package ollama

import "net/http"

type OllamaClient struct {
	baseUrl    string
	httpClient *http.Client
}

func NewOllamaClient(baseUrl string) *OllamaClient {
	return &OllamaClient{
		baseUrl:    baseUrl,
		httpClient: http.DefaultClient,
	}
}
