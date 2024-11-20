package main

import (
	"fmt"
	"gpt4cli-server/db"
	"gpt4cli-server/host"
	"gpt4cli-server/model/plan"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func main() {

	err := host.LoadIp()
	if err != nil {
		log.Fatal("Error loading IP: ", err)
	}

	err = db.Connect()
	if err != nil {
		log.Fatal("Error initializing database: ", err)
	}

	err = db.MigrationsUp()
	if err != nil {
		log.Fatal("Error running migrations: ", err)
	}

	if os.Getenv("GOENV") == "development" {
		log.Println("In development mode.")
	}

	// Get externalPort from the environment variable or default to 8080
	externalPort := os.Getenv("PORT")
	if externalPort == "" {
		externalPort = "8080"
	}

	go startServer(externalPort, routes())
	log.Println("Started server on port " + externalPort)

	sigTermChan := make(chan os.Signal, 1)
	signal.Notify(sigTermChan, syscall.SIGTERM)

	go func() {
		<-sigTermChan

		for {
			l := plan.NumActivePlans()
			if l == 0 {
				break
			}
			log.Printf("Waiting for %d active plans to finish...\n", l)
			time.Sleep(1 * time.Second)
		}

		os.Exit(0)
	}()

	select {}
}

func startServer(port string, routes *mux.Router) {
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), routes)
	if err != nil {
		log.Fatalf("Failed to start server on port %s: %v", port, err)
	}
}
