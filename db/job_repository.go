package db

import (
	"time"

	"github.com/futosawaguchi/go-job-queue/internal/job"
)

// Jobを保存する
func (db *DB) SaveJob(j job.Job) error {
	return db.Conn.Create(&j).Error
}

// IDでJobを取得する
func (db *DB) GetJob(id string) (*job.Job, error) {
	var j job.Job
	result := db.Conn.First(&j, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &j, nil
}

// JobのStatusを更新する
func (db *DB) UpdateJobStatus(id string, status job.Status) error {
	return db.Conn.Model(&job.Job{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		}).Error
}
