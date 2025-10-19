package db

import (
	"fmt"
	"log/slog"
	"net/url"
	"strconv"
	"time"
)

var (
	defaultPingTimeout = 4 * time.Second
	sslRootCert        = "/etc/ssl/cert.pem" // Certificate path for Alpine linux. This is the default path for the ca-certificates package.
)

// DBConfig contains all parameters needed initializing db
type DBConfig struct {
	CC ConnectionConfig `yaml:"conn"`
	SS SQLConfig        `yaml:"conf"`
}

// ConnectionConfig contains all parameters needed to connect to the database
type ConnectionConfig struct {
	Username string `yaml:"user" envconfig:"DB_USER" default:"pguser"`
	Password string `yaml:"pass" envconfig:"DB_PASS" default:"pgpass"`
	Name     string `yaml:"name" envconfig:"DB_NAME" default:"postgres"`
	Host     string `yaml:"host" envconfig:"DB_HOST" default:"127.0.0.1"`
	Port     int    `yaml:"port" envconfig:"DB_PORT" default:"5432"`
	SSL      bool   `yaml:"useSSL" envconfig:"USE_SSL_DB" default:"false"`
}

func (c *ConnectionConfig) LogValue() slog.Value {
	noPass := c
	c.Password = "*******"
	return slog.GroupValue(
		slog.String("username", c.Username),
		slog.String("password", "*********"),
		slog.String("database", c.Name),
		slog.String("host", c.Host),
		slog.Int("port", c.Port),
		slog.Bool("SSL", c.SSL),
		slog.String("URL", URLForConfig(*noPass)),
	)
}

// SQLConfig contains all parameters needed congigure connections to the database
type SQLConfig struct {
	MaxOpenConns    int           `yaml:"maxOpenConns" envconfig:"DB_MAX_OPEN_CONNS" default:"5"`             // strconv.Atoi
	MaxIdleConns    int           `yaml:"idleOpenConns" envconfig:"DB_IDLE_OPEN_CONNS" default:"3"`           // strconv.Atoi
	ConnMaxLifetime time.Duration `yaml:"maxLifetimeConns" envconfig:"DB_MAX_LIFETIME_CONNS" default:"1800s"` // time.Duration(maxSec) * time.Second
}

// config to url string
func URLForConfig(cc ConnectionConfig) string {
	sslMode := "disable"
	if cc.SSL {
		sslMode = "verify-ca&sslrootcert=" + sslRootCert
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		url.QueryEscape(cc.Username), cc.Password, url.QueryEscape(cc.Host),
		url.QueryEscape(strconv.Itoa(cc.Port)), url.QueryEscape(cc.Name), sslMode)
}
