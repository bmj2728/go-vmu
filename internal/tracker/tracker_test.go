package tracker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProgressTracker(t *testing.T) {
	// Test with different total files
	testCases := []struct {
		name       string
		totalFiles int
	}{
		{
			name:       "Zero files",
			totalFiles: 0,
		},
		{
			name:       "Single file",
			totalFiles: 1,
		},
		{
			name:       "Multiple files",
			totalFiles: 10,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tracker := NewProgressTracker(tc.totalFiles)

			assert.NotNil(t, tracker)
			assert.Equal(t, tc.totalFiles, tracker.totalFiles)
			assert.Equal(t, 0, tracker.completedFiles)
			assert.NotNil(t, tracker.currentFiles)
			assert.Empty(t, tracker.currentFiles)
			assert.NotNil(t, tracker.bar)
		})
	}
}

func TestProgressTracker_UpdateStage(t *testing.T) {
	// Create a tracker
	tracker := NewProgressTracker(2)

	// Test updating stage for a new file
	t.Run("New file", func(t *testing.T) {
		tracker.UpdateStage("/path/to/file1.mkv", StageBackup)

		// Verify the file is tracked
		assert.Len(t, tracker.currentFiles, 1)
		progress, exists := tracker.currentFiles["/path/to/file1.mkv"]
		assert.True(t, exists)
		assert.Equal(t, "/path/to/file1.mkv", progress.filename)
		assert.Equal(t, StageBackup, progress.stage)
		assert.False(t, progress.done)
	})

	// Test updating stage for an existing file
	t.Run("Existing file", func(t *testing.T) {
		tracker.UpdateStage("/path/to/file1.mkv", StageProcess)

		// Verify the file's stage is updated
		assert.Len(t, tracker.currentFiles, 1)
		progress, exists := tracker.currentFiles["/path/to/file1.mkv"]
		assert.True(t, exists)
		assert.Equal(t, "/path/to/file1.mkv", progress.filename)
		assert.Equal(t, StageProcess, progress.stage)
		assert.False(t, progress.done)
	})

	// Test updating stage for a second file
	t.Run("Second file", func(t *testing.T) {
		tracker.UpdateStage("/path/to/file2.mkv", StageBackup)

		// Verify both files are tracked
		assert.Len(t, tracker.currentFiles, 2)
		progress1, exists1 := tracker.currentFiles["/path/to/file1.mkv"]
		assert.True(t, exists1)
		assert.Equal(t, StageProcess, progress1.stage)
		progress2, exists2 := tracker.currentFiles["/path/to/file2.mkv"]
		assert.True(t, exists2)
		assert.Equal(t, StageBackup, progress2.stage)
	})
}

func TestProgressTracker_CompleteFile(t *testing.T) {
	// Create a tracker
	tracker := NewProgressTracker(2)

	// Add files to track
	tracker.UpdateStage("/path/to/file1.mkv", StageBackup)
	tracker.UpdateStage("/path/to/file2.mkv", StageProcess)

	// Verify initial state
	assert.Equal(t, 0, tracker.completedFiles)
	assert.Len(t, tracker.currentFiles, 2)

	// Complete the first file
	t.Run("Complete first file", func(t *testing.T) {
		tracker.CompleteFile("/path/to/file1.mkv")

		// Verify the file is removed from tracking and completedFiles is incremented
		assert.Equal(t, 1, tracker.completedFiles)
		assert.Len(t, tracker.currentFiles, 1)
		_, exists := tracker.currentFiles["/path/to/file1.mkv"]
		assert.False(t, exists)
	})

	// Complete the second file
	t.Run("Complete second file", func(t *testing.T) {
		tracker.CompleteFile("/path/to/file2.mkv")

		// Verify the file is removed from tracking and completedFiles is incremented
		assert.Equal(t, 2, tracker.completedFiles)
		assert.Empty(t, tracker.currentFiles)
	})

	// Complete a non-existent file
	t.Run("Complete non-existent file", func(t *testing.T) {
		tracker.CompleteFile("/path/to/nonexistent.mkv")

		// Verify completedFiles is still incremented
		assert.Equal(t, 3, tracker.completedFiles)
	})
}

func TestProgressTracker_updateDescription(t *testing.T) {
	// Create a tracker
	tracker := NewProgressTracker(5)

	// Test with no active files
	t.Run("No active files", func(t *testing.T) {
		tracker.updateDescription()
		// We can't easily verify the description, but at least ensure it doesn't panic
	})

	// Test with some active files
	t.Run("Some active files", func(t *testing.T) {
		tracker.UpdateStage("/path/to/file1.mkv", StageBackup)
		tracker.UpdateStage("/path/to/file2.mkv", StageProcess)
		tracker.updateDescription()
		// We can't easily verify the description, but at least ensure it doesn't panic
	})

	// Test with many active files (more than 3)
	t.Run("Many active files", func(t *testing.T) {
		tracker.UpdateStage("/path/to/file3.mkv", StageValidate)
		tracker.UpdateStage("/path/to/file4.mkv", StageCleanup)
		tracker.UpdateStage("/path/to/file5.mkv", StageBackup)
		tracker.updateDescription()
		// We can't easily verify the description, but at least ensure it doesn't panic
	})
}