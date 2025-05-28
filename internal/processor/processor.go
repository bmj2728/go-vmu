package processor

import (
	"github.com/rs/zerolog/log"
	"go-vmu/internal/pool"
)

type Processor struct {
	Pool *pool.Pool
}

func NewProcessor(workers int) *Processor {
	return &Processor{
		Pool: pool.NewPool(workers),
	}
}

func (p *Processor) ProcessDirectory(dir string) ([]*pool.ProcessResult, error) {

	//get jobs
	log.Debug().Msg("Getting jobs")
	err := p.Pool.GenerateJobs(dir)
	if err != nil {
		log.Error().Err(err).Msg("Error getting jobs")
		return nil, err
	}
	log.Debug().Msg("Jobs retrieved")

	// start workers
	log.Debug().Msg("Starting workers")
	p.Pool.Start()

	// wait for all workers to finish
	results := p.Pool.Wait()
	// end process directory

	return results, nil
}
