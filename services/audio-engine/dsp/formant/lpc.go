package formant

/*
 Pre-emphasis
 Autocorrelation
 Levinson–Durbin recursion

 Input: windowed audio frame
 Output: LPC coefficients representing vocal tract filter

 PreEmphasis applies a first-order high-pass filter
 y[n] = x[n] - alpha * x[n-1]
*/

func PreEmphasis(x []float64, alpha float64) []float64 {
	y := make([]float64, len(x))
	if len(x) == 0 {
		return y
	}

	y[0] = x[0]
	for i := 1; i < len(x); i++ {
		y[i] = x[i] - alpha*x[i-1]
	}
	return y
}

// Autocorrelation computes autocorrelation coefficients
// up to the given LPC order.
func Autocorrelation(x []float64, order int) []float64 {
	r := make([]float64, order+1)

	for lag := 0; lag <= order; lag++ {
		var sum float64
		for n := lag; n < len(x); n++ {
			sum += x[n] * x[n-lag]
		}
		r[lag] = sum
	}
	return r
}

// LevinsonDurbin solves the normal equations for LPC
// using Levinson–Durbin recursion.
//
// Returns:
//  - a: LPC coefficients (a[0] unused, a[1..order] valid)
//  - e: final prediction error energy
func LevinsonDurbin(r []float64, order int) ([]float64, float64) {
	a := make([]float64, order+1)

	if len(r) == 0 || r[0] == 0 {
		return a, 0
	}

	e := r[0]

	for i := 1; i <= order; i++ {
		var acc float64
		for j := 1; j < i; j++ {
			acc += a[j] * r[i-j]
		}

		k := (r[i] - acc) / e
		a[i] = k

		for j := 1; j < i; j++ {
			a[j] = a[j] - k*a[i-j]
		}

		e *= (1 - k*k)
		if e <= 0 {
			break
		}
	}

	return a, e
}

// ComputeLPC is a convenience wrapper that performs
// pre-emphasis, autocorrelation, and LPC estimation.
func ComputeLPC(frame []float64, order int, preEmph float64) ([]float64, float64) {
	if len(frame) == 0 {
		return nil, 0
	}

	// 1. Pre-emphasis
	emphasized := PreEmphasis(frame, preEmph)

	// 2. Autocorrelation
	r := Autocorrelation(emphasized, order)

	// 3. Levinson–Durbin
	a, e := LevinsonDurbin(r, order)

	return a, e
}
