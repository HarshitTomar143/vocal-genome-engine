package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
	"vocal-genome-engine/services/audio-engine/dsp/pitch"

)

// This struct define the actual synthetic signal specification

type SineConfig struct {
	Frequency float64 //hz
	SampleRate float64 // hz
	Duration float64 // second
	Amplitude float64 // 0.0 - 1.0
	NoiseLevel float64 // 0.0 = no noise	
}

// This function generates Sine waves with optional white noise
func GenerateSine(cfg SineConfig)[]float64 {
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

		fmt.Println("\nRunning YIN pitch detection on first frame (1024 samples)...")

	frameSize := 1024
	if len(signal) < frameSize {
		panic("signal shorter than frame size")
	}

	result := pitch.DetectYIN(
		signal[:frameSize],
		cfg.SampleRate,
		70.0,   // min frequency (Hz)
		500.0,  // max frequency (Hz)
		0.1,    // YIN threshold
	)

	fmt.Println("\nYIN RESULT:")
	if result.Voiced {
		fmt.Printf("Detected F0: %.2f Hz\n", result.F0)
		fmt.Printf("Confidence: %.3f\n", result.Confidence)
	} else {
		fmt.Println("Frame detected as unvoiced")
	}

	return signal
}


// Main function

func main(){
	cfg := SineConfig{
		Frequency: 440.0,
		SampleRate: 44100.0,
		Duration: 1.0,
		Amplitude: 0.8,
		NoiseLevel: 0.0,
	}

	signal := GenerateSine(cfg)


	fmt.Printf("Generated sine wave: \n")
	fmt.Printf("Frequency: %.2f Hz\n", cfg.Frequency)
	fmt.Printf("Samples: %d\n", len(signal))
	fmt.Printf("First 10 samples:\n")

	for i:= 0; i< 10 ; i++ {
		fmt.Printf("%d: %6f\n",i,signal[i])
	}
}