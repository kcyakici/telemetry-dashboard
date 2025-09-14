"use client";

import { useState, useEffect } from "react";
import FilterBar from "../../components/filters/FilterBar";
import KpiChart, { KpiProps } from "@/components/charts/KpiChart";

export default function KpisPage() {
  const [vehicle, setVehicle] = useState("B183");
  const [from, setFrom] = useState("2019-06-24T03:16:13Z");
  const [to, setTo] = useState("2019-06-24T03:20:00Z");
  const [kpis, setKpis] = useState<KpiProps | null>(null);

  const loadKpis = async () => {
    const url = new URL("http://localhost:8080/kpis");
    url.searchParams.append("vehicle_id", vehicle);
    url.searchParams.append("start", from);
    url.searchParams.append("end", to);

    const res = await fetch(url.toString());
    const json = await res.json();
    setKpis(json);
  };

  useEffect(() => {
    loadKpis();
  }, []);

  return (
    <div className="p-6 space-y-6 text-white bg-gray-900 min-h-screen">
      <h1 className="text-2xl font-bold">Vehicle KPIs</h1>

      <FilterBar
        vehicle={vehicle}
        setVehicle={setVehicle}
        from={from}
        setFrom={setFrom}
        to={to}
        setTo={setTo}
        onApply={loadKpis}
      />

      {kpis && <KpiChart {...kpis}></KpiChart>}
    </div>
  );
}
