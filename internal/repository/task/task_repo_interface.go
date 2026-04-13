package task

import (
	"ai-task-processor/internal/model"
	"context"
)

type TaskRepository interface {
	CreateTask(ctx context.Context, newTask model.Task) (int, error)
	GetTaskById(ctx context.Context, taskId int, userId int) (model.Task, error)
	GetTasksByUser(ctx context.Context, userId int) ([]model.Task, error)
	UpdateTaskStatus(ctx context.Context, taskId int, taskStatus model.TaskStatus) (bool, error)
	GetAllTasks(ctx context.Context) ([]model.Task, error)
	FetchBatch() ([]model.BatchFetcher, error)
}
