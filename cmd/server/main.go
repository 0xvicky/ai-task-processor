package main

import (
	"ai-task-processor/internal/db"
	"ai-task-processor/internal/handler"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	println("AI-TASK-PROCESSOR")

	//Connect with DB
	db.Init()
	defer db.Db.Close()
	//Register Routes
	http.HandleFunc("/", handler.RootHandler)
	http.HandleFunc("/health", handler.Health)
	http.HandleFunc("/createuser", handler.CreateUserHandler) //post req

	//Server
	err := http.ListenAndServe(":6969", nil)

	if err != nil {
		log.Fatal(err)
	}

	// println("Server running at 6969")
}
