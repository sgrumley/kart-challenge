package logger

import (
	"log/slog"
	"os"
)

const (
	HandlerJSON Handler = "json"
	HandlerText Handler = "text"
)

type (
	Handler       string
	Option        func(*LoggerOptions)
	LoggerOptions struct {
		level  slog.Level
		format Handler
	}
)

func WithLevel(level slog.Level) Option {
	return func(opts *LoggerOptions) {
		opts.level = level
	}
}

func WithFormat(format Handler) Option {
	return func(opts *LoggerOptions) {
		opts.format = format
	}
}

func NewLogger(options ...Option) *slog.Logger {
	// default
	opts := LoggerOptions{
		level:  slog.LevelInfo,
		format: HandlerText,
	}

	for _, opt := range options {
		opt(&opts)
	}

	return slog.New(getHandler(opts))
}

func getHandler(opts LoggerOptions) slog.Handler {
	baseOpts := &slog.HandlerOptions{
		Level: opts.level,
	}
	output := os.Stdout

	switch opts.format {
	case HandlerJSON:
		return slog.NewJSONHandler(output, baseOpts)

	case HandlerText:
		return slog.NewTextHandler(output, baseOpts)
	}

	return slog.NewTextHandler(output, baseOpts)
}
