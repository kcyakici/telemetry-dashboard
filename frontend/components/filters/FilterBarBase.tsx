"use client";

import { useState } from "react";
import DateInput from "./DateInput";

export type FiltersBase = {
  vehicle: string;
  from: string;
  to: string;
};

type FilterBarBaseProps = {
  initialVehicle: string;
  initialFrom: string;
  initialTo: string;
  onApply: (filters: { vehicle: string; from: string; to: string }) => void;
};

export default function FilterBarBase({
  initialVehicle,
  initialFrom,
  initialTo,
  onApply,
}: FilterBarBaseProps) {
  const [vehicle, setVehicle] = useState(initialVehicle);
  const [from, setFrom] = useState(initialFrom);
  const [to, setTo] = useState(initialTo);

  const isInvalidRange = new Date(from) > new Date(to);

  return (
    <div className="flex flex-wrap gap-4 bg-gray-800 p-4 rounded-xl">
      <div>
        <label className="block text-sm mb-1">Vehicle</label>
        <select
          value={vehicle}
          onChange={(e) => setVehicle(e.target.value)}
          className="border rounded p-2 bg-gray-700 text-white"
        >
          <option value="B183">B183</option>
          <option value="B208">B208</option>
        </select>
      </div>
      <DateInput label="From" date={from} handleDateChange={setFrom} />
      <DateInput label="To" date={to} handleDateChange={setTo} />
      <div className="flex items-end">
        <button
          onClick={() => onApply({ vehicle, from, to })}
          disabled={isInvalidRange}
          className={`px-4 py-2 rounded ${
            isInvalidRange
              ? "bg-gray-500 cursor-not-allowed"
              : "bg-blue-600 hover:bg-blue-700"
          }`}
        >
          Apply
        </button>
        {isInvalidRange && (
          <p className="text-red-400 text-sm">
            Start date must be before end date
          </p>
        )}
      </div>
    </div>
  );
}
