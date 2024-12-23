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

type PostgresDB struct {
    conn *sql.DB
    cfg  config.Config
}

// NewPostgresDB creates a new PostgresDB instance with the given configuration.
func NewPostgresDB(cfg config.Config) *PostgresDB {
    return &PostgresDB{cfg: cfg}
}

func (p *PostgresDB) Connect() error {
    dsnString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?%s", p.cfg.Database.User, p.cfg.Database.Password, p.cfg.Database.Host, p.cfg.Database.Port, p.cfg.Database.DBName, p.cfg.Database.Params)
    var err error
    p.conn, err = sql.Open("postgres", dsnString)
    return err
}

// Parse the SQL template.
func (db *PostgresDB) parseQueryTemplate(queryTmplPath string, tables config.Tables) (string, error) {
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

func (p *PostgresDB) Setup() error {
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

func (p *PostgresDB) Close() error {
    return p.conn.Close()
}
