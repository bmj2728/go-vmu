package pool

// ProcessResult contains the result of processing a file
type ProcessResult struct {
	FilePath string
	Success  bool
	Error    error
}

func (r *ProcessResult) WithResult(success bool, err error) *ProcessResult {
	return &ProcessResult{
		FilePath: r.FilePath,
		Success:  success,
		Error:    err,
	}
}
