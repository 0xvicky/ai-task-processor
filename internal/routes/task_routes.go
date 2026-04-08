package routes

import (
	"ai-task-processor/internal/handler"
	m "ai-task-processor/internal/middleware"
	"net/http"
)

func TaskRoutes(th *handler.TaskHandler) {
	http.Handle("/newtask", m.Public(http.HandlerFunc(th.CreateTaskHandler)))
}
