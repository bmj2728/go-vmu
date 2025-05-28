package main

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"go-vmu/internal/config"
	"go-vmu/internal/logger"
	"go-vmu/internal/processor"
)

//TODO - status - kinda works! I'm able to locally process a file - nfs share is an issue
// what doesn't work right:
// - nfs share

func main() {
	// load config
	cfg, err := config.Load("././config.toml")
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
	}
	// setup logger
	logger.Setup(&cfg.Logger)

	log.Debug().Str("startup", "logger").Msg("Logger started")
	log.Debug().Msgf("Config: %v", cfg)

	proc := processor.NewProcessor(cfg.Workers)
	results, err := proc.ProcessDirectory(cfg.ScanFolder)
	if err != nil {
		log.Error().Err(err).Msg("Error processing directory")
		return
	}

	// print results
	for _, result := range results {
		if result.Success {
			log.Debug().Msgf("Result: %s - %v", result.FilePath, result.Success)
		} else {
			log.Error().Msgf("Result: %s - %v", result.FilePath, result.Error)
		}
	}

}
