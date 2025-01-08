package main

import (
	"fmt"
	"os"

	"github.com/media-luna/eureka/configs"
	"github.com/media-luna/eureka/internal/eureka"
)

func main() {
    // Load configuration
    configFilePath := "../configs/config.yaml"
	config, err := config.LoadConfig(configFilePath)
	if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
	}

	// Get Eureka app
    app, err := eureka.NewEureka(*config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	app.Save("/home/daniel/projects/jamaivu/media/musicbox161/Good_Things_Go.wav")
}