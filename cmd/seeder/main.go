package main

import (
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"

	"github.com/sgrumley/kart-challenge/cmd/seeder/seeding"
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

	db, err := db.InitDBConnForApp(log, &cfg.Database.PostgreSQL.CC, &cfg.Database.PostgreSQL.SS)
	if err != nil {
		log.Error("error", "unable to create DB connection for app", err)
		return
	}

	sqlxDB := sqlx.NewDb(db, "postgres")
	if sqlxDB == nil {
		log.Error("error", "could not initialize database: ", err)
		return
	}

	seeding.InsertAll(sqlxDB)
	fmt.Println("successfully seeded")
}
