package fingerprint

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"math/cmplx"
	"os"
)

const (
	PEAK_THRESHOLD         = 0.2  // Samples will be considered as peak when reaching this value
	MIN_HASH_TIME_DELTA    = 0    // Min milliseconds between 2 peaks to considered fingerprint
	MAX_HASH_TIME_DELTA    = 2000 // Max milliseconds between 2 peaks to considered fingerprint
	FAN_VALUE              = 5    // Size of the target zone for peak pairing in the fingerprinting process
	WINDOW_SIZE            = 1024 // Size of the window used for the STFT
	DOWNSAMPLE_RATIO       = 1    // Downsampling ratio for the audio samples(devide the amount of samples by N)
	MIN_WAV_BYTES          = 44   // Minimum number of bytes required for a valid WAV file
	HEADER_BITS_PER_SAMPLE = 16   // Number of bits per sample in the WAV file header
)

// Fingerprint represents a single audio fingerprint
type Fingerprint struct {
	Hash   string
	SongID int
	Offset int
}

// Fingerprints
type Fingerprints struct {
	TimeMs float64
	Hash   string
}

type Peak struct {
	Time      float64
	TimeMS    float64
	Magnitude float64
	Freq      complex128
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
					timeMS := float64(t) * float64(WINDOW_SIZE) / float64(sampleRate) * 1000
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
//
//	spectrogram [][]complex128 - A 2D slice of complex128 numbers representing the spectrogram.
//
// Returns:
//
//	[][]float64 - A 2D slice of float64 numbers representing the magnitudes of the spectrogram.
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

// CalculateFileHash generates a SHA1 hash of the file contents
func CalculateFileHash(filePath string) string {
	f, err := os.Open(filePath)
	if err != nil {
		return ""
	}
	defer f.Close()

	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		return ""
	}

	return hex.EncodeToString(h.Sum(nil))
}

// GenerateFingerprints generates fingerprints from spectrogram peaks
func GenerateFingerprints(peaks []Peak) []Fingerprint {
	var fingerprints []Fingerprint

	// Fan out from each peak
	for i, anchor := range peaks {
		// Look at the next few peaks as target points
		for j := i + 1; j < i+FAN_VALUE && j < len(peaks); j++ {
			target := peaks[j]

			// Create hash using frequency and time delta
			timeDelta := target.TimeMS - anchor.TimeMS
			if timeDelta <= MIN_HASH_TIME_DELTA || timeDelta > MAX_HASH_TIME_DELTA {
				continue
			}

			// Use real part of complex frequency as the frequency value
			anchorFreq := int(real(anchor.Freq))
			targetFreq := int(real(target.Freq))

			hashStr := fmt.Sprintf("%d|%d|%d",
				anchorFreq,
				targetFreq,
				int(timeDelta))

			fingerprints = append(fingerprints, Fingerprint{
				Hash:   hashStr,
				Offset: int(anchor.TimeMS),
			})
		}
	}

	return fingerprints
}
