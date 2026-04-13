package task

import (
	"ai-task-processor/internal/model"
	"context"
	"database/sql"
	"encoding/json"
	"log"
)

type PostgresTaskRepository struct {
	db *sql.DB
}

func NewPostgresTaskRepo(db *sql.DB) TaskRepository {
	return &PostgresTaskRepository{db: db}
}

func (t *PostgresTaskRepository) CreateTask(ctx context.Context, newTask model.Task) (int, error) {
	var newTaskId int
	newTaskQuery := `INSERT INTO tasks (user_id, task_type, payload) VALUES( $1, $2, $3) RETURNING task_id;`
	row := t.db.QueryRowContext(ctx, newTaskQuery, newTask.UserId, newTask.TaskType, newTask.Payload)
	rowErr := row.Scan(&newTaskId)
	if rowErr != nil {
		log.Print(rowErr)
		return 0, rowErr
	}

	log.Print(newTaskId)
	return newTaskId, nil
}

func (t *PostgresTaskRepository) GetTaskById(ctx context.Context, taskId int, userId int) (model.Task, error) {
	var task model.Task
	var resultBytes []byte
	fetchTaskQuery := `SELECT task_id,user_id,  task_type, status, result, error, created_at, updated_at from tasks WHERE task_id = $1 AND user_id = $2;`

	taskRow := t.db.QueryRowContext(ctx, fetchTaskQuery, taskId, userId)
	scanErr := taskRow.Scan(&task.TaskId, &task.UserId, &task.TaskType, &task.Status, &resultBytes, &task.Error, &task.CreatedAt, &task.UpdatedAt)

	if scanErr != nil {

		return model.Task{}, scanErr
	}

	if task.Result != nil {
		task.Result = json.RawMessage(resultBytes)
	} else {
		task.Result = json.RawMessage(`{}`)
	}

	return task, nil
}

func (t *PostgresTaskRepository) GetTasksByUser(ctx context.Context, userId int) ([]model.Task, error) {
	var tasks []model.Task
	// log.Print("userid:%w", userId)

	fetchTasksQuery := `SELECT task_id,user_id, task_type, status, result, error, created_at, updated_at from tasks WHERE user_id = $1;`

	taskRows, rowErr := t.db.QueryContext(ctx, fetchTasksQuery, userId)
	if rowErr != nil {
		return nil, rowErr
	}
	for taskRows.Next() {
		var newTask model.Task
		var resultBytes []byte
		scanErr := taskRows.Scan(&newTask.TaskId, &newTask.UserId, &newTask.TaskType, &newTask.Status, &resultBytes, &newTask.Error, &newTask.CreatedAt, &newTask.UpdatedAt)
		if newTask.Result != nil {
			newTask.Result = json.RawMessage(resultBytes)
		} else {
			newTask.Result = json.RawMessage(`{}`)
		}
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

func (t *PostgresTaskRepository) GetAllTasks(ctx context.Context) ([]model.Task, error) {
	var tasks []model.Task
	tasksQuery := `SELECT task_id,user_id, task_type, status, result, error, created_at, updated_at from tasks;`
	tasksRows, rowsErr := t.db.QueryContext(ctx, tasksQuery)
	if rowsErr != nil {
		return nil, rowsErr
	}

	if tasksRows.Next() {
		var newTask model.Task
		var resultBytes []byte
		scanErr := tasksRows.Scan(&newTask.TaskId, &newTask.UserId, &newTask.TaskType, &newTask.Status, &resultBytes, &newTask.Error, &newTask.CreatedAt, &newTask.UpdatedAt)
		if newTask.Result != nil {
			newTask.Result = json.RawMessage(resultBytes)
		} else {
			newTask.Result = json.RawMessage(`{}`)
		}
		if scanErr != nil {
			return nil, scanErr
		}

		tasks = append(tasks, newTask)
	}

	if tasksRows.Err() != nil {
		return nil, tasksRows.Err()
	}
	return tasks, nil
}

func (t *PostgresTaskRepository) FetchBatch() ([]model.BatchFetcher, error) {
	var batch []model.BatchFetcher
	batchQuery := `UPDATE tasks SET status='RUNNING' WHERE task_id IN(SELECT task_id from tasks WHERE status='PENDING' ORDER BY task_id LIMIT 5) RETURNING task_id, user_id, payload`
	//later we'll add the logic to fetch the tasks having difference between last updated and time started greater than threshodl
	rows, rowsErr := t.db.Query(batchQuery)
	if rowsErr != nil {
		return nil, rowsErr
	}

	if rows.Next() {

		var task model.BatchFetcher
		scanErr := rows.Scan(&task.TaskId, &task.UserId, &task.Payload)
		if scanErr != nil {
			return nil, scanErr
		}
		batch = append(batch, task)
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		return nil, rowsErr
	}

	return batch, nil
}
