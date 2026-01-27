package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"vocal-genome-engine/services/audio-engine/features"
)

type SineConfig struct {
	Frequency  float64
	SampleRate float64
	Duration   float64
	Amplitude  float64
	NoiseLevel float64
}

func GenerateSine(cfg SineConfig) []float64 {
	totalSamples := int(cfg.SampleRate * cfg.Duration)
	signal := make([]float64, totalSamples)

	rand.Seed(time.Now().UnixNano())

	for n := 0; n < totalSamples; n++ {
		t := float64(n) / cfg.SampleRate
		value := cfg.Amplitude * math.Sin(2*math.Pi*cfg.Frequency*t)

		if cfg.NoiseLevel > 0 {
			noise := cfg.NoiseLevel * (rand.Float64()*2 - 1)
			value += noise
		}

		signal[n] = value
	}

	return signal
}

func main() {
	cfg := SineConfig{
		Frequency:  440.0,
		SampleRate: 44100.0,
		Duration:   0.1,
		Amplitude:  0.8,
		NoiseLevel: 0.0,
	}

	signal := GenerateSine(cfg)

	fmt.Println("Generated sine wave")
	fmt.Printf("Frequency: %.2f Hz\n", cfg.Frequency)
	fmt.Printf("Samples:   %d\n", len(signal))

	trackerCfg := features.TrackerConfig{
		SampleRate: cfg.SampleRate,
		FrameSize:  1024,
		HopSize:    512,
		MinFreq:    70.0,
		MaxFreq:    500.0,
		Threshold:  0.1,
	}

	pitchFrames := features.TrackPitch(signal, trackerCfg)
	pitchFrames = features.MedianSmoothPitch(pitchFrames, 3)

	fmt.Println("\nPitch track (first voiced frames):")
	fmt.Printf("Total pitch frames: %d\n", len(pitchFrames))

	printed := 0
	for _, f := range pitchFrames {
		if f.Voiced {
			fmt.Printf(
				"t = %.3f s → F0 = %.2f Hz (confidence %.3f)\n",
				f.Time, f.F0, f.Confidence,
			)
			printed++
			if printed >= 10 {
				break
			}
		}
	}

	traits := features.ComputePitchTraits(pitchFrames)

	fmt.Println("\nPitch traits:")
	fmt.Printf("Mean pitch:     %.2f Hz\n", traits.MeanHz)
	fmt.Printf("Pitch range:    %.4f Hz\n", traits.RangeHz)
	fmt.Printf("Stability (σ):  %.6f Hz\n", traits.StabilityHz)
	fmt.Printf("Mean glide:     %.6f Hz\n", traits.MeanGlideHz)

	fmt.Println("\nProgram completed successfully.")
}
