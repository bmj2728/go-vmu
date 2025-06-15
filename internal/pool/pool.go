package pool

import (
	"context"
	"github.com/bmj2728/go-vmu/internal/tracker"
	"github.com/rs/zerolog/log"
	"sync"
)

// Pool manages a collection of workers
type Pool struct {
	Workers         int
	Jobs            chan string
	Results         chan *tracker.ProcessResult
	Wg              sync.WaitGroup
	Ctx             context.Context
	CancelFunc      context.CancelFunc
	ProgressTracker *tracker.ProgressTracker
}

// NewPool creates a new worker pool
func NewPool(workerCount int) *Pool {
	ctx, cancel := context.WithCancel(context.Background())
	return &Pool{
		Workers:    workerCount,
		Ctx:        ctx,
		CancelFunc: cancel,
	}
}

// Start launches the worker pool
func (p *Pool) Start(tracker *tracker.ProgressTracker) {
	for i := 0; i < p.Workers; i++ {
		worker := NewWorker(i, p.Jobs, p.Results, &p.Wg, p.Ctx, tracker)
		log.Debug().Msgf("Starting worker %d", i)
		p.Wg.Add(1)
		go worker.Start()
	}
}

// Submit adds a job to the pool
func (p *Pool) Submit(filePath string) {
	log.Debug().Msgf("Submitting job for %s", filePath)
	p.Jobs <- filePath
}

//TODO update to return done bool and use tracker's Result param

// Wait waits for all jobs to complete and returns results
func (p *Pool) Wait() []*tracker.ProcessResult {
	// Use a defer to recover from panic if the channel is already closed
	defer func() {
		if r := recover(); r != nil {
			// Channel was already closed, ignore the panic
			log.Debug().Msg("Jobs channel was already closed")
		}
	}()

	// Try to close the Jobs channel
	close(p.Jobs)

	p.Wg.Wait() // Wait for all workers to finish

	// Collect all results that are already in the channel
	// changed way this works, but leaving in for now to avoid potential breaking changes
	var processResults []*tracker.ProcessResult
	// Keep reading from Results channel until it's empty
	for {
		select {
		case result, ok := <-p.Results:
			if !ok {
				// Channel is closed
				return processResults
			}
			processResults = append(processResults, result)
		default:
			// No more results available right now
			close(p.Results)
			// Drain any remaining results
			for result := range p.Results {
				processResults = append(processResults, result)
			}
			return processResults
		}
	}
}

// Stop cancels all workers
func (p *Pool) Stop() {
	p.CancelFunc()
}

func (p *Pool) SubmitJobs(paths []string) int {
	bufferSize := len(paths) // Use actual file count for buffer
	log.Debug().Msgf("Creating job channel with buffer size %d for %d files", bufferSize, len(paths))
	p.Jobs = make(chan string, bufferSize)
	p.Results = make(chan *tracker.ProcessResult, bufferSize)
	// Submit all jobs
	for _, path := range paths {
		p.Submit(path)
	}
	return len(p.Jobs)
}
