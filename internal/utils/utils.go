package utils

import (
	"ai-task-processor/internal/model"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

// jwt secret key from .env
var secretKey = []byte(os.Getenv("JWT_SECRET"))

func JWTInit(userId int) (string, error) {
	// println("%s", secretKey)
	claims := jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(time.Hour * 6).Unix(),
		"iss":    "ai-task-processor",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, signErr := token.SignedString(secretKey)
	if signErr != nil {
		return "", fmt.Errorf("Error while signing token:%w", signErr)
	}

	return tokenString, nil
}
