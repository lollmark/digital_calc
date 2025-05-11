package main

import (
	"log"

	"github.com/lollmark/digital_calc/internal"
)

func main() {
	agent := application.NewAgent()
	log.Println("Starting Agent...")
	agent.Run()
}
