package utils

import (
	"ai-task-processor/internal/model"
	"encoding/json"
	"net/http"
)

func WriteJsonResponse(w http.ResponseWriter, statusCode int, status bool, message string, payload any) {
	apiRes := model.APIResponse{
		Success: status,
		Message: message,
		Payload: payload,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	encodeErr := json.NewEncoder(w).Encode(apiRes)
	if encodeErr != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}

}
