package processor

import (
	"github.com/bmj2728/go-vmu/internal/pool"
	"github.com/bmj2728/go-vmu/internal/tracker"
	"github.com/bmj2728/go-vmu/internal/utils"
	"github.com/rs/zerolog/log"
)

type Processor struct {
	Pool            *pool.Pool
	ProgressTracker *tracker.ProgressTracker
}

func NewProcessor(workers int) *Processor {
	return &Processor{
		Pool: pool.NewPool(workers),
	}
}

func (p *Processor) ProcessDirectory(dir string) ([]*tracker.ProcessResult, error) {

	//get jobs
	log.Debug().Msg("Getting jobs")

	files, jobs, err := utils.GetFiles(dir)
	if err != nil {
		log.Error().Err(err).Msg("Error getting files")
		return nil, err
	}
	log.Debug().Msgf("Got %d files", jobs)

	p.ProgressTracker = tracker.NewProgressTracker(jobs)

	// add jobs to pool to avoid closing channel before submissions complete
	p.Pool.SubmitJobs(files)

	// start workers
	log.Debug().Msg("Starting workers")
	p.Pool.Start(p.ProgressTracker)

	// wait for all workers to finish
	results := p.Pool.Wait()
	// end process directory

	return results, nil
}
