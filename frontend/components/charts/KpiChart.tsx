import { Cell, Pie, PieChart, ResponsiveContainer, Tooltip } from "recharts";

export type KpiProps = {
  avg_speed: number;
  max_temp: number;
  total_power: number;
  avg_brake_pressure: number;
  door_open_ratio: number;
};

export default function KpiChart(props: KpiProps) {
  return (
    <div className="grid grid-cols-2 md:grid-cols-3 gap-6">
      <div className="bg-gray-800 p-4 rounded-xl shadow">
        <h2 className="text-lg">Avg Speed</h2>
        <p className="text-2xl font-bold">{props.avg_speed?.toFixed(2)} km/h</p>
      </div>
      <div className="bg-gray-800 p-4 rounded-xl shadow">
        <h2 className="text-lg">Max Temp</h2>
        <p className="text-2xl font-bold">
          {(props.max_temp - 273.15).toFixed(1)} °C
        </p>
      </div>
      <div className="bg-gray-800 p-4 rounded-xl shadow">
        <h2 className="text-lg">Total Power</h2>
        <p className="text-2xl font-bold">{props.total_power?.toFixed(1)} W</p>
      </div>
      <div className="bg-gray-800 p-4 rounded-xl shadow">
        <h2 className="text-lg">Avg Brake Pressure</h2>
        <p className="text-2xl font-bold">
          {props.avg_brake_pressure?.toFixed(1)} Pa
        </p>
      </div>

      <div className="bg-gray-800 p-4 rounded-xl shadow col-span-2">
        <h2 className="text-lg mb-2">Door Open Ratio</h2>
        <ResponsiveContainer width="100%" height={250}>
          <PieChart>
            <Pie
              data={[
                { name: "Open", value: props.door_open_ratio * 100 },
                { name: "Closed", value: 100 - props.door_open_ratio * 100 },
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
  );
}
