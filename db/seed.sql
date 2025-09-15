-- Enable TimescaleDB extension
CREATE EXTENSION IF NOT EXISTS timescaledb;

-- Create the telemetry table
CREATE TABLE IF NOT EXISTS telemetry (
    vehicle_id TEXT NOT NULL,
    time_iso TIMESTAMPTZ NOT NULL,
    time_unix BIGINT,
    electric_power_demand DOUBLE PRECISION,
    gnss_altitude DOUBLE PRECISION,
    gnss_course DOUBLE PRECISION,
    gnss_latitude DOUBLE PRECISION,
    gnss_longitude DOUBLE PRECISION,
    itcs_bus_route TEXT,
    itcs_number_of_passengers DOUBLE PRECISION,
    itcs_stop_name TEXT,
    odometry_articulation_angle DOUBLE PRECISION,
    odometry_steering_angle DOUBLE PRECISION,
    odometry_vehicle_speed DOUBLE PRECISION,
    odometry_wheel_speed_fl DOUBLE PRECISION,
    odometry_wheel_speed_fr DOUBLE PRECISION,
    odometry_wheel_speed_ml DOUBLE PRECISION,
    odometry_wheel_speed_mr DOUBLE PRECISION,
    odometry_wheel_speed_rl DOUBLE PRECISION,
    odometry_wheel_speed_rr DOUBLE PRECISION,
    status_door_is_open INT,
    status_grid_is_available INT,
    status_halt_brake_is_active INT,
    status_park_brake_is_active INT,
    temperature_ambient DOUBLE PRECISION,
    traction_brake_pressure DOUBLE PRECISION,
    traction_traction_force DOUBLE PRECISION,
    PRIMARY KEY (vehicle_id, time_iso)
);

-- Convert the table into a hypertable
SELECT create_hypertable('telemetry', 'time_iso', if_not_exists => TRUE);

-- Enable compression for old chunks (compress after 7 days)
-- Might be problematic for static old data
-- ALTER TABLE telemetry SET (
--     timescaledb.compress,
--     timescaledb.compress_orderby = 'time_iso',
--     timescaledb.compress_segmentby = 'vehicle_id'
-- );
-- SELECT add_compression_policy('telemetry', INTERVAL '7 days');

-- Continuous aggregates

-- Speed
CREATE MATERIALIZED VIEW IF NOT EXISTS trend_speed_1min
WITH (timescaledb.continuous) AS
SELECT time_bucket('1 minute', time_iso) AS bucket,
       vehicle_id,
       AVG(odometry_vehicle_speed) AS avg_speed
FROM telemetry
GROUP BY bucket, vehicle_id;

-- Temp
CREATE MATERIALIZED VIEW IF NOT EXISTS trend_temp_1min
WITH (timescaledb.continuous) AS
SELECT time_bucket('1 minute', time_iso) AS bucket,
       vehicle_id,
       AVG(temperature_ambient) AS avg_temp
FROM telemetry
GROUP BY bucket, vehicle_id;

-- Power
CREATE MATERIALIZED VIEW IF NOT EXISTS trend_power_1min
WITH (timescaledb.continuous) AS
SELECT time_bucket('1 minute', time_iso) AS bucket,
       vehicle_id,
       AVG(electric_power_demand) AS avg_power
FROM telemetry
GROUP BY bucket, vehicle_id;

-- Add refresh policies (refresh every 5 minutes, look back 1 hour)
SELECT add_continuous_aggregate_policy('trend_speed_1min',
    start_offset => NULL,
    -- start_offset => INTERVAL '1 hour',
    end_offset   => INTERVAL '1 minute',
    schedule_interval => INTERVAL '1 minute');

SELECT add_continuous_aggregate_policy('trend_temp_1min',
    start_offset => NULL,
    -- start_offset => INTERVAL '1 hour',
    end_offset   => INTERVAL '1 minute',
    schedule_interval => INTERVAL '1 minute');

SELECT add_continuous_aggregate_policy('trend_power_1min',
    start_offset => NULL,
    -- start_offset => INTERVAL '1 hour',
    end_offset   => INTERVAL '1 minute',
    schedule_interval => INTERVAL '1 minute');
