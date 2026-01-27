package audio

// FrameConfig defines framing parameters
type FrameConfig struct {
	FrameSize int // samples per frame
	HopSize   int // samples between frames
}

// FrameSignal slices a signal into overlapping frames
func FrameSignal(signal []float64, cfg FrameConfig) [][]float64 {
	if cfg.FrameSize <= 0 {
		panic("FrameSize must be > 0")
	}
	if cfg.HopSize <= 0 {
		panic("HopSize must be > 0 (otherwise infinite loop)")
	}

	var frames [][]float64

	for start := 0; start+cfg.FrameSize <= len(signal); start += cfg.HopSize {
		frame := signal[start : start+cfg.FrameSize]
		frames = append(frames, frame)
	}

	return frames
}

// ApplyWindow multiplies a frame by a window (in-place)
func ApplyWindow(frame []float64, window []float64) {
	if len(frame) != len(window) {
		panic("Frame and window length mismatch")
	}

	for i := 0; i < len(frame); i++ {
		frame[i] *= window[i]
	}
}
