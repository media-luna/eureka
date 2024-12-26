package fingerprint

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/cmplx"
	"os"

	"github.com/maddyblue/go-dsp/fft"
	"github.com/maddyblue/go-dsp/window"
)

const (
	downsampleRatio = 2 		// downsampling ratio for the audio samples
	freqBinSize = 1024 			// size of the frequency bins used in the spectrogram
	maxFreq = 22050.0 			// maximum frequency to be considered (22.05kHz for 44.1kHz sample rate)
	stepSize = freqBinSize / 32 // step size for the frequency bins
	windowSize = 1024			// size of the window used for the STFT
	overlap = 512				// overlap between consecutive windows in the STFT
	maxFreqBits = 9				// number of bits used to represent the frequency in the fingerprint
	maxDeltaBits = 14			// number of bits used to represent the time difference in the fingerprint
	targetZoneSize = 5			// size of the target zone for peak pairing in the fingerprinting process
)

// Spectrogram computes the spectrogram of a WAV file.
func SamplesToSpectrogram(samples []float64, sampleRate int) ([][]complex128, error) {
	// Apply Hamming window
	window := window.Hamming(len(samples))
	for i := range samples {
		samples[i] *= window[i]
	}

	// Apply low-pass filter (optional)
	filteredSamples := lowPassFilter(samples, windowSize)

	// Downsample
	downsampleSampleRate := sampleRate / downsampleRatio
	downsampledSamples, err := downsample(filteredSamples, sampleRate, downsampleSampleRate)
	if err != nil {
		return nil, err
	}

	// Compute FFT
	spectrogram := [][]complex128{}
	winSize := 1024 // Adjust window size as needed
	for i := 0; i < len(downsampledSamples)-winSize; i += winSize {
		frame := downsampledSamples[i : i+winSize]
		fftOut := fft.FFTReal(frame)
		spectrogram = append(spectrogram, fftOut)
	}

	return spectrogram, nil
}

// GenerateSpectrogramImage generates a spectrogram image from the given spectrogram data.
func SpectrogramToImage(spectrogram [][]complex128, path string) error {
	// calculate RMS
	rms := calculateRMS(spectrogram)

	// Calculate dimensions
	numFrames := len(spectrogram)
	numFreqs := len(spectrogram[0]) / 2 // Consider only positive frequencies
	imgWidth := numFrames
	imgHeight := numFreqs

	// Create image
	img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))

	// Normalize magnitudes
	maxMagnitude := 0.0
	for _, frame := range spectrogram {
		for i := 0; i < numFreqs; i++ {
			mag := cmplx.Abs(frame[i])
			if mag > maxMagnitude {
				maxMagnitude = mag
			}
		}
	}

	// Map magnitudes to colors
	for x, frame := range spectrogram {
		for y := 0; y < numFreqs; y++ {
			mag := cmplx.Abs(frame[y]) / maxMagnitude
			gray := uint8(255 * mag)
			img.Set(x, imgHeight-y-1, color.RGBA{gray, gray, gray, 255}) // Invert y-axis
		}
	}
	
	// Normalize magnitudes using RMS
	for x, frame := range spectrogram {
		for y := 0; y < numFreqs; y++ {
				mag := cmplx.Abs(frame[y]) / rms 
				gray := uint8(255 * mag) 
				img.Set(x, imgHeight-y-1, color.RGBA{gray, gray, gray, 255}) // Invert y-axis
		}
	}

	// Save file
	f, err := os.Create(path)
    if err != nil {
            fmt.Println(err)
            return err
    }
    defer f.Close()

    err = png.Encode(f, img)
    if err != nil {
            fmt.Println("Error encoding image:", err)
            return err
    }

    fmt.Println("Spectrogram image saved to spectrogram.png")
	
	return nil
}

// lowPassFilter applies a low-pass filter to the given samples using a moving average filter.
// It smooths the input signal by averaging each sample with its neighboring samples within the specified window size.
//
// Parameters:
// - samples: A slice of float64 representing the input signal samples.
// - windowSize: An integer specifying the size of the moving average window. It should be an odd number.
//
// Returns:
// - A slice of float64 containing the filtered samples.
func lowPassFilter(samples []float64, windowSize int) []float64 {
	filteredSamples := make([]float64, len(samples))
	for i := windowSize / 2; i < len(samples)-windowSize/2; i++ {
			sum := 0.0
			for j := -windowSize / 2; j <= windowSize/2; j++ {
					sum += samples[i+j]
			}
			filteredSamples[i] = sum / float64(windowSize)
	}
	return filteredSamples
}

// calculateRMS calculates the Root Mean Square (RMS) value of a given spectrogram.
// The spectrogram is represented as a 2D slice of complex128 values, where each
// element represents a frequency bin in a specific time frame.
//
// The RMS value is computed by taking the square root of the average of the
// squared magnitudes of the complex values in the spectrogram.
//
// Parameters:
// - spectogram: A 2D slice of complex128 values representing the spectrogram.
//
// Returns:
// - A float64 value representing the RMS of the spectrogram.
func calculateRMS(spectogram [][]complex128) float64{
	rms := 0.0
	for _, frame := range spectogram {
		for _, complexVal := range frame {
			rms += cmplx.Abs(complexVal) * cmplx.Abs(complexVal)
		}
	}
	
	rms = math.Sqrt(rms / float64(len(spectogram) * len(spectogram[0])))
	return rms
}

// Downsample downsamples the input audio from originalSampleRate to targetSampleRate
// downsample reduces the sample rate of the input signal to the target sample rate.
// It takes an input slice of float64 representing the original signal, the original sample rate,
// and the target sample rate. It returns a new slice of float64 representing the downsampled signal
// and an error if the sample rates are invalid.
//
// Parameters:
// - input: []float64 - the original signal to be downsampled
// - originalSampleRate: int - the sample rate of the original signal
// - targetSampleRate: int - the desired sample rate after downsampling
//
// Returns:
// - []float64 - the downsampled signal
// - error - an error if the sample rates are invalid or if the target sample rate is greater than the original sample rate
//
// The function ensures that the target sample rate is less than or equal to the original sample rate
// and that both sample rates are positive. It calculates the ratio of the original sample rate to the
// target sample rate and uses this ratio to average the input signal over intervals, producing the downsampled signal.
func downsample(input []float64, originalSampleRate, targetSampleRate int) ([]float64, error) {
	if targetSampleRate <= 0 || originalSampleRate <= 0 {
		return nil, errors.New("sample rates must be positive")
	}
	if targetSampleRate > originalSampleRate {
		return nil, errors.New("target sample rate must be less than or equal to original sample rate")
	}

	ratio := originalSampleRate / targetSampleRate
	if ratio <= 0 {
		return nil, errors.New("invalid ratio calculated from sample rates")
	}

	var resampled []float64
	for i := 0; i < len(input); i += ratio {
		end := i + ratio
		if end > len(input) {
			end = len(input)
		}

		sum := 0.0
		for j := i; j < end; j++ {
			sum += input[j]
		}
		avg := sum / float64(end-i)
		resampled = append(resampled, avg)
	}

	return resampled, nil
}