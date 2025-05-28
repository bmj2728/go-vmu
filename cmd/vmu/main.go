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

	//// create pool --> NewProcessor(cfg.Workers)
	//log.Debug().Msg("Creating pool")
	//workers := pool.NewPool(cfg.Workers)
	//log.Debug().Msg("Pool created")
	//
	//// get files to process --> ProcessDirectory(cfg.ScanFolder)
	//log.Debug().Msg("Getting jobs")
	//err = workers.GenerateJobs(cfg.ScanFolder)
	//if err != nil {
	//	log.Error().Err(err).Msg("Error getting jobs")
	//	return
	//}
	//log.Debug().Msg("Jobs retrieved")
	//// start workers
	//log.Debug().Msg("Starting workers")
	//workers.Start()
	//// wait for all workers to finish
	//results := workers.Wait()
	//// end process directory

	proc := processor.NewProcessor(cfg.Workers)
	results, err := proc.ProcessDirectory(cfg.ScanFolder)
	if err != nil {
		log.Error().Err(err).Msg("Error processing directory")
		return
	}

	// print results
	log.Debug().Msgf("Results: %v", results)

}
