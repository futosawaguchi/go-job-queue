# go-job-queue

Go言語の学習を目的として作成した、Job Queue / Task Worker システムです。
Goの並行処理（goroutine / channel）を活用し、複数のJobを並行して処理します。

## システム構成

```
フロントエンド（Next.js :3000）
        ↕ HTTP
バックエンド（Go :8080）
        ↕ GORM
データベース（PostgreSQL）
```

## 主な機能

- Job の投入（POST /jobs）
- Job のステータス確認（GET /jobs/:id）
- Job の一覧取得（GET /jobs）
- Worker Pool による並行処理（goroutine × 3）
- Graceful Shutdown（処理中の Job を安全に完了してから終了）
- リアルタイムステータス可視化（フロントエンド）

## 使用技術

### バックエンド

| 技術 | 用途 |
|---|---|
| Go 1.26 | メイン言語 |
| echo | HTTP フレームワーク |
| GORM | ORM |
| PostgreSQL | データベース |
| UUID | Job ID 生成 |

### フロントエンド

| 技術 | 用途 |
|---|---|
| Next.js 16 | フレームワーク |
| TypeScript | 言語 |
| Tailwind CSS | スタイリング |

## Goらしい設計ポイント

### goroutine / channel による Worker Pool

```go
// Worker 3人を goroutine で並行起動
for i := 0; i < wp.workerCount; i++ {
    wp.wg.Add(1)
    go wp.runWorker(i)
}

// channel 経由で Job を安全に渡す
func (wp *WorkerPool) Submit(j job.Job) {
    wp.jobQueue <- j
}
```

### Graceful Shutdown

```go
// OS シグナルを待ち受けて安全に終了
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit

// 処理中の Job が全部終わるまで待つ
pool.Stop()
```

### interface によるテストのモック化

```go
type JobDB interface {
    UpdateJobStatus(id string, status job.Status) error
}
// 本物の DB もモック DB も同じ interface で扱える
```

## ディレクトリ構成

```
go-job-queue/
├── cmd/
│   └── server/
│       └── main.go          # エントリーポイント
├── internal/
│   ├── job/
│   │   └── job.go           # Job の構造体・ステータス定義
│   ├── worker/
│   │   ├── worker.go        # Worker Pool（並行処理の核心）
│   │   └── worker_test.go   # 並行処理のテスト
│   └── handler/
│       └── handler.go       # HTTP ハンドラ
├── db/
│   ├── db.go                # DB 接続
│   └── job_repository.go    # DB 操作（GORM）
├── frontend/                # Next.js フロントエンド
├── .env.example             # 環境変数の見本
└── go.mod
```

## セットアップ

### 必要なもの

- Go 1.26+
- PostgreSQL 16+
- Node.js 22+

### 手順

**1. リポジトリをクローン**

```bash
git clone https://github.com/futosawaguchi/go-job-queue.git
cd go-job-queue
```

**2. 環境変数を設定**

```bash
cp .env.example .env
# .env を編集して DB のパスワードを設定
```

**3. データベースを作成**

```bash
psql postgres
```

```sql
CREATE DATABASE jobqueue;
CREATE USER jobuser WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE jobqueue TO jobuser;
ALTER DATABASE jobqueue OWNER TO jobuser;
ALTER SCHEMA public OWNER TO jobuser;
\q
```

**4. テーブルを作成**

```bash
psql -U jobuser -d jobqueue
```

```sql
CREATE TABLE jobs (
    id          VARCHAR(36) PRIMARY KEY,
    type        VARCHAR(50) NOT NULL,
    payload     TEXT,
    status      VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP NOT NULL DEFAULT NOW()
);
\q
```

**5. バックエンドを起動**

```bash
go run cmd/server/main.go
```

**6. フロントエンドを起動**

```bash
cd frontend
cp .env.local.example .env.local
npm install
npm run dev
```

**7. ブラウザで確認**

http://localhost:3000 を開いて Job を投入してみてください。

## テスト

```bash
go test -v ./internal/worker/...
```

並行処理のテスト結果例：

```
--- PASS: TestWorkerPool_ProcessJob (0.10s)
    worker_test.go: 最大同時実行数: 3
--- PASS: TestWorkerPool_Concurrency (0.50s)
```

## Job のライフサイクル

```
POST /jobs
    ↓
DB に保存（status: pending）
    ↓
Worker Pool に投入
    ↓
Worker が処理開始（status: running）
    ↓
処理完了（status: completed）
```
