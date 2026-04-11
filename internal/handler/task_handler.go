package handler

import (
	"ai-task-processor/internal/model"
	"ai-task-processor/internal/service"
	"ai-task-processor/internal/utils"
	"context"
	"encoding/json"
	"log"
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
	defer r.Body.Close()
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	var newTask model.CreateTask
	decoder := json.NewDecoder(r.Body)
	decodeErr := decoder.Decode(&newTask)
	if decodeErr != nil {
		log.Print("Error in handler here")
		utils.ErrorHandler(w, decodeErr)
		return
	}
	log.Print(newTask)
	userId, userIdErr := utils.ExtractUserId(r)
	if userIdErr != nil {
		utils.ErrorHandler(w, userIdErr)
		return
	}
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

func (h *TaskHandler) TaskByIdHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	userId, userIdErr := utils.ExtractUserId(r)
	if userIdErr != nil {
		utils.ErrorHandler(w, userIdErr)
		return
	}
	var req struct {
		TaskId int `json:"taskId"`
	}
	decoder := json.NewDecoder(r.Body)
	decodeErr := decoder.Decode(&req)
	if decodeErr != nil {

		utils.ErrorHandler(w, decodeErr)
		return
	}

	task, taskErr := h.service.GetTaskById(ctx, req.TaskId, userId)
	if taskErr != nil {
		utils.ErrorHandler(w, taskErr)
		return
	}
	utils.WriteJsonResponse(w, 200, true, "task fetched", task)
}

func (h *TaskHandler) TasksByUserHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	userId, userIdErr := utils.ExtractUserId(r)
	if userIdErr != nil {
		utils.ErrorHandler(w, userIdErr)
		return
	}

	tasks, tasksErr := h.service.GetTasksByUser(ctx, userId)
	if tasksErr != nil {
		utils.ErrorHandler(w, tasksErr)
		return
	}

	utils.WriteJsonResponse(w, 200, true, "tasks fetched", tasks)

}

func (h *TaskHandler) AllTasksHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	tasks, taskErr := h.service.GetAllTasks(ctx)
	if taskErr != nil {
		utils.ErrorHandler(w, taskErr)
		return
	}

	utils.WriteJsonResponse(w, 200, true, "all tasks fetched", tasks)

}
