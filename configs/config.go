package config

type DBConfig struct {
    Host     string
    Port     int
    User     string
    Password string
    DBName   string
    Params   string // Optional: for additional DSN parameters
}