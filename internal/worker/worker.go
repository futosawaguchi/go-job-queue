package worker

import (
	"fmt"
	"sync"

	"github.com/futosawaguchi/go-job-queue/internal/job"
)

// WorkerPoolの構造体
type WorkerPool struct {
	workerCount int
	jobQueue    chan job.Job
	wg          sync.WaitGroup
}

// WorkerPoolを作成する関数
func NewWorkerPool(workerCount int) *WorkerPool {
	return &WorkerPool{
		workerCount: workerCount,
		jobQueue:    make(chan job.Job, 100),
	}
}

// WorkerPoolを起動する
func (wp *WorkerPool) Start() {
	for i := 0; i < wp.workerCount; i++ {
		wp.wg.Add(1)
		go wp.runWorker(i)
	}
}

// 個々のWorkerの処理
func (wp *WorkerPool) runWorker(id int) {
	defer wp.wg.Done()
	for j := range wp.jobQueue {
		fmt.Printf("Worker %d: Job %s を処理中...\n", id, j.ID)
		// ここに実際の処理が入る
		fmt.Printf("Worker %d: Job %s 完了!\n", id, j.ID)
	}
}

// JobをQueueに追加する
func (wp *WorkerPool) Submit(j job.Job) {
	wp.jobQueue <- j
}

// 全Jobの完了を待って停止する
func (wp *WorkerPool) Stop() {
	close(wp.jobQueue)
	wp.wg.Wait()
}
