package audio

// Defining Framing parameters
type FrameConfig struct {
	FrameSize int // sample per frame
	HopSize   int // samples between frames
}

func FrameSignal(signal []float64, cfg FrameConfig) [][]float64 {
	if cfg.FrameSize <= 0 || cfg.HopSize <=0 {
		panic("Invalid Frame Configuration")
	}

	var frames [][]float64

	for start := 0; start+ cfg.FrameSize <= len(signal); start += cfg.HopSize{
		frame := signal[start : start+cfg.FrameSize]
		frames = append(frames, frame)
	}

	return  frames
}

// Applying multiplies a frame by a window 
func ApplyWindow(frame []float64, window []float64){
	if len(frame) != len(frame){
		panic("Frame and window length mismatch")
	}

	for i := 0; i< len(frame); i++{
		frame[i] *= window[i]
	}
}