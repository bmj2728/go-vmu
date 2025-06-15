package pool

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/bmj2728/go-vmu/internal/tracker"
	"github.com/stretchr/testify/assert"
)

func TestNewWorker(t *testing.T) {
	// Create channels
	jobs := make(chan string, 1)
	results := make(chan *tracker.ProcessResult, 1)

	// Create WaitGroup
	var wg sync.WaitGroup

	// Create context
	ctx := context.Background()

	// Create tracker
	tracker := tracker.NewProgressTracker(1)

	// Create worker
	worker := NewWorker(1, jobs, results, &wg, ctx, tracker)

	// Verify worker properties
	assert.Equal(t, 1, worker.Id)
	// Just check that the channels are not nil
	assert.NotNil(t, worker.Jobs)
	assert.NotNil(t, worker.Results)
	assert.Equal(t, &wg, worker.Wg)
	assert.Equal(t, ctx, worker.Ctx)
	assert.Equal(t, tracker, worker.ProgressTracker)
}

func TestWorker_Start(t *testing.T) {
	// Create a temporary test directory
	tmpDir, err := os.MkdirTemp("", "worker-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create a test file
	testFile := filepath.Join(tmpDir, "test.txt")
	err = os.WriteFile(testFile, []byte("test data"), 0644)
	assert.NoError(t, err)

	// Test with a cancelled context
	t.Run("Cancelled context", func(t *testing.T) {
		// Create channels
		jobs := make(chan string, 1)
		results := make(chan *tracker.ProcessResult, 1)

		// Create WaitGroup
		var wg sync.WaitGroup

		// Create context with cancel
		ctx, cancel := context.WithCancel(context.Background())

		// Create worker
		worker := NewWorker(1, jobs, results, &wg, ctx, nil)

		// Start worker in a goroutine
		wg.Add(1)
		go worker.Start()

		// Cancel the context immediately
		cancel()

		// Wait for the worker to finish
		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()

		// Verify the worker exits
		select {
		case <-done:
			// Expected
		case <-time.After(1 * time.Second):
			t.Fatal("Worker did not exit after context cancellation")
		}
	})

	// Test with a closed jobs channel
	t.Run("Closed jobs channel", func(t *testing.T) {
		// Create channels
		jobs := make(chan string, 1)
		results := make(chan *tracker.ProcessResult, 1)

		// Create WaitGroup
		var wg sync.WaitGroup

		// Create context
		ctx := context.Background()

		// Create worker
		worker := NewWorker(1, jobs, results, &wg, ctx, nil)

		// Start worker in a goroutine
		wg.Add(1)
		go worker.Start()

		// Close the jobs channel
		close(jobs)

		// Wait for the worker to finish
		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()

		// Verify the worker exits
		select {
		case <-done:
			// Expected
		case <-time.After(1 * time.Second):
			t.Fatal("Worker did not exit after jobs channel was closed")
		}
	})

	// Test with a job that fails (file doesn't exist)
	t.Run("Job with non-existent file", func(t *testing.T) {
		// Create channels
		jobs := make(chan string, 1)
		results := make(chan *tracker.ProcessResult, 1)

		// Create WaitGroup
		var wg sync.WaitGroup

		// Create context
		ctx := context.Background()

		// Create tracker
		tracker := tracker.NewProgressTracker(1)

		// Create worker
		worker := NewWorker(1, jobs, results, &wg, ctx, tracker)

		// Start worker in a goroutine
		wg.Add(1)
		go worker.Start()

		// Submit a job with a non-existent file
		jobs <- "/non/existent/file.mkv"
		close(jobs)

		// Wait for the result
		result := <-results

		// Verify the result
		assert.Equal(t, "/non/existent/file.mkv", result.FilePath)
		assert.False(t, result.Success)
		assert.Error(t, result.Error)

		// Wait for the worker to finish
		wg.Wait()
	})
}
