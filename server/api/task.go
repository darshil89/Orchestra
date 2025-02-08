package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	models "server/models"
)

func responseWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func responseWithError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	var data models.Tasks
	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Invalid request payload")
		return
	}

	// print the request body
	fmt.Printf("Request Body: %v\n", data)

	responseWithJSON(w, http.StatusOK, "Task created successfully")

}
