package pool

import (
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"go-vmu/internal/ffmpeg"
	"go-vmu/internal/metadata"
	"go-vmu/internal/nfo"
	"go-vmu/internal/tracker"
	"go-vmu/internal/utils"
	"os"
	"sync"
)

// Worker processes jobs from the pool
type Worker struct {
	Id              int
	Jobs            <-chan string
	Results         chan<- *ProcessResult
	Wg              *sync.WaitGroup
	Ctx             context.Context
	ProgressTracker *tracker.ProgressTracker
}

// NewWorker creates a new worker
func NewWorker(id int, jobs <-chan string, results chan<- *ProcessResult, wg *sync.WaitGroup, ctx context.Context, tracker *tracker.ProgressTracker) *Worker {
	return &Worker{
		Id:              id,
		Jobs:            jobs,
		Results:         results,
		Wg:              wg,
		Ctx:             ctx,
		ProgressTracker: tracker,
	}
}

// Start begins the worker's processing loop
func (w *Worker) Start() {
	defer w.Wg.Done()

	for {
		select {
		case filePath, ok := <-w.Jobs:
			if !ok {
				// Channel closed, no more jobs
				log.Debug().Msgf("Worker %d finished. Channel closed & no more jobs", w.Id)
				return
			}
			result := w.processFile(filePath)
			w.Results <- result
			log.Debug().Msgf("Result sent to channel. Completed files: %d", len(w.Results))

		case <-w.Ctx.Done():
			// Context cancelled, stop worker
			log.Debug().Msg("Context cancelled, stopping worker")
			return
		}
		log.Debug().Msg("Worker finished loop?")
	}

}

// processFile handles the actual file processing
func (w *Worker) processFile(filePath string) *ProcessResult {
	result := ProcessResult{FilePath: filePath}
	var success bool
	var err error

	log.Debug().Strs(fmt.Sprintf("Processing file %s", filePath), []string{"worker", "id", fmt.Sprintf("%d", w.Id)}).Msg("Processing file")

	//validate existence
	_, err = os.Stat(filePath)
	if err != nil {
		log.Error().Err(err).Msg("File does not exist")
		success = false
		if w.ProgressTracker != nil {
			w.ProgressTracker.CompleteFile(filePath)
		}
		return result.WithResult(success, err)
	}

	//get nfo file
	nfoPath, err := nfo.MatchEpisodeFile(filePath)
	if err != nil {
		log.Error().Err(err).Msg("Error matching NFO file")
		success = false
		if w.ProgressTracker != nil {
			w.ProgressTracker.CompleteFile(filePath)
		}
		return result.WithResult(success, err)
	}
	if nfoPath == "" {
		log.Error().Msg("NFO file not found")
		success = false
		if w.ProgressTracker != nil {
			w.ProgressTracker.CompleteFile(filePath)
		}
		return result.WithResult(success, errors.New("NFO file not found"))
	}

	//parse nfo file
	data, err := nfo.ParseEpisodeNFO(nfoPath)
	if err != nil {
		log.Error().Err(err).Msg("Error parsing NFO file")
		success = false
		if w.ProgressTracker != nil {
			w.ProgressTracker.CompleteFile(filePath)
		}
		return result.WithResult(success, err)
	}

	//process into metadata
	adapter := metadata.NewNFOAdapter(data)
	meta, err := adapter.TranslateNFO()
	if err != nil {
		log.Error().Err(err).Msg("Error translating NFO file")
		success = false
		if w.ProgressTracker != nil {
			w.ProgressTracker.CompleteFile(filePath)
		}
		return result.WithResult(success, err)
	}

	//create ffmpeg command
	outputFile := utils.InsertTagToFileName(filePath, "govmu-edit")
	cmd, err := ffmpeg.NewFFmpegCommand().WithInput(filePath).WithOutput(outputFile).WithMetadata(*meta)
	if err != nil {
		log.Error().Err(err).Msg("Error creating ffmpeg command")
		success = false
		if w.ProgressTracker != nil {
			w.ProgressTracker.CompleteFile(filePath)
		}
		return result.WithResult(success, err)
	}
	cmd = cmd.GenerateArgs()
	log.Debug().Msgf("FFmpeg command: %v", cmd)

	//create executor
	executor := ffmpeg.NewExecutor(cmd, w.ProgressTracker)

	//execute
	err = executor.Execute()
	if err != nil {
		log.Error().Err(err).Msg("Error executing ffmpeg command")
		success = false
		if w.ProgressTracker != nil {
			w.ProgressTracker.CompleteFile(filePath)
		}
		return result.WithResult(success, err)
	}

	//validate
	ok, err := executor.ValidateNewFile()
	if err != nil {
		log.Error().Err(err).Msg("Error validating new file")
		success = false
		if w.ProgressTracker != nil {
			w.ProgressTracker.CompleteFile(filePath)
		}
		return result.WithResult(success, err)
	}
	if !ok {
		log.Error().Msg("New file is invalid")
		success = false
		if w.ProgressTracker != nil {
			w.ProgressTracker.CompleteFile(filePath)
		}
		return result.WithResult(success, err)
	}

	//cleanup
	err = executor.Cleanup()
	if err != nil {
		log.Error().Err(err).Msg("Error cleaning up")
		success = false
		if w.ProgressTracker != nil {
			w.ProgressTracker.CompleteFile(filePath)
		}
		return result.WithResult(success, err)
	}

	success = true

	log.Debug().Msgf("Worker %d processed file successfully: %s", w.Id, filePath)

	if w.ProgressTracker != nil {
		w.ProgressTracker.CompleteFile(filePath)
	}

	//share results
	return result.WithResult(success, err)
}
