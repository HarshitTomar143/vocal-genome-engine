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

// Etracting interpretable traits from pitch frames
func 