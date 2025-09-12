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
