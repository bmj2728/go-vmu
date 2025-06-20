package validator

import (
	"context"
	"errors"
	"github.com/rs/zerolog/log"
	"gopkg.in/vansante/go-ffprobe.v2"
	"time"
)

// MediaProberInterface defines the interface for media probing operations
type MediaProberInterface interface {
	Probe(path string) error
	DurationMinutes() float64
	VideoCodec() string
	VideoBitrate() string
	VideoHeight() int
	VideoWidth() int
	VideoAspectRatio() string
	AudioCodec() string
	AudioBitrate() string
	AudioChannels() int
	Size() string
}

// Ensure MediaProber implements MediaProberInterface
var _ MediaProberInterface = (*MediaProber)(nil)

type MediaProber struct {
	Context     context.Context
	CancelFn    context.CancelFunc
	Data        *ffprobe.ProbeData
	ProbeFailed bool
}

func NewMediaProber(timeout time.Duration) *MediaProber {
	ctx, cancelFn := context.WithTimeout(context.Background(), timeout)
	return &MediaProber{
		Context:  ctx,
		CancelFn: cancelFn,
	}
}

func (m *MediaProber) Probe(path string) error {
	defer m.CancelFn()
	data, err := ffprobe.ProbeURL(m.Context, path)
	if err != nil {
		m.ProbeFailed = true
		return err
	}
	m.Data = data
	return nil
}

func (m *MediaProber) DurationMinutes() float64 {
	if m.Data == nil || m.Data.Format == nil {
		return 0.0
	}
	return m.Data.Format.Duration().Minutes()
}

func (m *MediaProber) VideoCodec() string {
	if m.Data == nil {
		return ""
	}
	video := m.Data.FirstVideoStream()
	if video == nil {
		return ""
	}
	return video.CodecName
}

func (m *MediaProber) VideoBitrate() string {
	if m.Data == nil {
		return ""
	}
	video := m.Data.FirstVideoStream()
	if video == nil {
		return ""
	}

	return video.BitRate
}

func (m *MediaProber) VideoHeight() int {
	if m.Data == nil {
		return 0
	}
	video := m.Data.FirstVideoStream()
	if video == nil {
		return 0
	}
	return video.Height
}

func (m *MediaProber) VideoWidth() int {
	if m.Data == nil {
		return 0
	}
	video := m.Data.FirstVideoStream()
	if video == nil {
		return 0
	}
	return video.Width
}

func (m *MediaProber) VideoAspectRatio() string {
	if m.Data == nil {
		return ""
	}
	video := m.Data.FirstVideoStream()
	if video == nil {
		return ""
	}
	return video.DisplayAspectRatio
}

func (m *MediaProber) AudioCodec() string {
	if m.Data == nil {
		return ""
	}
	audio := m.Data.FirstAudioStream()
	if audio == nil {
		return ""
	}
	return audio.CodecName
}

func (m *MediaProber) AudioBitrate() string {
	if m.Data == nil {
		return ""
	}
	audio := m.Data.FirstAudioStream()
	if audio == nil {
		return ""
	}

	return audio.BitRate
}

func (m *MediaProber) AudioChannels() int {
	if m.Data == nil {
		return 0
	}
	audio := m.Data.FirstAudioStream()
	if audio == nil {
		return 0
	}
	return audio.Channels
}

func (m *MediaProber) Size() string {
	if m.Data == nil || m.Data.Format == nil {
		return ""
	}
	return m.Data.Format.Size
}

func (m *MediaProber) Tags() (ffprobe.Tags, error) {
	if m.Data == nil || m.Data.Format.TagList == nil {
		return nil, errors.New("No tags")
	}
	log.Debug().Msgf("Tags: %v", m.Data.Format.TagList)
	return m.Data.Format.TagList, nil
}
