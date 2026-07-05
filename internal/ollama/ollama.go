package ollama

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type ChatRequestBody struct {
	Model    string      `json:"model"`
	Messages []Message   `json:"messages"` // chat history
	Stream   bool        `json:"stream"`
	Options  ChatOptions `json:"options"`
}

type ChatOptions struct {
	Temperature float32 `json:"temperature,omitempty"`
	NumCtx      uint    `json:"num_ctx,omitempty"`
	NumPredict  uint    `json:"num_predict,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatResponseBody struct {
	Model   string  `json:"model"`
	Message Message `json:"message"`
	Done    bool    `json:"done"`
}

type OllamaModelsResponse struct {
	Models []OllamaModel `json:"models"`
}

type OllamaModel struct {
	Name    string            `json:"name"`
	Details OllamaModelDetail `json:"details"`
}

type OllamaModelDetail struct {
	ParameterSize string `json:"parameter_size"`
}

func (oc *OllamaClient) GetModels() ([]OllamaModel, error) {
	resp, err := oc.httpClient.Get(oc.baseUrl + "/api/tags")
	if err != nil {
		return nil, fmt.Errorf("Error in GetModels: %s", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		rb, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Ollama returned status: %s / body: %s", resp.Status, string(rb))
	}

	rb, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	var modelResponse OllamaModelsResponse
	if err := json.Unmarshal(rb, &modelResponse); err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	return modelResponse.Models, nil
}

func (oc *OllamaClient) Chat(chat *OllamaChat, newMessage Message) (string, error) {
	// Append the user message
	chat.AppendUserMessage(newMessage.Content)

	// Prepare request body
	reqBody := ChatRequestBody{
		Model:    chat.Model,
		Messages: chat.Messages,
		Stream:   false,
		Options: ChatOptions{
			Temperature: chat.Temperature,
			NumCtx:      chat.CtxWindow,
			NumPredict:  chat.MaxPredictSize,
		},
	}

	jb, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("Failed to marshal request body: %v", err)
	}

	url := fmt.Sprintf("%s/api/chat", oc.baseUrl)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jb))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := oc.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}

	defer resp.Body.Close()

	var respBody ChatResponseBody
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return "", fmt.Errorf("failed to decode the response body: %v", err)
	}

	chat.AppendAssistantMessage(respBody.Message.Content)
	return respBody.Message.Content, nil
}

func (oc *OllamaClient) ChatStream(chat *OllamaChat, newMessage Message) (<-chan string, <-chan error, error) {
	// Append the user message
	chat.AppendUserMessage(newMessage.Content)

	// Prepare request body
	reqBody := ChatRequestBody{
		Model:    chat.Model,
		Messages: chat.Messages,
		Stream:   true,
		Options: ChatOptions{
			Temperature: chat.Temperature,
			NumCtx:      chat.CtxWindow,
			NumPredict:  chat.MaxPredictSize,
		},
	}

	jb, err := json.Marshal(reqBody)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to marshal request body: %v", err)
	}

	url := fmt.Sprintf("%s/api/chat", oc.baseUrl)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jb))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := oc.httpClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to send request: %v", err)
	}

	out := make(chan string)
	errChan := make(chan error)

	go func() {
		defer resp.Body.Close()
		defer close(out)
		defer close(errChan)

		if resp.StatusCode != http.StatusOK {
			errChan <- fmt.Errorf("Ollama returned status: %s", resp.Status)
			return
		}

		scanner := bufio.NewScanner(resp.Body)
		scanner.Split(bufio.ScanLines)

		var chatChunk ChatResponseBody
		var fullMessage strings.Builder

		for scanner.Scan() {
			if err = json.Unmarshal([]byte(scanner.Text()), &chatChunk); err != nil {
				errChan <- err
				return
			}

			fullMessage.WriteString(chatChunk.Message.Content)
			out <- chatChunk.Message.Content
		}

		chat.AppendAssistantMessage(fullMessage.String())
	}()

	return out, errChan, nil
}
