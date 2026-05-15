package worker

import (
	"fmt"
	"sync"
	"time"

	"github.com/futosawaguchi/go-job-queue/internal/job"
)

// DBのインターフェース（本物とモックを切り替えられる）
type JobDB interface {
	UpdateJobStatus(id string, status job.Status) error
}

type WorkerPool struct {
	workerCount int
	jobQueue    chan job.Job
	wg          sync.WaitGroup
	db          JobDB
	processor   func(job.Job)
}

func NewWorkerPool(workerCount int, db JobDB) *WorkerPool {
	return &WorkerPool{
		workerCount: workerCount,
		jobQueue:    make(chan job.Job, 100),
		db:          db,
	}
}

// テスト用：DBなしで作成
func NewWorkerPoolWithDB(workerCount int, db JobDB) *WorkerPool {
	return NewWorkerPool(workerCount, db)
}

// テスト用：カスタム処理を設定
func (wp *WorkerPool) SetProcessor(fn func(job.Job)) {
	wp.processor = fn
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
		wp.db.UpdateJobStatus(j.ID, job.StatusRunning)

		// カスタム処理があれば実行
		if wp.processor != nil {
			wp.processor(j)
		}
		time.Sleep(3 * time.Second)

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
