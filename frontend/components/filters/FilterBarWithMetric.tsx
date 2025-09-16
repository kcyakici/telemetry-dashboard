"use client";

import { useState } from "react";
import FilterBarBase from "./FilterBarBase";
import { metricsConfig } from "@/config/metrics"; // <-- import config

export type Filters = {
  vehicle: string;
  metric: string;
  from: string;
  to: string;
};

type FilterBarWithMetricProps = {
  initialFilters: Filters;
  onApply: (filters: Filters) => void;
};

export default function FilterBarWithMetric({
  initialFilters,
  onApply,
}: FilterBarWithMetricProps) {
  const [metric, setMetric] = useState(initialFilters.metric);

  return (
    <div className="space-y-4">
      <FilterBarBase
        initialVehicle={initialFilters.vehicle}
        initialFrom={initialFilters.from}
        initialTo={initialFilters.to}
        onApply={(baseFilters) => onApply({ ...baseFilters, metric })}
      />
      <div className="bg-gray-800 p-4 rounded-xl">
        <label className="block text-sm mb-1">Metric</label>
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
    </div>
  );
}
