"use client";

import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from "recharts";

type BarChartProps = {
  data: { range: string; count: number }[];
  metric: string;
  vehicle: string;
};

export default function BarChartComponent({
  data,
  metric,
  vehicle,
}: BarChartProps) {
  return (
    <div className="bg-gray-800 p-6 rounded-xl shadow">
      <h2 className="text-lg mb-4">
        {metric.toUpperCase()} Distribution ({vehicle})
      </h2>
      <ResponsiveContainer width="100%" height={400}>
        <BarChart data={data}>
          <CartesianGrid strokeDasharray="3 3" stroke="#374151" />
          <XAxis dataKey="range" tick={{ fill: "white" }} />
          <YAxis tick={{ fill: "white" }} />
          <Tooltip
            contentStyle={{ backgroundColor: "#1f2937", border: "none" }}
            labelStyle={{ color: "#e5e7eb" }}
          />
          <Bar dataKey="count" fill="#38bdf8" />
        </BarChart>
      </ResponsiveContainer>
    </div>
  );
}
