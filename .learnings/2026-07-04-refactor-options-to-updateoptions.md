# Learning: Refactored Options Management via UpdateOptions Method

## Requirement
Refactor option parsing logic out of `NewChatWithOptions` in `internal/ollama/chat.go` into a dedicated helper/method `UpdateOptions(opts map[string]any)` on `OllamaChat`.

## Refactoring Decisions
1. **Chat Options Structuring**:
   - Extracted option parsing for `temperature`, `num_ctx`, and `num_predict` into `(oc *OllamaChat) UpdateOptions(opts map[string]any)`.
   - Updated `NewChatWithOptions` to instantiate the `OllamaChat` with defaults and delegate parsing to `chat.UpdateOptions(opts)`.
2. **Simplified Handler Integration**:
   - Refactored `internal/api/handler.go`'s parameter updating logic from building an ephemeral `OllamaChat` to calling `chat.UpdateOptions(payload.Options)` directly. This makes option updates during active chat sessions extremely clean and efficient.
