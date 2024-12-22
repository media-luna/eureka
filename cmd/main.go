package main

import (
    "github.com/media-luna/eureka/internal/eureka"
    "github.com/media-luna/eureka/configs"
)

func main() {
    // Create the system object
    config := eureka.Config{
            DatabaseType: "mysql",
            FingerprintLimit: 0,
            Database: config.DBConfig{
                User: "mysql",
                Password: "password",
                DBName: "dejavu",
                Host: "daniel-server.local",
                Port: 3306,
                Params: "",
            },
    }

    app := eureka.NewEureka("Eureka", "0.0.1", config)

    // Start the system
    println(app.Name)
}