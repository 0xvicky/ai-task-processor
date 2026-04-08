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

func (t *PostgresTaskRepository) CreateTask(ctx context.Context, newTask model.Task) (int, error) {
	var newTaskId int
	newTaskQuery := `INSERT INTO tasks (user_id, task_type, status, payload) VALUES( $1, $2, $3, $4) RETURNING task_id;`
	row := t.db.QueryRowContext(ctx, newTaskQuery, newTask.TaskType, newTask.Payload)
	rowErr := row.Scan(&newTaskId)
	if rowErr != nil {
		return 0, rowErr
	}
	return newTaskId, nil
}

func (t *PostgresTaskRepository) GetTaskById(ctx context.Context, taskId int, userId int) (model.Task, error) {
	var task model.Task
	fetchTaskQuery := `SELECT task_id,user_id,  task_type, status, result, error, created_at, updated_at from tasks WHERE task_id = $1 AND user_id = $2;`

	taskRow := t.db.QueryRowContext(ctx, fetchTaskQuery, taskId, userId)
	scanErr := taskRow.Scan(&task.TaskId, &task.UserId, &task.TaskType, &task.Status, &task.Result, &task.Error, &task.CreatedAt, &task.UpdatedAt)

	if scanErr != nil {
		return model.Task{}, scanErr
	}
	return task, nil
}

func (t *PostgresTaskRepository) GetTasksByUser(ctx context.Context, userId int) ([]model.Task, error) {
	var tasks []model.Task

	fetchTasksQuery := `SELECT task_id, task_type, status, result, error, created_at, updated_at from tasks WHERE user_id = $1;`

	taskRows, rowErr := t.db.QueryContext(ctx, fetchTasksQuery, userId)
	if rowErr != nil {
		return nil, rowErr
	}
	for taskRows.Next() {
		var newTask model.Task
		scanErr := taskRows.Scan(&newTask.TaskId, &newTask.TaskType, &newTask.Status, &newTask.Result, &newTask.Error, &newTask.CreatedAt, &newTask.UpdatedAt)
		if scanErr != nil {
			return nil, scanErr
		}

		tasks = append(tasks, newTask)
	}
	if rowsErr := taskRows.Err(); rowsErr != nil {
		return nil, rowsErr
	}
	return tasks, nil
}

func (t *PostgresTaskRepository) UpdateTaskStatus(ctx context.Context, taskId int, taskStatus model.TaskStatus) (bool, error) {
	updateQuery := `UPDATE tasks SET status = $1 WHERE task_id=$2;`
	result, err := t.db.ExecContext(ctx, updateQuery, taskStatus, taskId)
	if err != nil {
		return false, err
	}

	rowAff, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	if rowAff == 0 {
		return false, sql.ErrNoRows
	}

	return true, nil
}
