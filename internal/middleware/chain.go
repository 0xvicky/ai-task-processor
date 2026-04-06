package middleware

import (
	"net/http"
)

func Public(h http.Handler) http.Handler {
	return RecoverMiddleware(h)
}

func Protected(h http.Handler) http.Handler {
	return RecoverMiddleware(AuthMiddleware(h))
}

func Admin(h http.Handler) http.Handler {
	return RecoverMiddleware(AuthMiddleware(RoleMiddleware("ADMIN")(h)))
}
