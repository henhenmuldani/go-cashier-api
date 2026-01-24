package handler

import (
	"encoding/json" //Encode/decode JSON  API response
	"net/http"      //HTTP server & request handling
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	// Set response header to JSON
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{
		"status":  "OK",
		"message": "API Running",
	})
}
