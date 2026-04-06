package middleware

import (
	"ai-task-processor/internal/apperrors"
	"ai-task-processor/internal/constants"
	"ai-task-processor/internal/model"
	"ai-task-processor/internal/utils"
	"context"
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
	"slices"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// jwt secret key from .env
var secretKey = []byte(os.Getenv("JWT_SECRET"))

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//extracting token
		reqHeader := r.Header.Get("Authorization")
		// fmt.Print(reqHeader)
		if reqHeader == "" || !strings.HasPrefix(reqHeader, "Bearer ") {
			//Error throw
			utils.WriteJsonResponse(w, 401, false, "Unauthorized", nil)
			return
		}

		tokenString := strings.TrimPrefix(reqHeader, "Bearer ")
		// fmt.Print(tokenString)
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
		// println("==============================")
		// println("claims here\n")
		// fmt.Print(claims)

		if parseErr != nil {
			utils.WriteJsonResponse(w, 401, false, "Unauthorized", nil)
			return
		}
		if token == nil || !token.Valid {
			utils.WriteJsonResponse(w, 401, false, "Unauthorized", nil)
			return
		}

		// fmt.Print(claims.UserId)

		if token.Valid {
			//move ahead logic
			userId := claims.UserId
			userRole := claims.UserRole
			ctx := r.Context()
			ctx = context.WithValue(ctx, constants.UserKey, userId)
			ctx = context.WithValue(ctx, constants.UserRole, userRole)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})

	//passing user_id forward

}

func RoleMiddleware(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// allowedRoles := []string{"ADMIN", "MOD"}
			userRole, ok := r.Context().Value(constants.UserRole).(string)
			if !ok {
				utils.WriteJsonResponse(w, 401, false, "Unauthorized", nil)
				return
			}
			if !slices.Contains(allowedRoles, userRole) {
				utils.WriteJsonResponse(w, 403, false, "Forbidden", nil)
				return
			}

			next.ServeHTTP(w, r)
		})
	}

}

func RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("PANIC RECOVERED:", err)
				fmt.Println(string(debug.Stack())) //to know where panic occured
				utils.ErrorHandler(w, apperrors.ErrInternal)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
