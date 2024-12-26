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

	// "github.com/hajimehoshi/go-mp3/internal/frame"
	"github.com/maddyblue/go-dsp/fft"
	"github.com/maddyblue/go-dsp/window"
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
	downsampleFactor := 2 // Downsample by x times
	downsampleSampleRate := sampleRate / downsampleFactor
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