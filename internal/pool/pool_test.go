package pool

import (
	"testing"

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
