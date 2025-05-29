package logger

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go-vmu/internal/config"
	"io"
	"os"
	"path/filepath"
)

// Setup configures the global logger
func Setup(cfg *config.LoggerConfig) {
	// Set the default time format
	timeFormatString := "2006-01-02 15:04:05"

	// Set the logger level
	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	// Configure output
	var output io.Writer = os.Stdout

	// If LogFile is specified, also write logs to a file with rotation
	if cfg.LogFile != "" {
		// Ensure directory exists
		logDir := filepath.Dir(cfg.LogFile)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			fmt.Printf("Failed to create log directory: %v\n", err)
		} else {
			// Configure log rotation
			rotateLogger := &lumberjack.Logger{
				Filename:   cfg.LogFile,
				MaxSize:    cfg.MaxSize,    // megabytes
				MaxBackups: cfg.MaxBackups, // number of backups
				MaxAge:     cfg.MaxAge,     // days
				Compress:   cfg.Compress,
			}

			// Use MultiWriter to write to both stdout and rotated file
			output = zerolog.MultiLevelWriter(output, rotateLogger)
		}
	}

	// If Pretty is enabled, use ConsoleWriter for human-readable output
	if cfg.Pretty {
		output = zerolog.ConsoleWriter{
			Out:        output,
			TimeFormat: timeFormatString,
			NoColor:    false,
			FormatLevel: func(i interface{}) string {
				return fmt.Sprintf("| %-6s|", i)
			},
			FormatCaller: func(i interface{}) string {
				return filepath.Base(fmt.Sprintf("%s", i))
			},
		}
	}

	// Set global logger
	log.Logger = zerolog.New(output).
		With().
		Timestamp().
		Caller().
		Logger()
}
