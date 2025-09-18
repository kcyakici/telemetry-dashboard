"use client";
import { TrendPoint } from "@/types";
import { useCallback, useEffect, useRef, useState } from "react";

export type LiveMessage =
  | { type: "point"; timestamp: string; value: number }
  | { type: "error"; error: string }
  | { type: "connected"; vehicle: string; metric: string };

export function useLiveTrend(
  vehicle: string,
  metric: string,
  url = "ws://localhost:8080/live-trend"
) {
  const [points, setPoints] = useState<TrendPoint[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [isConnected, setIsConnected] = useState(false);
  const wsRef = useRef<WebSocket | null>(null);

  const connect = useCallback(() => {
    // Don't create new connection if one exists and is connecting/open
    if (
      wsRef.current?.readyState === WebSocket.CONNECTING ||
      wsRef.current?.readyState === WebSocket.OPEN
    ) {
      return;
    }

    const wsUrl = new URL(url);
    wsUrl.searchParams.set("vehicle_id", vehicle);
    wsUrl.searchParams.set("metric", metric);

    const ws = new WebSocket(wsUrl.toString());
    wsRef.current = ws;

    // Reset state when starting new connection
    setPoints([]);
    setError(null);
    setIsConnected(false);

    ws.onopen = () => {
      setIsConnected(true);
      setError(null);
    };

    ws.onmessage = (evt) => {
      try {
        const msg: LiveMessage = JSON.parse(evt.data);

        switch (msg.type) {
          case "connected":
            setIsConnected(true);
            setError(null);
            break;
          case "point":
            setPoints((prev) => [
              ...prev.slice(-49), // Keep last 50 points
              { timestamp: msg.timestamp, value: msg.value },
            ]);
            break;
          case "error":
            setError(msg.error);
            break;
        }
      } catch (err) {
        console.error("Failed to parse WebSocket message:", err);
        setError("Failed to parse server message");
      }
    };

    ws.onerror = () => {
      setError("WebSocket connection error");
      setIsConnected(false);
    };

    ws.onclose = (event) => {
      setIsConnected(false);

      // Only set error if it wasn't a clean close
      if (event.code !== 1000) {
        setError(`Connection closed: ${event.reason || "Unknown error"}`);
      }

      console.log("WebSocket closed:", event.code, event.reason);
    };
  }, [vehicle, metric, url]);

  const disconnect = useCallback(() => {
    if (wsRef.current) {
      wsRef.current.onclose = null; // Prevent onclose from triggering
      wsRef.current.close(1000, "Component unmounted");
      wsRef.current = null;
    }
    setIsConnected(false);
  }, []);

  useEffect(() => {
    connect();
    return disconnect;
  }, [connect, disconnect]);

  return {
    points,
    error,
    isConnected,
    reconnect: connect,
  };
}
