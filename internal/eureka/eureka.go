package eureka

import (
	"fmt"
	"os"
	"path/filepath"

	config "github.com/media-luna/eureka/configs"
	"github.com/media-luna/eureka/internal/database"
	"github.com/media-luna/eureka/internal/database/mysql"
	fingerprint "github.com/media-luna/eureka/internal/fingerprint"
	"github.com/media-luna/eureka/utils/logger"
	"github.com/schollz/progressbar/v3"
)

// Eureka represents the main structure for the Eureka service,
// containing the configuration settings required for its operation.
type Eureka struct {
	Config   config.Config
	database database.Database
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
	db, err := database.NewDatabase(config)
	if err != nil {
		return nil, err
	}

	// Setup DB
	if err := db.Setup(); err != nil {
		return nil, err
	}

	return &Eureka{
		Config:   config,
		database: db,
	}, nil
}

// Save processes an audio file, generates its spectrogram, and extracts fingerprints.
func (e *Eureka) Save(path string) error {
	// Check if path is dir or file
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("error stating path: %v", err)
	}

	if info.IsDir() {
		return fmt.Errorf("path is a directory not supported, expected a file")
	}

	logger.Info(fmt.Sprintf("Processing audio file: %s", filepath.Base(path)))

	// Convert any file type to WAV
	filePath, err := fingerprint.ConvertToWAV(path, "output.wav")
	if err != nil {
		return fmt.Errorf("error converting to WAV: %v", err)
	}
	logger.Info("Audio file converted to WAV format")

	// Read wav info
	wavInfo, err := fingerprint.ReadWavInfo(filePath)
	if err != nil {
		return fmt.Errorf("error reading WAV info: %v", err)
	}

	logger.Info("Generating spectrogram...")
	// Generate spectrogram
	spectrogram, err := fingerprint.SamplesToSpectrogram(wavInfo.Samples, wavInfo.SampleRate)
	if err != nil {
		return fmt.Errorf("error creating spectrogram: %v", err)
	}

	// Collect spectrogram peaks
	peaks := fingerprint.PickPeaks(spectrogram, wavInfo.SampleRate)
	logger.Info(fmt.Sprintf("Found %d peaks in spectrogram", len(peaks)))

	// Save spectrogram image with peaks
	if err := fingerprint.SpectrogramToImage(spectrogram, peaks, wavInfo.SampleRate, "spectrogram.png"); err != nil {
		return fmt.Errorf("error saving spectrogram image: %v", err)
	}

	// Generate fingerprints
	logger.Info("Generating fingerprints...")
	fingerprints := fingerprint.GenerateFingerprints(peaks)
	logger.Info(fmt.Sprintf("Generated %d fingerprints", len(fingerprints)))

	// Calculate file hash
	fileHash := fingerprint.CalculateFileHash(path)

	// Store song in database
	songName := filepath.Base(path)
	songID, err := e.database.InsertSong(songName, "", fileHash, len(fingerprints))
	if err != nil {
		return fmt.Errorf("error inserting song: %v", err)
	}

	// Store fingerprints with progress bar
	logger.Info("Storing fingerprints in database...")
	bar := progressbar.Default(int64(len(fingerprints)))
	for _, fp := range fingerprints {
		if err := e.database.InsertFingerprints(fp.Hash, songID, fp.Offset); err != nil {
			return fmt.Errorf("error inserting fingerprint: %v", err)
		}
		bar.Add(1)
	}
	logger.Info(fmt.Sprintf("Successfully processed %s", songName))

	return nil
}

// ListSongs returns all songs from the database
func (e *Eureka) ListSongs() ([]mysql.Song, error) {
	if db, ok := e.database.(*mysql.DB); ok {
		return db.ListSongs()
	}
	return nil, fmt.Errorf("database type does not support listing songs")
}

// CleanupDuplicates removes duplicate songs from the database
func (e *Eureka) CleanupDuplicates() error {
	if db, ok := e.database.(*mysql.DB); ok {
		return db.CleanupDuplicates()
	}
	return fmt.Errorf("database type does not support cleanup")
}
