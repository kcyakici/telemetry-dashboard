"use client";

import DateInput from "./DateInput";

type FilterBarBaseProps = {
  vehicle: string;
  setVehicle: (v: string) => void;
  from: string;
  setFrom: (f: string) => void;
  to: string;
  setTo: (t: string) => void;
  onApply: () => void;
};

export default function FilterBarBase({
  vehicle,
  setVehicle,
  from,
  setFrom,
  to,
  setTo,
  onApply,
}: FilterBarBaseProps) {
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
      <DateInput
        label="From"
        date={from}
        handleDateChange={setFrom}
      ></DateInput>
      <DateInput label="To" date={to} handleDateChange={setTo}></DateInput>
      <div className="flex items-end">
        <button
          onClick={onApply}
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
