package ffmpeg

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"go-vmu/internal/tracker"
	"go-vmu/internal/utils"
	"go-vmu/internal/validator"
	"io"
	"os"
	"os/exec"
	"strings"
)

type Executor struct {
	FFmpegCommand   *FFmpegCommand
	Validator       *validator.Validator
	ProgressTracker *tracker.ProgressTracker
	backup          string
}

func NewExecutor(cmd *FFmpegCommand, tracker *tracker.ProgressTracker) *Executor {
	return &Executor{
		FFmpegCommand:   cmd,
		ProgressTracker: tracker,
	}
}

func (e *Executor) Execute() error {
	//validate args
	log.Debug().Msgf("Validating args: %v", e.FFmpegCommand.args)
	ok, err := e.validArgs()
	if !ok || err != nil {
		log.Error().Err(err).Msg("Invalid arguments")
		return err
	}

	//backup the file
	log.Debug().Msg("Backing up file")

	//update the tracker
	if e.ProgressTracker != nil {
		e.ProgressTracker.UpdateStage(e.FFmpegCommand.inputFile, tracker.StageBackup)
	}

	err = e.backupFile()
	if err != nil {
		log.Error().Err(err).Msg("Error copying file")
		return err
	}
	log.Debug().Msgf("File backed up to %s", e.backup)

	//execute the command
	command := exec.Command("ffmpeg", e.FFmpegCommand.args...)
	//log.Debug().Msgf("Executing command: %v", command.Args)

	quotedArgs := make([]string, len(e.FFmpegCommand.args))
	for i, arg := range e.FFmpegCommand.args {
		quotedArgs[i] = fmt.Sprintf("'%s'", strings.ReplaceAll(arg, "'", `'\''`)) // Basic shell-safe quoting for display
	}
	log.Info().Msgf("Executing FFmpeg command: ffmpeg %s", strings.Join(quotedArgs, " "))

	// Capture stdout and stderr
	var stdout, stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr

	//update the tracker
	if e.ProgressTracker != nil {
		e.ProgressTracker.UpdateStage(e.FFmpegCommand.inputFile, tracker.StageProcess)
	}

	err = command.Run()
	if err != nil {
		log.Error().Err(err).Msg("Error running command\n")
		//needs cleanup to revert file
		clErr := e.revertToBackup()
		if clErr != nil {
			log.Error().Err(clErr).Msg("Error reverting backup\n")
			return errors.Join(err, clErr)
		}
		clErr = e.removeOutputFile()
		if clErr != nil {
			log.Error().Err(clErr).Msg("Error removing output\n")
			return errors.Join(err, clErr)
		}
		clErr = e.removeBackupFile()
		if clErr != nil {
			log.Error().Err(clErr).Msg("Error removing backup file\n")
			return errors.Join(err, clErr)
		}
		return err
	}
	//everything went well
	log.Debug().Msg("FFmpegCommand executed successfully")
	return nil
}

func (e *Executor) ValidateNewFile() (bool, error) {
	e.Validator = validator.NewValidator(e.FFmpegCommand.inputFile, e.FFmpegCommand.outputFile, 300)
	//update the tracker
	if e.ProgressTracker != nil {
		e.ProgressTracker.UpdateStage(e.FFmpegCommand.inputFile, tracker.StageValidate)
	}
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
	//update the tracker
	if e.ProgressTracker != nil {
		e.ProgressTracker.UpdateStage(e.FFmpegCommand.inputFile, tracker.StageCleanup)
	}

	log.Debug().Msgf("Renaming %s to %s for cleanup.", e.FFmpegCommand.outputFile, e.FFmpegCommand.inputFile)
	err := os.Rename(e.FFmpegCommand.outputFile, e.FFmpegCommand.inputFile)
	if err != nil {
		log.Error().Err(err).Msgf("Error renaming output file to input file during cleanup: %s to %s", e.FFmpegCommand.outputFile, e.FFmpegCommand.inputFile)
		return fmt.Errorf("failed to replace original file with new file during cleanup: %w", err)
	}
	log.Debug().Msg("Original file replaced with updated file.")

	// remove the backup file
	err = e.removeBackupFile()
	if err != nil {
		log.Error().Err(err).Msg("Error removing backup file")
		return err
	}

	// reset the backup path
	e.backup = ""

	return nil
}

func (e *Executor) validArgs() (bool, error) {
	//check if args is nil
	if e.FFmpegCommand.args == nil {
		log.Error().Msg("args is nil")
		return false, fmt.Errorf("args is nil")
	}

	//check arg format
	// The ArgsString() method might still print quoted strings if it internally uses QuoteString,
	// but the args slice passed to exec.Command will be unquoted after the fix in command.go.
	argsString, err := e.FFmpegCommand.ArgsString()
	if err != nil {
		log.Error().Err(err).Msg("Error generating args string")
		return false, err
	}

	//check for dangerous characters and commands - we can probs do better
	// With exec.Command directly executing, shell injection via args is less of a concern.
	// These checks are still good for general input validation if the inputs come from untrusted sources.
	dangerousChars := []string{"&&", "rm -Rf", "|", ">", "<", "`", "$("} // Added more shell meta-characters
	for _, char := range dangerousChars {
		if strings.Contains(argsString, char) { // This check should be on the *raw* arguments, not the quoted string.
			// It's better to check individual arguments in the slice directly.
			// For simplicity, keeping it on argsString for now, but be aware of its limitations.
			return false, fmt.Errorf("dangerous character detected in arguments: %s", char)
		}
	}

	//confirm the input file was passed and is real
	if e.FFmpegCommand.inputFile == "" {
		return false, fmt.Errorf("input file is empty")
	}
	_, err = os.Stat(e.FFmpegCommand.inputFile)
	if err != nil {
		return false, fmt.Errorf("input file does not exist or is inaccessible: %w", err)
	}

	//output file dose not exist yet, so just ensure we have the value
	if e.FFmpegCommand.outputFile == "" {
		return false, fmt.Errorf("output file is empty")
	}

	//metadata is required
	if e.FFmpegCommand.metadata == nil {
		return false, fmt.Errorf("metadata is nil")
	}
	//woohoo we're good to go
	return true, nil
}

func (e *Executor) backupFile() error {

	newPath := utils.InsertTagToFileName(e.FFmpegCommand.inputFile, "backup")

	// Open the source file for reading
	sourceFile, err := os.Open(e.FFmpegCommand.inputFile)
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

	// Explicitly close files immediately after copy to release handles as early as possible.
	// The defers are still present as a fallback for other exit paths.
	if err := sourceFile.Close(); err != nil {
		log.Warn().Err(err).Msgf("Failed to explicitly close source file %s after backup copy.", sourceFile.Name())
	}
	if err := destFile.Close(); err != nil {
		log.Warn().Err(err).Msgf("Failed to explicitly close destination file %s after backup copy.", destFile.Name())
	}

	e.backup = newPath

	return nil
}

func closeFile(file *os.File, fileType string) {
	if err := file.Close(); err != nil && !errors.Is(err, os.ErrClosed) {
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
	err := os.Remove(e.FFmpegCommand.outputFile)
	if err != nil {
		log.Error().Err(err).Msg("Error removing output file")
		return err
	}
	return nil
}

func (e *Executor) revertToBackup() error {
	// Use os.Rename for atomic replacement if possible, safer than copy+overwrite
	log.Debug().Msgf("Reverting to backup: Renaming %s to %s", e.backup, e.FFmpegCommand.inputFile)
	err := os.Rename(e.backup, e.FFmpegCommand.inputFile)
	if err != nil {
		log.Error().Err(err).Msg("Error reverting to backup (rename failed). Attempting copy...")
		// Fallback to copy if rename fails (e.g., cross-device rename)
		// Open the source file for reading - the backup
		sourceFile, openErr := os.Open(e.backup)
		if openErr != nil {
			return fmt.Errorf("failed to open backup source file for revert: %w", errors.Join(err, openErr))
		}
		defer closeFile(sourceFile, "backup_source_for_revert")

		// open the destination file for writing
		destFile, createErr := os.OpenFile(e.FFmpegCommand.inputFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644) // Added O_CREATE in case original was deleted
		if createErr != nil {
			return fmt.Errorf("failed to open destination file for revert: %w", errors.Join(err, createErr))
		}
		defer closeFile(destFile, "original_dest_for_revert")

		// Copy data from source to destination
		if _, copyErr := io.Copy(destFile, sourceFile); copyErr != nil {
			return fmt.Errorf("error copying file during revert: %w", errors.Join(err, copyErr))
		}
		log.Debug().Msg("Successfully reverted to backup via copy.")

		// Try to remove backup after successful copy
		removeErr := e.removeBackupFile()
		if removeErr != nil {
			return errors.Join(err, removeErr) // Return original error and cleanup error
		}

		e.backup = "" // Reset backup path
		return nil    // Revert succeeded via copy
	}

	log.Debug().Msg("Successfully reverted to backup via rename.")
	e.backup = "" // Reset backup path since the file was moved/renamed

	return nil
}
