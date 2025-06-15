package tracker

type ProcessStatus int

const (
	StatusSuccess ProcessStatus = iota
	StatusFileNotFound
	StatusNFONotFound
	StatusNFOParseError
	StatusFFmpegError
	StatusValidationError
	StatusCleanupError
	StatusNetworkError // For those temporary blips
	StatusUnknownError
	StatusSkipped
)

func (ps ProcessStatus) String() string {
	switch ps {
	case StatusSuccess:
		return "Success"
	case StatusFileNotFound:
		return "FileNotFound"
	case StatusNFONotFound:
		return "NFONotFound"
	case StatusNFOParseError:
		return "NFOParseError"
	case StatusFFmpegError:
		return "FFmpegError"
	case StatusValidationError:
		return "ValidationError"
	case StatusCleanupError:
		return "CleanupError"
	case StatusNetworkError:
		return "NetworkError"
	case StatusUnknownError:
		return "UnknownError"
	case StatusSkipped:
		return "Skipped"
	default:
		return "UnknownStatus"
	}
}

// ProcessResult contains the result of processing a file
type ProcessResult struct {
	FilePath string        `json:"file_path"`
	Retries  int           `json:"retries,omitempty"` //to be used to determine if should retry this file
	Status   ProcessStatus `json:"status"`
	Success  bool          `json:"success"`
	Error    error         `json:"error,omitempty"`
}

type HumanReadableResult struct {
	FilePath string `json:"file_path"`
	Retries  int    `json:"retries,omitempty"`
	Status   string `json:"status"`
	Success  bool   `json:"success"`
	Error    string `json:"error,omitempty"`
}

func (r *ProcessResult) WithRetries(retries int) *ProcessResult {
	return &ProcessResult{
		FilePath: r.FilePath,
		Status:   r.Status,
		Success:  r.Success,
		Error:    r.Error,
		Retries:  retries,
	}
}

//add granular Withs for fields - can extend existing with WithStatus
//consider refactor of WithResult to maintain one-shot but add Status

func (r *ProcessResult) WithResult(success bool, err error) *ProcessResult {
	return &ProcessResult{
		FilePath: r.FilePath,
		Status:   r.Status,
		Success:  success,
		Error:    err,
		Retries:  r.Retries,
	}
}

func (r *ProcessResult) WithStatus(status ProcessStatus) *ProcessResult {
	return &ProcessResult{
		FilePath: r.FilePath,
		Status:   status,
		Success:  r.Success,
		Error:    r.Error,
		Retries:  r.Retries,
	}
}

func (r *ProcessResult) WithSuccess(success bool) *ProcessResult {
	return &ProcessResult{
		FilePath: r.FilePath,
		Status:   r.Status,
		Success:  success,
		Error:    r.Error,
		Retries:  r.Retries,
	}
}

func (r *ProcessResult) WithError(err error) *ProcessResult {
	return &ProcessResult{
		FilePath: r.FilePath,
		Status:   r.Status,
		Success:  r.Success,
		Error:    err,
		Retries:  r.Retries,
	}
}

func (r *ProcessResult) MakeHumanReadable() *HumanReadableResult {
	statusString := r.Status.String()
	errorString := ""
	if r.Error != nil {
		errorString = r.Error.Error()
	}

	return &HumanReadableResult{
		FilePath: r.FilePath,
		Retries:  r.Retries,
		Status:   statusString,
		Success:  r.Success,
		Error:    errorString,
	}
}
