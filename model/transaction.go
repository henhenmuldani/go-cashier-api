package model

import (
	"time"
)

type Transaction struct {
	ID          int                 `json:"id"`
	TotalAmount int                 `json:"total_amount"`
	CreatedAt   time.Time           `json:"created_at"`
	Details     []TransactionDetail `json:"details"`
}

type TransactionDetail struct {
	ID            int    `json:"id,omitempty"`
	TransactionID int    `json:"transaction_id"`
	ProductID     int    `json:"product_id"`
	ProductName   string `json:"product_name,omitempty"`
	Quantity      int    `json:"quantity"`
	Subtotal      int    `json:"subtotal"`
}

type TransactionResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Data    *Transaction `json:"data"`
}

type TransactionsResponse struct {
	Success            bool                `json:"success"`
	Message            string              `json:"message"`
	Data               []Transaction       `json:"data"`
	TotalTransactions  int                 `json:"total_transactions"`
	TotalRevenue       int                 `json:"total_revenue"`
	BestSellingProduct *BestSellingProduct `json:"best_selling_product"`
}

type CheckoutItem struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type CheckoutRequest struct {
	Items []CheckoutItem `json:"items"`
}

type BestSellingProduct struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	TotalSold int    `json:"total_sold"`
}
