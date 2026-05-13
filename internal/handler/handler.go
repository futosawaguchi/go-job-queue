package handler

import (
	"net/http"
	"time"

	"github.com/futosawaguchi/go-job-queue/db"
	"github.com/futosawaguchi/go-job-queue/internal/job"
	"github.com/futosawaguchi/go-job-queue/internal/worker"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	pool *worker.WorkerPool
	db   *db.DB
}

// dbを受け取るように変更
func NewHandler(pool *worker.WorkerPool, db *db.DB) *Handler {
	return &Handler{pool: pool, db: db}
}

type SubmitJobRequest struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

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
		ID:        uuid.New().String(),
		Type:      req.Type,
		Payload:   req.Payload,
		Status:    job.StatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// DBに保存
	if err := h.db.SaveJob(j); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "jobの保存に失敗しました",
		})
	}

	// Workerに投入
	h.pool.Submit(j)

	return c.JSON(http.StatusAccepted, SubmitJobResponse{
		JobID:  j.ID,
		Status: string(j.Status),
	})
}

// JobのレスポンスのStruct
type GetJobResponse struct {
	JobID     string `json:"job_id"`
	Type      string `json:"type"`
	Payload   string `json:"payload"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// GET /jobs/:id
func (h *Handler) GetJob(c echo.Context) error {
	id := c.Param("id")

	j, err := h.db.GetJob(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "jobが見つかりません",
		})
	}

	return c.JSON(http.StatusOK, GetJobResponse{
		JobID:     j.ID,
		Type:      j.Type,
		Payload:   j.Payload,
		Status:    string(j.Status),
		CreatedAt: j.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: j.UpdatedAt.Format("2006-01-02 15:04:05"),
	})
}
