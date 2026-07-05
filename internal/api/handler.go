package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/devakxhay/ollama-chat-ui-go/internal/ollama"
)

type ChatAPIRequest struct {
	Model   string         `json:"model"`
	Message ollama.Message `json:"message"`
	Options map[string]any `json:"options"`
}

type ChatAPIResponse struct {
	Message  ollama.Message   `json:"message"`
	Messages []ollama.Message `json:"messages"`
}

func (cs *ChatServer) handleChat(w http.ResponseWriter, r *http.Request) {

	var payload ChatAPIRequest

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Initialize the Chat
	var chat *ollama.OllamaChat
	if cs.activeChat == nil {
		chat = ollama.NewChatWithOptions(payload.Model, payload.Options)
		cs.activeChat = chat
	} else {
		chat = cs.activeChat
		chat.UpdateOptions(payload.Options)
	}

	message, err := cs.ollamaClient.Chat(chat, payload.Message)
	if err != nil {
		http.Error(w, "Failed to generate response", http.StatusInternalServerError)
		return
	}

	resp := ChatAPIResponse{
		Message: ollama.Message{
			Role:    "assistant",
			Content: message,
		},
		Messages: chat.Messages,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (cs *ChatServer) handleChatStream(w http.ResponseWriter, r *http.Request) {

	var payload ChatAPIRequest

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Initialize the Chat
	var chat *ollama.OllamaChat
	if cs.activeChat == nil {
		chat = ollama.NewChatWithOptions(payload.Model, payload.Options)
		cs.activeChat = chat
	} else {
		chat = cs.activeChat
		chat.UpdateOptions(payload.Options)
	}

	out, errCh, err := cs.ollamaClient.ChatStream(chat, payload.Message)
	if err != nil {
		http.Error(w, "Failed to generate response", http.StatusInternalServerError)
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Transfer-Encoding", "chunked")

	// Trigger sending headers before entering stream loop
	w.WriteHeader(http.StatusOK)
	flusher.Flush()

	for {
		select {
		case <-r.Context().Done():
			return
		case msg, ok := <-out:
			if !ok {
				out = nil
				break
			}
			if chunkBytes, err := json.Marshal(map[string]string{"content": msg}); err == nil {
				fmt.Fprintf(w, "data: %s\n\n", string(chunkBytes))
				flusher.Flush()
			}
		case err, ok := <-errCh:
			if !ok {
				errCh = nil
				break
			}
			fmt.Fprintf(w, "event: error\ndata: %s\n\n", err.Error())
			flusher.Flush()
			return
		}

		if out == nil && errCh == nil {
			break
		}
	}

	// Emit done event after stream ends normally
	fmt.Fprint(w, "event: done\ndata: {}\n\n")
	flusher.Flush()
}

func (cs *ChatServer) handleIndexPage(w http.ResponseWriter, r *http.Request) {
	cs.activeChat = nil

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.ServeFile(w, r, "ui/index.html")
}

func (cs *ChatServer) serverAppJs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	http.ServeFile(w, r, "ui/app.js")
}

func (cs *ChatServer) handleModels(w http.ResponseWriter, r *http.Request) {
	models, err := cs.ollamaClient.GetModels()
	if err != nil {
		http.Error(w, "Failed to fetch models", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models)
}
