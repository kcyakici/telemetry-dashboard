"use client";

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
      <div>
        <label className="block text-sm mb-1">From</label>
        <input
          type="datetime-local"
          value={from.slice(0, 16)}
          onChange={(e) => setFrom(e.target.value)}
          className="border rounded p-2 bg-gray-700 text-white"
        />
      </div>
      <div>
        <label className="block text-sm mb-1">To</label>
        <input
          type="datetime-local"
          value={to.slice(0, 16)}
          onChange={(e) => setTo(e.target.value)}
          className="border rounded p-2 bg-gray-700 text-white"
        />
      </div>
      <div className="flex items-end">
        <button
          onClick={onApply}
          className="px-4 py-2 bg-blue-600 rounded hover:bg-blue-700"
        >
          Apply
        </button>
      </div>
    </div>
  );
}
