package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sgrumley/kart-challenge/pkg/config"
	"github.com/sgrumley/kart-challenge/pkg/db"
	"github.com/sgrumley/kart-challenge/pkg/logger"
)

type Config struct {
	Database *DataConfig `yaml:"database"`
}

type DataConfig struct {
	PostgreSQL *db.DBConfig `yaml:"postgres"`
}

type Uploader struct {
	pool *pgxpool.Pool
}

func NewUploader(pool *pgxpool.Pool) *Uploader {
	return &Uploader{pool: pool}
}

func main() {
	log := logger.NewLogger(
		logger.WithLevel(slog.LevelDebug),
		logger.WithFormat(logger.HandlerJSON),
	)

	fp := "./config/local.yaml"
	cfg, err := config.LoadYAMLDocument[Config](fp)
	if err != nil {
		log.Error("error", "failed to configure environment", err)
		return
	}

	connStr := db.URLForConfig(cfg.Database.PostgreSQL.CC)
	pgxpool, err := NewPool(context.Background(), connStr)
	if err != nil {
		log.Error("error", "unable to create DB pool with pgx", err)
		return
	}
	defer pgxpool.Close()

	uploader := NewUploader(pgxpool)
	couponFiles := []string{"couponbase1", "couponbase2", "couponbase3"}
	uploader.processFiles(couponFiles)

	fmt.Println("✓ All coupon files processed")
}

func NewPool(ctx context.Context, connStr string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}

	config.MaxConns = int32(maxWorkers + 2)
	config.MinConns = int32(maxWorkers)
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute
	config.HealthCheckPeriod = time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.NewWithConfig: %w", err)
	}

	ctxPing, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := pool.Ping(ctxPing); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping failed: %w", err)
	}

	return pool, nil
}
