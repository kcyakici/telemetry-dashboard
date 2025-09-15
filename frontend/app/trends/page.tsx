"use client";

import { useCallback, useEffect, useState } from "react";
import TrendChart from "../../components/charts/TrendChart";
import { TrendPoint } from "@/types";
import FilterBarWithMetric from "../../components/filters/FilterBarWithMetric";

export default function TrendsPage() {
  const [data, setData] = useState<TrendPoint[]>([]);
  const [vehicle, setVehicle] = useState("B183");
  const [metric, setMetric] = useState("temp");
  const [from, setFrom] = useState("2019-06-24T03:16:00Z");
  const [to, setTo] = useState("2019-06-24T03:20:00Z");

  const loadTrend = async () => {
    const url = new URL("http://localhost:8080/trend");
    url.searchParams.append("vehicle_id", vehicle);
    url.searchParams.append("metric", metric);
    url.searchParams.append("start", from);
    url.searchParams.append("end", to);

    const res = await fetch(url.toString());
    const json: TrendPoint[] = await res.json();
    setData(json);
  };

  useEffect(() => {
    loadTrend();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return (
    <div className="p-6 space-y-6 text-white bg-gray-900 min-h-screen">
      <h1 className="text-2xl font-bold">Telemetry Trend Dashboard</h1>

      <FilterBarWithMetric
        vehicle={vehicle}
        setVehicle={setVehicle}
        from={from}
        setFrom={setFrom}
        to={to}
        setTo={setTo}
        metric={metric}
        setMetric={setMetric}
        onApply={loadTrend}
      />

      <div className="bg-white shadow p-4 rounded-xl">
        {data?.length > 0 ? (
          <TrendChart data={data} metric={metric} />
        ) : (
          <p className="text-gray-500">No data available.</p>
        )}
      </div>
    </div>
  );
}
