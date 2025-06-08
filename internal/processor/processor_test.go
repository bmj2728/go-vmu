package processor

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bmj2728/go-vmu/internal/pool"
	"github.com/stretchr/testify/assert"
)

func TestNewProcessor(t *testing.T) {
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
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			processor := NewProcessor(tc.workerCount)

			assert.NotNil(t, processor)
			assert.NotNil(t, processor.Pool)
			assert.Equal(t, tc.workerCount, processor.Pool.Workers)
			assert.Nil(t, processor.ProgressTracker)
		})
	}
}

func TestProcessor_ProcessDirectory(t *testing.T) {
	// Create a temporary test directory
	tmpDir, err := os.MkdirTemp("", "processor-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Test with an empty directory
	t.Run("Empty directory", func(t *testing.T) {
		processor := NewProcessor(1)
		results, err := processor.ProcessDirectory(tmpDir)

		assert.NoError(t, err)
		assert.Empty(t, results)
	})

	// Test with a directory containing a file
	t.Run("Directory with file", func(t *testing.T) {
		// Create a test file
		testFile := filepath.Join(tmpDir, "test.mkv")
		err := os.WriteFile(testFile, []byte("test data"), 0644)
		assert.NoError(t, err)

		processor := NewProcessor(1)
		results, err := processor.ProcessDirectory(tmpDir)

		assert.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, testFile, results[0].FilePath)
		assert.False(t, results[0].Success) // Will fail because no NFO file
		assert.Error(t, results[0].Error)
	})

	// Test with a non-existent directory
	t.Run("Non-existent directory", func(t *testing.T) {
		processor := NewProcessor(1)
		results, err := processor.ProcessDirectory("/non/existent/directory")

		assert.Error(t, err)
		assert.Nil(t, results)
	})
}

// TestProcessor_Integration tests the integration of the processor with the pool
func TestProcessor_Integration(t *testing.T) {
	// Create a mock pool for testing
	mockPool := &pool.Pool{
		Workers: 2,
	}

	// Create a processor with the mock pool
	processor := &Processor{
		Pool: mockPool,
	}

	// Verify the processor uses the pool
	assert.Equal(t, mockPool, processor.Pool)
	assert.Equal(t, 2, processor.Pool.Workers)
}
