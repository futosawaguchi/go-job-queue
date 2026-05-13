package worker

import (
	"fmt"
	"sync"
	"time"

	"github.com/futosawaguchi/go-job-queue/db"
	"github.com/futosawaguchi/go-job-queue/internal/job"
)

type WorkerPool struct {
	workerCount int
	jobQueue    chan job.Job
	wg          sync.WaitGroup
	db          *db.DB
}

// dbを受け取るように変更
func NewWorkerPool(workerCount int, db *db.DB) *WorkerPool {
	return &WorkerPool{
		workerCount: workerCount,
		jobQueue:    make(chan job.Job, 100),
		db:          db,
	}
}

func (wp *WorkerPool) Start() {
	for i := 0; i < wp.workerCount; i++ {
		wp.wg.Add(1)
		go wp.runWorker(i)
	}
}

func (wp *WorkerPool) runWorker(id int) {
	defer wp.wg.Done()
	for j := range wp.jobQueue {
		fmt.Printf("Worker %d: Job %s を処理中...\n", id, j.ID)

		// 処理中ステータスに更新
		wp.db.UpdateJobStatus(j.ID, job.StatusRunning)

		// ここに実際の処理が入る
		time.Sleep(5 * time.Second)

		// 完了ステータスに更新
		wp.db.UpdateJobStatus(j.ID, job.StatusCompleted)

		fmt.Printf("Worker %d: Job %s 完了!\n", id, j.ID)
	}
}

func (wp *WorkerPool) Submit(j job.Job) {
	wp.jobQueue <- j
}

func (wp *WorkerPool) Stop() {
	close(wp.jobQueue)
	wp.wg.Wait()
}
