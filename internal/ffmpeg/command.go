package ffmpeg

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"go-vmu/internal/metadata"
)

type FFmpegCommand struct {
	inputFile  string
	outputFile string
	metadata   map[string]interface{}
	// Other options if we need them
}

func NewFFmpegCommand() *FFmpegCommand {
	return &FFmpegCommand{}
}

func (cmd *FFmpegCommand) WithInput(input string) *FFmpegCommand {
	cmd.inputFile = input
	return &FFmpegCommand{
		inputFile:  input,
		outputFile: cmd.outputFile,
		metadata:   cmd.metadata,
	}
}

func (cmd *FFmpegCommand) WithOutput(output string) *FFmpegCommand {
	cmd.outputFile = output
	return &FFmpegCommand{
		inputFile:  cmd.inputFile,
		outputFile: output,
		metadata:   cmd.metadata,
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
	}, nil
}

func (cmd *FFmpegCommand) Build() []string {
	args := []string{"-i", cmd.inputFile, "-c", "copy"}

	for key, value := range cmd.metadata {
		args = append(args, "-metadata", fmt.Sprintf("%s=%s", key, value))
	}

	args = append(args, cmd.outputFile)
	return args
}
