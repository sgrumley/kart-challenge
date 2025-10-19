package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sgrumley/kart-challenge/pkg/config"
	"github.com/sgrumley/kart-challenge/pkg/db"
	"github.com/sgrumley/kart-challenge/pkg/graceful"
	"github.com/sgrumley/kart-challenge/pkg/logger"
)

var (
	localPort = "8080"
	localHost = "0.0.0.0"
)

func main() {
	ctx := context.Background()
	log := logger.NewLogger(
		logger.WithLevel(slog.LevelDebug),
		logger.WithFormat(logger.HandlerJSON),
	)
	ctx = logger.AddLoggerContext(ctx, log)

	env, err := LoadEnvVar()
	if err != nil {
		log.Error("failed to load environment config", slog.String("error", err.Error()))
		return
	}

	cfg, err := config.LoadYAMLDocument[Config](env.ConfigFilePath)
	if err != nil {
		log.Error("failed to load yaml config", slog.String("error", err.Error()))
		return
	}

	log.Info("started")
	if err := run(ctx, log, cfg); err != nil {
		log.Error("service terminated", slog.String("error", err.Error()))
		return
	}
}

func run(ctx context.Context, log *slog.Logger, cfg *Config) error {
	// configure database
	db, err := db.InitDBConnForApp(log, &cfg.Database.PostgreSQL.CC, &cfg.Database.PostgreSQL.SS)
	if err != nil {
		return fmt.Errorf("unable to create DB connection for app use with err: %w", err)
	}

	sqlxDB := sqlx.NewDb(db, "postgres")
	if sqlxDB == nil {
		return fmt.Errorf("could not initialize database: %w", err)
	}

	newAPI := NewHandler(ctx, *log, sqlxDB)

	svr := &http.Server{
		ReadHeaderTimeout: 30 * time.Second,
		Addr:              fmt.Sprintf("%s:%s", localHost, localPort),
		Handler:           newAPI,
	}

	log.Info("server started", slog.String("host", localHost), slog.String("port", localPort))
	return graceful.ListenAndServe(ctx, svr)
}
