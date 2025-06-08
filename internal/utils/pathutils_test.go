package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsertTagToFileName(t *testing.T) {
	testCases := []struct {
		name     string
		path     string
		tag      string
		expected string
	}{
		{
			name:     "Simple file",
			path:     "/home/user/video.mkv",
			tag:      "tagged",
			expected: "/home/user/video.tagged.mkv",
		},
		{
			name:     "File with dots in name",
			path:     "/home/user/video.file.mkv",
			tag:      "tagged",
			expected: "/home/user/video.file.tagged.mkv",
		},
		{
			name:     "File with no extension",
			path:     "/home/user/video",
			tag:      "tagged",
			expected: "/home/user/video.tagged",
		},
		{
			name:     "File with path and spaces",
			path:     "/home/user/my videos/video file.mkv",
			tag:      "tagged",
			expected: "/home/user/my videos/video file.tagged.mkv",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := InsertTagToFileName(tc.path, tc.tag)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestNFOPath(t *testing.T) {
	// Create a temporary test directory
	tmpDir, err := os.MkdirTemp("", "nfo-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Test case 1: Valid NFO file exists
	t.Run("Valid NFO exists", func(t *testing.T) {
		// Create a test video file
		videoPath := filepath.Join(tmpDir, "test-video.mkv")
		_, err = os.Create(videoPath)
		assert.NoError(t, err)

		// Create a matching NFO file
		nfoPath := filepath.Join(tmpDir, "test-video.nfo")
		_, err = os.Create(nfoPath)
		assert.NoError(t, err)

		// Test NFOPath function
		result, err := NFOPath(videoPath)
		assert.NoError(t, err)
		assert.Equal(t, nfoPath, result)
	})

	// Test case 2: NFO file doesn't exist
	t.Run("NFO doesn't exist", func(t *testing.T) {
		// Create a test video file without a matching NFO
		videoPath := filepath.Join(tmpDir, "test-video-no-nfo.mkv")
		_, err = os.Create(videoPath)
		assert.NoError(t, err)

		// Test NFOPath function
		result, err := NFOPath(videoPath)
		assert.Error(t, err)
		assert.Empty(t, result)
		assert.Contains(t, err.Error(), "does not exist")
	})

	// Test case 3: NFO path is a directory
	t.Run("NFO path is a directory", func(t *testing.T) {
		// Create a test video file
		videoPath := filepath.Join(tmpDir, "test-video-dir.mkv")
		_, err = os.Create(videoPath)
		assert.NoError(t, err)

		// Create a directory with the same name as the expected NFO file
		nfoDir := filepath.Join(tmpDir, "test-video-dir.nfo")
		err = os.Mkdir(nfoDir, 0755)
		assert.NoError(t, err)

		// Test NFOPath function
		result, err := NFOPath(videoPath)
		assert.Error(t, err)
		assert.Empty(t, result)
		assert.Contains(t, err.Error(), "is a directory")
	})
}

func TestGetFiles(t *testing.T) {
	// Create a temporary test directory
	tmpDir, err := os.MkdirTemp("", "files-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create a subdirectory
	subDir := filepath.Join(tmpDir, "subdir")
	err = os.Mkdir(subDir, 0755)
	assert.NoError(t, err)

	// Create various test files
	videoFiles := []string{
		filepath.Join(tmpDir, "video1.mkv"),
		filepath.Join(tmpDir, "video2.mp4"),
		filepath.Join(tmpDir, "video3.avi"),
		filepath.Join(subDir, "video4.wmv"),
		filepath.Join(subDir, "video5.mov"),
	}

	nonVideoFiles := []string{
		filepath.Join(tmpDir, "document.txt"),
		filepath.Join(tmpDir, "image.jpg"),
		filepath.Join(subDir, "data.json"),
	}

	// Create all the files
	for _, file := range append(videoFiles, nonVideoFiles...) {
		_, err = os.Create(file)
		assert.NoError(t, err)
	}

	// Test GetFiles function
	files, count, err := GetFiles(tmpDir)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, len(videoFiles), count)
	assert.Len(t, files, len(videoFiles))

	// Check that all video files are included
	for _, videoFile := range videoFiles {
		found := false
		for _, file := range files {
			if file == videoFile {
				found = true
				break
			}
		}
		assert.True(t, found, "Video file %s should be in the results", videoFile)
	}

	// Check that non-video files are not included
	for _, nonVideoFile := range nonVideoFiles {
		found := false
		for _, file := range files {
			if file == nonVideoFile {
				found = true
				break
			}
		}
		assert.False(t, found, "Non-video file %s should not be in the results", nonVideoFile)
	}

	// Test with non-existent directory
	files, count, err = GetFiles("/path/to/nonexistent/dir")
	assert.Error(t, err)
	assert.Nil(t, files)
	assert.Zero(t, count)
}
