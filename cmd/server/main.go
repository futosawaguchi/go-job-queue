package main

import (
	"fmt"

	"github.com/futosawaguchi/go-job-queue/internal/handler"
	"github.com/futosawaguchi/go-job-queue/internal/worker"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	fmt.Println("Job Queue Server starting...")

	// Worker 3人でPoolを作成・起動
	pool := worker.NewWorkerPool(3)
	pool.Start()

	// Handlerを作成
	h := handler.NewHandler(pool)

	// echoのインスタンスを作成
	e := echo.New()

	// ログとリカバリのミドルウェアを追加
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	// ルーティング
	e.POST("/jobs", h.SubmitJob)

	// サーバー起動
	e.Logger.Fatal(e.Start(":8080"))
}
