"use client";

import FilterBarBase from "./FilterBar";

type FilterBarWithMetricProps = {
  vehicle: string;
  setVehicle: (v: string) => void;
  from: string;
  setFrom: (f: string) => void;
  to: string;
  setTo: (t: string) => void;
  metric: string;
  setMetric: (m: string) => void;
  onApply: () => void;
};

export default function FilterBarWithMetric({
  vehicle,
  setVehicle,
  from,
  setFrom,
  to,
  setTo,
  metric,
  setMetric,
  onApply,
}: FilterBarWithMetricProps) {
  return (
    <div className="space-y-4">
      <FilterBarBase
        vehicle={vehicle}
        setVehicle={setVehicle}
        from={from}
        setFrom={setFrom}
        to={to}
        setTo={setTo}
        onApply={onApply}
      />
      <div className="bg-gray-800 p-4 rounded-xl">
        <label className="block text-sm mb-1">Metric</label>
        <select
          value={metric}
          onChange={(e) => setMetric(e.target.value)}
          className="border rounded p-2 bg-gray-700 text-white"
        >
          <option value="speed">Speed</option>
          <option value="temp">Temperature</option>
          <option value="power">Power Demand</option>
          <option value="traction">Traction Force</option>
        </select>
      </div>
    </div>
  );
}
