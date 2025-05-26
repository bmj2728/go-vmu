package main

import (
	"github.com/rs/zerolog/log"
	"go-vmu/internal/config"
	"go-vmu/internal/logger"
)

func main() {
	cfg, err := config.Load("././config.toml")
	if err != nil {
		panic(err)
	}
	logger.Setup(&cfg.Logger)

	log.Info().Str("startup", "logger").Msg("Logger started")
	log.Info().Msgf("Config: %v", cfg)

}
