import { useAuthStore } from '@/stores';

type MessageHandler = (data: unknown) => void;

// Get WebSocket URL from environment variable
const WS_URL = import.meta.env.VITE_WS_URL || 'ws://localhost:8080/ws';

// Heartbeat configuration
const HEARTBEAT_INTERVAL = 30000; // 30 seconds
const HEARTBEAT_TIMEOUT = 10000;  // 10 seconds to wait for pong

export class WebSocketService {
    private socket: WebSocket | null = null;
    private handlers: Map<string, Set<MessageHandler>> = new Map();
    private reconnectAttempts = 0;
    private maxReconnectAttempts = 5;
    private reconnectTimeout: ReturnType<typeof setTimeout> | null = null;
    private url: string;

    // Heartbeat state
    private heartbeatInterval: ReturnType<typeof setInterval> | null = null;
    private heartbeatTimeout: ReturnType<typeof setTimeout> | null = null;

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

            // Start heartbeat after connection
            this.startHeartbeat();
        };

        this.socket.onmessage = (event) => {
            try {
                const message = JSON.parse(event.data);
                const { type, payload } = message;

                // Handle pong response
                if (type === 'pong') {
                    this.handlePong();
                    return;
                }

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
            this.stopHeartbeat();
            this.attemptReconnect();
        };

        this.socket.onerror = (error) => {
            console.error('WebSocket error:', error);
        };
    }

    public disconnect() {
        this.stopHeartbeat();
        if (this.socket) {
            this.socket.close();
            this.socket = null;
        }
    }

    public send(type: string, payload: unknown) {
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

    public isConnected(): boolean {
        return this.socket?.readyState === WebSocket.OPEN;
    }

    private emit(type: string, payload: unknown) {
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

    // ============ Heartbeat Methods ============

    private startHeartbeat() {
        this.stopHeartbeat(); // Clear any existing heartbeat

        // Send ping at regular intervals
        this.heartbeatInterval = setInterval(() => {
            this.sendPing();
        }, HEARTBEAT_INTERVAL);

        // Send initial ping
        this.sendPing();
    }

    private stopHeartbeat() {
        if (this.heartbeatInterval) {
            clearInterval(this.heartbeatInterval);
            this.heartbeatInterval = null;
        }
        if (this.heartbeatTimeout) {
            clearTimeout(this.heartbeatTimeout);
            this.heartbeatTimeout = null;
        }
    }

    private sendPing() {
        if (this.socket?.readyState !== WebSocket.OPEN) return;

        // Send ping message
        this.send('ping', { timestamp: Date.now() });

        // Set timeout for pong response
        if (this.heartbeatTimeout) clearTimeout(this.heartbeatTimeout);

        this.heartbeatTimeout = setTimeout(() => {
            // No pong received within timeout - connection is dead
            console.warn('WebSocket heartbeat timeout - reconnecting...');
            this.socket?.close();
        }, HEARTBEAT_TIMEOUT);
    }

    private handlePong() {
        // Clear the timeout since we received pong
        if (this.heartbeatTimeout) {
            clearTimeout(this.heartbeatTimeout);
            this.heartbeatTimeout = null;
        }
    }
}

export const wsService = new WebSocketService();
