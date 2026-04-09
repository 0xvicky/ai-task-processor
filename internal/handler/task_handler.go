package handler

import (
	"ai-task-processor/internal/model"
	"ai-task-processor/internal/service"
	"ai-task-processor/internal/utils"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type TaskHandler struct {
	service *service.TaskService
}

func NewTaskHandler(svc *service.TaskService) *TaskHandler {
	return &TaskHandler{
		service: svc,
	}
}

func (h *TaskHandler) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	r.Body.Close()
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	var newTask model.CreateTask
	decoder := json.NewDecoder(r.Body)
	decodeErr := decoder.Decode(&newTask)
	if decodeErr != nil {
		utils.ErrorHandler(w, decodeErr)
		return
	}

	userId := utils.ExtractUserId(r)
	taskId, createTaskErr := h.service.CreateTask(ctx, newTask, userId)

	if createTaskErr != nil {
		utils.ErrorHandler(w, createTaskErr)
		return
	}
	newTaskPayload := struct {
		TaskId int `json:"taskId"`
	}{
		TaskId: taskId,
	}

	utils.WriteJsonResponse(w, 201, true, "task created successfully", newTaskPayload)

}

func (h *TaskHandler) GetTaskById(w http.ResponseWriter, r *http.Request) {

}
