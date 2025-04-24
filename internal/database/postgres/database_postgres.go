package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	config "github.com/media-luna/eureka/configs"
)

// DB represents a PostgreSQL database connection.
type DB struct {
	conn *sql.DB
	cfg  config.Config
}

const (
	createSongsTableSQL = `
		CREATE TABLE IF NOT EXISTS %s (
			%s SERIAL PRIMARY KEY,
			%s VARCHAR(250) NOT NULL,
			%s SMALLINT DEFAULT 0,
			%s BYTEA NOT NULL,
			%s INTEGER NOT NULL DEFAULT 0,
			date_created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			date_modified TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`

	createFingerprintsTableSQL = `
		CREATE TABLE IF NOT EXISTS %s (
			%s BYTEA NOT NULL,
			%s INTEGER NOT NULL REFERENCES %s(%s) ON DELETE CASCADE,
			%s INTEGER NOT NULL,
			date_created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			date_modified TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE (%s, %s, %s)
		);
		CREATE INDEX IF NOT EXISTS ix_%s_%s ON %s (%s);`

	deleteUnfingerprintedSQL = `DELETE FROM %s WHERE %s = 0;`
)

// NewDB creates a new DB instance with the given configuration.
func NewDB(cfg config.Config) (*DB, error) {
	db := &DB{cfg: cfg}
	if err := db.connect(); err != nil {
		return nil, err
	}
	return db, nil
}

// Connect to the PostgreSQL database.
func (p *DB) connect() error {
	var err error
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s %s",
		p.cfg.Database.Host,
		p.cfg.Database.Port,
		p.cfg.Database.User,
		p.cfg.Database.Password,
		p.cfg.Database.DBName,
		p.cfg.Database.Params)
	p.conn, err = sql.Open("postgres", connStr)
	return err
}

// Setup initializes the database tables.
func (p *DB) Setup() error {
	// Create songs table
	songsSQL := fmt.Sprintf(createSongsTableSQL,
		p.cfg.Tables.Songs.Name,
		p.cfg.Tables.Songs.Fields.ID,
		p.cfg.Tables.Songs.Fields.Name,
		p.cfg.Tables.Songs.Fields.Fingerprinted,
		p.cfg.Tables.Songs.Fields.FileSHA1,
		p.cfg.Tables.Songs.Fields.TotalHashes)

	if _, err := p.conn.Exec(songsSQL); err != nil {
		return fmt.Errorf("error creating songs table: %w", err)
	}

	// Create fingerprints table
	fpSQL := fmt.Sprintf(createFingerprintsTableSQL,
		p.cfg.Tables.Fingerprints.Name,
		p.cfg.Tables.Fingerprints.Fields.Hash,
		p.cfg.Tables.Songs.Fields.ID,
		p.cfg.Tables.Songs.Name,
		p.cfg.Tables.Songs.Fields.ID,
		p.cfg.Tables.Fingerprints.Fields.Offset,
		p.cfg.Tables.Songs.Fields.ID,
		p.cfg.Tables.Fingerprints.Fields.Hash,
		p.cfg.Tables.Fingerprints.Fields.Offset,
		p.cfg.Tables.Fingerprints.Name,
		p.cfg.Tables.Fingerprints.Fields.Hash,
		p.cfg.Tables.Fingerprints.Name,
		p.cfg.Tables.Fingerprints.Fields.Hash)

	if _, err := p.conn.Exec(fpSQL); err != nil {
		return fmt.Errorf("error creating fingerprints table: %w", err)
	}

	// Delete unfingerprinted songs
	cleanupSQL := fmt.Sprintf(deleteUnfingerprintedSQL,
		p.cfg.Tables.Songs.Name,
		p.cfg.Tables.Songs.Fields.Fingerprinted)

	if _, err := p.conn.Exec(cleanupSQL); err != nil {
		return fmt.Errorf("error cleaning up unfingerprinted songs: %w", err)
	}

	return nil
}

// Close closes the database connection.
func (p *DB) Close() error {
	return p.conn.Close()
}

// Insert song metadata into songs table
func (p *DB) InsertSong(songName string, artistName string, fileHash string, totalHashes int) (int, error) {
	query := fmt.Sprintf("INSERT INTO %s (%s, %s, %s) VALUES ($1, decode($2, 'hex'), $3) RETURNING %s",
		p.cfg.Tables.Songs.Name,
		p.cfg.Tables.Songs.Fields.Name,
		p.cfg.Tables.Songs.Fields.FileSHA1,
		p.cfg.Tables.Songs.Fields.TotalHashes,
		p.cfg.Tables.Songs.Fields.ID)

	var id int
	err := p.conn.QueryRow(query, songName, fileHash, totalHashes).Scan(&id)
	return id, err
}

// Insert fingerprints into fingerprints table
func (p *DB) InsertFingerprints(fingerprint string, songID int, offset int) error {
	query := fmt.Sprintf("INSERT INTO %s (%s, %s, %s) VALUES ($1, decode($2, 'hex'), $3) ON CONFLICT DO NOTHING",
		p.cfg.Tables.Fingerprints.Name,
		p.cfg.Tables.Songs.Fields.ID,
		p.cfg.Tables.Fingerprints.Fields.Hash,
		p.cfg.Tables.Fingerprints.Fields.Offset)

	_, err := p.conn.Exec(query, songID, fingerprint, offset)
	return err
}

// UpdateSongFingerprinted marks a song as fingerprinted in the database
func (p *DB) UpdateSongFingerprinted(songID int) error {
	updateQuery := fmt.Sprintf("UPDATE %s SET %s = 1 WHERE %s = $1",
		p.cfg.Tables.Songs.Name,
		p.cfg.Tables.Songs.Fields.Fingerprinted,
		p.cfg.Tables.Songs.Fields.ID)

	result, err := p.conn.Exec(updateQuery, songID)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("song with ID %d not found", songID)
	}

	return nil
}
