package main

import (
	"fmt"

	"github.com/rifkiadrn/cassandra-explore/config"
)

func main() {
	viperConfig := config.NewViper()
	log := config.NewLogger(viperConfig)
	noSQLDB := config.NewNoSQLDatabase(viperConfig, log)
	validate := config.NewValidator(viperConfig)
	app := config.NewFiber(viperConfig)

	fmt.Printf("noSQLDB: %d \n", noSQLDB)

	config.Bootstrap(&config.BootstrapConfig{
		NoSQLDB:  noSQLDB,
		App:      app,
		Log:      log,
		Validate: validate,
		Config:   viperConfig,
	})

	webPort := viperConfig.GetInt("app.port")
	fmt.Println(webPort)
	err := app.Listen(fmt.Sprintf(":%d", webPort))
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
