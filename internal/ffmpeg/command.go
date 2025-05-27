package ffmpeg

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"go-vmu/internal/metadata"
	"strings"
)

type FFmpegCommand struct {
	inputFile  string
	outputFile string
	metadata   map[string]interface{}
	args       []string
	// Other options if we need them
}

func NewFFmpegCommand() *FFmpegCommand {
	return &FFmpegCommand{}
}

func (cmd *FFmpegCommand) WithInput(input string) *FFmpegCommand {
	return &FFmpegCommand{
		inputFile:  input,
		outputFile: cmd.outputFile,
		metadata:   cmd.metadata,
		args:       cmd.args,
	}
}

func (cmd *FFmpegCommand) WithOutput(output string) *FFmpegCommand {
	return &FFmpegCommand{
		inputFile:  cmd.inputFile,
		outputFile: output,
		metadata:   cmd.metadata,
		args:       cmd.args,
	}
}

func (cmd *FFmpegCommand) WithMetadata(meta metadata.Metadata) (*FFmpegCommand, error) {

	metaFields, err := meta.ToMap()

	if err != nil {
		log.Error().Err(err).Msg("Error converting metadata to map")
		return nil, fmt.Errorf("error converting metadata to map: %v", err)
	}

	return &FFmpegCommand{
		inputFile:  cmd.inputFile,
		outputFile: cmd.outputFile,
		metadata:   metaFields,
		args:       cmd.args,
	}, nil
}

func (cmd *FFmpegCommand) GenerateArgs() *FFmpegCommand {
	args := []string{"-i", cmd.inputFile, "-c", "copy"}

	for key, value := range cmd.metadata {
		args = append(args, "-metadata", fmt.Sprintf("%s=%s", key, value))
	}

	args = append(args, cmd.outputFile)

	return &FFmpegCommand{
		inputFile:  cmd.inputFile,
		outputFile: cmd.outputFile,
		metadata:   cmd.metadata,
		args:       args,
	}
}

func (cmd *FFmpegCommand) ArgsString() (string, error) {

	if cmd.args == nil {
		return "", fmt.Errorf("args is nil")
	}

	argString := strings.Join(cmd.args, " ")

	return argString, nil
}
