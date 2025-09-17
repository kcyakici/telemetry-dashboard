"use client";

import { useState } from "react";
import { useLiveTrend } from "@/hooks/useLiveTrend";
import TrendChart from "@/components/charts/TrendChart";
import { metricsConfig } from "@/config/metrics";

export default function LiveTrendPage() {
  const [vehicle, setVehicle] = useState("B183");
  const [metric, setMetric] = useState("speed");

  const { points, error } = useLiveTrend(vehicle, metric);

  return (
    <div className="p-6 space-y-6 text-white bg-gray-900 min-h-screen">
      <h1 className="text-2xl font-bold">Live Telemetry Trend</h1>

      {/* Filters */}
      <div className="flex gap-4 bg-gray-800 p-4 rounded-xl">
        <div>
          <label className="block text-sm mb-1">Vehicle</label>
          <select
            value={vehicle}
            onChange={(e) => setVehicle(e.target.value)}
            className="border rounded p-2 bg-gray-700 text-white"
          >
            <option value="B183">B183</option>
            <option value="B208">B208</option>
          </select>
        </div>

        <select
          value={metric}
          onChange={(e) => setMetric(e.target.value)}
          className="border rounded p-2 bg-gray-700 text-white"
        >
          {metricsConfig.map((m) => (
            <option key={m.value} value={m.value}>
              {m.label}
            </option>
          ))}
        </select>
      </div>

      {/* Error banner */}
      {error && (
        <div className="bg-red-700 p-3 rounded text-white">
          WebSocket error: {error}
        </div>
      )}

      {/* Live chart */}
      {points.length > 0 ? (
        <TrendChart
          data={points}
          header={`${vehicle} ${metric} Live Telemetry`}
          isAnimationActive={false}
        />
      ) : (
        <p className="text-gray-400">Waiting for live telemetry...</p>
      )}
    </div>
  );
}
