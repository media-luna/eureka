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
)

// LowPassFilter is a first-order low-pass filter using H(p) = 1 / (1 + pRC)
type LowPassFilter struct {
	alpha float64 // Filter coefficient
	yPrev float64 // Previous output value
}

const (
	dspRatio    = 4
	freqBinSize = 1024
	maxFreq     = 22050.0 // 22.05kHz, typical max frequency for 44.1kHz sample rate
	stepSize    = freqBinSize / 32
	windowSize  = 1024
	overlap     = 512

	maxFreqBits    = 9
	maxDeltaBits   = 14
	targetZoneSize = 5
)

func SpectrogramFromSamples(samples []float64, sampleRate int) ([][]complex128, error) {
	lpf, err := newLowPassFilter(maxFreq, float64(sampleRate))
	if err != nil {
		return nil, fmt.Errorf("could not generate low-pass filter")
	}

	filteredSamples := lpf.Filter(samples)

	downsampledSamples, err := downsample(filteredSamples, sampleRate, sampleRate/dspRatio)
	if err != nil {
		return nil, fmt.Errorf("couldn't downsample audio samples: %v", err)
	}

	numOfWindows := (len(downsampledSamples) - windowSize) / (windowSize - overlap) + 1
	spectrogram := make([][]complex128, numOfWindows)

	// Apply Hamming window function
	window := make([]float64, windowSize)
	for i := range window {
		window[i] = 0.54 - 0.46*math.Cos(2*math.Pi*float64(i)/(float64(windowSize)-1))
	}

	// Perform STFT
	for i := 0; i < numOfWindows; i++ {
		start := i * (windowSize - overlap)
		end := start + windowSize
		if end > len(downsampledSamples) {
			end = len(downsampledSamples)
		}

		bin := make([]float64, windowSize)
		copy(bin, downsampledSamples[start:end])

		// Apply Hamming window
		for j := range window {
			bin[j] *= window[j]
		}

		spectrogram[i] = FFT(bin)
	}

	return spectrogram, nil
}

func SpectrogramToImage(spectrogram [][]complex128, path string) (error) {
	height := len(spectrogram[0])
	width := len(spectrogram)

	// Validate input data
	if len(spectrogram) == 0 {
		return fmt.Errorf("spectrogram data is empty")
	}

	// Check spectrogram inconsistensy
	for _, column := range spectrogram {
		if len(column) != height {
			return fmt.Errorf("spectrogram data has inconsistent column heights")
		}
	}

	// Find maximum magnitude for normalization
	maxMagnitude := getMaxMagnitude(spectrogram)
	if maxMagnitude == 0 {
		return fmt.Errorf("maximum magnitude is zero, possibly all-zero input")
	}

	// Initialize the image
	img := image.NewRGBA(image.Rect(0, 0, height, width))

	// Process each time-frequency point
	for timeStep := 0; timeStep < width; timeStep++ {
		for freq := 0; freq < height; freq++ {
			// Compute magnitude and normalize
			magnitude := cmplx.Abs(spectrogram[timeStep][freq])
			normalized := math.Log(1 + magnitude) / math.Log(1 + maxMagnitude)

			// Color mapping
			gray := uint8(normalized * 255)

			// Set pixel color, flipping timeStep and freq for correct orientation
			img.Set(freq, timeStep, color.RGBA{R: gray, G: gray, B: gray, A: 255})
		}
	}

	// Create spectogram image
	outputFile, err := os.Create(path)
	if err != nil {
		fmt.Println("Error creating image file:", err)
		return err
	}
	defer outputFile.Close()

	// Save spectogram data to file
	err = png.Encode(outputFile, img)
	if err != nil {
		fmt.Println("Error saving image:", err)
		return err
	}
	
	fmt.Println("Spectrogram image saved as 'spectrogram.png'")
	return nil
}

// getMaxMagnitude finds the maximum magnitude in the spectrogram data for normalization.
func getMaxMagnitude(spectrogramData [][]complex128) float64 {
	maxMagnitude := 0.0
	for _, column := range spectrogramData {
		for _, value := range column {
			magnitude := cmplx.Abs(value)
			if magnitude > maxMagnitude {
				maxMagnitude = magnitude
			}
		}
	}
	return maxMagnitude
}

// newLowPassFilter creates a new low-pass filter
func newLowPassFilter(cutoffFrequency, sampleRate float64) (*LowPassFilter, error) {
	if cutoffFrequency <= 0 || sampleRate <= 0 {
		return nil, errors.New("cutoff frequency and sample rate must be positive")
	}
	rc := 1.0 / (2 * math.Pi * cutoffFrequency)
	dt := 1.0 / sampleRate
	alpha := dt / (rc + dt)
	return &LowPassFilter{
		alpha: alpha,
		yPrev: 0,
	}, nil
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

// Fft performs the Fast Fourier Transform on the input signal.
func FFT(input []float64) []complex128 {
	// Convert input to complex128
	complexArray := make([]complex128, len(input))
	for i, v := range input {
		complexArray[i] = complex(v, 0)
	}

	fftResult := make([]complex128, len(complexArray))
	copy(fftResult, complexArray) // Copy input to result buffer
	return recursiveFFT(fftResult)
}

// recursiveFFT performs the recursive FFT algorithm.
func recursiveFFT(complexArray []complex128) []complex128 {
	N := len(complexArray)
	if N <= 1 {
		return complexArray
	}

	even := make([]complex128, N/2)
	odd := make([]complex128, N/2)
	for i := 0; i < N/2; i++ {
		even[i] = complexArray[2*i]
		odd[i] = complexArray[2*i+1]
	}

	even = recursiveFFT(even)
	odd = recursiveFFT(odd)

	fftResult := make([]complex128, N)
	for k := 0; k < N/2; k++ {
		t := complex(math.Cos(-2*math.Pi*float64(k)/float64(N)), math.Sin(-2*math.Pi*float64(k)/float64(N)))
		fftResult[k] = even[k] + t*odd[k]
		fftResult[k+N/2] = even[k] - t*odd[k]
	}

	return fftResult
}

// Filter processes the input signal through the low-pass filter
func (lpf *LowPassFilter) Filter(input []float64) []float64 {
	filtered := make([]float64, len(input))
	for i, x := range input {
		if i == 0 {
			filtered[i] = x * lpf.alpha
		} else {
			filtered[i] = lpf.alpha*x + (1-lpf.alpha)*lpf.yPrev
		}
		lpf.yPrev = filtered[i]
	}
	return filtered
}

// Peak extraction
type Peak struct {
	Time float64
	Freq complex128
}

// ExtractPeaks analyzes a spectrogram and extracts significant peaks in the frequency domain over time.
func ExtractPeaks(spectrogram [][]complex128, audioDuration float64) []Peak {
	if len(spectrogram) < 1 {
		return []Peak{}
	}

	type maxies struct {
		maxMag  float64
		maxFreq complex128
		freqIdx int
	}

	bands := []struct{ min, max int }{{0, 10}, {10, 20}, {20, 40}, {40, 80}, {80, 160}, {160, 512}}

	var peaks []Peak
	binDuration := audioDuration / float64(len(spectrogram))

	for binIdx, bin := range spectrogram {
		var maxMags []float64
		var maxFreqs []complex128
		var freqIndices []float64

		binBandMaxies := []maxies{}
		for _, band := range bands {
			var maxx maxies
			var maxMag float64
			for idx, freq := range bin[band.min:band.max] {
				magnitude := cmplx.Abs(freq)
				if magnitude > maxMag {
					maxMag = magnitude
					freqIdx := band.min + idx
					maxx = maxies{magnitude, freq, freqIdx}
				}
			}
			binBandMaxies = append(binBandMaxies, maxx)
		}

		for _, value := range binBandMaxies {
			maxMags = append(maxMags, value.maxMag)
			maxFreqs = append(maxFreqs, value.maxFreq)
			freqIndices = append(freqIndices, float64(value.freqIdx))
		}

		// Calculate the average magnitude
		var maxMagsSum float64
		for _, max := range maxMags {
			maxMagsSum += max
		}
		avg := maxMagsSum / float64(len(maxFreqs)) // * coefficient

		// Add peaks that exceed the average magnitude
		for i, value := range maxMags {
			if value > avg {
				peakTimeInBin := freqIndices[i] * binDuration / float64(len(bin))

				// Calculate the absolute time of the peak
				peakTime := float64(binIdx)*binDuration + peakTimeInBin

				peaks = append(peaks, Peak{Time: peakTime, Freq: maxFreqs[i]})
			}
		}
	}

	return peaks
}

// Fingerprint
type Couple struct {
	AnchorTimeMs uint32
	SongID       uint32
}

// Fingerprint generates fingerprints from a list of peaks and stores them in an array.
// The fingerprints are encoded using a 32-bit integer format and stored in an array.
// Each fingerprint consists of an address and a couple.
// The address is a hash. The couple contains the anchor time and the song ID.
func Fingerprint(peaks []Peak, songID uint32) map[uint32]Couple {
	fingerprints := map[uint32]Couple{}

	for i, anchor := range peaks {
		for j := i + 1; j < len(peaks) && j <= i+targetZoneSize; j++ {
			target := peaks[j]

			address := createAddress(anchor, target)
			anchorTimeMs := uint32(anchor.Time * 1000)

			fingerprints[address] = Couple{anchorTimeMs, songID}
		}
	}

	return fingerprints
}

// createAddress generates a unique address for a pair of anchor and target points.
// The address is a 32-bit integer where certain bits represent the frequency of
// the anchor and target points, and other bits represent the time difference (delta time)
// between them. This function combines these components into a single address (a hash).
func createAddress(anchor, target Peak) uint32 {
	anchorFreq := int(real(anchor.Freq))
	targetFreq := int(real(target.Freq))
	deltaMs := uint32((target.Time - anchor.Time) * 1000)

	// Combine the frequency of the anchor, target, and delta time into a 32-bit address
	address := uint32(anchorFreq<<23) | uint32(targetFreq<<14) | deltaMs

	return address
}