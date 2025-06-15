package pool

import (
	"errors"
	"testing"
	"time"

	"github.com/bmj2728/go-vmu/internal/tracker"
	"github.com/stretchr/testify/assert"
)

func TestNewPool(t *testing.T) {
	// Test with different worker counts
	testCases := []struct {
		name        string
		workerCount int
	}{
		{
			name:        "Single worker",
			workerCount: 1,
		},
		{
			name:        "Multiple workers",
			workerCount: 5,
		},
		{
			name:        "Zero workers",
			workerCount: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pool := NewPool(tc.workerCount)

			assert.NotNil(t, pool)
			assert.Equal(t, tc.workerCount, pool.Workers)
			assert.NotNil(t, pool.Ctx)
			assert.NotNil(t, pool.CancelFunc)
			assert.Nil(t, pool.Jobs)
			assert.Nil(t, pool.Results)
		})
	}
}

func TestPool_SubmitJobs(t *testing.T) {
	// Test with empty paths
	t.Run("Empty paths", func(t *testing.T) {
		pool := NewPool(1)
		count := pool.SubmitJobs([]string{})

		assert.Equal(t, 0, count)
		assert.NotNil(t, pool.Jobs)
		assert.NotNil(t, pool.Results)
		assert.Equal(t, 0, len(pool.Jobs))
	})

	// Test with single path
	t.Run("Single path", func(t *testing.T) {
		pool := NewPool(1)
		paths := []string{"/path/to/file.mkv"}
		count := pool.SubmitJobs(paths)

		assert.Equal(t, 1, count)
		assert.NotNil(t, pool.Jobs)
		assert.NotNil(t, pool.Results)

		// Verify job was submitted
		job := <-pool.Jobs
		assert.Equal(t, "/path/to/file.mkv", job)
	})

	// Test with multiple paths
	t.Run("Multiple paths", func(t *testing.T) {
		pool := NewPool(1)
		paths := []string{
			"/path/to/file1.mkv",
			"/path/to/file2.mkv",
			"/path/to/file3.mkv",
		}
		count := pool.SubmitJobs(paths)

		assert.Equal(t, 3, count)
		assert.NotNil(t, pool.Jobs)
		assert.NotNil(t, pool.Results)

		// Verify jobs were submitted in order
		job1 := <-pool.Jobs
		assert.Equal(t, "/path/to/file1.mkv", job1)

		job2 := <-pool.Jobs
		assert.Equal(t, "/path/to/file2.mkv", job2)

		job3 := <-pool.Jobs
		assert.Equal(t, "/path/to/file3.mkv", job3)
	})
}

func TestPool_Submit(t *testing.T) {
	pool := NewPool(1)
	pool.Jobs = make(chan string, 1)

	// Submit a job
	pool.Submit("/path/to/file.mkv")

	// Verify job was submitted
	job := <-pool.Jobs
	assert.Equal(t, "/path/to/file.mkv", job)
}

func TestPool_Stop(t *testing.T) {
	pool := NewPool(1)

	// Verify context is not done before stopping
	select {
	case <-pool.Ctx.Done():
		t.Fatal("Context should not be done before stopping")
	default:
		// Expected
	}

	// Stop the pool
	pool.Stop()

	// Verify context is done after stopping
	select {
	case <-pool.Ctx.Done():
		// Expected
	default:
		t.Fatal("Context should be done after stopping")
	}
}

func TestPool_Start(t *testing.T) {
	// Create a pool with 2 workers
	pool := NewPool(2)
	pool.Jobs = make(chan string, 2)
	pool.Results = make(chan *tracker.ProcessResult, 2)

	// Create a tracker
	tracker := tracker.NewProgressTracker(2)

	// Start the pool
	pool.Start(tracker)

	// We can't directly verify the WaitGroup, but we can verify that
	// the pool has started by checking that the workers are running
	// by submitting a job and getting a result

	// Submit a job
	pool.Jobs <- "/path/to/file.mkv"

	// Wait for the result to be sent to the Results channel
	// This is similar to how the worker test does it
	var result *tracker.ProcessResult
	select {
	case result = <-pool.Results:
		// Got a result
	case <-time.After(1 * time.Second):
		t.Fatal("Timed out waiting for result")
	}

	// Close the jobs channel to signal no more jobs
	close(pool.Jobs)

	// Wait for the workers to finish
	pool.Wait()

	// Verify we got a result (even though it will be an error since the file doesn't exist)
	assert.NotNil(t, result)
	assert.Equal(t, "/path/to/file.mkv", result.FilePath)
	assert.False(t, result.Success)
	assert.Error(t, result.Error)
}

func TestPool_Wait(t *testing.T) {
	// Create a pool with 2 workers
	pool := NewPool(2)

	// Create jobs and results channels
	pool.Jobs = make(chan string, 2)
	pool.Results = make(chan *tracker.ProcessResult, 2)

	// Add some results to the channel
	result1 := &tracker.ProcessResult{FilePath: "/path/to/file1.mkv", Success: true}
	result2 := &tracker.ProcessResult{FilePath: "/path/to/file2.mkv", Success: false, Error: errors.New("test error")}

	pool.Results <- result1
	pool.Results <- result2

	// Call Wait
	results := pool.Wait()

	// Verify the results
	assert.Len(t, results, 2)
	assert.Contains(t, results, result1)
	assert.Contains(t, results, result2)
}
