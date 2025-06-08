package validator

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewMediaProber(t *testing.T) {
	// Test with different timeouts
	testCases := []struct {
		name    string
		timeout time.Duration
	}{
		{
			name:    "Zero timeout",
			timeout: 0,
		},
		{
			name:    "Short timeout",
			timeout: 1 * time.Second,
		},
		{
			name:    "Long timeout",
			timeout: 5 * time.Minute,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			prober := NewMediaProber(tc.timeout)

			assert.NotNil(t, prober)
			assert.NotNil(t, prober.Context)
			assert.NotNil(t, prober.CancelFn)
			assert.Nil(t, prober.Data)
			assert.False(t, prober.ProbeFailed)

			// Verify the context has the correct timeout
			deadline, ok := prober.Context.Deadline()
			assert.True(t, ok)
			expectedDeadline := time.Now().Add(tc.timeout)
			assert.WithinDuration(t, expectedDeadline, deadline, 100*time.Millisecond)
		})
	}
}

func TestMediaProber_Probe_Error(t *testing.T) {
	// Create a prober with a short timeout
	prober := NewMediaProber(1 * time.Second)

	// Test with a non-existent file
	err := prober.Probe("/non/existent/file.mkv")

	// Verify the error and state
	assert.Error(t, err)
	assert.True(t, prober.ProbeFailed)
	assert.Nil(t, prober.Data)
}

func TestMediaProber_Methods_NilData(t *testing.T) {
	// Create a prober
	prober := NewMediaProber(1 * time.Second)

	// Test methods with nil data
	// These should not panic, but return zero values

	// We need to set up a deferred recover to catch any panics
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Method panicked: %v", r)
		}
	}()

	// Test all methods
	assert.Equal(t, 0.0, prober.DurationMinutes())
	assert.Equal(t, "", prober.VideoCodec())
	assert.Equal(t, "", prober.VideoBitrate())
	assert.Equal(t, 0, prober.VideoHeight())
	assert.Equal(t, 0, prober.VideoWidth())
	assert.Equal(t, "", prober.VideoAspectRatio())
	assert.Equal(t, "", prober.AudioCodec())
	assert.Equal(t, "", prober.AudioBitrate())
	assert.Equal(t, 0, prober.AudioChannels())
	assert.Equal(t, "", prober.Size())
}
