package api

import (
	"net/http"

	"github.com/devakxhay/ollama-chat-ui-go/internal/ollama"
)

type ChatServer struct {
	ollamaClient *ollama.OllamaClient
	activeChat   *ollama.OllamaChat
}

func NewChatServer(client *ollama.OllamaClient) *ChatServer {
	return &ChatServer{
		ollamaClient: client,
		activeChat:   nil,
	}
}

func (cs *ChatServer) SetupRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /", cs.handleIndexPage)
	router.HandleFunc("GET /app.js", cs.serverAppJs)
	router.HandleFunc("GET /api/models", cs.handleModels)
	router.HandleFunc("POST /api/chat", cs.handleChatStream)
}
