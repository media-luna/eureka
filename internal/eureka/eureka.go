package eureka

import (
	"fmt"

	"github.com/media-luna/eureka/configs"
    "github.com/media-luna/eureka/internal/database"
)

type Config struct {
    Database         config.DBConfig
    DatabaseType     string
    FingerprintLimit int
}

type Eureka struct {
    Config  Config
    Name    string
    Version string
}

// NewSystem creates and initializes a new System.
func NewEureka(name, version string, config Config) *Eureka {
    // Init DB object
    database, err := database.NewDatabase(config.DatabaseType, config.Database)
    if err != nil {
        fmt.Println("error initializing database:", err)
        return nil
    }

    // Connect to DB
    if err := database.Connect(); err != nil {
        fmt.Println("error connecting to database:", err)
        return nil
    }
    defer database.Close()

    fmt.Println("Database connected successfully!")

    return &Eureka{
        Name:    name,
        Version: version,
        Config: config,
    }
}