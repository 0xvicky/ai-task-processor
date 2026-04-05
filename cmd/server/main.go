package main

import (
	"ai-task-processor/internal/config"
	"ai-task-processor/internal/db"
	"ai-task-processor/internal/handler"
	"ai-task-processor/internal/middleware"
	"ai-task-processor/internal/repository"
	"ai-task-processor/internal/service"
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

	repo := repository.NewPostgresUserRepo(db.Db)
	svc := service.NewUserService(repo)
	h := handler.NewUserHandler(svc)

	http.HandleFunc("/createuser", h.CreateUserHandler) //post req
	http.HandleFunc("/login", h.LoginUserHandler)
	http.Handle("/update", middleware.AuthMiddleware(http.HandlerFunc(h.UpdateUserHandler))) //PATCH req
	http.Handle("/delete", middleware.AuthMiddleware(http.HandlerFunc(h.DeleteUserHandler))) //DELETE req
	http.Handle("/me", middleware.AuthMiddleware(http.HandlerFunc(h.MeHandler)))
	http.Handle("/admin/users", middleware.AuthMiddleware(middleware.RoleMiddleware("ADMIN")(http.HandlerFunc(h.FetchAllUsersHandler))))

	//Server
	err := http.ListenAndServe(":6969", nil)

	if err != nil {
		log.Fatal(err)
	}

	// println("Server running at 6969")
}
