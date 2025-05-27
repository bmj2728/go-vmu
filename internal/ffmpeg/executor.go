package ffmpeg

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"go-vmu/internal/validator"
	"io"
	"os"
	"os/exec"
	"strings"
)

type Executor struct {
	Command   *FFmpegCommand
	Validator *validator.Validator
	backup    string
}

func NewExecutor(cmd *FFmpegCommand) *Executor {
	newValidator := validator.NewValidator(cmd.inputFile, cmd.outputFile, 10)
	return &Executor{
		Command:   cmd,
		Validator: newValidator,
	}
}

func (e *Executor) Execute() error {
	//validate args
	ok, err := e.validArgs()
	if !ok || err != nil {
		log.Error().Err(err).Msg("Invalid arguments")
		return err
	}

	//backup the file
	err = e.backupFile()
	if err != nil {
		log.Error().Err(err).Msg("Error copying file")
		return err
	}

	//execute the command
	command := exec.Command("ffmpeg", e.Command.args...)
	log.Info().Msgf("Executing command: %v", command.Args)
	err = command.Run()
	if err != nil {
		log.Error().Err(err).Msg("Error running command")
		//needs cleanup to revert file
		clErr := e.revertToBackup()
		if clErr != nil {
			log.Error().Err(clErr).Msg("Error cleaning up")
			return errors.Join(err, clErr)
		}
		clErr = e.removeOutputFile()
		if clErr != nil {
			log.Error().Err(clErr).Msg("Error cleaning up")
			return errors.Join(err, clErr)
		}
		clErr = e.removeBackupFile()
		if clErr != nil {
			log.Error().Err(clErr).Msg("Error cleaning up")
			return errors.Join(err, clErr)
		}
		return err
	}
	return nil
}

func (e *Executor) ValidateNewFile() (bool, error) {
	err := e.Validator.Validate()
	if err != nil {
		log.Error().Err(err).Msg("Error validating new file")
		//needs cleanup to revert file
		clErr := e.revertToBackup()
		if clErr != nil {
			log.Error().Err(clErr).Msg("Error cleaning up")
		}
		clErr = e.removeOutputFile()
		if clErr != nil {
			log.Error().Err(clErr).Msg("Error cleaning up")
		}
		clErr = e.removeBackupFile()
		if clErr != nil {
			log.Error().Err(clErr).Msg("Error cleaning up")
		}
		return false, errors.Join(err, clErr)
	}
	return true, nil
}

func (e *Executor) Cleanup() error {

	//open source file output
	sourceFile, err := os.Open(e.Command.outputFile)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer closeFile(sourceFile, "source")

	//open the destination file for writing
	destFile, err := os.Open(e.Command.inputFile)
	if err != nil {
		return fmt.Errorf("failed to open destination file: %w", err)
	}
	defer closeFile(destFile, "destination")

	//copy data from source to destination
	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return fmt.Errorf("error copying file: %w", err)
	}

	// remove the backup file
	err = e.removeBackupFile()
	if err != nil {
		log.Error().Err(err).Msg("Error removing backup file")
		return err
	}

	//remove the new file
	err = e.removeOutputFile()
	if err != nil {
		log.Error().Err(err).Msg("Error removing output file")
		return err
	}

	//reset the backup
	e.backup = ""

	return nil
}

func (e *Executor) validArgs() (bool, error) {
	//check if args is nil
	if e.Command.args == nil {
		log.Error().Msg("args is nil")
		return false, fmt.Errorf("args is nil")
	}

	//check arg format
	argsString, err := e.Command.ArgsString()
	if err != nil {
		log.Error().Err(err).Msg("Error generating args string")
		return false, err
	}

	//check for dangerous characters and commands - we can probs do better
	dangerousChars := []string{"|", ";", "&&", ">", "<", "`", "$", "(", ")"}
	for _, char := range dangerousChars {
		if strings.Contains(argsString, char) {
			return false, fmt.Errorf("dangerous character detected in arguments: %s", char)
		}
	}

	//confirm the input file was passed and is real
	if e.Command.inputFile == "" {
		return false, fmt.Errorf("input file is empty")
	}
	_, err = os.Stat(e.Command.inputFile)
	if err != nil {
		return false, err
	}

	//output file dose not exist yet, so just ensure we have the value
	if e.Command.outputFile == "" {
		return false, fmt.Errorf("output file is empty")
	}

	//metadata is required
	if e.Command.metadata == nil {
		return false, fmt.Errorf("metadata is nil")
	}
	//woohoo we're good to go
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

func (e *Executor) removeBackupFile() error {
	err := os.Remove(e.backup)
	if err != nil {
		log.Error().Err(err).Msg("Error removing backup file")
		return err
	}
	return nil
}

func (e *Executor) removeOutputFile() error {
	err := os.Remove(e.Command.outputFile)
	if err != nil {
		log.Error().Err(err).Msg("Error removing output file")
		return err
	}
	return nil
}

func (e *Executor) revertToBackup() error {

	// Open the source file for reading - the backup
	sourceFile, err := os.Open(e.backup)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer closeFile(sourceFile, "source")

	// open the destination file for writing
	destFile, err := os.Open(e.Command.inputFile)
	if err != nil {
		return fmt.Errorf("failed to open destination file: %w", err)
	}
	defer closeFile(destFile, "destination")

	// Copy data from source to destination
	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return fmt.Errorf("error copying file: %w", err)
	}

	err = e.removeBackupFile()
	if err != nil {
		log.Error().Err(err).Msg("Error removing backup file")
		return err
	}

	e.backup = ""

	return nil
}
