"use client";

import { useEffect, useState } from "react";
import DistributionChart from "../../components/charts/DistributionChart";
import FilterBarWithMetric from "../../components/filters/FilterBarWithMetric";

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

export default function DistributionPage() {
  const [vehicle, setVehicle] = useState("B183");
  const [metric, setMetric] = useState("temp");
  const [from, setFrom] = useState("2019-06-24T03:16:00Z");
  const [to, setTo] = useState("2019-06-24T03:20:00Z");
  const [bins, setBins] = useState(10);
  const [data, setData] = useState<DistributionResponse | null>(null);

  const loadDistribution = async () => {
    const url = new URL("http://localhost:8080/distribution");
    url.searchParams.append("vehicle_id", vehicle);
    url.searchParams.append("metric", metric);
    url.searchParams.append("from", from);
    url.searchParams.append("to", to);
    url.searchParams.append("bins", String(bins));

    const res = await fetch(url.toString());
    const json = await res.json();
    setData(json);
  }; // all dependencies;

  useEffect(() => {
    loadDistribution();
    // eslint-disable-next-line react-hooks/exhaustive-deps
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
        vehicle={vehicle}
        setVehicle={setVehicle}
        from={from}
        setFrom={setFrom}
        to={to}
        setTo={setTo}
        metric={metric}
        setMetric={setMetric}
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
          value={bins}
          onChange={(e) => setBins(Number(e.target.value))}
          className="border rounded p-2 bg-gray-700 text-white"
        />
      </div>

      {chartData.length > 0 ? (
        <DistributionChart data={chartData} metric={metric} vehicle={vehicle} />
      ) : (
        <p className="text-gray-400">No distribution data available.</p>
      )}
    </div>
  );
}
