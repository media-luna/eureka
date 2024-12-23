package eureka

import (
	"fmt"

	"github.com/media-luna/eureka/configs"
	"github.com/media-luna/eureka/internal/database"
)

// Eureka represents the main structure for the Eureka service,
// containing the configuration settings required for its operation.
type Eureka struct {
    Config  config.Config
}

// NewSystem creates and initializes a new System.
// NewEureka initializes a new Eureka instance with the provided configuration.
// It performs the following steps:
// 1. Initializes the database object using the provided configuration.
// 2. Connects to the database.
// 3. Sets up the database.
//
// If any of these steps fail, it logs the error and returns nil.
//
// Parameters:
//   - config: The configuration object used to initialize the database.
//
// Returns:
//   - A pointer to the initialized Eureka instance, or nil if an error occurred.
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