package window

import "math"

// This function generates a Hann window of length N
func Hann(N int) []float64 {
	w := make([]float64, N)

	if N == 1 {
		w[0] = 1.0
		return w
	}

	for n := 0; n < N; n++ {
		w[n] = 0.5 * (1.0 - math.Cos(2.0*math.Pi*float64(n)/float64(N-1)))
	}

	return w
}
