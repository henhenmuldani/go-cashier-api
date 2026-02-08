package service

import (
	"errors"
	"fmt"
	"time"

	"go-cashier-api/model"
	"go-cashier-api/repository"
)

// Business logic interface
type TransactionService interface {
	Checkout(request model.CheckoutRequest) (*model.TransactionResponse, error)
	GetTransactionsByDate(startDateStr, endDateStr string) (*model.TransactionsResponse, error)
	GetTransactionsToday() (*model.TransactionsResponse, error)
}

// Service implementation with dependencies
type TransactionServiceImpl struct {
	repo        repository.TransactionRepository // Transaction operations
	productRepo repository.ProductRepository     // Product operations
}

// Constructor with dependency injection
func NewTransactionService(repo repository.TransactionRepository,
	productRepo repository.ProductRepository) TransactionService {
	return &TransactionServiceImpl{
		repo:        repo,
		productRepo: productRepo,
	}
}

func (s *TransactionServiceImpl) Checkout(request model.CheckoutRequest) (*model.TransactionResponse, error) {
	// Validate request has at least one item
	if len(request.Items) == 0 {
		return nil, fmt.Errorf("items cannot be empty")
	}

	// Validate each item
	for _, item := range request.Items {
		if item.ProductID <= 0 {
			return nil, fmt.Errorf("invalid product id")
		}
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("quantity must be greater than 0")
		}
	}

	// Call repository to create transaction
	transaction, err := s.repo.CreateTransaction(request.Items)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Create success response
	response := &model.TransactionResponse{
		Success: true,
		Message: "Transaction created successfully",
		Data:    transaction,
	}

	return response, nil
}

func (s *TransactionServiceImpl) GetTransactionsByDate(startDateStr, endDateStr string) (*model.TransactionsResponse, error) {
	startDate, err := time.Parse("2026-01-31", startDateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid start date format. Use YYYY-MM-DD")
	}

	endDate, err := time.Parse("2026-01-31", endDateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid end date format. Use YYYY-MM-DD")
	}

	// Add one day to end date to include the entire day
	endDate = endDate.Add(24 * time.Hour)

	transactions, totalTransactions, totalRevenue, bestSellingProduct, err := s.repo.GetTransactionsByDate(startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions by date: %w", err)
	}

	return &model.TransactionsResponse{
		Success:            true,
		Message:            "Transactions retrieved successfully",
		Data:               transactions,
		TotalTransactions:  totalTransactions,
		TotalRevenue:       totalRevenue,
		BestSellingProduct: bestSellingProduct,
	}, nil
}

func (s *TransactionServiceImpl) GetTransactionsToday() (*model.TransactionsResponse, error) {
	now := time.Now().UTC()

	today := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		0, 0, 0, 0,
		now.Location(),
	)

	// Add one day to end date to include the entire day
	endDate := today.Add(24 * time.Hour)

	transactions, totalTransactions, totalRevenue, bestSellingProduct, err := s.repo.GetTransactionsByDate(today, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions by date: %w", err)
	}

	if transactions == nil {
		return nil, errors.New("Transactions not found")
	}

	return &model.TransactionsResponse{
		Success:            true,
		Message:            "Transactions retrieved successfully",
		Data:               transactions,
		TotalTransactions:  totalTransactions,
		TotalRevenue:       totalRevenue,
		BestSellingProduct: bestSellingProduct,
	}, nil
}
