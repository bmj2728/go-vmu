package validator

import (
	"context"
	"gopkg.in/vansante/go-ffprobe.v2"
	"time"
)

type MediaProber struct {
	Context     context.Context
	CancelFn    context.CancelFunc
	Data        *ffprobe.ProbeData
	ProbeFailed bool
}

func NewMediaProber(timeout time.Duration) *MediaProber {
	ctx, cancelFn := context.WithTimeout(context.Background(), timeout*time.Second)
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
	return m.Data.Format.Duration().Minutes()
}

func (m *MediaProber) VideoCodec() string {
	video := m.Data.FirstVideoStream()
	if video == nil {
		return ""
	}
	return video.CodecName
}

func (m *MediaProber) VideoBitrate() string {
	video := m.Data.FirstVideoStream()
	if video == nil {
		return ""
	}

	return video.BitRate
}

func (m *MediaProber) VideoHeight() int {
	video := m.Data.FirstVideoStream()
	if video == nil {
		return 0
	}
	return video.Height
}

func (m *MediaProber) VideoWidth() int {
	video := m.Data.FirstVideoStream()
	if video == nil {
		return 0
	}
	return video.Width
}

func (m *MediaProber) VideoAspectRatio() string {
	video := m.Data.FirstVideoStream()
	if video == nil {
		return ""
	}
	return video.DisplayAspectRatio
}

func (m *MediaProber) AudioCodec() string {
	audio := m.Data.FirstAudioStream()
	if audio == nil {
		return ""
	}
	return audio.CodecName
}

func (m *MediaProber) AudioBitrate() string {
	audio := m.Data.FirstAudioStream()
	if audio == nil {
		return ""
	}

	return audio.BitRate
}

func (m *MediaProber) AudioChannels() int {
	audio := m.Data.FirstAudioStream()
	if audio == nil {
		return 0
	}
	return audio.Channels
}

func (m *MediaProber) Size() string {
	return m.Data.Format.Size
}
