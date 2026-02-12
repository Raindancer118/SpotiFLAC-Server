// WebSocket Client for real-time updates from SpotiFLAC Server

const WS_URL = import.meta.env.VITE_WS_URL || 'ws://localhost:8080/ws';

export interface DownloadProgress {
    is_downloading: boolean;
    mb_downloaded: number;
    mb_total: number;
    percentage: number;
    speed_mbps: number;
}

export interface QueueUpdate {
    queue: any[];
    is_downloading: boolean;
}

type MessageHandler = (data: any) => void;

/**
 * WebSocket client for real-time updates
 * Replaces Wails EventsOn/EventsOff
 */
class WebSocketClient {
    private ws: WebSocket | null = null;
    private reconnectTimer: NodeJS.Timeout | null = null;
    private handlers: Map<string, Set<MessageHandler>> = new Map();
    private reconnectAttempts = 0;
    private maxReconnectAttempts = 10;
    private reconnectDelay = 1000;

    constructor() {
        this.connect();
    }

    /**
     * Connect to WebSocket server
     */
    private connect(): void {
        try {
            this.ws = new WebSocket(WS_URL);

            this.ws.onopen = () => {
                console.log('WebSocket connected');
                this.reconnectAttempts = 0;
                this.reconnectDelay = 1000;
            };

            this.ws.onmessage = (event) => {
                try {
                    const message = JSON.parse(event.data);
                    this.handleMessage(message);
                } catch (error) {
                    console.error('Failed to parse WebSocket message:', error);
                }
            };

            this.ws.onclose = () => {
                console.log('WebSocket disconnected');
                this.scheduleReconnect();
            };

            this.ws.onerror = (error) => {
                console.error('WebSocket error:', error);
            };
        } catch (error) {
            console.error('Failed to create WebSocket:', error);
            this.scheduleReconnect();
        }
    }

    /**
     * Schedule reconnection attempt
     */
    private scheduleReconnect(): void {
        if (this.reconnectAttempts >= this.maxReconnectAttempts) {
            console.error('Max reconnection attempts reached');
            return;
        }

        if (this.reconnectTimer) {
            clearTimeout(this.reconnectTimer);
        }

        this.reconnectTimer = setTimeout(() => {
            this.reconnectAttempts++;
            console.log(`Reconnecting... (attempt ${this.reconnectAttempts})`);
            this.connect();
        }, this.reconnectDelay);

        // Exponential backoff
        this.reconnectDelay = Math.min(this.reconnectDelay * 2, 30000);
    }

    /**
     * Handle incoming WebSocket message
     */
    private handleMessage(message: { type: string; data: any }): void {
        const { type, data } = message;

        // Call all registered handlers for this event type
        const handlers = this.handlers.get(type);
        if (handlers) {
            handlers.forEach((handler) => handler(data));
        }
    }

    /**
     * Register event handler (replaces Wails EventsOn)
     */
    on(eventName: string, callback: MessageHandler): void {
        if (!this.handlers.has(eventName)) {
            this.handlers.set(eventName, new Set());
        }
        this.handlers.get(eventName)!.add(callback);
    }

    /**
     * Unregister event handler (replaces Wails EventsOff)
     */
    off(eventName: string, callback?: MessageHandler): void {
        if (!callback) {
            // Remove all handlers for this event
            this.handlers.delete(eventName);
        } else {
            // Remove specific handler
            const handlers = this.handlers.get(eventName);
            if (handlers) {
                handlers.delete(callback);
                if (handlers.size === 0) {
                    this.handlers.delete(eventName);
                }
            }
        }
    }

    /**
     * Send message to server
     */
    send(type: string, data?: any): void {
        if (this.ws && this.ws.readyState === WebSocket.OPEN) {
            this.ws.send(JSON.stringify({ type, data }));
        } else {
            console.warn('WebSocket not connected, cannot send message');
        }
    }

    /**
     * Request current status from server
     */
    requestStatus(): void {
        this.send('request_status');
    }

    /**
     * Send ping to server
     */
    ping(): void {
        this.send('ping');
    }

    /**
     * Close WebSocket connection
     */
    close(): void {
        if (this.reconnectTimer) {
            clearTimeout(this.reconnectTimer);
            this.reconnectTimer = null;
        }

        if (this.ws) {
            this.ws.close();
            this.ws = null;
        }

        this.handlers.clear();
    }

    /**
     * Check if connected
     */
    isConnected(): boolean {
        return this.ws?.readyState === WebSocket.OPEN;
    }
}

// Create singleton instance
export const wsClient = new WebSocketClient();

// Convenience functions that match Wails API
export const EventsOn = (eventName: string, callback: MessageHandler): void => {
    wsClient.on(eventName, callback);
};

export const EventsOff = (eventName: string, callback?: MessageHandler): void => {
    wsClient.off(eventName, callback);
};

// Export for use in components
export default wsClient;
