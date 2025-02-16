package main

import (
	"log"
	"gpt4cli-server/routes"
	"gpt4cli-server/setup"

	"github.com/gorilla/mux"
)

func main() {
	// Configure the default logger to include milliseconds in timestamps
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	r := mux.NewRouter()
	routes.AddHealthRoutes(r)
	routes.AddApiRoutes(r)

	setup.MustLoadIp()
	setup.MustInitDb()
	setup.StartServer(r)
}
