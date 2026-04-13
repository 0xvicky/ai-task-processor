package service

import (
	"ai-task-processor/internal/model"
	"ai-task-processor/internal/repository/task"
)

type TaskProcessingService struct {
	taskRepo task.TaskRepository
}

func NewTaskProcessing(taskRepo task.TaskRepository) *TaskProcessingService {
	return &TaskProcessingService{
		taskRepo: taskRepo,
	}
}

func (t *TaskProcessingService) FetchBatch() ([]model.BatchFetcher, error) {
	batch, batchErr := t.taskRepo.FetchBatch()

	if batchErr != nil {
		return nil, batchErr
	}
	return batch, nil
}
