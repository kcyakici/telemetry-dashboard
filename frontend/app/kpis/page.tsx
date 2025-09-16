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
  avg_speed: number;
  max_temp: number;
  total_power: number;
  avg_brake_pressure: number;
  door_open_ratio: number;
};

export default function KpisPage() {
  const [kpis, setKpis] = useState<KpiResponse | null>(null);
  const [appliedFilters, setAppliedFilters] =
    useState<FiltersBase>(initialFilters);

  const loadKpis = async (filters: {
    vehicle: string;
    from: string;
    to: string;
  }) => {
    const url = new URL("http://localhost:8080/kpis");
    url.searchParams.append("vehicle_id", filters.vehicle);
    url.searchParams.append("start", filters.from);
    url.searchParams.append("end", filters.to);

    const res = await fetch(url.toString());
    const json = await res.json();
    setAppliedFilters(filters);
    setKpis(json);
  };

  useEffect(() => {
    const fetchData = async (filters: {
      vehicle: string;
      from: string;
      to: string;
    }) => {
      const url = new URL("http://localhost:8080/kpis");
      url.searchParams.append("vehicle_id", filters.vehicle);
      url.searchParams.append("start", filters.from);
      url.searchParams.append("end", filters.to);

      const res = await fetch(url.toString());
      const json = await res.json();
      setAppliedFilters(filters);
      setKpis(json);
    };

    fetchData(initialFilters);
  }, []);

  return (
    <div className="p-6 space-y-6 text-white bg-gray-900 min-h-screen">
      <h1 className="text-2xl font-bold">Vehicle KPIs</h1>

      <FilterBarBase
        initialVehicle={initialFilters.vehicle}
        initialFrom={initialFilters.from}
        initialTo={initialFilters.to}
        onApply={loadKpis}
      />

      {kpis && (
        <KpiChart
          {...kpis}
          header={`${appliedFilters.vehicle} from ${appliedFilters.from} to ${appliedFilters.to}`}
        />
      )}
    </div>
  );
}
