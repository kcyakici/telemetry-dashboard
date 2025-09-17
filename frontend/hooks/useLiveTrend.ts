"use client";

import { TrendPoint } from "@/types";
import { useEffect, useRef, useState } from "react";

export type LiveMessage =
  | { type: "point"; timestamp: string; value: number }
  | { type: "error"; error: string };

export function useLiveTrend(
  vehicle: string,
  metric: string,
  url = "ws://localhost:8080/live-trend"
) {
  const [points, setPoints] = useState<TrendPoint[]>([]);
  const [error, setError] = useState<string | null>(null);
  const wsRef = useRef<WebSocket | null>(null);

  useEffect(() => {
    const cleanup = (reason: string) => {
      if (wsRef.current) {
        wsRef.current.onclose = null;
        wsRef.current.close(1000, reason);
        wsRef.current = null;
      }
    };

    cleanup("reconnect");

    const wsUrl = new URL(url);
    wsUrl.searchParams.append("vehicle_id", vehicle);
    wsUrl.searchParams.append("metric", metric);

    const ws = new WebSocket(wsUrl.toString());
    wsRef.current = ws;

    setPoints([]);
    setError(null);

    ws.onmessage = (evt) => {
      const msg: LiveMessage = JSON.parse(evt.data);
      if (msg.type === "point") {
        setPoints((prev) => [
          ...prev.slice(-49),
          { timestamp: msg.timestamp, value: msg.value },
        ]);
      } else if (msg.type === "error") {
        setError(msg.error);
      }
    };

    ws.onerror = () => setError("WebSocket connection error");
    ws.onclose = () => console.log("ws closed");

    return () => cleanup("component unmount / metric change");
  }, [vehicle, metric, url]);

  return { points, error };
}
