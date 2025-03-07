package utils

import (
	"encoding/json"
	"net/http"
)

func SendErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	responseData := map[string]string{
		"message": message,
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(responseData)
}
