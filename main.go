package main

import (
	"log"
	"net/http"

	"go-cashier-api/handler"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", handler.HealthHandler)
	mux.HandleFunc("/api/products", handler.ProductHandler)
	mux.HandleFunc("/api/products/", handler.ProductHandler)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
