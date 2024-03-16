package main

import (
	"fmt"
	"github.com/BukhryakovVladimir/vkTest/internal/handlers/filmotekahandler"
	"github.com/BukhryakovVladimir/vkTest/internal/routes"
	"log"
	"net/http"
	"os"
)

func main() {
	err := routes.InitConnPool()

	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	mux := http.NewServeMux()

	filmotekahandler.SetupRoutes(mux)

	strPort := os.Getenv("PORT")
	if strPort == "" {
		log.Fatalf("Environment variable PORT is empty.")
	}
	port := fmt.Sprintf(":%s", strPort)

	err = http.ListenAndServe(port, mux)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
