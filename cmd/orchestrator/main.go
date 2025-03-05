package main

import (
	"log"
	"net/http"
	"digitalcalc/internal/orchestrator"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("Failed to initialize logger", err)
	}
	defer logger.Sync()

	log.Println("Starting Orchestrator...")
	router := orchestrator.NewRouter(logger)
	log.Fatal(http.ListenAndServe(":8080", router))
}