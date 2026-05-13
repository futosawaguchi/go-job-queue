package main

import (
	"fmt"
	"os"

	"github.com/futosawaguchi/go-job-queue/db"
	"github.com/futosawaguchi/go-job-queue/internal/handler"
	"github.com/futosawaguchi/go-job-queue/internal/worker"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	fmt.Println("Job Queue Server starting...")

	// .envファイルを読み込む
	if err := godotenv.Load(); err != nil {
		panic("'.envファイルが見つかりません")
	}

	// DB接続
	database, err := db.New(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("DB接続成功！")

	// Worker 3人でPoolを作成・起動
	pool := worker.NewWorkerPool(3, database)
	pool.Start()

	// Handlerを作成
	h := handler.NewHandler(pool, database)

	// echoのインスタンスを作成
	e := echo.New()
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	// ルーティング
	e.POST("/jobs", h.SubmitJob)
	e.GET("/jobs/:id", h.GetJob)

	// サーバー起動
	e.Logger.Fatal(e.Start(":8080"))
}
