package service

import (
	"ai-task-processor/internal/apperrors"
	"ai-task-processor/internal/model"
	"ai-task-processor/internal/repository/task"
	"context"
	"database/sql"
	"errors"
)

type TaskService struct {
	taskRepo task.TaskRepository
}

func NewTaskService(taskRepo task.TaskRepository) *TaskService {
	return &TaskService{
		taskRepo: taskRepo,
	}
}

func (t *TaskService) CreateTask(ctx context.Context, newTask model.CreateTask, userId int) (int, error) {

	if userId == 0 {
		return 0, apperrors.ErrUnauthorized
	}
	if newTask.TaskType == "" {
		return 0, apperrors.ErrBadRequest
	}
	task := model.Task{
		UserId:   userId,
		TaskType: newTask.TaskType,
		Payload:  newTask.Payload,
	}

	taskId, createTaskErr := t.taskRepo.CreateTask(ctx, task)
	if createTaskErr != nil {
		//add context deadline and cancelled errors
		return 0, apperrors.ErrInternal
	}

	return taskId, nil
}

func (t *TaskService) GetTaskById(ctx context.Context, taskId int, userId int) (model.Task, error) {
	task, fetchErr := t.taskRepo.GetTaskById(ctx, taskId, userId)
	if fetchErr != nil {
		//context errors

		if errors.Is(fetchErr, sql.ErrNoRows) {
			return model.Task{}, apperrors.ErrTaskNotFound
		}
		return model.Task{}, apperrors.ErrInternal
	}

	return task, nil
}

func (t *TaskService) GetTasksByUser(ctx context.Context, userId int) ([]model.Task, error) {
	tasks, tasksErr := t.taskRepo.GetTasksByUser(ctx, userId)
	if tasksErr != nil {
		//context errors

		return nil, apperrors.ErrInternal
	}

	return tasks, nil

}
