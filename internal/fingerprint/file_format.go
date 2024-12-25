package fingerprint

import (
	"fmt"
	"os"
	"strings"

	"github.com/faiface/beep"
	"github.com/faiface/beep/flac"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/wav"
)

// monoStreamer combines multiple channels into a single mono channel
type monoStreamer struct {
	streamer beep.Streamer
	format   beep.Format
}

func (m *monoStreamer) Stream(samples [][2]float64) (n int, ok bool) {
	n, ok = m.streamer.Stream(samples)
	for i := range samples[:n] {
		// Average left and right channels
		monoValue := (samples[i][0] + samples[i][1]) / 2
		// Set both channels to the same value for mono
		samples[i][0], samples[i][1] = monoValue, monoValue
	}
	return n, ok
}

func (m *monoStreamer) Err() error {
	return m.streamer.Err()
}

// ConvertToWAV decodes the input audio file and saves it as a WAV file
func ConvertToWAV(inputPath string, outputPath string, channels int) error {
	// Open the input file
	file, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("error opening input file: %v", err)
	}
	defer file.Close()

	// Create a decoder based on file extension
	var streamer beep.StreamSeekCloser
	var format beep.Format

	// Which file format
	switch ext := strings.ToLower(getFileExtension(inputPath)); ext {
	case "mp3":
		streamer, format, err = mp3.Decode(file)
	case "flac":
		streamer, format, err = flac.Decode(file)
	case "wav":
		fmt.Println("File already WAV type")
		return nil
	default:
		return fmt.Errorf("unsupported format: %s", ext)
	}

	// Error handling
	if err != nil {
		return fmt.Errorf("error decoding file: %v", err)
	}
	defer streamer.Close()

	// Create the WAV output file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer outputFile.Close()
	
	// Validate num channels
	if channels > 2 || channels < 1 {
		channels = 1
	}

	// whether decode and save as 1 or 2 channels(stereo/mono)
	switch channels {
	case 1:
		// Mono streamer
		mono := &monoStreamer{streamer: streamer, format: format}
		format.NumChannels = 1

		err = wav.Encode(outputFile, mono, format)
		if err != nil {
			return fmt.Errorf("error encoding WAV: %v", err)
		}
	case 2:
		// Encode the audio as WAV format
		err = wav.Encode(outputFile, streamer, format)
		if err != nil {
			return fmt.Errorf("error encoding WAV: %v", err)
		}
	}

	fmt.Println("Conversion completed:", outputPath)
	return nil
}

func getFileExtension(filename string) string {
	parts := strings.Split(filename, ".")
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}
	return ""
}