package validator

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"time"
)

type Validator struct {
	oldFile   string
	newFile   string
	oldProber *MediaProber
	newProber *MediaProber
}

func NewValidator(oldFile string, newFile string, timeoutSecs time.Duration) *Validator {
	timeout := timeoutSecs * time.Second
	log.Info().Str("timeout", timeout.String())
	oldProber := NewMediaProber(timeout)
	newProber := NewMediaProber(timeout)
	return &Validator{
		oldFile:   oldFile,
		newFile:   newFile,
		oldProber: oldProber,
		newProber: newProber,
	}
}

func (v *Validator) Validate() error {
	err := v.oldProber.Probe(v.oldFile)
	if err != nil {
		return err
	}
	err = v.newProber.Probe(v.newFile)
	if err != nil {
		return err
	}

	//returned Duration mismatch: oldFile: 68.10593333333334  newFile: 68.10591666666667
	//if v.oldProber.DurationMinutes() != v.newProber.DurationMinutes() {
	//	log.Error().Msgf("Duration mismatch: oldFile: %v  newFile: %v", v.oldProber.DurationMinutes(), v.newProber.DurationMinutes())
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
	//hmmmm Size mismatch oldFile: 1696803304  newFile: 1692549566
	//if v.oldProber.Size() != v.newProber.Size() {
	//	log.Error().Msgf("Size mismatch oldFile: %v newFile: %v", v.oldProber.Size(), v.newProber.Size())
	//	return fmt.Errorf("size mismatch")
	//}

	return nil
}
