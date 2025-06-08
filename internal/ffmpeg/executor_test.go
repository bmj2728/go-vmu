package ffmpeg

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bmj2728/go-vmu/internal/metadata"
	"github.com/bmj2728/go-vmu/internal/tracker"
	"github.com/bmj2728/go-vmu/internal/utils"
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

// TestExecutor_Execute tests the Execute method
func TestExecutor_Execute(t *testing.T) {
	// Create a temporary test directory
	tmpDir, err := os.MkdirTemp("", "ffmpeg-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create test files
	inputFile := filepath.Join(tmpDir, "input.mkv")
	outputFile := filepath.Join(tmpDir, "output.mkv")
	err = os.WriteFile(inputFile, []byte("test data"), 0644)
	assert.NoError(t, err)

	// Test with invalid args
	t.Run("Invalid args", func(t *testing.T) {
		cmd := NewFFmpegCommand().WithInput("").WithOutput(outputFile)
		executor := NewExecutor(cmd, nil)

		err := executor.Execute()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "args is nil")
	})

	// Mock the command execution for the valid case
	// This is a simplified test since we can't easily mock os/exec
	t.Run("Valid args", func(t *testing.T) {
		meta := &metadata.Metadata{
			Title: "Test Title",
		}
		cmd := NewFFmpegCommand().WithInput(inputFile).WithOutput(outputFile)
		cmdWithMeta, _ := cmd.WithMetadata(*meta)
		cmdWithArgs := cmdWithMeta.GenerateArgs()

		tracker := tracker.NewProgressTracker(1)
		executor := NewExecutor(cmdWithArgs, tracker)

		// We can't fully test the execution without mocking os/exec
		// But we can at least verify that the method doesn't panic
		// and handles the expected error case

		// The command will fail because ffmpeg isn't available in the test environment
		err := executor.Execute()
		assert.Error(t, err)
	})
}

// We skip TestExecutor_ValidateNewFile because it requires a real validator
// and we can't easily mock it without changing the code structure

// TestExecutor_Cleanup tests the Cleanup method
func TestExecutor_Cleanup(t *testing.T) {
	// Create a temporary test directory
	tmpDir, err := os.MkdirTemp("", "ffmpeg-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create test files
	inputFile := filepath.Join(tmpDir, "input.mkv")
	outputFile := filepath.Join(tmpDir, "output.mkv")
	backupFile := utils.InsertTagToFileName(inputFile, "backup")

	err = os.WriteFile(inputFile, []byte("test data"), 0644)
	assert.NoError(t, err)
	err = os.WriteFile(outputFile, []byte("test output data"), 0644)
	assert.NoError(t, err)
	err = os.WriteFile(backupFile, []byte("backup data"), 0644)
	assert.NoError(t, err)

	// Create a command
	meta := &metadata.Metadata{
		Title: "Test Title",
	}
	cmd := NewFFmpegCommand().WithInput(inputFile).WithOutput(outputFile)
	cmdWithMeta, _ := cmd.WithMetadata(*meta)
	cmdWithArgs := cmdWithMeta.GenerateArgs()

	// Test cleanup
	t.Run("Cleanup success", func(t *testing.T) {
		executor := NewExecutor(cmdWithArgs, nil)

		// Set the backup file path
		executor.backup = backupFile

		err := executor.Cleanup()
		assert.NoError(t, err)

		// Verify the backup file was removed
		_, err = os.Stat(backupFile)
		assert.True(t, os.IsNotExist(err))
	})
}

// TestExecutor_backupFile tests the backupFile method
func TestExecutor_backupFile(t *testing.T) {
	// Create a temporary test directory
	tmpDir, err := os.MkdirTemp("", "ffmpeg-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create a test file
	inputFile := filepath.Join(tmpDir, "input.mkv")
	err = os.WriteFile(inputFile, []byte("test data"), 0644)
	assert.NoError(t, err)

	// Create a command
	cmd := NewFFmpegCommand().WithInput(inputFile).WithOutput(filepath.Join(tmpDir, "output.mkv"))
	executor := NewExecutor(cmd, nil)

	// Test backup
	err = executor.backupFile()
	assert.NoError(t, err)

	// Verify the backup file was created
	backupFile := utils.InsertTagToFileName(inputFile, "backup")
	_, err = os.Stat(backupFile)
	assert.NoError(t, err)

	// Verify the backup file has the correct content
	content, err := os.ReadFile(backupFile)
	assert.NoError(t, err)
	assert.Equal(t, "test data", string(content))
}

// TestExecutor_removeBackupFile tests the removeBackupFile method
func TestExecutor_removeBackupFile(t *testing.T) {
	// Create a temporary test directory
	tmpDir, err := os.MkdirTemp("", "ffmpeg-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create a test file
	inputFile := filepath.Join(tmpDir, "input.mkv")
	backupFile := utils.InsertTagToFileName(inputFile, "backup")
	err = os.WriteFile(backupFile, []byte("backup data"), 0644)
	assert.NoError(t, err)

	// Create a command
	cmd := NewFFmpegCommand().WithInput(inputFile).WithOutput(filepath.Join(tmpDir, "output.mkv"))
	executor := NewExecutor(cmd, nil)
	executor.backup = backupFile

	// Test removing backup
	err = executor.removeBackupFile()
	assert.NoError(t, err)

	// Verify the backup file was removed
	_, err = os.Stat(backupFile)
	assert.True(t, os.IsNotExist(err))
}

// TestExecutor_removeOutputFile tests the removeOutputFile method
func TestExecutor_removeOutputFile(t *testing.T) {
	// Create a temporary test directory
	tmpDir, err := os.MkdirTemp("", "ffmpeg-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create a test file
	outputFile := filepath.Join(tmpDir, "output.mkv")
	err = os.WriteFile(outputFile, []byte("output data"), 0644)
	assert.NoError(t, err)

	// Create a command
	cmd := NewFFmpegCommand().WithInput(filepath.Join(tmpDir, "input.mkv")).WithOutput(outputFile)
	executor := NewExecutor(cmd, nil)

	// Test removing output
	err = executor.removeOutputFile()
	assert.NoError(t, err)

	// Verify the output file was removed
	_, err = os.Stat(outputFile)
	assert.True(t, os.IsNotExist(err))
}

// TestExecutor_revertToBackup tests the revertToBackup method
func TestExecutor_revertToBackup(t *testing.T) {
	// Create a temporary test directory
	tmpDir, err := os.MkdirTemp("", "ffmpeg-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create test files
	inputFile := filepath.Join(tmpDir, "input.mkv")
	backupFile := utils.InsertTagToFileName(inputFile, "backup")
	err = os.WriteFile(inputFile, []byte("original data"), 0644)
	assert.NoError(t, err)
	err = os.WriteFile(backupFile, []byte("backup data"), 0644)
	assert.NoError(t, err)

	// Create a command
	cmd := NewFFmpegCommand().WithInput(inputFile).WithOutput(filepath.Join(tmpDir, "output.mkv"))
	executor := NewExecutor(cmd, nil)
	executor.backup = backupFile

	// Test reverting to backup
	err = executor.revertToBackup()
	assert.NoError(t, err)

	// Verify the input file now has the backup content
	content, err := os.ReadFile(inputFile)
	assert.NoError(t, err)
	assert.Equal(t, "backup data", string(content))
}
