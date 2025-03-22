package config

import (
	"database/sql"
	"log/slog"
	"time"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

func InitializeDBConnection() (*sql.DB, error) {
	cfg := mysql.Config{
		User:      AppConfig.DbUser,
		Passwd:    AppConfig.DbPassword,
		Net:       "tcp",
		Addr:      AppConfig.DbUrl,
		DBName:    "forex",
		ParseTime: true,
		Loc:       time.UTC,
	}

	// var err error
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		slog.Error("fail initialize database connection: %v", err)
		panic("fail to initalize database connection")
		return nil, err
	}

	return db, nil
}

func CloseDBConnection() {
	db.Close()
}
