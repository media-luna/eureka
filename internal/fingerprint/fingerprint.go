package fingerprint

import "math/cmplx"

// Fingerprint
type Couple struct {
	AnchorTimeMs uint32
	SongID       uint32
}

type Peak struct {
	Time float64
	Freq complex128
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
func PickPeaks(spectrogram [][]complex128, threshold float64) []Peak {
    magnitudes := getMagnitudes(spectrogram)
    var peaks []Peak

    for t, frame := range magnitudes {
        for f, magnitude := range frame {
            if magnitude > threshold && isLocalPeak(magnitudes, t, f) {
				peaks = append(peaks, Peak{Time: float64(t), Freq: spectrogram[t][f]})
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

    for _, dt := range deltaT {
        for _, df := range deltaF {
            if dt == 0 && df == 0 {
                continue
            }
            tt, ff := t+dt, f+df
            if tt >= 0 && tt < len(magnitudes) && ff >= 0 && ff < len(magnitudes[0]) {
                if magnitudes[t][f] <= magnitudes[tt][ff] {
                    return false
                }
            }
        }
    }
    return true
}