package ffmpeg

import (
	"testing"

	"github.com/bmj2728/go-vmu/internal/metadata"
	"github.com/stretchr/testify/assert"
)

func TestNewFFmpegCommand(t *testing.T) {
	cmd := NewFFmpegCommand()

	assert.NotNil(t, cmd)
	assert.Empty(t, cmd.inputFile)
	assert.Empty(t, cmd.outputFile)
	assert.Nil(t, cmd.metadata)
	assert.Nil(t, cmd.args)
}

func TestFFmpegCommand_WithInput(t *testing.T) {
	// Test with empty command
	cmd := NewFFmpegCommand()
	result := cmd.WithInput("/path/to/input.mkv")

	assert.NotNil(t, result)
	assert.Equal(t, "/path/to/input.mkv", result.inputFile)
	assert.Empty(t, result.outputFile)
	assert.Nil(t, result.metadata)
	assert.Nil(t, result.args)

	// Test with existing values
	cmd = &FFmpegCommand{
		outputFile: "/path/to/output.mkv",
		metadata:   map[string]interface{}{"title": "Test"},
		args:       []string{"-test"},
	}

	result = cmd.WithInput("/path/to/new_input.mkv")

	assert.NotNil(t, result)
	assert.Equal(t, "/path/to/new_input.mkv", result.inputFile)
	assert.Equal(t, "/path/to/output.mkv", result.outputFile)
	assert.Equal(t, map[string]interface{}{"title": "Test"}, result.metadata)
	assert.Equal(t, []string{"-test"}, result.args)
}

func TestFFmpegCommand_WithOutput(t *testing.T) {
	// Test with empty command
	cmd := NewFFmpegCommand()
	result := cmd.WithOutput("/path/to/output.mkv")

	assert.NotNil(t, result)
	assert.Empty(t, result.inputFile)
	assert.Equal(t, "/path/to/output.mkv", result.outputFile)
	assert.Nil(t, result.metadata)
	assert.Nil(t, result.args)

	// Test with existing values
	cmd = &FFmpegCommand{
		inputFile: "/path/to/input.mkv",
		metadata:  map[string]interface{}{"title": "Test"},
		args:      []string{"-test"},
	}

	result = cmd.WithOutput("/path/to/new_output.mkv")

	assert.NotNil(t, result)
	assert.Equal(t, "/path/to/input.mkv", result.inputFile)
	assert.Equal(t, "/path/to/new_output.mkv", result.outputFile)
	assert.Equal(t, map[string]interface{}{"title": "Test"}, result.metadata)
	assert.Equal(t, []string{"-test"}, result.args)
}

func TestFFmpegCommand_WithMetadata(t *testing.T) {
	// Create test metadata
	meta := &metadata.Metadata{
		Title:     "Test Title",
		Plot:      "Test Plot",
		Runtime:   120,
		ShowTitle: "Test Show",
		Season:    1,
		Episode:   2,
	}

	// Test with empty command
	cmd := NewFFmpegCommand()
	result, err := cmd.WithMetadata(*meta)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result.inputFile)
	assert.Empty(t, result.outputFile)
	assert.NotNil(t, result.metadata)
	assert.Nil(t, result.args)

	// Verify metadata fields
	assert.Equal(t, "Test Title", result.metadata["title"])
	assert.Equal(t, "Test Plot", result.metadata["plot"])
	assert.Equal(t, 120, result.metadata["runtime"])
	assert.Equal(t, "Test Show", result.metadata["showtitle"])
	assert.Equal(t, 1, result.metadata["season"])
	assert.Equal(t, 2, result.metadata["episode"])

	// Test with existing values
	cmd = &FFmpegCommand{
		inputFile:  "/path/to/input.mkv",
		outputFile: "/path/to/output.mkv",
		args:       []string{"-test"},
	}

	result, err = cmd.WithMetadata(*meta)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "/path/to/input.mkv", result.inputFile)
	assert.Equal(t, "/path/to/output.mkv", result.outputFile)
	assert.NotNil(t, result.metadata)
	assert.Equal(t, []string{"-test"}, result.args)

	// Test with empty metadata
	emptyMeta := metadata.Metadata{}
	result, err = cmd.WithMetadata(emptyMeta)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result.metadata)
}

func TestFFmpegCommand_GenerateArgs(t *testing.T) {
	// Test with input, output, and metadata
	cmd := &FFmpegCommand{
		inputFile:  "/path/to/input.mkv",
		outputFile: "/path/to/output.mkv",
		metadata: map[string]interface{}{
			"title":     "Test Title",
			"showtitle": "Test Show",
		},
	}

	result := cmd.GenerateArgs()

	assert.NotNil(t, result)
	assert.Equal(t, "/path/to/input.mkv", result.inputFile)
	assert.Equal(t, "/path/to/output.mkv", result.outputFile)
	assert.Equal(t, cmd.metadata, result.metadata)
	assert.NotNil(t, result.args)

	// Verify args contain expected values
	expectedArgs := []string{
		"-loglevel", "debug",
		"-i", "/path/to/input.mkv",
		"-c", "copy",
		"-metadata", "title=Test Title",
		"-metadata", "showtitle=Test Show",
		"/path/to/output.mkv",
	}

	// Since map iteration order is not guaranteed, we need to check that all expected args are present
	// rather than checking the exact order
	assert.Len(t, result.args, len(expectedArgs))
	assert.Contains(t, result.args, "-loglevel")
	assert.Contains(t, result.args, "debug")
	assert.Contains(t, result.args, "-i")
	assert.Contains(t, result.args, "/path/to/input.mkv")
	assert.Contains(t, result.args, "-c")
	assert.Contains(t, result.args, "copy")
	assert.Contains(t, result.args, "-metadata")
	assert.Contains(t, result.args, "/path/to/output.mkv")

	// Test with no metadata
	cmd = &FFmpegCommand{
		inputFile:  "/path/to/input.mkv",
		outputFile: "/path/to/output.mkv",
	}

	result = cmd.GenerateArgs()

	assert.NotNil(t, result)
	assert.NotNil(t, result.args)

	// Verify args contain expected values
	expectedArgs = []string{
		"-loglevel", "debug",
		"-i", "/path/to/input.mkv",
		"-c", "copy",
		"/path/to/output.mkv",
	}

	assert.Equal(t, expectedArgs, result.args)
}

func TestFFmpegCommand_ArgsString(t *testing.T) {
	// Test with args
	cmd := &FFmpegCommand{
		args: []string{"-loglevel", "debug", "-i", "/path/to/input.mkv", "-c", "copy", "/path/to/output.mkv"},
	}

	result, err := cmd.ArgsString()

	assert.NoError(t, err)
	assert.Equal(t, "-loglevel debug -i /path/to/input.mkv -c copy /path/to/output.mkv", result)

	// Test with nil args
	cmd = &FFmpegCommand{}

	result, err = cmd.ArgsString()

	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "args is nil")
}

func TestFFmpegCommand_Integration(t *testing.T) {
	// Test the full command building process
	meta := &metadata.Metadata{
		Title:   "Test Title",
		Runtime: 120,
		Season:  1,
		Episode: 2,
	}

	cmd := NewFFmpegCommand().
		WithInput("/path/to/input.mkv").
		WithOutput("/path/to/output.mkv")

	cmdWithMeta, err := cmd.WithMetadata(*meta)
	assert.NoError(t, err)

	cmdWithArgs := cmdWithMeta.GenerateArgs()

	argsString, err := cmdWithArgs.ArgsString()
	assert.NoError(t, err)

	// Verify the final command string contains all expected parts
	assert.Contains(t, argsString, "-loglevel debug")
	assert.Contains(t, argsString, "-i /path/to/input.mkv")
	assert.Contains(t, argsString, "-c copy")
	assert.Contains(t, argsString, "-metadata title=Test Title")
	assert.Contains(t, argsString, "-metadata runtime=120")
	assert.Contains(t, argsString, "-metadata season=1")
	assert.Contains(t, argsString, "-metadata episode=2")
	assert.Contains(t, argsString, "/path/to/output.mkv")
}
