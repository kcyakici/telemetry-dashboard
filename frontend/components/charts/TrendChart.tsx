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
  metric: string;
};

export default function TrendChart({ data, metric }: TrendChartProps) {
  const chartData = data.map((d) => ({
    time: new Date(d.timestamp).toLocaleTimeString(),
    value: d.value,
  }));

  const values = chartData.map((d) => d.value);
  const minValue = Math.min(...values);
  const maxValue = Math.max(...values);
  const domain =
    minValue === maxValue
      ? [minValue - 1, maxValue + 1]
      : [minValue * 0.95, maxValue * 1.05];

  return (
    <ResponsiveContainer width="100%" height={350}>
      <LineChart data={chartData}>
        <CartesianGrid strokeDasharray="3 3" />
        <XAxis dataKey="time" />
        <YAxis
          allowDecimals={false}
          tickFormatter={(tick) => Math.round(tick).toString()}
          domain={domain}
          label={{
            value: metric,
            angle: -90,
            position: "insideLeft",
            offset: 10,
            style: { textAnchor: "middle", fill: "#000" },
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
