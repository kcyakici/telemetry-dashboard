"use client";

import { useEffect, useState } from "react";
import TrendChart from "../../components/TrendChart";
import { TrendPoint } from "../../types";

export default function Home() {
  const [data, setData] = useState<TrendPoint[]>([]);
  const [vehicle, setVehicle] = useState("B183");
  const [metric, setMetric] = useState("temp");
  const [from, setFrom] = useState("2019-06-24T03:16:13Z");
  const [to, setTo] = useState("2019-06-24T03:20:00Z");

  const loadTrend = async () => {
    const url = new URL("http://localhost:8080/trend");
    url.searchParams.append("vehicle_id", vehicle);
    url.searchParams.append("metric", metric);
    url.searchParams.append("start", from);
    url.searchParams.append("end", to);
    console.log("Here");

    const res = await fetch(url.toString());
    console.log("There");
    const json: TrendPoint[] = await res.json();
    console.log(json);
    setData(json);
  };

  useEffect(() => {
    loadTrend();
  }, []); // load initial data once

  return (
    <div className="p-6 space-y-6">
      <h1 className="text-2xl font-bold">Telemetry Trend Dashboard</h1>

      {/* Filter Controls */}
      <div className="flex flex-wrap gap-4 bg-gray-800 p-4 rounded-xl text-white">
        <div>
          <label className="block text-sm font-medium mb-1">Vehicle</label>
          <select
            value={vehicle}
            onChange={(e) => setVehicle(e.target.value)}
            className="border rounded p-2 bg-gray-700 text-white"
          >
            <option value="B183">B183</option>
            <option value="B208">B208</option>
          </select>
        </div>

        <div>
          <label className="block text-sm font-medium mb-1">Metric</label>
          <select
            value={metric}
            onChange={(e) => setMetric(e.target.value)}
            className="border rounded p-2 bg-gray-700 text-white"
          >
            <option value="temp">Temperature</option>
            <option value="speed">Speed</option>
            <option value="power">Power Demand</option>
          </select>
        </div>

        <div>
          <label className="block text-sm font-medium mb-1">From</label>
          <input
            type="datetime-local"
            value={from.slice(0, 16)}
            onChange={(e) => setFrom(e.target.value)}
            className="border rounded p-2 bg-gray-700 text-white"
          />
        </div>

        <div>
          <label className="block text-sm font-medium mb-1">To</label>
          <input
            type="datetime-local"
            value={to.slice(0, 16)}
            onChange={(e) => setTo(e.target.value)}
            className="border rounded p-2 bg-gray-700 text-white"
          />
        </div>

        <div className="flex items-end">
          <button
            onClick={loadTrend}
            className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
          >
            Apply
          </button>
        </div>
      </div>

      {/* Chart */}
      <div className="bg-white shadow p-4 rounded-xl">
        {data?.length > 0 ? (
          <TrendChart data={data} metric={metric} />
        ) : (
          <p className="text-gray-500">No data available.</p>
        )}
      </div>
    </div>
  );
}
