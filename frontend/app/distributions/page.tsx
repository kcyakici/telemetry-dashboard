"use client";

import { useEffect, useState } from "react";
import DistributionChart from "../../components/charts/DistributionChart";
import FilterBarWithMetric, {
  Filters,
} from "../../components/filters/FilterBarWithMetric";
import { metricsConfig } from "@/config/metrics";

const initialFilters: Filters = {
  vehicle: "B183",
  metric: metricsConfig[0].value,
  from: "2019-06-24T03:16:00Z",
  to: "2019-06-24T03:20:00Z",
};

export default function DistributionPage() {
  const [data, setData] = useState<DistributionResponse | null>(null);
  const [appliedFilters, setAppliedFilters] = useState<Filters>(initialFilters);
  const [draftBin, setDraftBin] = useState(10);
  const [appliedBin, setAppliedBin] = useState(10);

  const loadDistribution = async (filters: Filters) => {
    const url = new URL("http://localhost:8080/distribution");
    url.searchParams.append("vehicle_id", filters.vehicle);
    url.searchParams.append("metric", filters.metric);
    url.searchParams.append("start", filters.from);
    url.searchParams.append("end", filters.to);
    url.searchParams.append("bins", String(appliedBin));

    const res = await fetch(url.toString());
    const json: DistributionResponse = await res.json();
    setAppliedFilters(filters);
    setAppliedBin(draftBin);
    setData(json);
  };

  useEffect(() => {
    const fetchData = async (filters: Filters) => {
      const url = new URL("http://localhost:8080/distribution");
      url.searchParams.append("vehicle_id", filters.vehicle);
      url.searchParams.append("metric", filters.metric);
      url.searchParams.append("start", filters.from);
      url.searchParams.append("end", filters.to);
      url.searchParams.append("bins", String(10));

      const res = await fetch(url.toString());
      const json: DistributionResponse = await res.json();
      setData(json);
      setAppliedFilters(filters);
    };

    fetchData(initialFilters);
  }, []);

  const chartData =
    data?.buckets.map((b) => ({
      range: `${b.range_min.toFixed(1)} - ${b.range_max.toFixed(1)}`,
      count: b.count,
    })) || [];

  const minBinCount = 5;
  const maxBinCount = 20;

  return (
    <div className="p-6 space-y-6 text-white bg-gray-900 min-h-screen">
      <h1 className="text-2xl font-bold">Distribution</h1>

      <FilterBarWithMetric
        initialFilters={initialFilters}
        onApply={loadDistribution}
      />

      <div>
        <label className="block text-sm mb-1">Bins</label>
        <input
          type="number"
          name="bins"
          step="1"
          min={minBinCount}
          max={maxBinCount}
          value={draftBin}
          onChange={(e) => setDraftBin(Number(e.target.value))}
          className="border rounded p-2 bg-gray-700 text-white"
        />
      </div>

      {chartData.length > 0 ? (
        <DistributionChart
          header={`${appliedFilters.metric.toUpperCase()} Distribution - 
          Vehicle:  ${appliedFilters.vehicle} - FROM ${
            appliedFilters.from
          } TO ${appliedFilters.to} - Bins: ${appliedBin}`}
          data={chartData}
        />
      ) : (
        <p className="text-gray-400">No distribution data available.</p>
      )}
    </div>
  );
}

type Bucket = {
  bucket: number;
  count: number;
  range_min: number;
  range_max: number;
};

type DistributionResponse = {
  metric: string;
  vehicle: string;
  bins: number;
  min: number;
  max: number;
  from?: string;
  to?: string;
  buckets: Bucket[];
};
