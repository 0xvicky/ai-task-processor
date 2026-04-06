package main

import (
	"ai-task-processor/internal/config"
	"ai-task-processor/internal/db"
	"ai-task-processor/internal/handler"
	m "ai-task-processor/internal/middleware"
	"ai-task-processor/internal/repository"
	"ai-task-processor/internal/service"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	println("AI-TASK-PROCESSOR")

	//load the env
	envErr := config.EnvInit()
	if envErr != nil {
		fmt.Println("Env load error")
		return
	}
	//Db intialization and connection
	dbConn, dbErr := db.Init()
	if dbErr != nil {
		log.Fatal("Db connection failed:%w", dbErr)
	}
	defer dbConn.Close()

	repo := repository.NewPostgresUserRepo(dbConn)
	svc := service.NewUserService(repo)
	h := handler.NewUserHandler(svc)

	http.Handle("/createuser", m.Public(http.HandlerFunc(h.CreateUserHandler))) //post req
	http.Handle("/login", m.Public(http.HandlerFunc(h.LoginUserHandler)))

	http.Handle("/update", m.Protected(http.HandlerFunc(h.UpdateUserHandler))) //PATCH req
	http.Handle("/delete", m.Protected(http.HandlerFunc(h.DeleteUserHandler))) //DELETE req
	http.Handle("/me", m.Protected(http.HandlerFunc(h.MeHandler)))
	http.Handle("/admin/users", m.Admin(http.HandlerFunc(h.FetchAllUsersHandler)))

	//Server
	err := http.ListenAndServe(":6969", nil)

	if err != nil {
		log.Fatal(err)
	}

	// println("Server running at 6969")
}
