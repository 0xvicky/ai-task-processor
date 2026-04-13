package engine

import (
	"ai-task-processor/internal/model"
	"ai-task-processor/internal/service"
	"encoding/json"
	"log"
	"sync"
	"time"
)

func StartEngine(tps *service.TaskProcessingService, nWorkers int, wg *sync.WaitGroup) {
	//create a channel
	taskChannel := make(chan model.BatchFetcher)
	//spawn the workers
	for i := 1; i <= nWorkers; i++ {
		wg.Add(1)
		go workerExec(i, wg, taskChannel, tps)
	}

	//run the batch fetcher
	go batchFetcherInit(tps, taskChannel)

}

//each worker polls db to check for pending tasks

func workerExec(workerId int, wg *sync.WaitGroup, workerQueue <-chan model.BatchFetcher, tps *service.TaskProcessingService) {
	defer wg.Done()
	log.Printf("Tracking the channel and executing with worker:%d", workerId)
	//if found then safely lock it and update the state to running and if not found, sleep

	//track the channel
	for task := range workerQueue {
		//if task found, execute
		// taskPayload := task.Payload
		log.Println("Picked task no:%w", task.TaskId)
		taskId := task.TaskId
		//simulate the processing
		time.Sleep(2 * time.Second)
		rawResult := `{"result":"This is result for the task payload"}`
		result := json.RawMessage(rawResult)

		var taskUpdate = model.TaskUpdate{
			TaskId:    taskId,
			Result:    &result,
			Error:     nil,
			UpdatedAt: time.Now(),
		}

		taskUpdate.Status = model.StatusCompleted

		res, updateErr := tps.UpdateTask(taskUpdate)
		if updateErr != nil {
			println("Failed to update task")
			continue
		}
		println("Task updated:%w", res)

	}

	//process the task
	//store results in db of that particular task and update the status
}

func batchFetcherInit(tps *service.TaskProcessingService, workerQueue chan<- model.BatchFetcher) {
	for {
		batch, batchErr := tps.FetchBatch()
		if batchErr != nil {
			log.Println("Batch error")
			continue
		}
		if len(batch) <= 0 {
			log.Println("no pending tasks present")
			time.Sleep(5 * time.Second)
			continue
		}
		for _, task := range batch {
			log.Println("Put Task:%w in the queue", task.TaskId)
			workerQueue <- task
		}
	}

}
