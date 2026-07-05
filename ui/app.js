const api = {
  async sendChat(payload) {
    const res = await fetch('/api/chat', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    });
    if (!res.ok) throw new Error(`HTTP Error ${res.status}`);
    return res;
  },
  async getModels() {
    const res = await fetch('/api/models')
    if (!res.ok) throw new Error(`HTTP Error ${res.status}`);
    return res;
  }
}


document.addEventListener('alpine:init', () => {
  Alpine.data('chatApp', () => ({
    userInput: '',
    messages: [],
    model: '',
    models: [],
    maxToken: 64, // num_predict
    temperature: 0.5,
    contextWindow: 2048, // num_ctx
    loading: false,
    showSettings: false,

    init() {
      api.getModels()
        .then(resp => resp.json())
        .then(data => {
          this.models = data;
          if (data.length > 0) this.model = data[0].name;
        })
    },
    async sendMessage() {
      if (!this.userInput.trim() || this.loading) return;

      const userMessage = { role: 'user', content: this.userInput };
      this.messages.push(userMessage);
      this.userInput = '';
      this.loading = true;

      try {
        const payload = {
          model: this.model,
          message: userMessage,
          options: {
            num_predict: parseInt(this.maxToken) || 64,
            temperature: parseFloat(this.temperature) ?? 0.5,
            num_ctx: parseInt(this.contextWindow) || 2048
          }
        };

        const resp = await api.sendChat(payload);

        let assistantMessage = { role: 'assistant', content: '' };
        this.messages.push(assistantMessage);
        const targetIdx = this.messages.length - 1;

        await consumeSSEStream(
          resp,
          (content) => {
            assistantMessage.content += content;
            this.messages[targetIdx] = { ...assistantMessage };
          },
          (errorMsg) => {
            assistantMessage.content = 'Error: ' + errorMsg;
            this.messages[targetIdx] = { ...assistantMessage };
          }
        );

        this.loading = false;
      } catch (err) {
        this.loading = false;
        console.error(err);
      }
    }
  }));
});

async function consumeSSEStream(response, onChunk, onError) {
  const reader = response.body.getReader();
  const decoder = new TextDecoder('utf-8');
  let buffer = '';
  let currentEvent = 'message';

  try {
    while (true) {
      const { done, value } = await reader.read();
      if (done) break;

      buffer += decoder.decode(value, { stream: true });
      const lines = buffer.split('\n\n');
      buffer = lines.pop(); // keep partial line in buffer

      for (const line of lines) {
        const cleanLine = line.trim();
        if (!cleanLine) continue;

        // Parse event headers
        if (cleanLine.startsWith('event:')) {
          currentEvent = cleanLine.replace('event:', '').trim();
          continue;
        }

        // Process data payloads
        if (cleanLine.startsWith('data:')) {
          const dataVal = cleanLine.slice(5).trim();

          if (currentEvent === 'error') {
            onError(dataVal);
            reader.cancel();
            return;
          }

          try {
            const parsed = JSON.parse(dataVal);
            onChunk(parsed.content);
          } catch (e) {
            onChunk(dataVal); // Fallback if raw text
          }

          // Reset status tracker
          currentEvent = 'message';
        }
      }
    }
  } catch (err) {
    onError(err.message);
  }
}