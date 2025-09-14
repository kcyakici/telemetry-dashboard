"use client";

import Link from "next/link";

export default function Navbar() {
  return (
    <nav className="bg-gray-800 text-white px-6 py-3 shadow-md">
      <div className="flex items-center justify-between max-w-7xl mx-auto">
        <div className="text-xl font-bold">Telemetry Dashboard</div>

        <div className="flex space-x-6">
          <Link href="/" className="hover:text-blue-400 transition">
            Home
          </Link>
          <Link href="/ingestion" className="hover:text-blue-400 transition">
            Upload
          </Link>
          <Link href="/trends" className="hover:text-blue-400 transition">
            Trends
          </Link>
          <Link href="/distribution" className="hover:text-blue-400 transition">
            Distribution
          </Link>
          <Link href="/kpis" className="hover:text-blue-400 transition">
            KPIs
          </Link>
        </div>
      </div>
    </nav>
  );
}
