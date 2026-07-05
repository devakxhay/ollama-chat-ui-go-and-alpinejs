# Learning: Sidebar Configuration Panel & Active Session Options Update

## Issue / Requirement
The user requested a minimal sidebar configuration panel to control:
- `max-predict-tokens` (default: 64)
- `temperature` (default: 0.5)
- `context window` (default: 2048)

## Design Decisions
1. **Sidebar layout**:
   - Added a collapsible sidebar settings panel toggleable via a settings gear/control button in the header.
   - Built modern, minimal range slider and number inputs leveraging Tailwind's sleek dark theme color palette.
   - Restructured the top-level body grid/flex to support a responsive `max-w-5xl` container, positioning settings next to the main chat pane.
2. **AlpineJS binding & JSON Payload Fix**:
   - Mapped settings state to the `maxToken`, `temperature`, and `contextWindow` values inside AlpineJS data structure.
   - Updated the `app.js` request builder to send parameter overrides in the request body `options` field (`num_predict`, `temperature`, `num_ctx`).
   - Fixed the API payload key from plural `messages` to singular `message` to match the Go backend `ChatAPIRequest` structure.
3. **Go Backend Option Propagation**:
   - Fixed the type-assertion logic in the Go backend (`internal/ollama/chat.go`) to check for `float64` first since JSON decoding into `map[string]any` automatically parses numbers as `float64`.
   - Updated the chat handler (`internal/api/handler.go`) to dynamic-update settings parameters mid-session if they are changed by the user.
