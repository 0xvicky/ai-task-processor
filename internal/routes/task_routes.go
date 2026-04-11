package routes

import (
	"ai-task-processor/internal/handler"
	m "ai-task-processor/internal/middleware"
	"net/http"
)

func TaskRoutes(th *handler.TaskHandler) {
	http.Handle("/newTask", m.Protected(http.HandlerFunc(th.CreateTaskHandler)))
	http.Handle("/getTask", m.Protected(http.HandlerFunc(th.TaskByIdHandler)))
	http.Handle("/taskByUser", m.Protected(http.HandlerFunc(th.TasksByUserHandler)))
	http.Handle("/allTasks", m.Admin(http.HandlerFunc(th.AllTasksHandler)))
}
