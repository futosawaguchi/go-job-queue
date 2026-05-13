package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

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
		panic(".envファイルが見つかりません")
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

	// Worker Pool作成・起動
	pool := worker.NewWorkerPool(3, database)
	pool.Start()

	// Handler作成
	h := handler.NewHandler(pool, database)

	// echoのインスタンス作成
	e := echo.New()
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// ルーティング
	e.POST("/jobs", h.SubmitJob)
	e.GET("/jobs/:id", h.GetJob)

	// 別goroutineでサーバーを起動
	go func() {
		if err := e.Start(":8080"); err != nil {
			fmt.Println("サーバーを停止しました")
		}
	}()

	// OSシグナルを待ち受ける
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // Ctrl+C が押されるまでここで待機

	fmt.Println("シャットダウン開始...")

	// 新しいリクエストの受付を停止（10秒でタイムアウト）
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		fmt.Println("HTTPサーバーの停止に失敗:", err)
	}

	// 処理中のJobが全部終わるまで待つ
	pool.Stop()

	fmt.Println("サーバーを安全に停止しました")
}
