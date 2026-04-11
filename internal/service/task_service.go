package service

import (
	"ai-task-processor/internal/apperrors"
	"ai-task-processor/internal/model"
	"ai-task-processor/internal/repository/task"
	"context"
	"database/sql"
	"errors"
	"log"
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
		if errors.Is(createTaskErr, context.DeadlineExceeded) {
			return 0, apperrors.ErrDeadlineExceeded
		}
		if errors.Is(createTaskErr, context.Canceled) {
			return 0, apperrors.ErrCanceled
		}
		log.Print("Error here")
		return 0, apperrors.ErrInternal
	}
	log.Print(taskId)
	return taskId, nil
}

func (t *TaskService) GetTaskById(ctx context.Context, taskId int, userId int) (model.Task, error) {
	task, fetchErr := t.taskRepo.GetTaskById(ctx, taskId, userId)
	if fetchErr != nil {
		//context errors
		if errors.Is(fetchErr, context.DeadlineExceeded) {
			return model.Task{}, apperrors.ErrDeadlineExceeded
		}
		if errors.Is(fetchErr, context.Canceled) {
			return model.Task{}, apperrors.ErrCanceled
		}
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
		if errors.Is(tasksErr, context.DeadlineExceeded) {
			return nil, apperrors.ErrDeadlineExceeded
		}
		if errors.Is(tasksErr, context.Canceled) {
			return nil, apperrors.ErrCanceled
		}
		return nil, apperrors.ErrInternal
	}

	return tasks, nil

}
