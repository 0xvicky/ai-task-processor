package routes

import (
	"ai-task-processor/internal/handler"
	m "ai-task-processor/internal/middleware"
	"net/http"
)

func TaskRoutes(th *handler.TaskHandler) {
	http.Handle("/newTask", m.Protected(http.HandlerFunc(th.CreateTaskHandler)))
	http.Handle("/getTask", m.Protected(http.HandlerFunc(th.GetTaskById)))
	http.Handle("/getAllTasks", m.Protected(http.HandlerFunc(th.GetTasksByUser)))
}
