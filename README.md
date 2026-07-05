# Ollama Chat UI

A lightweight, responsive web-based chat interface for Ollama LLMs. It features a Go backend and an Alpine.js/Tailwind CSS frontend with streaming responses and live parameter configuration.

## Features

- **Model Selector**: Automatically fetches and lists your locally installed Ollama models.
- **Streaming Responses**: Real-time markdown/text streaming directly from the LLM.
- **Sidebar Configuration**: Customize model parameters on the fly:
  - Temperature (control randomness)
  - Max Predict Tokens (control maximum response length)
  - Context Window (adjust memory buffer size)
- **Sleek UI**: Minimalist, dark-themed responsive design.

## Prerequisites

- **Go** (1.20 or newer)
- **Ollama** installed and running on your local machine (default: `http://localhost:11434`).

## Getting Started

1. **Start Ollama**
   Make sure the Ollama service is up and running:
   ```bash
   ollama serve
   ```

2. **Run the Application**
   Clone the repository, navigate to the directory, and run the server:
   ```bash
   # Set the OLLAMA_BASE_URL environment variable if your Ollama is running on a different port or host
   export OLLAMA_BASE_URL="http://localhost:11434"

   go run main.go
   ```

3. **Access the Chat Interface**
   Open your browser and navigate to:
   ```
   http://localhost:8080
   ```

## Project Structure

- `main.go` — Server entry point and configuration.
- `internal/` — Go backend code.
  - `api/` — HTTP server routes, handlers, and assets server.
  - `ollama/` — Ollama API client and chat management logic.
- `ui/` — Frontend assets.
  - `index.html` — Main interface structure styled with Tailwind CSS and Alpine.js logic.
  - `app.js` — Client-side app state and streaming request/response logic.
