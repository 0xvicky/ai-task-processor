package middleware

import (
	"ai-task-processor/internal/constants"
	"ai-task-processor/internal/model"
	"ai-task-processor/internal/utils"
	"context"
	"fmt"
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
		fmt.Print(reqHeader)
		if reqHeader == "" || !strings.HasPrefix(reqHeader, "Bearer ") {
			//Error throw
			utils.WriteJsonResponse(w, 401, false, "Unauthorized", nil)
			return
		}

		tokenString := strings.TrimPrefix(reqHeader, "Bearer ")
		fmt.Print(tokenString)
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
		println("==============================")
		println("claims here\n")
		fmt.Print(claims)

		if parseErr != nil {
			utils.WriteJsonResponse(w, 401, false, "Unauthorized", nil)
			return
		}
		if token == nil || !token.Valid {
			utils.WriteJsonResponse(w, 401, false, "Unauthorized", nil)
			return
		}

		fmt.Print(claims.UserId)

		if token.Valid {
			//move ahead logic
			userId := claims.UserId
			ctx := context.WithValue(r.Context(), constants.UserKey, userId)
			next.ServeHTTP(w, r.WithContext(ctx))
		}

	})

	//passing user_id forward

}
