package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"vocal-genome-engine/services/audio-engine/audio/decoder"
)

/*
	Response schema (for now)
*/
type AnalyzeResponse struct {
	MeanPitch float64 `json:"mean_pitch"`
	Frames    int     `json:"frames"`
	SampleRate int    `json:"sample_rate"`
}

/*
	Helper: always return JSON errors
*/
func writeJSONError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"error": msg,
	})
}

/*
	POST /analyze
*/
func analyzeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, "POST only", http.StatusMethodNotAllowed)
		return
	}

	// Read raw audio bytes (Next.js already strips multipart)
	audioBytes, err := io.ReadAll(r.Body)
	if err != nil {
		writeJSONError(w, "failed to read request body", 400)
		return
	}

	if len(audioBytes) == 0 {
		writeJSONError(w, "empty audio payload", 400)
		return
	}

	// Decode WAV → PCM
	pcm, sampleRate, err := decoder.DecodeWAV(audioBytes)
	if err != nil {
		log.Println("Decode error:", err)
		writeJSONError(w, err.Error(), 400)
		return
	}

	log.Printf(
		"Decoded audio: %d samples | %d Hz\n",
		len(pcm),
		sampleRate,
	)

	// ⚠️ DSP plugs in here next
	resp := AnalyzeResponse{
		MeanPitch:  0.0,            // placeholder
		Frames:     len(pcm),
		SampleRate: sampleRate,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

/*
	Server entrypoint
*/
func main() {
	http.HandleFunc("/analyze", analyzeHandler)

	log.Println("🎧 Go Audio Engine running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
