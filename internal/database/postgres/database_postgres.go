package postgres

import (
	"bytes"
	"database/sql"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	_ "github.com/lib/pq" // PostgreSQL driver

	"github.com/media-luna/eureka/configs"
)

// PostgresDB represents a PostgreSQL database connection and its configuration.
type DB struct {
    conn *sql.DB
    cfg  config.Config
}

// NewPostgresDB creates a new PostgresDB instance with the given configuration.
func NewDB(cfg config.Config) (*DB, error) {
	db := &DB{cfg: cfg}

	if err := db.connect(); err != nil {
		return nil, err
	}
	return db, nil
}

// Connect establishes a connection to the PostgreSQL database using the configuration
// provided in the DB struct. It constructs the DSN (Data Source Name) string from the
// configuration parameters and attempts to open a connection. If successful, it assigns
// the connection to the DB struct's conn field. It returns an error if the connection
// attempt fails.
func (p *DB) connect() error {
    dsnString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?%s", p.cfg.Database.User, p.cfg.Database.Password, p.cfg.Database.Host, p.cfg.Database.Port, p.cfg.Database.DBName, p.cfg.Database.Params)
    var err error
    p.conn, err = sql.Open("postgres", dsnString)
    return err
}

// Parse the SQL template.
func (p *DB) parseQueryTemplate(queryTmplPath string, tables config.Tables) (string, error) {
	// Read the file
	content, err := os.ReadFile(queryTmplPath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}
	
	// Parse the SQL template
	tmpl, err := template.New("sqlTemplate").Parse(string(content))
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	// Execute the template with data
	var output bytes.Buffer
	if err := tmpl.Execute(&output, tables); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return output.String(), nil
}

// Setup initializes the database by executing a series of SQL queries
// defined in the configuration templates. It parses each template,
// constructs the query string, and executes it against the database
// connection. If any step fails, it returns an error.
//
// Returns:
//   error: An error if any of the query templates fail to parse or execute.
func (p *DB) Setup() error {
    // Tebles templates query paths
	templates := []string{
		p.cfg.SQLTemplates.Template.CreateSongsTable,
		p.cfg.SQLTemplates.Template.CreateFingerprintsTable,
		p.cfg.SQLTemplates.Template.DeleteUnfingerprinted,
	}

	for _, tmpl := range templates {
		tmplPath := filepath.Join(p.cfg.SQLTemplates.Postgres, tmpl)
		queryString, err := p.parseQueryTemplate(tmplPath, p.cfg.Tables)
		if err != nil {
			return fmt.Errorf("failed to parse query template %s: %w", tmplPath, err)
		}

		// Execute the query
		_, err = p.conn.Exec(queryString)
		if err != nil {
			return fmt.Errorf("failed to execute query %s: %w", tmplPath, err)
		}
	}

	return nil
}

// Close closes the database connection.
// It returns an error if the connection cannot be closed.
func (p *DB) Close() error {
    return p.conn.Close()
}

func (m *DB) InsertSong(songName string, artistName string, fileHash string, totalHashes int) error {
	return nil
}

func (m *DB) InsertFingerprints(fingerprint string, songID int, offset int) error {
	return nil
}