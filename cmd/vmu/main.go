package main

import (
	"fmt"
	"github.com/bmj2728/go-vmu/internal/logger"
	"github.com/bmj2728/go-vmu/internal/processor"
	"github.com/bmj2728/go-vmu/internal/utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"runtime"
)

func main() {

	var workerCount int
	var verbose bool
	var retries int
	var resultsPath string
	var saveResults bool

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
			//set sane retry attempts val
			if retries < 0 {
				retries = 0
			}
			if retries > 5 {
				retries = 5
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

			//if no location don't try to save
			if saveResults && resultsPath == "" {
				resultsPath = directory
			}

			log.Info().Msgf("Processing directory: %s with %d workers\n", directory, workerCount)

			// Initialize processor
			proc := processor.NewProcessor(workerCount)

			// Process files
			results, err := proc.ProcessDirectory(directory, retries)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}

			// Report results
			fmt.Printf("Processed %d files. Success: %d, Failed: %d\n",
				len(results),
				utils.CountSuccesses(results),
				len(results)-utils.CountSuccesses(results))
			counts := utils.GetStatusCounts(results)
			utils.PrintStatusCounts(counts)

			if saveResults {
				err = utils.SaveResults(filepath.Join(resultsPath, "results.json"), results)
				if err != nil {
					log.Error().Msgf("Error saving results: %v", err)
				}
				err = utils.SaveFailures(filepath.Join(resultsPath, "failures.json"), results)
				if err != nil {
					log.Error().Msgf("Error saving results: %v", err)
				}
			}
		},
	}

	// Define flags
	rootCmd.Flags().IntVarP(&workerCount, "workers", "w", runtime.NumCPU(), "Number of concurrent workers(1-#CPUs)")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	rootCmd.Flags().IntVarP(&retries, "retries", "r", 3, "Number of retries (0-5)")
	rootCmd.Flags().BoolVarP(&saveResults, "save", "s", false, "Save results to file - results.json/failures.json in directory. If no path is specified, results will be saved to processed directory.")
	rootCmd.Flags().StringVarP(&resultsPath, "path", "p", "", "Path to directory to save results")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
