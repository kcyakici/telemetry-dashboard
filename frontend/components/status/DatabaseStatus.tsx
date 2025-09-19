"use client";

import { useEffect, useState } from "react";

export default function DatabaseStatus() {
  const [rowCount, setRowCount] = useState<number | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchStats = async () => {
      try {
        const res = await fetch("http://localhost:8080/stats");
        if (!res.ok) throw new Error("Failed to fetch stats");
        const data = await res.json();
        setRowCount(data.row_count);
      } catch (err) {
        setError("Could not fetch database stats");
      }
    };

    fetchStats();
  }, []);

  if (error) {
    return <div className="bg-red-800 p-4 rounded-xl text-white">{error}</div>;
  }

  if (rowCount === null) {
    return (
      <div className="bg-gray-800 p-4 rounded-xl text-white">
        Loading database status...
      </div>
    );
  }

  if (rowCount === 0) {
    return (
      <div className="bg-blue-800 p-4 rounded-xl text-white">
        You have <strong>0 rows</strong> in the database.
        <br />
        Start by uploading a ZTBus dataset sample to use the dashboard!
      </div>
    );
  }

  return (
    <div className="bg-green-800 p-4 rounded-xl text-white">
      You currently have <strong>{rowCount}</strong> rows in the database.
      <br />
      Continue exploring your data!
    </div>
  );
}
