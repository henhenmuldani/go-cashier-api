package handler

import (
	"encoding/json" // JSON parsing
	"net/http"      // HTTP operations

	"go-cashier-api/model"
	"go-cashier-api/pkg/response" // Alias the package
	"go-cashier-api/service"
)

type TransactionHandler struct {
	service service.TransactionService
}

func NewTransactionHandler(service service.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

func (h *TransactionHandler) HandleCheckout(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.Checkout(w, r)
	default:
		response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (h *TransactionHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	var request model.CheckoutRequest

	// Parse JSON from request body
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")

		return
	}

	// Validate using struct tags from model
	// if err := h.validate.Struct(request); err != nil {
	// 					response.Error(w, http.StatusBadRequest, err.Error())

	//     return
	// }

	// Call service layer
	responseData, err := h.service.Checkout(request)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Check if operation was successful
	if !responseData.Success {
		response.Error(w, http.StatusNotFound, responseData.Message)

		return
	}

	// Return 201 Created for successful creation
	response.JSON(w, http.StatusCreated, responseData)
}

func (h *TransactionHandler) GetTransactionsByDate(w http.ResponseWriter, r *http.Request) {
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	var responseData interface{}
	var err error

	if startDate != "" && endDate == "" || startDate == "" && endDate != "" {
		response.Error(w, http.StatusBadRequest, "start_date and end_date must be provided together")
		return
	}

	if startDate != "" && endDate != "" {
		responseData, err = h.service.GetTransactionsByDate(startDate, endDate)
	}
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())

		return
	}

	response.JSON(w, http.StatusOK, responseData)

}

func (h *TransactionHandler) GetTransactionsToday(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.GetTransactionsToday()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, data)

}
