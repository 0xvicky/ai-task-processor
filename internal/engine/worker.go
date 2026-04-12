package engine

import (
	"ai-task-processor/internal/repository/task"
	"sync"
)

func StartEngine(tRepo task.TaskRepository, nWorkers int) {
	//spawn the workers
	wgu * sync.WaitGroup

	for i := 0; i < nWorkers; i++ {

	}
	//each worker polls db to check for pending tasks

	//if found then safely lock it and update the state to running and if not found, sleep

	//process the task

	//store results in db of that particular task and update the status

}
