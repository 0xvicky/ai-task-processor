package task

import (
	"ai-task-processor/internal/model"
	"context"
	"database/sql"
)

type PostgresTaskRepository struct {
	db *sql.DB
}

func NewPostgresTaskRepo(db *sql.DB) TaskRepository {
	return &PostgresTaskRepository{db: db}
}

func (t *PostgresTaskRepository) CreateTask(ctx context.Context, newTaskDetail model.Task) (int, error) {
	return 0, nil
}

func (t *PostgresTaskRepository) GetTaskById(ctx context.Context, taskId int) (model.Task, error) {
	return model.Task{}, nil
}

func (t *PostgresTaskRepository) GetTasksByUser(ctx context.Context, userId int) ([]model.Task, error) {
	return nil, nil
}

func (t *PostgresTaskRepository) UpdateTaskStatus(ctx context.Context, taskStatus string) (string, error) {
	return "0", nil
}
