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

// MySQL is a type that implements MyInterface.
type MySQLDB struct {
    conn *sql.DB
	cfg config.Config
}

// NewMySQLDB creates a new MySQLDB instance with the given configuration.
func NewMySQLDB(cfg config.Config) *MySQLDB {
    return &MySQLDB{cfg: cfg}
}

// Connect to the MySQL database.
func (m *MySQLDB) Connect() error {
    var err error
	dsnString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", m.cfg.Database.User, m.cfg.Database.Password, m.cfg.Database.Host, m.cfg.Database.Port, m.cfg.Database.DBName, m.cfg.Database.Params)
    m.conn, err = sql.Open("mysql", dsnString)
    return err
}

// Parse the SQL template.
func (db *MySQLDB) parseQueryTemplate(queryTmplPath string, tables config.Tables) (string, error) {
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

// Setup the MySQL database with tables.
func (m *MySQLDB) Setup() error {
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
		println(queryString)
		_, err = m.conn.Exec(queryString)
		if err != nil {
			return fmt.Errorf("failed to execute query %s: %w", tmplPath, err)
		}
	}

	return nil
}

// Close the MySQL database connection.
func (m *MySQLDB) Close() error {
    return m.conn.Close()
}