CREATE EXTENSION IF NOT EXISTS timescaledb;

CREATE TABLE telemetry (
    id SERIAL PRIMARY KEY,
    vehicle_id TEXT NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL,
    speed DOUBLE PRECISION,
    temperature DOUBLE PRECISION,
    energy DOUBLE PRECISION
);

SELECT create_hypertable('telemetry', 'timestamp');
