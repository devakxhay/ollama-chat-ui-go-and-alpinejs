package ollama

const (
	DEFAULT_CONTEXT_WINDOW   uint    = 2048 //2k
	DEFAULT_TEMP             float32 = 0.5
	DEFAUTL_MAX_PREDICT_SIZE uint    = 128
)

type OllamaChat struct {
	Model          string
	Messages       []Message
	CtxWindow      uint // Context Window Size
	Temperature    float32
	MaxPredictSize uint
}

func NewChat(model string) *OllamaChat {
	return NewChatWithOptions(model, nil)
}

func NewChatWithOptions(model string, opts map[string]any) *OllamaChat {
	chat := &OllamaChat{
		Model:          model,
		Messages:       []Message{},
		CtxWindow:      DEFAULT_CONTEXT_WINDOW,
		Temperature:    DEFAULT_TEMP,
		MaxPredictSize: DEFAUTL_MAX_PREDICT_SIZE,
	}
	chat.UpdateOptions(opts)
	return chat
}

func (oc *OllamaChat) AppendUserMessage(content string) {
	oc.appendMessage("user", content)
}

func (oc *OllamaChat) AppendAssistantMessage(content string) {
	oc.appendMessage("assistant", content)
}

func (oc *OllamaChat) appendMessage(role, content string) {
	switch role {
	case "user":
		oc.Messages = append(oc.Messages, Message{
			Role:    "user",
			Content: content,
		})

	case "assistant":
		oc.Messages = append(oc.Messages, Message{
			Role:    "assistant",
			Content: content,
		})

	}
}

func (oc *OllamaChat) UpdateOptions(opts map[string]any) {
	if opts == nil {
		return
	}

	if temp, ok := opts["temperature"]; ok {
		if t, ok := temp.(float64); ok && t >= 0 && t <= 1 {
			oc.Temperature = float32(t)
		} else if t, ok := temp.(float32); ok && t >= 0 && t <= 1 {
			oc.Temperature = t
		}
	}
	if ctx, ok := opts["num_ctx"]; ok {
		if c, ok := ctx.(float64); ok && c > 0 {
			oc.CtxWindow = uint(c)
		} else if c, ok := ctx.(int); ok && c > 0 {
			oc.CtxWindow = uint(c)
		}
	}
	if predict, ok := opts["num_predict"]; ok {
		if p, ok := predict.(float64); ok && p > 0 {
			oc.MaxPredictSize = uint(p)
		} else if p, ok := predict.(int); ok && p > 0 {
			oc.MaxPredictSize = uint(p)
		}
	}
}
