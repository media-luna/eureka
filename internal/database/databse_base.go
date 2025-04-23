package database

import (
	"fmt"

	config "github.com/media-luna/eureka/configs"
	"github.com/media-luna/eureka/internal/database/mysql"
)

// Database defines the interface that all database implementations must satisfy
type Database interface {
	Setup() error
	Close() error
	// BeforeFork()
	// AfterFork()
	// Empty()
	// DeleteUnfingerprintedSongs()
	// GetNumSongs() int
	// GetNumFingerprints() int
	// SetSongFingerprinted(songID int)
	// GetSongs() []map[string]string
	// GetSongByID(songID int) map[string]string
	InsertFingerprints(fingerprint string, songID int, offset int) error
	InsertSong(songName string, artistName string, fileHash string, totalHashes int) (int, error)
	// Qurey(fingerprint string) []string
	// GetIterableKVPairs() []string
	// InsetHashes(songID int, hashes []map[string]int, batchSize int)
	// ReturnMatches(hashes []map[string]int, batchSize int) []map[string]string
	// DeleteSongById(songIDs []int, batchSize int)
}

// NewDatabase creates a new database instance based on the configuration
func NewDatabase(cfg config.Config) (Database, error) {
	switch cfg.Database.Type {
	case "mysql":
		return mysql.NewDB(cfg)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.Database.Type)
	}
}
