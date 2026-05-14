package worker

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/futosawaguchi/go-job-queue/internal/job"
)

// テスト用のダミーDB
type mockDB struct {
	mu       sync.Mutex
	statuses map[string]job.Status
}

func newMockDB() *mockDB {
	return &mockDB{
		statuses: make(map[string]job.Status),
	}
}

func (m *mockDB) UpdateJobStatus(id string, status job.Status) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.statuses[id] = status
	return nil
}

// テスト1: Jobが正しく処理されるか
func TestWorkerPool_ProcessJob(t *testing.T) {
	db := newMockDB()
	pool := NewWorkerPoolWithDB(3, db)
	pool.Start()

	j := job.Job{
		ID:     "test-job-1",
		Type:   "test",
		Status: job.StatusPending,
	}

	pool.Submit(j)

	// 処理完了を待つ
	time.Sleep(100 * time.Millisecond)
	pool.Stop()

	// completedになっているか確認
	db.mu.Lock()
	status := db.statuses["test-job-1"]
	db.mu.Unlock()

	if status != job.StatusCompleted {
		t.Errorf("期待値: completed, 実際: %s", status)
	}
}

// テスト2: 並行処理されているか証明する
func TestWorkerPool_Concurrency(t *testing.T) {
	db := newMockDB()
	pool := NewWorkerPoolWithDB(3, db)
	pool.Start()

	// 処理中のWorker数をカウント
	var concurrent int32
	var maxConcurrent int32

	pool.SetProcessor(func(j job.Job) {
		// 同時実行数をカウント
		current := atomic.AddInt32(&concurrent, 1)

		// 最大同時実行数を記録
		for {
			max := atomic.LoadInt32(&maxConcurrent)
			if current <= max {
				break
			}
			if atomic.CompareAndSwapInt32(&maxConcurrent, max, current) {
				break
			}
		}

		time.Sleep(50 * time.Millisecond)
		atomic.AddInt32(&concurrent, -1)
	})

	// 6つのJobを投入
	for i := 0; i < 6; i++ {
		pool.Submit(job.Job{
			ID:   fmt.Sprintf("job-%d", i),
			Type: "test",
		})
	}

	time.Sleep(500 * time.Millisecond)
	pool.Stop()

	// 最大同時実行数が2以上なら並行処理されている
	if maxConcurrent < 2 {
		t.Errorf("並行処理されていない: 最大同時実行数 = %d", maxConcurrent)
	}
	t.Logf("最大同時実行数: %d", maxConcurrent)
}
