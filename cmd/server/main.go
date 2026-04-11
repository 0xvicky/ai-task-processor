package main

import (
	"ai-task-processor/internal/config"
	"ai-task-processor/internal/db"
	"ai-task-processor/internal/engine"
	"ai-task-processor/internal/handler"
	"ai-task-processor/internal/repository/task"
	"ai-task-processor/internal/repository/user"
	"ai-task-processor/internal/routes"
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

	//User services and handlers
	repo := user.NewPostgresUserRepo(dbConn)
	svc := service.NewUserService(repo)
	h := handler.NewUserHandler(svc)

	//Task services and handlers
	tr := task.NewPostgresTaskRepo(dbConn)
	ts := service.NewTaskService(tr)
	th := handler.NewTaskHandler(ts)

	//routes
	routes.UserRoutes(h)

	routes.TaskRoutes(th)

	//Worker Engine
	engine.StartEngine(tr, 5)
	//Server
	err := http.ListenAndServe(":6969", nil)

	if err != nil {
		log.Fatal(err)
	}

	// println("Server running at 6969")
}
