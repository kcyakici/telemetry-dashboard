export default function HomePage() {
  return (
    <div className="max-w-4xl mx-auto text-center space-y-8">
      <h1 className="text-4xl font-bold">Welcome to Telemetry Dashboard</h1>
      <p className="text-gray-300 text-lg">
        This project provides a simple way to ingest, store, and visualize
        vehicle telemetry data using TimescaleDB, Go, Nextjs and Recharts.
      </p>

      <div className="grid md:grid-cols-3 gap-6 mt-8">
        <div className="bg-gray-800 p-6 rounded-xl shadow">
          <h2 className="text-xl font-semibold">📈 Trends</h2>
          <p className="text-gray-400 mt-2">
            View time-series charts of telemetry metrics like speed, temperature
            and power demand.
          </p>
        </div>

        <div className="bg-gray-800 p-6 rounded-xl shadow">
          <h2 className="text-xl font-semibold">📊 Distribution</h2>
          <p className="text-gray-400 mt-2">
            Explore value distributions for telemetry metrics with histograms
            and ranges.
          </p>
        </div>

        <div className="bg-gray-800 p-6 rounded-xl shadow">
          <h2 className="text-xl font-semibold">🚍 KPIs</h2>
          <p className="text-gray-400 mt-2">
            See aggregated statistics such as average speed, passenger count,
            and power usage.
          </p>
        </div>
      </div>

      <footer className="text-gray-500 mt-12 text-sm">
        Used{" "}
        <a href="https://www.research-collection.ethz.ch/entities/researchdata/61ac2f6e-2ca9-4229-8242-aed3b0c0d47c">
          ZTBus
        </a>{" "}
        dataset
      </footer>
    </div>
  );
}
