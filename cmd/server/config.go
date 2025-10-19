package main

import (
	"github.com/kelseyhightower/envconfig"

	"github.com/sgrumley/kart-challenge/pkg/db"
)

type EnvVar struct {
	ConfigFilePath string `envconfig:"CONFIG_FILE_PATH" default:"./config/local.yaml"`
	Port           string `envconfig:"APP_PORT" default:"8080"`
	Host           string `envconfig:"APP_HOST" default:""`
	Environment    string `envconfig:"ENVIRONMENT" default:"prod"`
}

func LoadEnvVar() (EnvVar, error) {
	c := EnvVar{}
	err := envconfig.Process("", &c)
	return c, err
}

type Config struct {
	Database *DataConfig `yaml:"database"`
}

type DataConfig struct {
	PostgreSQL *db.DBConfig `yaml:"postgres"`
}
