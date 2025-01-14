package mysql

import (
	"bytes"
	"database/sql"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql" // MySQL driver

	"github.com/media-luna/eureka/configs"
)

// DB is a type that implements MyInterface.
// DB represents a MySQL database connection.
// It holds a connection to the database and the configuration settings.
type DB struct {
	conn *sql.DB
	cfg config.Config
}

// NewDB creates a new DB instance with the given configuration.
func NewDB(cfg config.Config) (*DB, error) {
	db := &DB{cfg: cfg}

	if err := db.connect(); err != nil {
		return nil, err
	}
	return db, nil
}

// Connect to the MySQL database.
func (m *DB) connect() error {
	var err error
	dsnString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", m.cfg.Database.User, m.cfg.Database.Password, m.cfg.Database.Host, m.cfg.Database.Port, m.cfg.Database.DBName, m.cfg.Database.Params)
	m.conn, err = sql.Open("mysql", dsnString)
	return err
}

// Setup the MySQL database with tables.
func (m *DB) Setup() error {
	templates := []string{
		m.cfg.SQLTemplates.Template.CreateSongsTable,
		m.cfg.SQLTemplates.Template.CreateFingerprintsTable,
		m.cfg.SQLTemplates.Template.DeleteUnfingerprinted,
	}

	for _, tmpl := range templates {
		tmplPath := filepath.Join(m.cfg.SQLTemplates.MySQL, tmpl)
		queryString, err := m.parseQueryTemplate(tmplPath, m.cfg.Tables)
		if err != nil {
			return fmt.Errorf("failed to parse query template %s: %w", tmplPath, err)
		}

		// Execute the query
		_, err = m.conn.Exec(queryString)
		if err != nil {
			return fmt.Errorf("failed to execute query %s: %w", tmplPath, err)
		}
	}

	return nil
}

// Parse the SQL template.
func (m *DB) parseQueryTemplate(queryTmplPath string, tables config.Tables) (string, error) {
	// Step 1: Read the file
	content, err := os.ReadFile(queryTmplPath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}
	
	// Step 2: Parse the SQL template
	tmpl, err := template.New("sqlTemplate").Parse(string(content))
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}
	
	// Step 3: Execute the template with data
	var output bytes.Buffer
	if err := tmpl.Execute(&output, tables); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return output.String(), nil
}

// Close the MySQL database connection.
func (m *DB) Close() error {
    return m.conn.Close()
}

// Insert song metadata into songs table
func (m *DB) InsertSong(songName string, artistName string, fileHash string, totalHashes int) error {
	return nil
}

// Insert fingerprints into fingerprints table
func (m *DB) InsertFingerprints(fingerprint string, songID int, offset int) error {
	return nil
}