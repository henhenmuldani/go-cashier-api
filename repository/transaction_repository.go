package repository

import (
	"database/sql"
	"fmt"
	"time"

	"go-cashier-api/model"
)

// Interface defines what methods the repository must implement
// This allows for dependency injection and easier testing
type TransactionRepository interface {
	CreateTransaction(items []model.CheckoutItem) (*model.Transaction, error)
	GetTransactionsByDate(startDate, endDate time.Time) ([]model.Transaction, int, int, *model.BestSellingProduct, error)
	getTransactionDetails(transactionId int) ([]model.TransactionDetail, error)
}

// Implementation of the interface
type TransactionRepositoryImpl struct {
	db *sql.DB // Database connection pool
}

// Constructor function - creates new instance of repository
func NewTransactionRepository(db *sql.DB) TransactionRepository {
	return &TransactionRepositoryImpl{db: db}
}

func (repo *TransactionRepositoryImpl) CreateTransaction(items []model.CheckoutItem) (*model.Transaction, error) {
	// Start a database transaction - ensures all operations succeed or fail together
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	// Defer ensures rollback happens if we don't reach commit()
	defer tx.Rollback()

	totalAmount := 0 // Initialize total price counter
	// Pre-allocate slice with capacity equal to number of items (for better performance)
	details := make([]model.TransactionDetail, 0, len(items))

	for _, item := range items {
		// Validate quantity is positive
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("invalid quantity for product id %d", item.ProductID)
		}

		var productPrice, stock int
		var productName string

		// Query product details with FOR UPDATE to lock row during transaction
		err := tx.QueryRow("SELECT name, price, stock FROM products WHERE id = $1", item.ProductID).Scan(&productName, &productPrice, &stock)

		// Handle cases where product doesn't exist
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product id %d not found", item.ProductID)
		}
		if err != nil {
			return nil, err
		}

		// Check if we have enough stock
		if stock < item.Quantity {
			return nil, fmt.Errorf("insufficient stock for product %s. Available: %d, Requested: %d",
				productName, stock, item.Quantity)
		}

		// Calculate subtotal for this item
		subtotal := productPrice * item.Quantity
		totalAmount += subtotal // Add to running total

		// Update product stock (decrease by purchased quantity)
		_, err = tx.Exec(`
            UPDATE products 
            SET stock = stock - $1 
            WHERE id = $2
        `, item.Quantity, item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("failed to update product stock: %w", err)
		}

		// Create transaction detail object (without database ID yet)
		details = append(details, model.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	var transactionID int
	var createdAt time.Time
	// Insert main transaction record and get auto-generated ID and timestamp
	err = tx.QueryRow(`
        INSERT INTO transactions (total_amount) 
        VALUES ($1) 
        RETURNING id, created_at
    `, totalAmount).Scan(&transactionID, &createdAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Insert each transaction detail into database
	for i := range details {
		details[i].TransactionID = transactionID // Set foreign key
		_, err = tx.Exec(`
            INSERT INTO transaction_details 
            (transaction_id, product_id, quantity, subtotal) 
            VALUES ($1, $2, $3, $4)
        `, transactionID, details[i].ProductID, details[i].Quantity, details[i].Subtotal)
		if err != nil {
			return nil, fmt.Errorf("failed to create transaction detail: %w", err)
		}
	}

	// Commit all changes to database - if successful, transaction is permanent
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Return complete transaction object
	return &model.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		CreatedAt:   createdAt,
		Details:     details,
	}, nil
}

func (repo *TransactionRepositoryImpl) GetTransactionsByDate(startDate, endDate time.Time) ([]model.Transaction, int, int, *model.BestSellingProduct, error) {
	rows, err := repo.db.Query(`
		SELECT id, total_amount, created_at
		FROM transactions
		WHERE created_at BETWEEN $1 AND $2
		ORDER BY created_at DESC
	`, startDate, endDate)

	if err != nil {
		return nil, 0, 0, nil, fmt.Errorf("failed to get transactions by date: %w", err)
	}
	defer rows.Close()

	var transactions []model.Transaction
	for rows.Next() {
		var transaction model.Transaction
		err := rows.Scan(&transaction.ID, &transaction.TotalAmount, &transaction.CreatedAt)
		if err != nil {
			return nil, 0, 0, nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		// 🔥 load details for each transaction
		transaction.Details, err = repo.getTransactionDetails(transaction.ID)
		if err != nil {
			return nil, 0, 0, nil, err
		}
		transactions = append(transactions, transaction)
	}

	// Get total transactions count for date range
	var totalTransactions int
	err = repo.db.QueryRow(`
		SELECT COUNT(*) 
		FROM transactions 
		WHERE created_at BETWEEN $1 AND $2
	`, startDate, endDate).Scan(&totalTransactions)

	if err != nil {
		return nil, 0, 0, nil, fmt.Errorf("failed to get total transactions count: %w", err)
	}

	// Get total revenue sum for date range
	var totalRevenue int
	err = repo.db.QueryRow(`
		SELECT COALESCE(SUM(total_amount), 0)
		FROM transactions 
		WHERE created_at BETWEEN $1 AND $2
	`, startDate, endDate).Scan(&totalRevenue)

	if err != nil {
		return nil, 0, 0, nil, fmt.Errorf("failed to get total revenue count: %w", err)
	}

	// Get best-selling product for date range
	var bestSellingProduct model.BestSellingProduct
	err = repo.db.QueryRow(`
		SELECT
			p.id,
			p.name,
			COALESCE(SUM(td.quantity), 0) AS total_sold
		FROM transaction_details td
		JOIN products p ON p.id = td.product_id
		JOIN transactions t ON t.id = td.transaction_id
		WHERE t.created_at BETWEEN $1 AND $2
		GROUP BY p.id, p.name
		ORDER BY total_sold DESC
		LIMIT 1
	`, startDate, endDate).Scan(&bestSellingProduct.ID,
		&bestSellingProduct.Name,
		&bestSellingProduct.TotalSold)

	if err != nil {
		if err == sql.ErrNoRows {
			// No sales in date range → return nil instead of error
			return transactions, totalTransactions, totalRevenue, nil, nil
		}
		return nil, 0, 0, nil, fmt.Errorf("failed to get best selling product: %w", err)
	}

	return transactions, totalTransactions, totalRevenue, &bestSellingProduct, nil
}

func (repo *TransactionRepositoryImpl) getTransactionDetails(transactionID int) ([]model.TransactionDetail, error) {
	// Get transaction details
	rows, err := repo.db.Query(`
		SELECT td.id, td.product_id, p.name, td.quantity, td.subtotal
		FROM transaction_details td
		JOIN products p ON td.product_id = p.id
		WHERE td.transaction_id = $1
	`, transactionID)

	if err != nil {
		return nil, fmt.Errorf("failed to get transaction details: %w", err)
	}
	defer rows.Close()

	var details []model.TransactionDetail
	for rows.Next() {
		var detail model.TransactionDetail
		err := rows.Scan(&detail.ID, &detail.ProductID, &detail.ProductName,
			&detail.Quantity, &detail.Subtotal)
		if err != nil {
			return nil, fmt.Errorf("failed to scan detail: %w", err)
		}
		detail.TransactionID = transactionID
		details = append(details, detail)
	}

	return details, nil

}
