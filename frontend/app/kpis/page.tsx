"use client";

import KpiChart, { KpiProps } from "@/components/charts/KpiChart";
import { useEffect, useState } from "react";
import FilterBarBase, { FiltersBase } from "@/components/filters/FilterBarBase";

const initialFilters = {
  vehicle: "B183",
  from: "2019-06-24T03:16:00Z",
  to: "2019-06-24T03:20:00Z",
};

type KpiResponse = {
  avg_speed: number | null;
  max_temp: number | null;
  total_power: number | null;
  avg_brake_pressure: number | null;
  door_open_ratio: number | null;
};

export default function KpisPage() {
  const [kpis, setKpis] = useState<KpiResponse | null>(null);
  const [appliedFilters, setAppliedFilters] =
    useState<FiltersBase>(initialFilters);

  const loadKpis = async (filters: FiltersBase) => {
    const url = new URL("http://localhost:8080/kpis");
    url.searchParams.append("vehicle_id", filters.vehicle);
    url.searchParams.append("start", filters.from);
    url.searchParams.append("end", filters.to);

    const res = await fetch(url.toString());
    const json = await res.json();
    setAppliedFilters(filters);
    setKpis(json);
  };

  // initial load
  useEffect(() => {
    loadKpis(initialFilters);
  }, []);

  const hasData =
    kpis &&
    (kpis.avg_speed !== null ||
      kpis.max_temp !== null ||
      kpis.total_power !== null ||
      kpis.avg_brake_pressure !== null ||
      kpis.door_open_ratio !== null);

  return (
    <div className="p-6 space-y-6 text-white bg-gray-900 min-h-screen">
      <h1 className="text-2xl font-bold">Vehicle KPIs</h1>

      <FilterBarBase
        initialVehicle={initialFilters.vehicle}
        initialFrom={initialFilters.from}
        initialTo={initialFilters.to}
        onApply={loadKpis}
      />

      {hasData ? (
        <KpiChart
          {...kpis!}
          header={`${appliedFilters.vehicle} from ${appliedFilters.from} to ${appliedFilters.to}`}
        />
      ) : (
        <p className="text-gray-400">No data available.</p>
      )}
    </div>
  );
}
