package features

import "math"

type PitchFrame struct {
	Time float64
	F0 float64
	Confidence float64
	Voiced bool
}

type PitchTraits struct {
	RangeHz float64
	StabilityHz float64
	MeanHz float64
	MeanGlideHz float64
}

// Extracting interpretable traits from pitch frames
func ComputePitchTraits(frames []PitchFrame) PitchTraits {
	var voiced []float64
	
	for _, f := range frames {
		if f.Voiced && f.Confidence >= 0.85 {
			voiced = append(voiced, f.F0)
		}
	}

	if len(voiced) < 2 {
		return PitchTraits{}
	}

	minF, maxF := voiced[0], voiced[0]
	var sum float64
	
	for _, f := range voiced{
		sum += f
		if f<minF {
			minF= f
		}
		if f>maxF {
			maxF= f
		}
	}

	mean := sum/ float64(len(voiced))

	var variance float64
	for _, f := range voiced {
		d := f - mean
		variance += d * d

	}

	variance /= float64(len(voiced))


	var glideSum float64
	for i := 1; i<len(voiced); i++ {
		glideSum += math.Abs(voiced[i]- voiced[i-1])
	}

	return PitchTraits{
		RangeHz: maxF - minF,
		StabilityHz: math.Sqrt(variance),
		MeanHz: mean,
		MeanGlideHz: glideSum / float64(len(voiced)-1),
	}

}