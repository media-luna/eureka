package eureka

import (
	"fmt"

	"github.com/media-luna/eureka/configs"
	"github.com/media-luna/eureka/internal/database"
)

// type Config struct {
//     Database         config.DBConfig
//     DatabaseType     string
//     FingerprintLimit int
// }

type Eureka struct {
    Config  config.Config
}

// NewSystem creates and initializes a new System.
func NewEureka(config config.Config) *Eureka {
    // Init DB object
    database, err := database.NewDatabase(config)
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

    // Setup DB
    if err := database.Setup(); err != nil {
        fmt.Println("error connecting to database:", err)
        return nil
    }
    
    return &Eureka{Config: config}
}