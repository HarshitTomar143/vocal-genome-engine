package pitch

type Result struct {
	F0 float64
	Confidence float64
	Voiced bool
}

// This function computes the yin difference function

func differece(signal []float64, maxLag int) []float64 {
	n := len(signal)
	diff := make([]float64, maxLag)

	for tau := 1; tau < maxLag ; tau++ {
		var sum float64
		for i := 0; i<n-tau; i++ {
			delta := signal[i] - signal[i+tau]
			sum += delta * delta
		}

		diff[tau] = sum
	}

	return diff
}

// This function normalizes the difference function
func cumulativeMeanNormalizedDifference(diff []float64) []float64 {
	cmnd := make([]float64, len(diff))
	cmnd[0] = 1.0

	var runningSum float64
	for tau := 1; tau < len(diff); tau++ {
		runningSum += diff[tau]
		cmnd[tau] = diff[tau] * float64(tau)/ runningSum
	}

	return cmnd
}

// This function finds the first minimum threshold
func absoluteThreshold(cmnd []float64, threshold float64) int {
	for tau := 2; tau < len(cmnd); tau++ {
		if cmnd[tau] < threshold {
			return tau
		}
	}
	return -1
}


// This function refines lag estimate
func parabolicInterpolation(cmnd []float64, tau int) float64 {
	if tau <= 0 || tau >= len(cmnd)-1 {
		return float64(tau)
	}

	s0 := cmnd[tau-1]
	s1 := cmnd[tau]
	s2 := cmnd[tau+1]

	denominator := 2 * (2*s1 - s0 - s2)
	if denominator == 0 {
		return  float64(tau)
	}

	return float64(tau) + (s2- s0)/denominator
}

//This function determines pitch using the YIN algorithm
func DetectYIN(
	signal []float64,
	sampleRate float64,
	minFreq float64,
	maxFreq float64,
	threshold float64,
)Result{
	maxLag := int(sampleRate/ minFreq)
	minLag := int(sampleRate/maxFreq)

	diff := differece(signal, maxLag)
	cmnd := cumulativeMeanNormalizedDifference(diff)

	tau := absoluteThreshold(cmnd, threshold)

	if tau == -1 || tau < minLag {
		return Result{
			Voiced: false,
		}
	}

	refinedTau := parabolicInterpolation(cmnd, tau)
	f0 := sampleRate / refinedTau
	confidence := 1.0 - cmnd[tau]


	return Result {
		F0: f0,
		Confidence: confidence,
		Voiced: true,
	}

}