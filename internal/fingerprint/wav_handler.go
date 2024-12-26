package fingerprint

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
)

// Get WAV file information
// WavHeader defines the structure of a WAV header
type WavHeader struct {
	ChunkID       [4]byte
	ChunkSize     uint32
	Format        [4]byte
	Subchunk1ID   [4]byte
	Subchunk1Size uint32
	AudioFormat   uint16
	NumChannels   uint16
	SampleRate    uint32
	BytesPerSec   uint32
	BlockAlign    uint16
	BitsPerSample uint16
	Subchunk2ID   [4]byte
	Subchunk2Size uint32
}

// WavInfo defines a struct containing information extracted from the WAV header
type WavInfo struct {
	Channels   int
	SampleRate int
	Data	   []byte
	Samples    []float64
	Duration   float64
	FileHash   string
}

const (
	minWavBytes = 44
	headerBitsPerSample = 16
)

func hashFile(filePath string) (string, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Create a SHA256 hash object
	hasher := sha256.New()

	// Copy the file content into the hasher
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	// Compute the final hash sum
	hashSum := hasher.Sum(nil)

	// Convert the hash to a hexadecimal string
	return fmt.Sprintf("%x", hashSum), nil
}

// ReadWavInfo reads and parses the WAV file specified by the given filename.
// It returns a pointer to a WavInfo struct containing the parsed information,
// or an error if the file could not be read or parsed.
//
// The function performs the following steps:
// 1. Loads the WAV file data.
// 2. Checks if the file size is at least 44 bytes (minimum WAV header size).
// 3. Parses the WAV header from the first 44 bytes of the file data.
// 4. Extracts additional WAV information from the remaining file data.
//
// Parameters:
// - filename: The path to the WAV file to be read.
//
// Returns:
// - *WavInfo: A pointer to a WavInfo struct containing the parsed WAV information.
// - error: An error if the file could not be read or parsed, or if the file size is invalid.
func ReadWavInfo(filename string) (*WavInfo, error) {
	data, err := loadWAVFile(filename)
	if err != nil {
		return nil, err
	}

	// Check file not too small
	if len(data) < minWavBytes {
		return nil, errors.New("invalid WAV file size (too small)")
	}

	// Get only header info
	header, err := parseWavHeader(data[:minWavBytes])
	if err != nil {
		return nil, err
	}

	// Exctract samples form file
	samples, err := bytesToSamples(data[minWavBytes:])
	if err != nil {
		return nil, err
	}

	// Generate hash string for file
	hash, err := hashFile(filename)
	if err != nil {
		return nil, err
	}

	// Exstract wav header values to struct
	info, err := extractWavInfo(header, data, samples, hash)
	if err != nil {
		return nil, err
	}

	return info, nil
}

// parseWavHeader parses the WAV file header from the provided byte slice.
// It returns a pointer to a WavHeader struct if the header is valid, or an error if the header is invalid or if there is an issue reading the header data.
//
// The function checks for the following conditions to validate the WAV header:
// - The ChunkID must be "RIFF"
// - The Format must be "WAVE"
// - The AudioFormat must be 1 (indicating PCM format)
//
// Parameters:
// - headerData: A byte slice containing the WAV header data.
//
// Returns:
// - A pointer to a WavHeader struct if the header is valid.
// - An error if the header is invalid or if there is an issue reading the header data.
func parseWavHeader(headerData []byte) (*WavHeader, error) {
	var header WavHeader
	err := binary.Read(bytes.NewReader(headerData), binary.LittleEndian, &header)
	if err != nil {
		return nil, err
	}

	if string(header.ChunkID[:]) != "RIFF" || string(header.Format[:]) != "WAVE" || header.AudioFormat != 1 {
		return nil, errors.New("invalid WAV header format")
	}

	return &header, nil
}

// extractWavInfo extracts information from a WAV file header and data.
// It returns a WavInfo struct containing the number of channels, sample rate,
// data, and duration of the audio if the bits per sample is 16. If the bits per
// sample format is unsupported, it returns an error.
//
// Parameters:
//   - header: A pointer to a WavHeader struct containing the WAV file header information.
//   - data: A byte slice containing the WAV file audio data.
//
// Returns:
//   - A pointer to a WavInfo struct containing the extracted information.
//   - An error if the bits per sample format is unsupported.
func extractWavInfo(header *WavHeader, data []byte, samples []float64, hash string) (*WavInfo, error) {
	info := &WavInfo{
		Channels:   int(header.NumChannels),
		SampleRate: int(header.SampleRate),
		Samples:    samples,
		Data:		data,
		FileHash: 	hash,
	}

	if header.BitsPerSample == headerBitsPerSample {
		info.Duration = float64(len(info.Data)) / float64(int(header.NumChannels)*2*int(header.SampleRate))
	} else {
		return nil, errors.New("unsupported bits per sample format")
	}

	return info, nil
}

// loadWAVFile loads a WAV file from the specified filename and returns its contents as a byte slice.
// It returns an error if there is any issue opening or reading the file.
//
// Parameters:
//   - filename: The path to the WAV file to be loaded.
//
// Returns:
//   - []byte: The contents of the WAV file.
//   - error: An error if there is an issue opening or reading the file.
func loadWAVFile(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var data bytes.Buffer
	buf := make([]byte, 4096)

	for {
		n, err := reader.Read(buf)
		if n > 0 {
			data.Write(buf[:n])
		}
		if err != nil {
			if err.Error() == "EOF" {
				break // EOF reached
			}
			return nil, fmt.Errorf("error reading file: %v", err)
		}
	}

	return data.Bytes(), nil
}

// BytesToSamples converts a byte slice containing 16-bit PCM WAV audio data
// into a slice of float64 samples scaled to the range [-1, 1].
//
// The input byte slice must have an even length, as each sample is represented
// by two bytes. If the length of the input is not even, an error is returned.
//
// Parameters:
//   - input: A byte slice containing 16-bit PCM WAV audio data.
//
// Returns:
//   - A slice of float64 samples scaled to the range [-1, 1].
//   - An error if the input length is invalid.
func bytesToSamples(input []byte) ([]float64, error) {
	if len(input)%2 != 0 {
		return nil, errors.New("invalid input length")
	}

	numSamples := len(input) / 2
	output := make([]float64, numSamples)

	for i := 0; i < numSamples; i++ {
		// Interpret bytes as a 16-bit signed integer (little-endian)
		sample := int16(binary.LittleEndian.Uint16(input[i*2 : i*2+2]))

		// Scale the sample to the range [-1, 1]
		output[i] = float64(sample) / 32768.0
	}

	return output, nil
}