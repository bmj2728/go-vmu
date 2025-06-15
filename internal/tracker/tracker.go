package tracker

import (
	"fmt"
	//"github.com/bmj2728/go-vmu/internal/pool"
	"github.com/schollz/progressbar/v3"
	"path/filepath"
	"sync"
)

// Define stages for better tracking
const (
	StageBackup   = "Backup"
	StageProcess  = "Process"
	StageValidate = "Validate"
	StageCleanup  = "Cleanup"
)

// Progress tracking structure
type ProgressTracker struct {
	totalFiles     int
	completedFiles int
	currentFiles   map[string]*FileProgress
	Results        []*ProcessResult
	mu             sync.Mutex
	bar            *progressbar.ProgressBar
}

type FileProgress struct {
	filename string
	stage    string
	done     bool
}

// Create a new progress tracker
func NewProgressTracker(totalFiles int) *ProgressTracker {
	return &ProgressTracker{
		totalFiles:   totalFiles,
		currentFiles: make(map[string]*FileProgress),
		bar:          progressbar.Default(int64(totalFiles)),
	}
}

// Update stage for a file
func (p *ProgressTracker) UpdateStage(filename, stage string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if progress, exists := p.currentFiles[filename]; exists {
		progress.stage = stage
	} else {
		p.currentFiles[filename] = &FileProgress{
			filename: filename,
			stage:    stage,
			done:     false,
		}
	}

	// Update description to show active files and their stages
	p.updateDescription()
}

// Mark a file as complete
func (p *ProgressTracker) CompleteFile(filename string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if progress, exists := p.currentFiles[filename]; exists {
		progress.done = true
		delete(p.currentFiles, filename) // Remove from active tracking
	}

	p.completedFiles++
	err := p.bar.Add(1)
	if err != nil {
		return
	}
	p.updateDescription()
}

// Update the progress bar description to show active files
func (p *ProgressTracker) updateDescription() {
	// Build description showing active files (limit to 2-3 to avoid clutter)
	desc := fmt.Sprintf("%d/%d complete", p.completedFiles, p.totalFiles)

	activeCount := 0
	for _, progress := range p.currentFiles {
		if activeCount < 3 { // Show max 3 active files
			desc += fmt.Sprintf("\n%s: %s", filepath.Base(progress.filename), progress.stage)
			activeCount++
		}
	}

	if activeCount > 0 && len(p.currentFiles) > 3 {
		desc += fmt.Sprintf(" (+ %d more)", len(p.currentFiles)-3)
	}

	p.bar.Describe(desc)
}

func (p *ProgressTracker) AppendResult(result *ProcessResult) {
	p.Results = append(p.Results, result)
}
