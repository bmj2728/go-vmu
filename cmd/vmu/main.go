package main

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"go-vmu/internal/logger"
	"go-vmu/internal/pool"
	"go-vmu/internal/processor"
	"os"
	"runtime"
)

//TODO - status I'm able to locally process a batch of files - nfs share is an issue
// what doesn't work right:
// - nfs share causes failures
// - intermittent file errors locally
// what works:
// - local file processing
// - cobra cli
// - logging
// - progress tracker

func main() {

	var workerCount int
	var verbose bool

	rootCmd := &cobra.Command{
		Use:   "vmu [directory]",
		Short: "Video Metadata Updater",
		Long:  "Update metadata in video files based on NFO files",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			// setup logger
			logger.Setup(logger.NewLoggerConfig(verbose))
			log.Info().Msgf("is verbose - %v", verbose)
			log.Info().Msg("Starting vmu")

			// Validate arguments

			//ensure sane worker count
			if workerCount < 1 {
				workerCount = 1
				log.Warn().Msg("Worker count must be greater than 0, defaulting to 1")
			}
			if workerCount > runtime.NumCPU() {
				workerCount = runtime.NumCPU()
				log.Warn().Msgf("Worker count must be less than %d, defaulting to %d", runtime.NumCPU(), runtime.NumCPU())
			}

			directory := args[0]

			// Validate directory
			path, err := os.Stat(directory)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			if !path.IsDir() {
				fmt.Printf("Error: %s is not a directory\n", directory)
				os.Exit(1)
			}
			if len(directory) == 0 {
				fmt.Printf("Error: directory is empty\n")
				os.Exit(1)
			}

			log.Info().Msgf("Processing directory: %s with %d workers\n", directory, workerCount)

			// Initialize processor
			proc := processor.NewProcessor(workerCount)

			// Process files
			results, err := proc.ProcessDirectory(directory)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}

			// Report results
			fmt.Printf("Processed %d files. Success: %d, Failed: %d\n",
				len(results),
				countSuccesses(results),
				len(results)-countSuccesses(results))
		},
	}

	// Define flags
	rootCmd.Flags().IntVarP(&workerCount, "workers", "w", runtime.NumCPU(), "Number of concurrent workers")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func countSuccesses(results []*pool.ProcessResult) int {
	count := 0
	for _, r := range results {
		if r.Success {
			count++
		}
	}
	return count
}
