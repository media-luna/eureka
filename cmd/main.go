package main

import (
	"fmt"

	"github.com/media-luna/eureka/configs"
	"github.com/media-luna/eureka/internal/eureka"
)

func main() {
    // Load configuration
    configFilePath := "../configs/config.yaml"
	config, err := config.LoadConfig(configFilePath)
	if err != nil {
        fmt.Printf("Failed to load configuration: %v\n", err)
	}

    app := eureka.NewEureka(*config)
	app.Save("/home/daniel/projects/jamaivu/media/musicbox161/Good_Things_Go.wav") // From_Zero.wav
    println(app)
}