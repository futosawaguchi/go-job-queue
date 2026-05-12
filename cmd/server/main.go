package main

import (
	"fmt"
	"time"

	"github.com/futosawaguchi/go-job-queue/internal/job"
	"github.com/futosawaguchi/go-job-queue/internal/worker"
)

func main() {
	fmt.Println("Job Queue Server starting...")

	// Worker3人でPoolを作成
	pool := worker.NewWorkerPool(3)
	pool.Start()

	// テスト用にJobを3つ投入
	for i := 1; i <= 3; i++ {
		pool.Submit(job.Job{
			ID:     fmt.Sprintf("job-%d", i),
			Type:   "test",
			Status: job.StatusPending,
		})
	}

	// 少し待ってから停止
	time.Sleep(1 * time.Second)
	pool.Stop()

	fmt.Println("全Job完了！")
}
