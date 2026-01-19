package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"vocal-genome-engine/services/audio-engine/audio"
	"vocal-genome-engine/services/audio-engine/dsp/pitch"
	"vocal-genome-engine/services/audio-engine/dsp/window"
)

/*
SineConfig defines the synthetic signal specification.
This is our ground-truth generator.
*/
type SineConfig struct {
	Frequency  float64 // Hz
	SampleRate float64 // Hz
	Duration   float64 // seconds
	Amplitude  float64 // 0.0 – 1.0
	NoiseLevel float64 // 0.0 = no noise
}

// GenerateSine generates a sine wave with optional white noise.
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
	// ---- Ground truth configuration ----
	cfg := SineConfig{
		Frequency:  440.0,
		SampleRate: 44100.0,
		Duration:   1.0,
		Amplitude:  0.8,
		NoiseLevel: 0.0,
	}

	// ---- Generate synthetic signal ----
	signal := GenerateSine(cfg)

	fmt.Println("Generated sine wave")
	fmt.Printf("Frequency: %.2f Hz\n", cfg.Frequency)
	fmt.Printf("Samples:   %d\n", len(signal))

	// ---- Framing configuration ----
	frameCfg := audio.FrameConfig{
		FrameSize: 1024,
		HopSize:   512,
	}

	frames := audio.FrameSignal(signal, frameCfg)
	hann := window.Hann(frameCfg.FrameSize)

	fmt.Println("\nRunning YIN on windowed frames...\n")

	// ---- Run YIN on first few frames ----
	for i := 0; i < 5 && i < len(frames); i++ {
		frame := frames[i]

		// Copy frame to avoid mutating original signal
		buf := make([]float64, len(frame))
		copy(buf, frame)

		// Apply Hann window
		audio.ApplyWindow(buf, hann)

		result := pitch.DetectYIN(
			buf,
			cfg.SampleRate,
			70.0,  // min frequency
			500.0, // max frequency
			0.1,   // YIN threshold
		)

		if result.Voiced {
			fmt.Printf(
				"Frame %d → F0: %.2f Hz | confidence: %.3f\n",
				i,
				result.F0,
				result.Confidence,
			)
		} else {
			fmt.Printf("Frame %d → unvoiced\n", i)
		}
	}
}
