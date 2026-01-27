package features

import (
	"fmt"
	"vocal-genome-engine/services/audio-engine/audio"
	"vocal-genome-engine/services/audio-engine/dsp/pitch"
	"vocal-genome-engine/services/audio-engine/dsp/window"
)

type TrackerConfig struct {
	SampleRate float64
	FrameSize int
	HopSize int
	MinFreq float64
	MaxFreq float64
	Threshold float64
}

func TrackPitch(
	signal []float64,
	cfg TrackerConfig,
) []PitchFrame {

	fmt.Println("TrackPitch: entered")

	frameCfg := audio.FrameConfig{
		FrameSize: cfg.FrameSize,
		HopSize:   cfg.HopSize,
	}

	fmt.Println("TrackPitch: before FrameSignal")

	frames := audio.FrameSignal(signal, frameCfg)

	fmt.Println("TrackPitch: after FrameSignal, frames =", len(frames))

	win := window.Hann(cfg.FrameSize)

	fmt.Println("TrackPitch: window generated")

	var result []PitchFrame

	for i, frame := range frames {

		if i == 0 {
			fmt.Println("TrackPitch: processing first frame")
		}

		buf := make([]float64, len(frame))
		copy(buf, frame)

		audio.ApplyWindow(buf, win)

		yin := pitch.DetectYIN(
			buf,
			cfg.SampleRate,
			cfg.MinFreq,
			cfg.MaxFreq,
			cfg.Threshold,
		)

		result = append(result, PitchFrame{
			Time:       float64(i*cfg.HopSize+cfg.FrameSize/2) / cfg.SampleRate,
			F0:         yin.F0,
			Confidence: yin.Confidence,
			Voiced:     yin.Voiced,
		})
	}

	fmt.Println("TrackPitch: exiting")

	return result
}
