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

// NewSystem creates and initializes a new System.
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
func NewEureka(config config.Config) *Eureka {
	// audioDownloader , err := downloader.NewAudioDownloader("https://www.youtube.com/watch?v=s8QYxmpuyxg")
	// println(audioDownloader.GetTrack())

	// TODO: Load all fingerprinted songs and their hashes to memory
	// if possible to make the process a bit faster

	// Init DB object
	database, err := database.NewDatabase(config)
	if err != nil {
		fmt.Println("error initializing database:", err)
		return nil
	}

	// Connect to DB
	if err := database.Connect(); err != nil {
		fmt.Println("error connecting to database:", err)
		return nil
	}
	defer database.Close()

	fmt.Println("Database connected successfully!")

	// Setup DB
	if err := database.Setup(); err != nil {
		fmt.Println("error connecting to database:", err)
		return nil
	}

	return &Eureka{Config: config}
}

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
	err = fingerprint.ConvertToWAV(path, "output.wav", 2)
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Read wav info
	wavInfo, err := fingerprint.ReadWavInfo("output.wav")
	if err != nil {
		fmt.Println("Error:", err)
	}

	println(wavInfo.FileHash)

	// Wav bytes to samples
	samples, err := fingerprint.BytesToSamples(wavInfo.Data)
	if err != nil {
		return fmt.Errorf("error converting wav bytes to float64: %v", err)
	}

	println(samples)

	// make spectogram
	spectro, err := fingerprint.SpectrogramFromSamples(samples, wavInfo.SampleRate)
	if err != nil {
		return fmt.Errorf("error creating spectrogram: %v", err)
	}

	if err := fingerprint.SpectrogramToImage(spectro, "spectrogram.png"); err != nil {
		return fmt.Errorf("error creating spectrogram: %v", err)
	}

	// extract peaks
	peaks := fingerprint.ExtractPeaks(spectro, wavInfo.Duration)

	// fingerprint file
	fingerprints := fingerprint.Fingerprint(peaks, 1)
	println(fingerprints)

	// DB conn
	// store song info and hash to DB
	// 	songID, err := dbclient.RegisterSong(songTitle, songArtist, ytID)
	// store fingerprints to DB
	// 	err = dbclient.StoreFingerprints(fingerprints)

	return nil
}
