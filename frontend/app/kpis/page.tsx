"use client";

import { useState, useEffect } from "react";
import { PieChart, Pie, Cell, ResponsiveContainer, Tooltip } from "recharts";

type KpiData = {
  avg_speed: number;
  max_temp: number;
  total_power: number;
  avg_passengers: number;
  door_open_ratio: number;
};

export default function KpisPage() {
  const [vehicle, setVehicle] = useState("bus_1");
  const [from, setFrom] = useState("2019-06-24T03:16:13Z");
  const [to, setTo] = useState("2019-06-24T03:20:00Z");
  const [kpis, setKpis] = useState<KpiData | null>(null);

  const loadKpis = async () => {
    const url = new URL("http://localhost:8080/kpis");
    url.searchParams.append("vehicle_id", vehicle);
    url.searchParams.append("start", from);
    url.searchParams.append("end", to);

    const res = await fetch(url.toString());
    const json = await res.json();
    setKpis(json);
  };

  useEffect(() => {
    loadKpis();
  }, []);

  return (
    <div className="p-6 space-y-6 text-white bg-gray-900 min-h-screen">
      <h1 className="text-2xl font-bold">Vehicle KPIs</h1>

      {/* Filters */}
      <div className="flex gap-4 bg-gray-800 p-4 rounded-xl">
        <div>
          <label className="block text-sm mb-1">Vehicle</label>
          <select
            value={vehicle}
            onChange={(e) => setVehicle(e.target.value)}
            className="border rounded p-2 bg-gray-700 text-white"
          >
            <option value="bus_1">Bus 1</option>
            <option value="bus_2">Bus 2</option>
            <option value="B183">B183</option>
          </select>
        </div>
        <div>
          <label className="block text-sm mb-1">From</label>
          <input
            type="datetime-local"
            value={from.slice(0, 16)}
            onChange={(e) => setFrom(e.target.value)}
            className="border rounded p-2 bg-gray-700 text-white"
          />
        </div>
        <div>
          <label className="block text-sm mb-1">To</label>
          <input
            type="datetime-local"
            value={to.slice(0, 16)}
            onChange={(e) => setTo(e.target.value)}
            className="border rounded p-2 bg-gray-700 text-white"
          />
        </div>
        <div className="flex items-end">
          <button
            onClick={loadKpis}
            className="px-4 py-2 bg-blue-600 rounded hover:bg-blue-700"
          >
            Apply
          </button>
        </div>
      </div>

      {kpis && (
        <div className="grid grid-cols-2 md:grid-cols-3 gap-6">
          {/* Simple KPI cards */}
          <div className="bg-gray-800 p-4 rounded-xl shadow">
            <h2 className="text-lg">Avg Speed</h2>
            <p className="text-2xl font-bold">
              {kpis.avg_speed?.toFixed(2)} km/h
            </p>
          </div>
          <div className="bg-gray-800 p-4 rounded-xl shadow">
            <h2 className="text-lg">Max Temp</h2>
            <p className="text-2xl font-bold">
              {(kpis.max_temp - 273.15).toFixed(1)} °C
            </p>
          </div>
          <div className="bg-gray-800 p-4 rounded-xl shadow">
            <h2 className="text-lg">Total Power</h2>
            <p className="text-2xl font-bold">
              {kpis.total_power?.toFixed(1)} W
            </p>
          </div>
          <div className="bg-gray-800 p-4 rounded-xl shadow">
            <h2 className="text-lg">Avg Passengers</h2>
            <p className="text-2xl font-bold">
              {kpis.avg_passengers?.toFixed(1)}
            </p>
          </div>

          {/* Door open ratio chart */}
          <div className="bg-gray-800 p-4 rounded-xl shadow col-span-2">
            <h2 className="text-lg mb-2">Door Open Ratio</h2>
            <ResponsiveContainer width="100%" height={250}>
              <PieChart>
                <Pie
                  data={[
                    { name: "Open", value: kpis.door_open_ratio * 100 },
                    { name: "Closed", value: 100 - kpis.door_open_ratio * 100 },
                  ]}
                  dataKey="value"
                  cx="50%"
                  cy="50%"
                  outerRadius={80}
                  label
                >
                  <Cell fill="#38bdf8" />
                  <Cell fill="#64748b" />
                </Pie>
                <Tooltip />
              </PieChart>
            </ResponsiveContainer>
          </div>
        </div>
      )}
    </div>
  );
}
