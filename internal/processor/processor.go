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

func (p *Processor) ProcessDirectory(dir string, retries int) ([]*tracker.ProcessResult, error) {
	//account for the initial run
	retries = retries + 1

	//get an initial jobs list
	log.Debug().Msg("Getting jobs")
	files, jobs, err := utils.GetFiles(dir)
	if err != nil {
		log.Error().Err(err).Msg("Error getting files")
		return nil, err
	}
	log.Debug().Msgf("Got %d files", jobs)

	//create a variable to hold successes during later loops
	var trackedResults []*tracker.ProcessResult

	for retries > 0 {
		p.ProgressTracker = tracker.NewProgressTracker(jobs)

		// add jobs to pool to avoid closing channel before submissions complete
		p.Pool.SubmitJobs(files)

		// start workers
		log.Debug().Msg("Starting workers")
		p.Pool.Start(p.ProgressTracker)

		// wait for all workers to finish
		results := p.Pool.Wait()
		// end process directory
		log.Debug().Msgf("Ending process directory - %v", results)

		successes, failures := utils.SplitResults(p.ProgressTracker.Results)

		//add successes to global tracker
		trackedResults = append(trackedResults, successes...)

		log.Debug().Msgf("Failures: %d", len(failures))
		log.Debug().Msgf("Retries Remaining: %d", retries)

		//wipe files
		files = []string{}

		for _, failure := range failures {
			files = append(files, failure.FilePath)
		}
		jobs = len(files)
		//we now have new jobs count and files to process
		//Recreate the pool - copying worker setting
		p.Pool = pool.NewPool(p.Pool.Workers)
		retries-- //decrement retry counter
		if retries == 0 {
			trackedResults = append(trackedResults, failures...)
			break
		}
		//end of the loop?
	}

	//now
	return trackedResults, nil
}
