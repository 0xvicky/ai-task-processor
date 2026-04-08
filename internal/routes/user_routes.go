package routes

import (
	"ai-task-processor/internal/handler"
	m "ai-task-processor/internal/middleware"
	"net/http"
)

//User Handlers

func UserRoutes(h *handler.UserHandler) {

	http.Handle("/createuser", m.Public(http.HandlerFunc(h.CreateUserHandler))) //post req
	http.Handle("/login", m.Public(http.HandlerFunc(h.LoginUserHandler)))

	http.Handle("/update", m.Protected(http.HandlerFunc(h.UpdateUserHandler))) //PATCH req
	http.Handle("/delete", m.Protected(http.HandlerFunc(h.DeleteUserHandler))) //DELETE req
	http.Handle("/me", m.Protected(http.HandlerFunc(h.MeHandler)))
	http.Handle("/admin/users", m.Admin(http.HandlerFunc(h.FetchAllUsersHandler)))
}
