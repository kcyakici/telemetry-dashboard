"use client";

import { useEffect, useState } from "react";
import { TrendPoint } from "@/types";
import TrendChart from "../../components/charts/TrendChart";
import FilterBarWithMetric, {
  Filters,
} from "../../components/filters/FilterBarWithMetric";

const initialFilters: Filters = {
  vehicle: "B183",
  metric: "temp",
  from: "2019-06-24T03:16:00Z",
  to: "2019-06-24T03:20:00Z",
};

export default function TrendsPage() {
  const [data, setData] = useState<TrendPoint[]>([]);
  const [appliedFilters, setAppliedFilters] = useState<Filters>(initialFilters);

  const loadTrend = async (filters: Filters) => {
    const url = new URL("http://localhost:8080/trend");
    url.searchParams.append("vehicle_id", filters.vehicle);
    url.searchParams.append("metric", filters.metric);
    url.searchParams.append("start", filters.from);
    url.searchParams.append("end", filters.to);

    const res = await fetch(url.toString());
    const json: TrendPoint[] = await res.json();
    setData(json);
    setAppliedFilters(filters);
  };

  // initial load with defaults
  useEffect(() => {
    const fetchData = async (filters: Filters) => {
      const url = new URL("http://localhost:8080/trend");
      url.searchParams.append("vehicle_id", filters.vehicle);
      url.searchParams.append("metric", filters.metric);
      url.searchParams.append("start", filters.from);
      url.searchParams.append("end", filters.to);

      const res = await fetch(url.toString());
      const json: TrendPoint[] = await res.json();
      setData(json);
      setAppliedFilters(filters);
    };

    fetchData(initialFilters);
  }, []);

  return (
    <div className="p-6 space-y-6 text-white bg-gray-900 min-h-screen">
      <h1 className="text-2xl font-bold">Telemetry Trend Dashboard</h1>

      <FilterBarWithMetric
        initialFilters={initialFilters}
        onApply={loadTrend}
      />

      {data?.length > 0 ? (
        <TrendChart
          data={data}
          isAnimationActive={true}
          header={`${appliedFilters.vehicle} ${appliedFilters.metric} Trend - From ${appliedFilters.from} To ${appliedFilters.to}`}
        />
      ) : (
        <p className="text-gray-500">No data available.</p>
      )}
    </div>
  );
}
