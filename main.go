package main

import (
	"log"
	"net/http"
	"os"

	"github.com/devakxhay/ollama-chat-ui-go/internal/api"
	"github.com/devakxhay/ollama-chat-ui-go/internal/ollama"
)

const DEFAULT_OLLAMA_BASE_URL = "http://localhost:11434"

func main() {
	ollamaUrl := os.Getenv("OLLAMA_BASE_URL")
	if ollamaUrl == "" {
		log.Printf("Ollama base URL not set. Using the default: %s\n", DEFAULT_OLLAMA_BASE_URL)
		ollamaUrl = DEFAULT_OLLAMA_BASE_URL
	}

	log.Printf("Using Ollama URL: %s", ollamaUrl)

	// Initialize the Ollama client
	client := ollama.NewOllamaClient(ollamaUrl)

	// Initialise the Chat Server
	server := api.NewChatServer(client)

	// Add routes to the server
	router := http.NewServeMux()
	server.SetupRoutes(router)

	// Listen and Serve
	log.Println("Starting server on port :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
