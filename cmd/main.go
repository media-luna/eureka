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
    println(app)
}