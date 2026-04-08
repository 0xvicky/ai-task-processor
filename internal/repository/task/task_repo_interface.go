package task

import (
	"ai-task-processor/internal/model"
	"context"
)

type TaskRepository interface {
	CreateTask(ctx context.Context, newTaskDetail model.Task) (int, error)
	GetTaskById(ctx context.Context, taskId int) (model.Task, error)
	GetTasksByUser(ctx context.Context, userId int) ([]model.Task, error)
	UpdateTaskStatus(ctx context.Context, taskStatus string) (string, error)
}
