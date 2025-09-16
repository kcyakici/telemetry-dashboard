import { Cell, Pie, PieChart, ResponsiveContainer, Tooltip } from "recharts";

export type KpiProps = {
  avg_speed: number | null;
  max_temp: number | null;
  total_power: number | null;
  avg_brake_pressure: number | null;
  door_open_ratio: number | null;
  header: string;
};

export default function KpiChart({
  avg_speed,
  max_temp,
  total_power,
  avg_brake_pressure,
  door_open_ratio,
  header,
}: KpiProps) {
  return (
    <div className="grid grid-cols-2 md:grid-cols-3 gap-6">
      {/* Header across all columns */}
      <div className="col-span-full">
        <h2 className="text-2xl font-bold mb-6">{header.toUpperCase()}</h2>
      </div>

      <div className="bg-gray-800 p-4 rounded-xl shadow">
        <h2 className="text-lg">Avg Speed</h2>
        <p className="text-2xl font-bold">
          {avg_speed !== null ? `${avg_speed.toFixed(2)} km/h` : "N/A"}
        </p>
      </div>

      <div className="bg-gray-800 p-4 rounded-xl shadow">
        <h2 className="text-lg">Max Temp</h2>
        <p className="text-2xl font-bold">
          {max_temp !== null ? `${(max_temp - 273.15).toFixed(1)} °C` : "N/A"}
        </p>
      </div>

      <div className="bg-gray-800 p-4 rounded-xl shadow">
        <h2 className="text-lg">Total Power</h2>
        <p className="text-2xl font-bold">
          {total_power !== null ? `${total_power.toFixed(1)} W` : "N/A"}
        </p>
      </div>

      <div className="bg-gray-800 p-4 rounded-xl shadow">
        <h2 className="text-lg">Avg Brake Pressure</h2>
        <p className="text-2xl font-bold">
          {avg_brake_pressure !== null
            ? `${avg_brake_pressure.toFixed(1)} Pa`
            : "N/A"}
        </p>
      </div>

      <div className="bg-gray-800 p-4 rounded-xl shadow col-span-2">
        <h2 className="text-lg mb-2">Door Open Ratio</h2>
        {door_open_ratio !== null ? (
          <ResponsiveContainer width="100%" height={250}>
            <PieChart>
              <Pie
                data={[
                  { name: "Open", value: door_open_ratio * 100 },
                  { name: "Closed", value: 100 - door_open_ratio * 100 },
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
        ) : (
          <p className="text-gray-400">N/A</p>
        )}
      </div>
    </div>
  );
}
