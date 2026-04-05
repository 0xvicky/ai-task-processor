package main

import (
	"ai-task-processor/internal/config"
	"ai-task-processor/internal/db"
	"ai-task-processor/internal/handler"
	"ai-task-processor/internal/middleware"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	println("AI-TASK-PROCESSOR")

	envErr := config.EnvInit()
	if envErr != nil {
		fmt.Println("Env load error")
		return
	}
	//Connect with DB
	db.Init()
	defer db.Db.Close()
	//Register Routes
	http.HandleFunc("/", handler.RootHandler)
	http.HandleFunc("/health", handler.Health)
	http.HandleFunc("/createuser", handler.CreateUserHandler) //post req
	http.HandleFunc("/login", handler.LoginUserHandler)
	http.Handle("/update", middleware.AuthMiddleware(http.HandlerFunc(handler.UpdateUserHandler))) //PATCH req
	http.Handle("/delete", middleware.AuthMiddleware(http.HandlerFunc(handler.DeleteUserHandler))) //DELETE req
	http.Handle("/me", middleware.AuthMiddleware(http.HandlerFunc(handler.MeHandler)))
	http.Handle("/admin/users", middleware.AuthMiddleware(middleware.RoleMiddleware("ADMIN")(http.HandlerFunc(handler.FetchAllUsersHandler))))

	//Server
	err := http.ListenAndServe(":6969", nil)

	if err != nil {
		log.Fatal(err)
	}

	// println("Server running at 6969")
}
