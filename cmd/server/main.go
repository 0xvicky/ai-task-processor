package main

import (
	"ai-task-processor/internal/handler"
	"net/http"
)

func main() {
	println("AI-TASK-PROCESSOR")
	//Routes
	http.HandleFunc("/", handler.Handler)
	http.HandleFunc("/health", handler.Health)

	//Server
	println("Server running at 6969")
	http.ListenAndServe(":6969", nil)
}
