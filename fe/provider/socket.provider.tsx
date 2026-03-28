"use client";

import { createContext, useContext, useEffect, useState } from "react";

type EventCallback = (data: unknown) => void;

// Thin WebSocket wrapper with an event-emitter API compatible with how the backend
// gofiber/contrib/v3/socketio library exchanges messages:
//   client → server: {"event":"<name>","data":<payload>}
//   server → client: {"event":"<name>","data":<payload>}
export class SocketClient {
    private ws: WebSocket;
    private listeners = new Map<string, Set<EventCallback>>();
    private pending: string[] = [];

    constructor(url: string) {
        this.ws = new WebSocket(url);

        this.ws.addEventListener("open", () => {
            this.pending.forEach((msg) => this.ws.send(msg));
            this.pending = [];
        });

        this.ws.addEventListener("message", (event) => {
            try {
                const msg = JSON.parse(event.data as string) as {
                    event: string;
                    data: unknown;
                };
                if (msg.event) {
                    this.listeners.get(msg.event)?.forEach((cb) => cb(msg.data));
                }
            } catch {
                // ignore malformed messages
            }
        });
    }

    emit(event: string, data?: unknown): void {
        const msg = JSON.stringify({ event, data });
        if (this.ws.readyState === WebSocket.OPEN) {
            this.ws.send(msg);
        } else {
            this.pending.push(msg);
        }
    }

    on(event: string, callback: EventCallback): void {
        if (!this.listeners.has(event)) {
            this.listeners.set(event, new Set());
        }
        this.listeners.get(event)!.add(callback);
    }

    off(event: string, callback: EventCallback): void {
        this.listeners.get(event)?.delete(callback);
    }

    disconnect(): void {
        this.ws.close();
    }
}

export const SocketContext = createContext<SocketClient | null>(null);

export const SocketProvider = ({ children }: { children: React.ReactNode }) => {
    const [client] = useState<SocketClient>(() => {
        const token = localStorage.getItem("auth_token");
        const userId = token ? token.split("_")[1] : "anonymous";
        const username = localStorage.getItem("username") ?? "";

        const apiUrl = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";
        const wsBase = apiUrl.replace(/^http/, "ws");
        const url = `${wsBase}/socket.io/?user_id=${encodeURIComponent(userId)}&username=${encodeURIComponent(username)}`;

        return new SocketClient(url);
    });

    useEffect(() => {
        return () => {
            client.disconnect();
        };
    }, [client]);

    return (
        <SocketContext.Provider value={client}>
            {children}
        </SocketContext.Provider>
    );
};

export const useSocket = () => {
    const socket = useContext(SocketContext);
    if (!socket) {
        throw new Error("useSocket must be used within a SocketProvider");
    }
    return socket;
};

export default SocketProvider;
