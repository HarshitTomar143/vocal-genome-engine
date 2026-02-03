package decoder

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

// DecodeWAV decodes a 16-bit PCM WAV file into mono float64 PCM samples
// normalized to [-1, 1].
func DecodeWAV(data []byte) ([]float64, int, error) {
	r := bytes.NewReader(data)

	// --- RIFF header ---
	var riff [4]byte
	if _, err := io.ReadFull(r, riff[:]); err != nil {
		return nil, 0, err
	}
	if string(riff[:]) != "RIFF" {
		return nil, 0, errors.New("not RIFF")
	}

	// Skip chunk size
	r.Seek(4, io.SeekCurrent)

	var wave [4]byte
	io.ReadFull(r, wave[:])
	if string(wave[:]) != "WAVE" {
		return nil, 0, errors.New("not WAVE")
	}

	var sampleRate int
	var bitsPerSample int
	var numChannels int

	// --- Read chunks ---
	for {
		var chunkID [4]byte
		if _, err := io.ReadFull(r, chunkID[:]); err != nil {
			return nil, 0, err
		}

		var chunkSize uint32
		binary.Read(r, binary.LittleEndian, &chunkSize)

		switch string(chunkID[:]) {

		case "fmt ":
			var audioFormat uint16
			binary.Read(r, binary.LittleEndian, &audioFormat)

			binary.Read(r, binary.LittleEndian, &numChannels)
			binary.Read(r, binary.LittleEndian, &sampleRate)

			// Skip byteRate + blockAlign
			r.Seek(6, io.SeekCurrent)

			binary.Read(r, binary.LittleEndian, &bitsPerSample)

			// Skip rest of fmt chunk
			r.Seek(int64(chunkSize-16), io.SeekCurrent)

			if audioFormat != 1 {
				return nil, 0, errors.New("unsupported WAV format")
			}

		case "data":
			if bitsPerSample != 16 {
				return nil, 0, errors.New("only 16-bit PCM supported")
			}

			sampleCount := int(chunkSize) / 2
			raw := make([]int16, sampleCount)
			binary.Read(r, binary.LittleEndian, &raw)

			// Convert to mono float64
			pcm := make([]float64, 0, sampleCount)

			if numChannels == 1 {
				for _, v := range raw {
					pcm = append(pcm, float64(v)/32768.0)
				}
			} else {
				// Average channels
				for i := 0; i < len(raw); i += numChannels {
					var sum int
					for c := 0; c < numChannels; c++ {
						sum += int(raw[i+c])
					}
					pcm = append(pcm, float64(sum)/float64(numChannels)/32768.0)
				}
			}

			return pcm, sampleRate, nil

		default:
			// Skip unknown chunk
			r.Seek(int64(chunkSize), io.SeekCurrent)
		}
	}
}
