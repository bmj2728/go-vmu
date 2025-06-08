package nfo

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseEpisodeNFO_ValidFile(t *testing.T) {
	// Get the absolute path to the test NFO file
	testDataDir, err := filepath.Abs("../../test-data")
	assert.NoError(t, err)

	validNFOPath := filepath.Join(testDataDir, "test-video.nfo")

	// Parse the valid NFO file
	details, err := ParseEpisodeNFO(validNFOPath)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, details)

	// Verify some key fields
	assert.Equal(t, "The Hidden Hand", details.Title)
	assert.Equal(t, "Dune: Prophecy", details.ShowTitle)
	assert.Equal(t, 1, details.Season)
	assert.Equal(t, 1, details.Episode)
	assert.Equal(t, 66, details.Runtime)
	assert.Equal(t, "tt10467954", details.IMDBID)
	assert.Equal(t, 2024, details.Year)

	// Verify arrays
	assert.Len(t, details.Genre, 3)
	assert.Contains(t, details.Genre, "Action")
	assert.Contains(t, details.Genre, "Adventure")
	assert.Contains(t, details.Genre, "Drama")

	// Verify actors
	assert.NotEmpty(t, details.Actor)
	assert.Equal(t, "Emily Watson", details.Actor[0].Name)
	assert.Equal(t, "Mother Superior Valya Harkonnen", details.Actor[0].Role)
}

func TestParseEpisodeNFO_InvalidFile(t *testing.T) {
	// Get the absolute path to the test NFO file
	testDataDir, err := filepath.Abs("../../test-data")
	assert.NoError(t, err)

	invalidNFOPath := filepath.Join(testDataDir, "test-video-invalid.nfo")

	// Parse the invalid NFO file
	details, err := ParseEpisodeNFO(invalidNFOPath)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, details)
	assert.Contains(t, err.Error(), NFOUnMarshallError)
}

func TestParseEpisodeNFO_NonExistentFile(t *testing.T) {
	// Try to parse a non-existent file
	details, err := ParseEpisodeNFO("/path/to/nonexistent/file.nfo")

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, details)
	assert.Contains(t, err.Error(), NFOReadError)
}

func TestMatchEpisodeFile_ValidFile(t *testing.T) {
	// Create a temporary test file
	tmpDir, err := os.MkdirTemp("", "nfo-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create a test video file
	videoPath := filepath.Join(tmpDir, "test-video.mkv")
	_, err = os.Create(videoPath)
	assert.NoError(t, err)

	// Create a matching NFO file
	nfoPath := filepath.Join(tmpDir, "test-video.nfo")
	_, err = os.Create(nfoPath)
	assert.NoError(t, err)

	// Test matching
	matchedPath, err := MatchEpisodeFile(videoPath)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, nfoPath, matchedPath)
}

func TestMatchEpisodeFile_MissingNFO(t *testing.T) {
	// Create a temporary test file
	tmpDir, err := os.MkdirTemp("", "nfo-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create a test video file without a matching NFO
	videoPath := filepath.Join(tmpDir, "test-video-no-nfo.mkv")
	_, err = os.Create(videoPath)
	assert.NoError(t, err)

	// Test matching
	matchedPath, err := MatchEpisodeFile(videoPath)

	// Assertions
	assert.Error(t, err)
	assert.Empty(t, matchedPath)
	assert.Contains(t, err.Error(), "does not exist")
}

func TestMatchEpisodeFile_NonExistentFile(t *testing.T) {
	// Try to match a non-existent file
	matchedPath, err := MatchEpisodeFile("/path/to/nonexistent/file.mkv")

	// Assertions
	assert.Error(t, err)
	assert.Empty(t, matchedPath)
	assert.Contains(t, err.Error(), NFONotFoundError)
}
