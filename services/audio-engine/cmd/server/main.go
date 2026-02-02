package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type AnalyseResponse struct {
	MeanPitch float64 `json:"mean_pitch"`
	Frames int `json:"frames"`
}

func analyseHandler(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return 
	}

	/* Karna Kya h? 
		Reading audio bytes
		Decoding WAV/PCM
		Run pitch Tracker
		Compute traits
	*/

	resp := AnalyseResponse{
		MeanPitch: 440.0, // ise abhi ke liye hardcode kar diya h
		Frames: 100,
	}

	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	http.HandleFunc("/analyze", analyseHandler)
	log.Println("Go audio engine listening on : 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}