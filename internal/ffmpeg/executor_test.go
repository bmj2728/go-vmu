package ffmpeg

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bmj2728/go-vmu/internal/metadata"
	"github.com/stretchr/testify/assert"
)

// TestExecutor_validArgs tests the validArgs method indirectly through a test-only exported method
func TestExecutor_validArgs(t *testing.T) {
	// Create a temporary test file
	tmpDir, err := os.MkdirTemp("", "ffmpeg-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	inputFile := filepath.Join(tmpDir, "input.mkv")
	outputFile := filepath.Join(tmpDir, "output.mkv")

	// Create the input file
	err = os.WriteFile(inputFile, []byte("test data"), 0644)
	assert.NoError(t, err)

	// Create metadata
	meta := &metadata.Metadata{
		Title: "Test Title",
	}

	// Test cases
	testCases := []struct {
		name        string
		setupCmd    func() *FFmpegCommand
		expectValid bool
		expectError string
	}{
		{
			name: "Valid command",
			setupCmd: func() *FFmpegCommand {
				cmd := NewFFmpegCommand().WithInput(inputFile).WithOutput(outputFile)
				cmdWithMeta, _ := cmd.WithMetadata(*meta)
				return cmdWithMeta.GenerateArgs()
			},
			expectValid: true,
			expectError: "",
		},
		{
			name: "Nil args",
			setupCmd: func() *FFmpegCommand {
				cmd := NewFFmpegCommand().WithInput(inputFile).WithOutput(outputFile)
				cmdWithMeta, _ := cmd.WithMetadata(*meta)
				// Don't generate args
				return cmdWithMeta
			},
			expectValid: false,
			expectError: "args is nil",
		},
		{
			name: "Empty input file",
			setupCmd: func() *FFmpegCommand {
				cmd := NewFFmpegCommand().WithInput("").WithOutput(outputFile)
				cmdWithMeta, _ := cmd.WithMetadata(*meta)
				return cmdWithMeta.GenerateArgs()
			},
			expectValid: false,
			expectError: "input file is empty",
		},
		{
			name: "Non-existent input file",
			setupCmd: func() *FFmpegCommand {
				cmd := NewFFmpegCommand().WithInput("/non/existent/file.mkv").WithOutput(outputFile)
				cmdWithMeta, _ := cmd.WithMetadata(*meta)
				return cmdWithMeta.GenerateArgs()
			},
			expectValid: false,
			expectError: "input file does not exist",
		},
		{
			name: "Empty output file",
			setupCmd: func() *FFmpegCommand {
				cmd := NewFFmpegCommand().WithInput(inputFile).WithOutput("")
				cmdWithMeta, _ := cmd.WithMetadata(*meta)
				return cmdWithMeta.GenerateArgs()
			},
			expectValid: false,
			expectError: "output file is empty",
		},
		{
			name: "Nil metadata",
			setupCmd: func() *FFmpegCommand {
				// Create a command without metadata
				cmd := NewFFmpegCommand().WithInput(inputFile).WithOutput(outputFile)
				return cmd.GenerateArgs()
			},
			expectValid: false,
			expectError: "metadata is nil",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			cmd := tc.setupCmd()
			executor := NewExecutor(cmd, nil)

			// Test validArgs through our test wrapper
			valid, err := validateArgs(executor)

			// Assertions
			if tc.expectValid {
				assert.True(t, valid)
				assert.NoError(t, err)
			} else {
				assert.False(t, valid)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectError)
			}
		})
	}
}

// validateArgs is a test-only wrapper for the private validArgs method
func validateArgs(e *Executor) (bool, error) {
	return e.validArgs()
}
