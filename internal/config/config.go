package config

import "time"

type Config struct {
	NFOType    string       `toml:"nfo_type"`
	ScanFolder string       `toml:"scan_folder"`
	Workers    int          `toml:"workers"`
	Logger     LoggerConfig `toml:"logger"`
}

type LoggerConfig struct {
	Level      string    `toml:"level"`
	Pretty     bool      `toml:"pretty"`
	TimeFormat time.Time `toml:"time_format"`
	LogFile    string    `toml:"log_file"`
	MaxSize    int       `toml:"max_size"`
	MaxBackups int       `toml:"max_backups"`
	MaxAge     int       `toml:"max_age"`
	Compress   bool      `toml:"compress"`
}

func NewConfig() *Config {
	return &Config{}
}
