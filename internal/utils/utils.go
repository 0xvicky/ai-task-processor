package utils

import (
	"ai-task-processor/internal/apperrors"
	"ai-task-processor/internal/constants"
	"ai-task-processor/internal/model"
	"context"
	"encoding/json"
	"errors"
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

func ErrorHandler(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, apperrors.ErrBadRequest):
		WriteJsonResponse(w, 400, false, err.Error(), nil)
	case errors.Is(err, apperrors.ErrInternal):
		WriteJsonResponse(w, 401, false, err.Error(), nil)
	case errors.Is(err, apperrors.ErrUnauthorized):
		WriteJsonResponse(w, 401, false, err.Error(), nil)
	case errors.Is(err, apperrors.ErrUserAlreadyExists):
		WriteJsonResponse(w, 403, false, err.Error(), nil)
	case errors.Is(err, apperrors.ErrCanceled):
		WriteJsonResponse(w, 404, false, err.Error(), nil)
	case errors.Is(err, apperrors.ErrForbidden):
		WriteJsonResponse(w, 409, false, err.Error(), nil)
	case errors.Is(err, apperrors.ErrInvalidCredentials):
		WriteJsonResponse(w, 499, false, err.Error(), nil)
	case errors.Is(err, apperrors.ErrUserNotFound):
		WriteJsonResponse(w, 500, false, err.Error(), nil)
	case errors.Is(err, apperrors.ErrDeadlineExceeded):
		WriteJsonResponse(w, 504, false, err.Error(), nil)
	default:
		WriteJsonResponse(w, 500, false, "internal server error", nil)

	}
}

// jwt secret key from .env
var secretKey = []byte(os.Getenv("JWT_SECRET"))

func JWTInit(userId int, userRole string) (string, error) {
	// println("%s", secretKey)
	claims := jwt.MapClaims{
		"userId":   userId,
		"userRole": userRole,
		"exp":      time.Now().Add(time.Hour * 6).Unix(),
		"iss":      "ai-task-processor",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, signErr := token.SignedString(secretKey)
	if signErr != nil {
		return "", fmt.Errorf("Error while signing token:%w", signErr)
	}

	return tokenString, nil
}

func ExtractUserId(ctx context.Context) int {
	userId, ok := ctx.Value(constants.UserKey).(int)
	// fmt.Print(userId)
	if !ok {
		return 0
	}
	return userId
}
