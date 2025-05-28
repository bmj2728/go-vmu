package pool

import (
	"context"
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Pool manages a collection of workers
type Pool struct {
	Workers    int
	Jobs       chan string
	Results    chan *ProcessResult
	Wg         sync.WaitGroup
	Ctx        context.Context
	CancelFunc context.CancelFunc
}

// NewPool creates a new worker pool
func NewPool(workerCount int) *Pool {
	ctx, cancel := context.WithCancel(context.Background())
	return &Pool{
		Workers:    workerCount,
		Jobs:       make(chan string, workerCount*2), // Buffer for efficiency
		Results:    make(chan *ProcessResult),
		Ctx:        ctx,
		CancelFunc: cancel,
	}
}

// Start launches the worker pool
func (p *Pool) Start() {
	for i := 0; i < p.Workers; i++ {
		worker := NewWorker(i, p.Jobs, p.Results, &p.Wg, p.Ctx)
		p.Wg.Add(1)
		go worker.Start()
	}
}

// Submit adds a job to the pool
func (p *Pool) Submit(filePath string) {
	p.Jobs <- filePath
}

// Wait waits for all jobs to complete and returns results
func (p *Pool) Wait() []*ProcessResult {
	close(p.Jobs) // No more jobs will be submitted
	p.Wg.Wait()   // Wait for all workers to finish
	close(p.Results)

	var processResults []*ProcessResult
	for result := range p.Results {
		processResults = append(processResults, result)
	}
	return processResults
}

// Stop cancels all workers
func (p *Pool) Stop() {
	p.CancelFunc()
}

func (p *Pool) GetJobs(root string) error {

	//helper function
	getFiles := func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		} else if info.IsDir() {
			return err
		} else if strings.HasSuffix(info.Name(), ".avi") ||
			strings.HasSuffix(info.Name(), ".mp4") ||
			strings.HasSuffix(info.Name(), ".mkv") ||
			strings.HasSuffix(info.Name(), ".mpg") ||
			strings.HasSuffix(info.Name(), ".mov") ||
			strings.HasSuffix(info.Name(), ".wmv") ||
			strings.HasSuffix(info.Name(), ".flv") ||
			strings.HasSuffix(info.Name(), ".m4v") {
			p.Submit(path)
		}

		return nil
	}

	//walk the directory
	err := filepath.Walk(root, getFiles)

	//handle errors
	if err != nil {
		log.Error().Err(err).Msg("Error walking directory")
		return err
	}

	return nil
}
