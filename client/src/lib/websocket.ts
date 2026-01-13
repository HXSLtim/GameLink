import { useAuthStore } from '@/stores';

type MessageHandler = (data: any) => void;

class WebSocketService {
    private socket: WebSocket | null = null;
    private handlers: Map<string, Set<MessageHandler>> = new Map();
    private reconnectAttempts = 0;
    private maxReconnectAttempts = 5;
    private reconnectTimeout: ReturnType<typeof setTimeout> | null = null;
    private url: string;

    constructor() {
        // ws://localhost:8080/ws
        // In production, infer from window.location
        this.url = 'ws://localhost:8080/ws';
    }

    public connect() {
        if (this.socket?.readyState === WebSocket.OPEN) return;

        const token = useAuthStore.getState().token;
        if (!token) {
            console.warn("WebSocket: No token found, skipping connection");
            return;
        }

        // Remove token from query string (Security Fix: URL leakage)
        const wsUrl = this.url;

        // Use sub-protocol for authentication to prevent token logging in URL
        // Sub-protocol format: 'bearer-TOKEN' (Server must support this)
        const protocols = ['bearer-' + token];

        this.socket = new WebSocket(wsUrl, protocols);

        this.socket.onopen = () => {
            console.log('WebSocket connected');
            this.reconnectAttempts = 0;
        };

        this.socket.onmessage = (event) => {
            try {
                const message = JSON.parse(event.data);
                const { type, payload } = message;
                this.emit(type, payload);
            } catch (e) {
                console.error('WebSocket message parse error:', e);
            }
        };

        this.socket.onclose = () => {
            console.log('WebSocket disconnected');
            this.attemptReconnect();
        };

        this.socket.onerror = (error) => {
            console.error('WebSocket error:', error);
        };
    }

    public disconnect() {
        if (this.socket) {
            this.socket.close();
            this.socket = null;
        }
    }

    public send(type: string, payload: any) {
        if (this.socket?.readyState === WebSocket.OPEN) {
            this.socket.send(JSON.stringify({ type, payload }));
        } else {
            console.warn('WebSocket not connected, message dropped:', type);
        }
    }

    public on(type: string, handler: MessageHandler) {
        if (!this.handlers.has(type)) {
            this.handlers.set(type, new Set());
        }
        this.handlers.get(type)!.add(handler);
    }

    public off(type: string, handler: MessageHandler) {
        if (this.handlers.has(type)) {
            this.handlers.get(type)!.delete(handler);
        }
    }

    private emit(type: string, payload: any) {
        if (this.handlers.has(type)) {
            this.handlers.get(type)!.forEach((handler) => handler(payload));
        }
    }

    private attemptReconnect() {
        if (this.reconnectAttempts < this.maxReconnectAttempts) {
            this.reconnectAttempts++;
            const delay = Math.min(1000 * Math.pow(2, this.reconnectAttempts), 30000);
            console.log(`WebSocket reconnecting in ${delay}ms...`);

            if (this.reconnectTimeout) clearTimeout(this.reconnectTimeout);

            this.reconnectTimeout = setTimeout(() => {
                this.connect();
            }, delay);
        }
    }
}

export const wsService = new WebSocketService();
