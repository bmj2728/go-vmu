package tracker

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessResult_WithResult(t *testing.T) {
	// Test with success = true and no error
	result := &ProcessResult{FilePath: "/path/to/file.mkv"}
	newResult := result.WithResult(true, nil)

	assert.NotNil(t, newResult)
	assert.Equal(t, "/path/to/file.mkv", newResult.FilePath)
	assert.True(t, newResult.Success)
	assert.Nil(t, newResult.Error)

	// Test with success = false and an error
	testErr := errors.New("test error")
	result = &ProcessResult{FilePath: "/path/to/file.mkv"}
	newResult = result.WithResult(false, testErr)

	assert.NotNil(t, newResult)
	assert.Equal(t, "/path/to/file.mkv", newResult.FilePath)
	assert.False(t, newResult.Success)
	assert.Equal(t, testErr, newResult.Error)

	// Test that the original result is not modified
	result = &ProcessResult{FilePath: "/path/to/file.mkv", Success: true}
	newResult = result.WithResult(false, testErr)

	assert.True(t, result.Success)
	assert.Nil(t, result.Error)
	assert.False(t, newResult.Success)
	assert.Equal(t, testErr, newResult.Error)
}
