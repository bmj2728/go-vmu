package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/rs/zerolog/log"
	"os"
)

const (
	ConfNotFoundError = "error locating GoVMU config file"
	ConfDecodeError   = "error decoding GoVMU config file"
)

func Load(path string) (*Config, error) {
	confLogger := log.With().Str("config-loader", path).Logger()
	confLogger.Info().Msg("Loading config file %s")
	//Check if the config file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		confLogger.Error().Msg(ConfNotFoundError)
		return nil, fmt.Errorf(ConfNotFoundError+": %v", err)
	}
	conf := NewConfig()
	//Try to decode the config file
	_, err := toml.DecodeFile(path, &conf)
	if err != nil {
		confLogger.Error().Msg(ConfDecodeError)
		return nil, fmt.Errorf(ConfDecodeError+": %v", err)
	}

	return conf, nil
}
