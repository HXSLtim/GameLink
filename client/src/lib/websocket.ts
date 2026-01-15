import { useAuthStore } from '@/stores';

type MessageHandler = (data: any) => void;

// Get WebSocket URL from environment variable
const WS_URL = import.meta.env.VITE_WS_URL || 'ws://localhost:8080/ws';

export class WebSocketService {
    private socket: WebSocket | null = null;
    private handlers: Map<string, Set<MessageHandler>> = new Map();
    private reconnectAttempts = 0;
    private maxReconnectAttempts = 5;
    private reconnectTimeout: ReturnType<typeof setTimeout> | null = null;
    private url: string;

    constructor() {
        this.url = WS_URL;
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

        // Authentication via message to prevent token logging
        // Standard WebSocket connection without sensitive data in URL/Protocols
        this.socket = new WebSocket(wsUrl);

        this.socket.onopen = () => {
            // console.debug('WebSocket connected');
            this.reconnectAttempts = 0;

            // Send authentication message immediately after connection
            if (token) {
                this.send('auth', { token });
            }
        };

        this.socket.onmessage = (event) => {
            try {
                const message = JSON.parse(event.data);
                const { type, payload } = message;

                // Handle auth response if needed, otherwise emit
                if (type === 'auth_success') {
                    // console.debug('WebSocket authenticated successfully');
                } else if (type === 'auth_error') {
                    console.error('WebSocket authentication failed:', payload);
                    this.disconnect(); // Or handle token refresh logic
                } else {
                    this.emit(type, payload);
                }
            } catch (e) {
                console.error('WebSocket message parse error:', e);
            }
        };

        this.socket.onclose = () => {
            // console.debug('WebSocket disconnected');
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
            // console.debug(`WebSocket reconnecting in ${delay}ms...`);

            if (this.reconnectTimeout) clearTimeout(this.reconnectTimeout);

            this.reconnectTimeout = setTimeout(() => {
                this.connect();
            }, delay);
        }
    }
}

export const wsService = new WebSocketService();
