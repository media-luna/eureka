package database

import (
	"errors"

	"github.com/media-luna/eureka/internal/database/mysql"
	"github.com/media-luna/eureka/internal/database/postgres"
	"github.com/media-luna/eureka/configs"
)

// BaseDatabase defines methods for databse.
type BaseDatabase interface {
	// Called on creation or shortly afterwards.
	Setup() error
	Connect() error
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
	// Insert(fingerprint string, songID int, offset int)
	// InsertSong(songName string, fileHash string, totalHashes int)
	// Qurey(fingerprint string) []string
	// GetIterableKVPairs() []string
	// InsetHashes(songID int, hashes []map[string]int, batchSize int)
	// ReturnMatches(hashes []map[string]int, batchSize int) []map[string]string
	// DeleteSongById(songIDs []int, batchSize int)
}

// NewDatabase creates a new database connection based on the provided configuration.
// It supports "postgres" and "mysql" database types.
// Returns an instance of BaseDatabase or an error if the database type is invalid.
func NewDatabase(cfg config.Config) (BaseDatabase, error) {
    switch cfg.Database.Type {
    case "postgres":
        return postgres.NewDB(cfg), nil
    case "mysql":
        return mysql.NewDB(cfg), nil
    default:
        return nil, errors.New("invalid database type")
    }
}