package db

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/lib/pq"
)

func InitDBConnForApp(log *slog.Logger, cc *ConnectionConfig, dsc *SQLConfig) (*sql.DB, error) {
	db, err := setupDB(log, cc)
	if err != nil {
		return nil, err
	}

	if dsc != nil { // Ensure type here?
		db.SetMaxOpenConns(dsc.MaxOpenConns)
		db.SetMaxIdleConns(dsc.MaxIdleConns)
		db.SetConnMaxLifetime(dsc.ConnMaxLifetime)

		log.Info("DB connection data",
			slog.Group("connection pools",
				slog.Int("max open connections", dsc.MaxOpenConns),
				slog.Int("max idle connections", dsc.MaxIdleConns),
				slog.Duration("connection max life time", dsc.ConnMaxLifetime),
			),
		)
	}
	return db, nil
}

func setupDB(log *slog.Logger, cc *ConnectionConfig) (*sql.DB, error) {
	dsn := URLForConfig(*cc)
	log.Info("DB", "Connection", cc)

	var db *sql.DB
	var err error

	db, err = sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("sql Open failed:%v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), defaultPingTimeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("error pinging db: %v", err)
	}

	return db, nil
}
