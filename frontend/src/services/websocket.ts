export class WebSocketService {
  private ws: WebSocket | null = null;
  private messageHandlers: Set<(msg: any) => void> = new Set();
  private reconnectTimer: ReturnType<typeof setTimeout> | null = null;
  private intentionalClose = false;
  private url: string;

  constructor(token: string) {
    // Determine WS protocol based on HTTP protocol
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    // Remove trailing slash if present
    const host = 'localhost:8080';
    this.url = `${protocol}//${host}/ws?token=${token}`; // Assuming we pass token in query. Wait, Hertz JWTAuth uses Header Authorization: Bearer token usually. Or Query parameter?
  }
  
  connect() {
    if (this.ws) return;
    this.intentionalClose = false;
    
    // JWT auth middleware in Hertz usually checks both query "token" and header. Let's try passing in query "token"
    this.ws = new WebSocket(this.url);

    this.ws.onopen = () => {
      console.log('WebSocket connected');
      if (this.reconnectTimer) clearTimeout(this.reconnectTimer);
    };

    this.ws.onmessage = (event) => {
      try {
        const lines = event.data.split('\n');
        for (const line of lines) {
          if (!line.trim()) continue;
          const data = JSON.parse(line);
          this.messageHandlers.forEach(handler => handler(data));
        }
      } catch (e) {
        console.error('Failed to parse WS message', e);
      }
    };

    this.ws.onclose = () => {
      console.log('WebSocket disconnected');
      this.ws = null;
      if (!this.intentionalClose) {
        this.reconnectTimer = setTimeout(() => this.connect(), 3000);
      }
    };

    this.ws.onerror = (error) => {
      console.error('WebSocket error', error);
      this.ws?.close();
    };
  }

  sendMessage(to: number, message: string, type: string = 'text') {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify({ action: 'message', to, message, type }));
    } else {
      console.error('WebSocket is not connected');
    }
  }

  sendStatusUpdate(to: number, status: string) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify({ action: 'status_update', to, status }));
    }
  }

  onMessage(handler: (msg: any) => void) {
    this.messageHandlers.add(handler);
    return () => {
      this.messageHandlers.delete(handler);
    };
  }

  disconnect() {
    this.intentionalClose = true;
    if (this.reconnectTimer) clearTimeout(this.reconnectTimer);
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
  }
}
