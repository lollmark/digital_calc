package main

import (
	"log"

	"github.com/lollmark/calculator_go/internal"
)

func main() {
	agent := application.NewAgent()
	log.Println("Starting Agent...")
	agent.Run()
}
