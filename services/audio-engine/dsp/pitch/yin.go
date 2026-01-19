package pitch

import "math"

// Result represents the output of YIN pitch detection
type Result struct {
	F0         float64
	Confidence float64
	Voiced     bool
}

// difference computes the YIN difference function
// d(tau) = sum (x[n] - x[n+tau])^2
func difference(signal []float64, maxLag int) []float64 {
	n := len(signal)
	diff := make([]float64, maxLag)

	for tau := 1; tau < maxLag; tau++ {
		var sum float64
		for i := 0; i < n-tau; i++ {
			delta := signal[i] - signal[i+tau]
			sum += delta * delta
		}
		diff[tau] = sum
	}

	return diff
}

// cumulativeMeanNormalizedDifference computes the CMND function
func cumulativeMeanNormalizedDifference(diff []float64) []float64 {
	cmnd := make([]float64, len(diff))
	cmnd[0] = 1.0

	var runningSum float64
	for tau := 1; tau < len(diff); tau++ {
		runningSum += diff[tau]

		// Guard against division by zero
		if runningSum == 0 {
			cmnd[tau] = 1.0
		} else {
			cmnd[tau] = diff[tau] * float64(tau) / runningSum
		}
	}

	return cmnd
}

// absoluteThreshold finds the first tau below threshold
func absoluteThreshold(cmnd []float64, threshold float64) int {
	for tau := 2; tau < len(cmnd); tau++ {
		if cmnd[tau] < threshold {
			return tau
		}
	}
	return -1
}

// parabolicInterpolation refines lag estimate to sub-sample precision
func parabolicInterpolation(cmnd []float64, tau int) float64 {
	if tau <= 0 || tau >= len(cmnd)-1 {
		return float64(tau)
	}

	s0 := cmnd[tau-1]
	s1 := cmnd[tau]
	s2 := cmnd[tau+1]

	denominator := 2 * (2*s1 - s0 - s2)
	if denominator == 0 || math.IsNaN(denominator) {
		return float64(tau)
	}

	return float64(tau) + (s2-s0)/denominator
}

// DetectYIN estimates pitch using the YIN algorithm
func DetectYIN(
	signal []float64,
	sampleRate float64,
	minFreq float64,
	maxFreq float64,
	threshold float64,
) Result {

	if len(signal) == 0 {
		return Result{Voiced: false}
	}

	// Convert frequency bounds to lag bounds
	maxLag := int(sampleRate / minFreq)
	minLag := int(sampleRate / maxFreq)

	if maxLag >= len(signal) {
		maxLag = len(signal) - 1
	}

	diff := difference(signal, maxLag)
	cmnd := cumulativeMeanNormalizedDifference(diff)

	tau := absoluteThreshold(cmnd, threshold)
	if tau == -1 || tau < minLag {
		return Result{Voiced: false}
	}

	refinedTau := parabolicInterpolation(cmnd, tau)
	if refinedTau <= 0 || math.IsNaN(refinedTau) {
		return Result{Voiced: false}
	}

	f0 := sampleRate / refinedTau
	confidence := 1.0 - cmnd[tau]

	// Extra safety: reject impossible results
	if f0 < minFreq || f0 > maxFreq || confidence <= 0 {
		return Result{Voiced: false}
	}

	return Result{
		F0:         f0,
		Confidence: confidence,
		Voiced:     true,
	}
}
