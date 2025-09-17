"use client";

import {
  CartesianGrid,
  Line,
  LineChart,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";
import { TrendPoint } from "../../types";

type TrendChartProps = {
  data: TrendPoint[];
  isAnimationActive: boolean;
  header: string;
};

export default function TrendChart({
  data,
  isAnimationActive,
  header,
}: TrendChartProps) {
  const chartData = data.map((d) => ({
    time: new Date(d.timestamp).toLocaleTimeString(),
    value: d.value,
  }));

  const values = chartData.map((d) => d.value);
  const minValue = Math.min(...values);
  const maxValue = Math.max(...values);
  const domain = calculateDomainRange(minValue, maxValue);

  return (
    <div className="bg-gray-800 p-6 rounded-xl shadow">
      <h2 className="text-lg mb-4">{header}</h2>
      <ResponsiveContainer width="100%" height={400}>
        <LineChart data={chartData}>
          <CartesianGrid strokeDasharray="3 3" stroke="#374151" />
          <XAxis dataKey="time" tick={{ fill: "white" }} />
          <YAxis
            domain={domain}
            allowDecimals={false}
            tick={{ fill: "white" }}
            // tickFormatter={(tick) => Math.round(tick).toString()} strange behavior while rendering: eliminates smoothness
          />
          <Tooltip
            contentStyle={{ backgroundColor: "#1f2937", border: "none" }}
            labelStyle={{ color: "#e5e7eb" }}
          />
          <Line
            type="monotone"
            dataKey="value"
            stroke="#38bdf8"
            strokeWidth={2}
            dot={false}
            isAnimationActive={isAnimationActive}
          />
        </LineChart>
      </ResponsiveContainer>
    </div>
  );
}

function calculateDomainRange(minValue: number, maxValue: number): number[] {
  if (minValue === maxValue) {
    return [minValue - 1, maxValue + 1];
  } else if (minValue < 0 && maxValue > 0) {
    return [minValue * 1.05, maxValue * 1.05];
  } else if (minValue < 0 && maxValue < 0) {
    return [minValue * 1.05, maxValue * 0.95];
  } else {
    return [minValue * 0.95, maxValue * 1.05];
  }
}
