package ffmpeg

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"os/exec"
)

type Executor struct {
	Command *FFmpegCommand
	backup  string
}

func NewExecutor(cmd *FFmpegCommand) *Executor {
	return &Executor{
		Command: cmd,
	}
}

func (e *Executor) Execute() error {
	ok, err := e.validArgs()
	if !ok || err != nil {
		log.Error().Err(err).Msg("Invalid arguments")
		return err
	}

	err = e.backupFile()
	if err != nil {
		log.Error().Err(err).Msg("Error copying file")
		return err
	}

	argString, err := e.Command.ArgsString()
	if err != nil {
		log.Error().Err(err).Msg("Error generating args")
		return err
	}

	command := exec.Command("ffmpeg", argString)
	err = command.Run()
	if err != nil {
		log.Error().Err(err).Msg("Error running command")
		return err
	}
	return nil
}

func (e *Executor) Validate() (bool, error) {
	//validate new file path
	//compare to backup w/ ffprobe?
	return true, nil
}

func (e *Executor) Cleanup() error {
	//remove backup file
	//mv new file to use original path
	return nil
}

func (e *Executor) validArgs() (bool, error) {
	//uh figure it out
	return true, nil
}

func (e *Executor) backupFile() error {
	newPath := "temp-" + e.Command.inputFile

	// Open the source file for reading
	sourceFile, err := os.Open(e.Command.inputFile)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer closeFile(sourceFile, "source")

	// Create the destination file
	destFile, err := os.Create(newPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer closeFile(destFile, "destination")

	// Copy data from source to destination
	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return fmt.Errorf("error copying file: %w", err)
	}

	e.backup = newPath

	return nil
}

func closeFile(file *os.File, fileType string) {
	if err := file.Close(); err != nil {
		log.Error().Err(err).Msgf("Error closing %s file", fileType)
	}
}

func (e *Executor) removeFile() error {
	return nil
}
