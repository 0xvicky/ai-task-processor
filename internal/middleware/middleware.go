package middleware

import (
	"ai-task-processor/internal/model"
	"ai-task-processor/internal/utils"
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// jwt secret key from .env
var secretKey = []byte(os.Getenv("JWT_SECRET"))

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//extracting token
		reqHeader := r.Header.Get("Authorization")
		if reqHeader == "" || !strings.HasPrefix(reqHeader, "Bearer ") {
			//Error throw
			utils.WriteJsonResponse(w, 401, false, "Unauthorized", nil)
			return
		}

		tokenString := strings.TrimPrefix(reqHeader, "Bearer ")
		if tokenString == "" {
			utils.WriteJsonResponse(w, 401, false, "Unauthorized", nil)
			return
		}

		claims := &model.JwtClaims{}
		//validating token

		token, parseErr := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
			// Return the secret key for signature verification
			return secretKey, nil
		})

		if parseErr != nil {
			utils.WriteJsonResponse(w, 401, false, "Unauthorized", nil)
			return
		}
		if token == nil || !token.Valid {
			utils.WriteJsonResponse(w, 401, false, "Unauthorized", nil)
			return
		}

		type contextKey string
		const userKey contextKey = "userId"

		if token.Valid {
			//move ahead logic
			userId := claims.UserId
			ctx := context.WithValue(r.Context(), userKey, userId)
			next.ServeHTTP(w, r.WithContext(ctx))
		}

	})

	//passing user_id forward

}
