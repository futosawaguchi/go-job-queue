package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB接続の構造体
type DB struct {
	Conn *gorm.DB
}

// DB接続を作成する関数
func New(host, user, password, dbname string) (*DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s sslmode=disable",
		host, user, password, dbname,
	)

	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("DB接続失敗: %w", err)
	}

	return &DB{Conn: conn}, nil
}
