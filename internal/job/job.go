package job

import "time"

// JobのステータスをString型で定義
type Status string

const (
	StatusPending   Status = "pending"   // 待機中
	StatusRunning   Status = "running"   // 実行中
	StatusCompleted Status = "completed" // 完了
	StatusFailed    Status = "failed"    // 失敗
)

// Jobの構造体
type Job struct {
	ID        string
	Type      string
	Payload   string
	Status    Status
	CreatedAt time.Time
	UpdatedAt time.Time
}
