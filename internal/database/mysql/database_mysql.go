package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"

	"github.com/media-luna/eureka/configs"
)

// MySQL is a type that implements MyInterface.
type MySQLDB struct {
    conn *sql.DB
	cfg config.DBConfig
}

// NewMySQLDB creates a new MySQLDB instance with the given configuration.
func NewMySQLDB(cfg config.DBConfig) *MySQLDB {
    return &MySQLDB{cfg: cfg}
}

// Empty is a no-op function.
func (m *MySQLDB) Setup() error {
    fmt.Println("Empty function called")
	return fmt.Errorf("not implemented yet")
}

func (m *MySQLDB) Connect() error {
    var err error
	dsnString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", m.cfg.User, m.cfg.Password, m.cfg.Host, m.cfg.Port, m.cfg.DBName, m.cfg.Params)
    m.conn, err = sql.Open("mysql", dsnString)
    return err
}

func (m *MySQLDB) Close() error {
    return m.conn.Close()
}