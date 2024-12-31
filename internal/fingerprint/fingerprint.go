package fingerprint

import (
	"crypto/sha1"
	"fmt"
	"math/cmplx"
)

const (
	PEAK_THRESHOLD = 0.2 // Samples will be considered as peak when reaching this value
	MIN_HASH_TIME_DELTA = 0
	MAX_HASH_TIME_DELTA = 2000 // Max milliseconds between 2 peaks to considered fingerprint   
	FAN_VALUE = 5  // size of the target zone for peak pairing in the fingerprinting process
)

// Fingerprint
type Fingerprints struct {
	TimeMs float64
	Hash   string
}

type Peak struct {
	Time   float64
	TimeMS float64
	Magnitude float64
	Freq   complex128
}

// ExtractPeaks identifies and extracts peaks from a given spectrogram based on a specified threshold.
// A peak is defined as a local maximum in the magnitude of the spectrogram that exceeds the threshold.
//
// Parameters:
//   - spectrogram: A 2D slice of complex128 values representing the spectrogram data.
//   - threshold: A float64 value representing the minimum magnitude required for a peak.
//
// Returns:
//   - A slice of Peak structs, each representing a detected peak with its time and frequency.
func PickPeaks(spectrogram [][]complex128, sampleRate int) []Peak {
	magnitudes := getMagnitudes(spectrogram)
	var peaks []Peak
	freqMap := make(map[string]bool)

	for t, frame := range magnitudes {
		for f, magnitude := range frame {
			if magnitude > PEAK_THRESHOLD && isLocalPeak(magnitudes, t, f) {
				freqStr := fmt.Sprintf("%.10f", real(spectrogram[t][f]))
				if _, exists := freqMap[freqStr]; !exists {
					timeMS := float64(t) * float64(windowSize) / float64(sampleRate) * 1000
					peaks = append(peaks, Peak{Time: float64(t), TimeMS: timeMS, Freq: spectrogram[t][f], Magnitude: magnitude})
					freqMap[freqStr] = true
				}
			}
		}
	}
	return peaks
}

// getMagnitudes computes the magnitudes of a given 2D spectrogram.
// Each element in the spectrogram is a complex number, and the magnitude
// is calculated using the absolute value of the complex number.
//
// Parameters:
//   spectrogram [][]complex128 - A 2D slice of complex128 numbers representing the spectrogram.
//
// Returns:
//   [][]float64 - A 2D slice of float64 numbers representing the magnitudes of the spectrogram.
func getMagnitudes(spectrogram [][]complex128) [][]float64 {
    magnitudes := make([][]float64, len(spectrogram))
    for i, row := range spectrogram {
        magnitudes[i] = make([]float64, len(row))
        for j, val := range row {
            magnitudes[i][j] = cmplx.Abs(val)
        }
    }
    return magnitudes
}

// TODO: find neighbors dynamically
// isLocalPeak determines if the magnitude at a given time-frequency point (t, f)
// is a local peak in the spectrogram. A local peak is defined as a point that has
// a higher magnitude than all of its immediate neighbors.
//
// Parameters:
// - magnitudes: A 2D slice of float64 representing the magnitudes in the spectrogram.
// - t: The time index of the point to check.
// - f: The frequency index of the point to check.
//
// Returns:
// - bool: True if the point (t, f) is a local peak, false otherwise.
func isLocalPeak(magnitudes [][]float64, t, f int) bool {
    deltaT := []int{-1, 0, 1}
    deltaF := []int{-1, 0, 1}

    peakValue := magnitudes[t][f]
    for _, dt := range deltaT {
        for _, df := range deltaF {
            if dt == 0 && df == 0 {
                continue
            }
            tt, ff := t+dt, f+df
            if tt >= 0 && tt < len(magnitudes) && ff >= 0 && ff < len(magnitudes[0]) {
                if peakValue <= magnitudes[tt][ff] {
                    return false
                }
            }
        }
    }
    return true
}

func Fingerprint(peaks []Peak) []Fingerprints {
	var hashes []Fingerprints

    for i := 0; i < len(peaks); i++ {
        for j := 1; j <= FAN_VALUE; j++ {
            if i+j < len(peaks) {
                f1 := peaks[i].Magnitude
                f2 := peaks[i+j].Magnitude
                t1 := peaks[i].TimeMS
                t2 := peaks[i+j].TimeMS
				tDelta := int(t2-t1)

				if MIN_HASH_TIME_DELTA <= tDelta && tDelta <= MAX_HASH_TIME_DELTA {
					hashData := fmt.Sprintf("%f|%f|%d", f1, f2, tDelta)
					hash := fmt.Sprintf("%x", sha1.Sum([]byte(hashData)))
					hashes = append(hashes, Fingerprints{TimeMs: t1, Hash: hash[:20]})
				}
            }
        }
    }
    return hashes
}
