package database

import (
	"database/sql"
	"fmt"
	"product-service/config"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() error {
	cfg := config.LoadConfig()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	if err = db.Ping(); err != nil {
		return err
	}

	// 设置连接池参数
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	DB = db
	return nil
}

func CloseDB() {
	if DB != nil {
		err := DB.Close()
		if err != nil {
			return
		}
	}
}
