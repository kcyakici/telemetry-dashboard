# Telemetry Dashboard

A full-stack web application for exploring and visualizing vehicle telemetry data.  
The project demonstrates backend data handling with **Go + Gin + TimescaleDB** and a responsive frontend with **Next.js + Recharts**.

It was designed to showcase:

- Trend visualization of telemetry metrics
- KPI dashboards (aggregates like speed, temperature, power demand, etc.)
- Distribution charts with dynamic binning
- Performance optimizations with TimescaleDB features like **continuous aggregates**

---

## 🚀 Features

- **Trend Charts**  
  Line charts for metrics like speed, temperature, power demand, etc.  
  Supports automatic aggregation for large time ranges to keep performance smooth.

- **Distribution Charts**  
  Histograms for selected metrics, with adjustable bin counts.

- **KPI Dashboard**  
  Highlights average speed, maximum temperature, power totals, brake pressure, and door usage ratios.

- **Live Trends**
  Observe real-time data as new telemetry data are ingested.

- **TimescaleDB Continuous Aggregates**  
  Efficient queries on large telemetry datasets via pre-aggregated materialized views.

---

## ⚙️ Installation

This project is containerized with **Docker Compose** for easy setup.

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)

### Clone the repository

```bash
git clone https://github.com/kcyakici/telemetry-dashboard.git
cd telemetry-dashboard
```

### Run the application

```bash
docker compose --profile all up --build
```

This will start:

- Backend (Go + Gin API) on http://localhost:8080
- Frontend (Next.js) on http://localhost:3000
- TimescaleDB (Postgres extension for time-series data)

### Stop the application

```bash
docker compose --profile all down -v
```

## Dataset

- This project is done using [ZTBus: A Large Dataset of Time-Resolved City Bus Driving Missions](https://www.research-collection.ethz.ch/entities/researchdata/61ac2f6e-2ca9-4229-8242-aed3b0c0d47c). You can download dataset samples and use the "Upload" section in the web application to upload CSV files.

## ⚠️ Important Notes on TimescaleDB Aggregations

This project uses continuous aggregates to speed up queries on long time ranges.

Aggregates refresh on a 1-minute interval (configurable).

Since the telemetry dataset is static historical data (2019–2021), you may need to wait until the first refresh cycle completes before running wide time-range queries.

If aggregates haven’t been refreshed yet, certain queries (trend charts for large ranges) may return empty results.

---

## 🌟 Future Improvements

- More KPIs and customizable dashboards

- Advanced TimescaleDB features (manual compression, hypercore optimizations)

- Role-based authentication

---
