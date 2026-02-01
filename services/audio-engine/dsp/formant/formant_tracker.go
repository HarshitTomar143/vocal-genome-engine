package formant

import (
	"vocal-genome-engine/services/audio-engine/features"
)

// FormantFrame represents formant frequencies at a given time
type FormantFrame struct {
	Time     float64
	Formants []float64
}

// TrackFormants extracts formants for voiced frames only.
func TrackFormants(
	frames [][]float64,
	times []float64,
	pitchFrames []features.PitchFrame,
	sampleRate float64,
	lpcOrder int,
) []FormantFrame {

	var results []FormantFrame

	for i := 0; i < len(frames) && i < len(pitchFrames); i++ {

		// Use FEATURE-level voicing decision
		if !pitchFrames[i].Voiced {
			continue
		}

		lpc, errEnergy := ComputeLPC(frames[i], lpcOrder, 0.97)
		if errEnergy <= 0 {
			continue
		}

		formants := ExtractFormants(lpc, sampleRate)
		if len(formants) == 0 {
			continue
		}

		results = append(results, FormantFrame{
			Time:     times[i],
			Formants: formants,
		})
	}

	return results
}
