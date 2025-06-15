package utils

import (
	"encoding/json"
	"fmt"
	"github.com/bmj2728/go-vmu/internal/tracker"
	"os"
)

func CountSuccesses(results []*tracker.ProcessResult) int {
	count := 0
	for _, r := range results {
		if r.Success {
			count++
		}
	}
	return count
}

func GetStatusCounts(results []*tracker.ProcessResult) map[tracker.ProcessStatus]int {
	counts := make(map[tracker.ProcessStatus]int)
	for _, r := range results {
		if _, ok := counts[r.Status]; !ok {
			counts[r.Status] = 0
		}
		counts[r.Status]++
	}
	return counts
}

func PrintStatusCounts(counts map[tracker.ProcessStatus]int) {
	for status, count := range counts {
		fmt.Printf("Results:\n %s - %d\n", status, count)
	}
}

func SaveResults(filePath string, results []*tracker.ProcessResult) error {

	humanizedResults := make([]*tracker.HumanReadableResult, 0)
	for _, r := range results {
		humanizedResults = append(humanizedResults, r.MakeHumanReadable())
	}

	jsonData, err := json.MarshalIndent(humanizedResults, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal results: %w", err)
	}

	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write results file: %w", err)
	}

	return nil
}

func SaveFailures(filePath string, results []*tracker.ProcessResult) error {
	// Filter only failed results
	failures := make([]*tracker.HumanReadableResult, 0)
	for _, r := range results {
		if r.Status != tracker.StatusSuccess {
			failures = append(failures, r.MakeHumanReadable())
		}
	}

	jsonData, err := json.MarshalIndent(failures, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal failures: %w", err)
	}

	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write failures file: %w", err)
	}

	return nil
}
