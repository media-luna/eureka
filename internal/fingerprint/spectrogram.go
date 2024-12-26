package fingerprint

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/cmplx"
	"os"

	"github.com/maddyblue/go-dsp/fft"
	"github.com/maddyblue/go-dsp/window"
)

// Spectrogram computes the spectrogram of a WAV file.
func Spectrogram1(samples []float64, sampleRate int) ([][]complex128, error) {
	// Apply Hamming window
	window := window.Hamming(len(samples))
	for i := range samples {
		samples[i] *= window[i]
	}

	// Downsample
	downsampleFactor := 2 // Downsample by x times
	downsampleSampleRate := sampleRate / downsampleFactor

	samples, err := downsample(samples, sampleRate, downsampleSampleRate)
	if err != nil {
		return nil, err
	}

	// Apply low-pass filter (optional)
	// ...

	// Compute FFT
	spectrogram := [][]complex128{}
	winSize := 1024 // Adjust window size as needed
	for i := 0; i < len(samples)-winSize; i += winSize {
		frame := samples[i : i+winSize]
		fftOut := fft.FFTReal(frame)
		spectrogram = append(spectrogram, fftOut)
	}

	return spectrogram, nil
}

// GenerateSpectrogramImage generates a spectrogram image from the given spectrogram data.
func GenerateSpectrogramImage1(spectrogram [][]complex128, rms float64) *image.RGBA {
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
			img.Set(x, y, color.RGBA{gray, gray, gray, 255})
		}
	}

	// Normalize magnitudes using RMS
	for x, frame := range spectrogram {
        for y := 0; y < numFreqs; y++ {
                mag := cmplx.Abs(frame[y]) / rms 
                gray := uint8(255 * mag) 
                img.Set(x, y, color.RGBA{gray, gray, gray, 255})
        }
    }

	// Save file
	f, err := os.Create("spectrogram.png")
    if err != nil {
            fmt.Println(err)
            return nil
    }
    defer f.Close()

    err = png.Encode(f, img)
    if err != nil {
            fmt.Println("Error encoding image:", err)
            return nil
    }

    fmt.Println("Spectrogram image saved to spectrogram.png")
	
	return img
}