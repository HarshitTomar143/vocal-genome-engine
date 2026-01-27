package features

import "sort"



func MedianSmoothPitch(frames []PitchFrame, width int) []PitchFrame {
	if width < 3 || width%2 == 0 {
		return frames
	}

	half := width/2
	out := make([]PitchFrame, len(frames))
	copy(out, frames)

	for i:= half; i<len(frames)-half; i++ {
		var vals []float64

		for j := i-half; j <= i+half; j++ {
			if frames[j].Voiced {
				vals = append(vals, frames[j].F0)
			}
		}

		if len(vals) == 0 {
			out[i].Voiced = false
			continue
		}

		sort.Float64s(vals)
		out[i].F0 = vals[len(vals)/2]
		out[i].Voiced = true
	} 

	return out
}