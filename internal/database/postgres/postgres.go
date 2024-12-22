package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq" // PostgreSQL driver

	"github.com/media-luna/eureka/configs"
)

type PostgresDB struct {
    conn *sql.DB
    cfg  config.DBConfig
}

// NewPostgresDB creates a new PostgresDB instance with the given configuration.
func NewPostgresDB(cfg config.DBConfig) *PostgresDB {
    return &PostgresDB{cfg: cfg}
}

func (p *PostgresDB) Connect() error {
    dsnString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?%s", p.cfg.User, p.cfg.Password, p.cfg.Host, p.cfg.Port, p.cfg.DBName, p.cfg.Params)
    var err error
    p.conn, err = sql.Open("postgres", dsnString)
    return err
}

// func (p *PostgresDB) Query(query string, args ...interface{}) (interface{}, error) {
//     rows, err := p.conn.Query(query, args...)
//     if err != nil {
//         return nil, err
//     }
//     defer rows.Close()
//     var results []map[string]interface{}
//     for rows.Next() {
//         // Process rows
//     }
//     return results, nil
// }

func (p *PostgresDB) Setup() error {
    println("Empty function called")
    return nil
}

func (p *PostgresDB) Close() error {
    return p.conn.Close()
}
