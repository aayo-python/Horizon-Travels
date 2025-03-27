package main

import (
	"log"

	"github.com/janto-pee/Horizon-Travels.git/controllers"
	"github.com/janto-pee/Horizon-Travels.git/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load configuration environment", err)
	}
	runGinServer(config)
}

func runGinServer(config util.Config) {
	server, err := controllers.NewServer(config)
	if err != nil {
		log.Fatal("Could not create server")
	}
	server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("Could not start server")
	}
}
