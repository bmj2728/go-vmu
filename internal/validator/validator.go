package validator

import (
	"fmt"
	"github.com/rs/zerolog/log"
)

type Validator struct {
	old       string
	new       string
	oldProber *MediaProber
	newProber *MediaProber
}

func NewValidator(old string, new string) *Validator {
	oldProber := NewMediaProber(10)
	newProber := NewMediaProber(10)
	return &Validator{
		old:       old,
		new:       new,
		oldProber: oldProber,
		newProber: newProber,
	}
}

func (v *Validator) Validate() error {
	err := v.oldProber.Probe(v.old)
	if err != nil {
		return err
	}
	err = v.newProber.Probe(v.new)
	if err != nil {
		return err
	}

	if v.oldProber.DurationMinutes() != v.newProber.DurationMinutes() {
		log.Error().Msg("Duration mismatch")
		return fmt.Errorf("duration mismatch")
	}

	if v.oldProber.VideoCodec() != v.newProber.VideoCodec() {
		log.Error().Msg("Video codec mismatch")
		return fmt.Errorf("video codec mismatch")
	}

	if v.oldProber.VideoBitrate() != v.newProber.VideoBitrate() {
		log.Error().Msg("Video bitrate mismatch")
		return fmt.Errorf("video bitrate mismatch")
	}

	if v.oldProber.VideoHeight() != v.newProber.VideoHeight() {
		log.Error().Msg("Height mismatch")
		return fmt.Errorf("height mismatch")
	}

	if v.oldProber.VideoWidth() != v.newProber.VideoWidth() {
		log.Error().Msg("Width mismatch")
		return fmt.Errorf("width mismatch")
	}

	if v.oldProber.VideoAspectRatio() != v.newProber.VideoAspectRatio() {
		log.Error().Msg("Aspect ratio mismatch")
		return fmt.Errorf("aspect ratio mismatch")
	}

	if v.oldProber.AudioCodec() != v.newProber.AudioCodec() {
		log.Error().Msg("Audio codec mismatch")
		return fmt.Errorf("audio codec mismatch")
	}

	if v.oldProber.AudioBitrate() != v.newProber.AudioBitrate() {
		log.Error().Msg("Audio bitrate mismatch")
		return fmt.Errorf("audio bitrate mismatch")
	}

	if v.oldProber.AudioChannels() != v.newProber.AudioChannels() {
		log.Error().Msg("Audio channels mismatch")
		return fmt.Errorf("audio channels mismatch")
	}

	if v.oldProber.Size() != v.newProber.Size() {
		log.Error().Msg("Size mismatch")
		return fmt.Errorf("size mismatch")
	}

	return nil
}
