package validator

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockMediaProber is a mock implementation of MediaProber for testing
type MockMediaProber struct {
	mock.Mock
}

func NewMockMediaProber() *MockMediaProber {
	return &MockMediaProber{}
}

func (m *MockMediaProber) Probe(path string) error {
	args := m.Called(path)
	return args.Error(0)
}

func (m *MockMediaProber) DurationMinutes() float64 {
	args := m.Called()
	return args.Get(0).(float64)
}

func (m *MockMediaProber) VideoCodec() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockMediaProber) VideoBitrate() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockMediaProber) VideoHeight() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockMediaProber) VideoWidth() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockMediaProber) VideoAspectRatio() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockMediaProber) AudioCodec() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockMediaProber) AudioBitrate() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockMediaProber) AudioChannels() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockMediaProber) Size() string {
	args := m.Called()
	return args.String(0)
}

func TestNewValidator(t *testing.T) {
	validator := NewValidator("old.mkv", "new.mkv", 30)

	assert.NotNil(t, validator)
	assert.Equal(t, "old.mkv", validator.oldFile)
	assert.Equal(t, "new.mkv", validator.newFile)
	assert.NotNil(t, validator.oldProber)
	assert.NotNil(t, validator.newProber)
}

func TestValidator_Validate_ProbeError(t *testing.T) {
	// Create a validator with mock probers
	oldProber := NewMockMediaProber()
	newProber := NewMockMediaProber()
	validator := &Validator{
		oldFile:   "old.mkv",
		newFile:   "new.mkv",
		oldProber: oldProber,
		newProber: newProber,
	}

	// Setup the old prober to return an error
	probeErr := errors.New("probe error")
	oldProber.On("Probe", "old.mkv").Return(probeErr)

	// Test validation
	err := validator.Validate()

	// Verify the error is returned
	assert.Error(t, err)
	assert.Equal(t, probeErr, err)
	oldProber.AssertExpectations(t)

	// Reset mocks
	oldProber = NewMockMediaProber()
	newProber = NewMockMediaProber()
	validator.oldProber = oldProber
	validator.newProber = newProber

	// Setup the old prober to succeed but new prober to fail
	oldProber.On("Probe", "old.mkv").Return(nil)
	newProber.On("Probe", "new.mkv").Return(probeErr)

	// Test validation
	err = validator.Validate()

	// Verify the error is returned
	assert.Error(t, err)
	assert.Equal(t, probeErr, err)
	oldProber.AssertExpectations(t)
	newProber.AssertExpectations(t)
}

func TestValidator_Validate_Success(t *testing.T) {
	// Create a validator with mock probers
	oldProber := NewMockMediaProber()
	newProber := NewMockMediaProber()
	validator := &Validator{
		oldFile:   "old.mkv",
		newFile:   "new.mkv",
		oldProber: oldProber,
		newProber: newProber,
	}

	// Setup the probers to return matching values
	oldProber.On("Probe", "old.mkv").Return(nil)
	newProber.On("Probe", "new.mkv").Return(nil)

	oldProber.On("VideoCodec").Return("h264")
	newProber.On("VideoCodec").Return("h264")

	oldProber.On("VideoBitrate").Return("1000000")
	newProber.On("VideoBitrate").Return("1000000")

	oldProber.On("VideoHeight").Return(1080)
	newProber.On("VideoHeight").Return(1080)

	oldProber.On("VideoWidth").Return(1920)
	newProber.On("VideoWidth").Return(1920)

	oldProber.On("VideoAspectRatio").Return("16:9")
	newProber.On("VideoAspectRatio").Return("16:9")

	oldProber.On("AudioCodec").Return("aac")
	newProber.On("AudioCodec").Return("aac")

	oldProber.On("AudioBitrate").Return("128000")
	newProber.On("AudioBitrate").Return("128000")

	oldProber.On("AudioChannels").Return(2)
	newProber.On("AudioChannels").Return(2)

	// Test validation
	err := validator.Validate()

	// Verify no error is returned
	assert.NoError(t, err)
	oldProber.AssertExpectations(t)
	newProber.AssertExpectations(t)
}

func TestValidator_Validate_Mismatches(t *testing.T) {
	testCases := []struct {
		name           string
		setupMocks     func(oldProber, newProber *MockMediaProber)
		expectedErrMsg string
	}{
		{
			name: "Video codec mismatch",
			setupMocks: func(oldProber, newProber *MockMediaProber) {
				oldProber.On("VideoCodec").Return("h264")
				newProber.On("VideoCodec").Return("h265")
			},
			expectedErrMsg: "video codec mismatch",
		},
		{
			name: "Video bitrate mismatch",
			setupMocks: func(oldProber, newProber *MockMediaProber) {
				oldProber.On("VideoCodec").Return("h264")
				newProber.On("VideoCodec").Return("h264")
				oldProber.On("VideoBitrate").Return("1000000")
				newProber.On("VideoBitrate").Return("2000000")
			},
			expectedErrMsg: "video bitrate mismatch",
		},
		{
			name: "Video height mismatch",
			setupMocks: func(oldProber, newProber *MockMediaProber) {
				oldProber.On("VideoCodec").Return("h264")
				newProber.On("VideoCodec").Return("h264")
				oldProber.On("VideoBitrate").Return("1000000")
				newProber.On("VideoBitrate").Return("1000000")
				oldProber.On("VideoHeight").Return(1080)
				newProber.On("VideoHeight").Return(720)
			},
			expectedErrMsg: "height mismatch",
		},
		{
			name: "Video width mismatch",
			setupMocks: func(oldProber, newProber *MockMediaProber) {
				oldProber.On("VideoCodec").Return("h264")
				newProber.On("VideoCodec").Return("h264")
				oldProber.On("VideoBitrate").Return("1000000")
				newProber.On("VideoBitrate").Return("1000000")
				oldProber.On("VideoHeight").Return(1080)
				newProber.On("VideoHeight").Return(1080)
				oldProber.On("VideoWidth").Return(1920)
				newProber.On("VideoWidth").Return(1280)
			},
			expectedErrMsg: "width mismatch",
		},
		{
			name: "Aspect ratio mismatch",
			setupMocks: func(oldProber, newProber *MockMediaProber) {
				oldProber.On("VideoCodec").Return("h264")
				newProber.On("VideoCodec").Return("h264")
				oldProber.On("VideoBitrate").Return("1000000")
				newProber.On("VideoBitrate").Return("1000000")
				oldProber.On("VideoHeight").Return(1080)
				newProber.On("VideoHeight").Return(1080)
				oldProber.On("VideoWidth").Return(1920)
				newProber.On("VideoWidth").Return(1920)
				oldProber.On("VideoAspectRatio").Return("16:9")
				newProber.On("VideoAspectRatio").Return("4:3")
			},
			expectedErrMsg: "aspect ratio mismatch",
		},
		{
			name: "Audio codec mismatch",
			setupMocks: func(oldProber, newProber *MockMediaProber) {
				oldProber.On("VideoCodec").Return("h264")
				newProber.On("VideoCodec").Return("h264")
				oldProber.On("VideoBitrate").Return("1000000")
				newProber.On("VideoBitrate").Return("1000000")
				oldProber.On("VideoHeight").Return(1080)
				newProber.On("VideoHeight").Return(1080)
				oldProber.On("VideoWidth").Return(1920)
				newProber.On("VideoWidth").Return(1920)
				oldProber.On("VideoAspectRatio").Return("16:9")
				newProber.On("VideoAspectRatio").Return("16:9")
				oldProber.On("AudioCodec").Return("aac")
				newProber.On("AudioCodec").Return("mp3")
			},
			expectedErrMsg: "audio codec mismatch",
		},
		{
			name: "Audio bitrate mismatch",
			setupMocks: func(oldProber, newProber *MockMediaProber) {
				oldProber.On("VideoCodec").Return("h264")
				newProber.On("VideoCodec").Return("h264")
				oldProber.On("VideoBitrate").Return("1000000")
				newProber.On("VideoBitrate").Return("1000000")
				oldProber.On("VideoHeight").Return(1080)
				newProber.On("VideoHeight").Return(1080)
				oldProber.On("VideoWidth").Return(1920)
				newProber.On("VideoWidth").Return(1920)
				oldProber.On("VideoAspectRatio").Return("16:9")
				newProber.On("VideoAspectRatio").Return("16:9")
				oldProber.On("AudioCodec").Return("aac")
				newProber.On("AudioCodec").Return("aac")
				oldProber.On("AudioBitrate").Return("128000")
				newProber.On("AudioBitrate").Return("192000")
			},
			expectedErrMsg: "audio bitrate mismatch",
		},
		{
			name: "Audio channels mismatch",
			setupMocks: func(oldProber, newProber *MockMediaProber) {
				oldProber.On("VideoCodec").Return("h264")
				newProber.On("VideoCodec").Return("h264")
				oldProber.On("VideoBitrate").Return("1000000")
				newProber.On("VideoBitrate").Return("1000000")
				oldProber.On("VideoHeight").Return(1080)
				newProber.On("VideoHeight").Return(1080)
				oldProber.On("VideoWidth").Return(1920)
				newProber.On("VideoWidth").Return(1920)
				oldProber.On("VideoAspectRatio").Return("16:9")
				newProber.On("VideoAspectRatio").Return("16:9")
				oldProber.On("AudioCodec").Return("aac")
				newProber.On("AudioCodec").Return("aac")
				oldProber.On("AudioBitrate").Return("128000")
				newProber.On("AudioBitrate").Return("128000")
				oldProber.On("AudioChannels").Return(2)
				newProber.On("AudioChannels").Return(6)
			},
			expectedErrMsg: "audio channels mismatch",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a validator with mock probers
			oldProber := NewMockMediaProber()
			newProber := NewMockMediaProber()
			validator := &Validator{
				oldFile:   "old.mkv",
				newFile:   "new.mkv",
				oldProber: oldProber,
				newProber: newProber,
			}

			// Setup the probers to return values
			oldProber.On("Probe", "old.mkv").Return(nil)
			newProber.On("Probe", "new.mkv").Return(nil)

			// Setup the specific test case mocks
			tc.setupMocks(oldProber, newProber)

			// Test validation
			err := validator.Validate()

			// Verify the expected error is returned
			assert.Error(t, err)
			assert.Equal(t, tc.expectedErrMsg, err.Error())
			oldProber.AssertExpectations(t)
			newProber.AssertExpectations(t)
		})
	}
}
