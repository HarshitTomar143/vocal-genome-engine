package formant

import (
	"math"
	"sort"

	"gonum.org/v1/gonum/mat"
)

// ExtractFormants converts LPC coefficients into formant frequencies.
//
// Inputs:
//  - lpc: LPC coefficients (a[1..p], a[0] unused)
//  - sampleRate: sampling rate in Hz
//
// Output:
//  - sorted slice of formant frequencies in Hz
func ExtractFormants(lpc []float64, sampleRate float64) []float64 {
	if len(lpc) < 2 {
		return nil
	}

	order := len(lpc) - 1

	// Build companion matrix for LPC polynomial
	// A(z) = 1 - a1 z^-1 - a2 z^-2 - ... - ap z^-p
	data := make([]float64, order*order)

	// Subdiagonal ones
	for i := 1; i < order; i++ {
		data[i*order+(i-1)] = 1.0
	}

	// First row = -LPC coefficients
	for i := 0; i < order; i++ {
		data[i] = -lpc[i+1]
	}

	M := mat.NewDense(order, order, data)

	var eig mat.Eigen
	ok := eig.Factorize(M, mat.EigenRight)
	if !ok {
		return nil
	}

	values := eig.Values(nil)
	var formants []float64

	for _, v := range values {
		// We only care about complex conjugate roots
		if imag(v) <= 0 {
			continue
		}

		angle := math.Atan2(imag(v), real(v))
		freq := angle * sampleRate / (2 * math.Pi)

		// Physical vocal tract limits
		if freq > 90 && freq < 5000 {
			formants = append(formants, freq)
		}
	}

	sort.Float64s(formants)
	return formants
}
