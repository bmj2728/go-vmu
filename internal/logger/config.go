package logger

import "time"

type LoggerConfig struct {
	Level      string `toml:"level"`
	Pretty     bool   `toml:"pretty"`
	TimeFormat string `toml:"time_format"`
	LogFile    string `toml:"log_file"`
	MaxSize    int    `toml:"max_size"`
	MaxBackups int    `toml:"max_backups"`
	MaxAge     int    `toml:"max_age"`
	Compress   bool   `toml:"compress"`
}

func NewLoggerConfig(verbose bool) *LoggerConfig {
	level := "info"
	if verbose {
		level = "debug"
	}
	return &LoggerConfig{
		Level:      level,
		Pretty:     true,
		TimeFormat: time.RFC3339,
		LogFile:    "",
		MaxSize:    5,
		MaxBackups: 5,
		MaxAge:     14,
		Compress:   false,
	}
}
