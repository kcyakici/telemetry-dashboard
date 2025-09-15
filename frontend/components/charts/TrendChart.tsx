"use client";

import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from "recharts";
import { TrendPoint } from "../../types";

type TrendChartProps = {
  data: TrendPoint[];
  metric: string;
};

export default function TrendChart({ data, metric }: TrendChartProps) {
  const chartData = data.map((d) => ({
    time: new Date(d.timestamp).toLocaleTimeString(),
    value: d.value,
  }));

  return (
    <ResponsiveContainer width="100%" height={350}>
      <LineChart data={chartData}>
        <CartesianGrid strokeDasharray="3 3" />
        <XAxis dataKey="time" />
        <YAxis
          label={{
            value: metric,
            angle: -90,
            position: "insideLeft",
            offset: 10,
            style: { textAnchor: "middle", fill: "#fff" },
          }}
        />
        <Tooltip />
        <Line
          type="monotone"
          dataKey="value"
          stroke="#3182ce"
          strokeWidth={2}
        />
      </LineChart>
    </ResponsiveContainer>
  );
}
