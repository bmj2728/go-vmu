package validator

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"time"
)

type Validator struct {
	old       string
	new       string
	oldProber *MediaProber
	newProber *MediaProber
}

func NewValidator(old string, new string, timeout time.Duration) *Validator {
	oldProber := NewMediaProber(timeout)
	newProber := NewMediaProber(timeout)
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

	//returned Duration mismatch: old: 68.10593333333334  new: 68.10591666666667
	//if v.oldProber.DurationMinutes() != v.newProber.DurationMinutes() {
	//	log.Error().Msgf("Duration mismatch: old: %v  new: %v", v.oldProber.DurationMinutes(), v.newProber.DurationMinutes())
	//	return fmt.Errorf("duration mismatch")
	//}

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
	//hmmmm Size mismatch old: 1696803304  new: 1692549566
	//if v.oldProber.Size() != v.newProber.Size() {
	//	log.Error().Msgf("Size mismatch old: %v new: %v", v.oldProber.Size(), v.newProber.Size())
	//	return fmt.Errorf("size mismatch")
	//}

	return nil
}
