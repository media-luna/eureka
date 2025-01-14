package eureka

import (
	"fmt"
	"os"

	config "github.com/media-luna/eureka/configs"
	"github.com/media-luna/eureka/internal/database"
	fingerprint "github.com/media-luna/eureka/internal/fingerprint"
)

// Eureka represents the main structure for the Eureka service,
// containing the configuration settings required for its operation.
type Eureka struct {
	Config config.Config
}

// NewEureka initializes a new Eureka instance with the provided configuration.
// It performs the following steps:
// 1. Initializes the database object using the provided configuration.
// 2. Connects to the database.
// 3. Sets up the database.
//
// If any of these steps fail, it logs the error and returns nil.
//
// Parameters:
//   - config: The configuration object used to initialize the database.
//
// Returns:
//   - A pointer to the initialized Eureka instance, or nil if an error occurred.
func NewEureka(config config.Config) (*Eureka, error) {
	// audioDownloader , err := downloader.NewAudioDownloader("https://www.youtube.com/watch?v=s8QYxmpuyxg")
	// println(audioDownloader.GetTrack())

	// TODO: Load all fingerprinted songs and their hashes to memory
	// if possible to make the process a bit faster

	// Init DB object
	database, err := database.NewDatabase(config)
	if err != nil {
		return nil, err
	}

	// Setup DB
	if err := database.Setup(); err != nil {
		return nil, err
	}

	return &Eureka{Config: config}, nil
}

// Save processes an audio file, generates its spectrogram, and extracts fingerprints.
// It performs the following steps:
// 1. Checks if the provided path is a file (directories are not supported).
// 2. Converts the file to WAV format.
// 3. Reads WAV file information.
// 4. Generates a spectrogram from the WAV samples (this step may take a long time).
// 5. Collects peaks from the spectrogram.
// 6. Saves the spectrogram image with peaks.
// 7. Extracts fingerprints from the peaks.
// 8. (TODO) Stores song information, file hash, and fingerprints to the database.
//
// Parameters:
// - path: The file path of the audio file to be processed.
//
// Returns:
// - error: An error if any step fails, otherwise nil.
func (e *Eureka) Save(path string) error {
	// Check if path is dir or file
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("error stating path: %v", err)
	}

	if info.IsDir() {
		return fmt.Errorf("path is a directory not supported, expected a file")
	}

	// Convert any file type to WAV
	filePath, err := fingerprint.ConvertToWAV(path, "output.wav")
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Read wav info
	wavInfo, err := fingerprint.ReadWavInfo(filePath)
	if err != nil {
		fmt.Println("Error:", err)
	}

	// TODO: this func takes too long like 20 seconds!!!
	// Generate spectrogram
	spectrogram, err := fingerprint.SamplesToSpectrogram(wavInfo.Samples, wavInfo.SampleRate)
	if err != nil {
		return fmt.Errorf("error creating spectrogram: %v", err)
	}

	// Collect spectrogram peaks
	peaks := fingerprint.PickPeaks(spectrogram, wavInfo.SampleRate)

	// Save spectrogram image with peaks
	if err := fingerprint.SpectrogramToImage(spectrogram, peaks, wavInfo.SampleRate, "spectrogram.png"); err != nil {
		return fmt.Errorf("error creating spectrogram: %v", err)
	}

	fingerprints := fingerprint.Fingerprint(peaks)
	print(fingerprints)

	// Get DB conn
	// store song info and file hash to DB
	// store fingerprints to DB
	// set song fingerfrinted

	return nil
}
