package handler

import (
	"net/http"
	"time"

	"github.com/futosawaguchi/go-job-queue/internal/job"
	"github.com/futosawaguchi/go-job-queue/internal/worker"
	"github.com/labstack/echo/v4"
)

// Handlerの構造体
type Handler struct {
	pool *worker.WorkerPool
}

// Handlerを作成する関数
func NewHandler(pool *worker.WorkerPool) *Handler {
	return &Handler{pool: pool}
}

// JobをSubmitするAPIのリクエスト
type SubmitJobRequest struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

// JobをSubmitするAPIのレスポンス
type SubmitJobResponse struct {
	JobID  string `json:"job_id"`
	Status string `json:"status"`
}

// POST /jobs
func (h *Handler) SubmitJob(c echo.Context) error {
	req := new(SubmitJobRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request",
		})
	}

	j := job.Job{
		ID:        time.Now().Format("20060102150405"),
		Type:      req.Type,
		Payload:   req.Payload,
		Status:    job.StatusPending,
		CreatedAt: time.Now(),
	}

	h.pool.Submit(j)

	return c.JSON(http.StatusAccepted, SubmitJobResponse{
		JobID:  j.ID,
		Status: string(j.Status),
	})
}
