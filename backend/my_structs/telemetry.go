package my_structs

import "time"

type Telemetry struct {
	VehicleID                 *string   `json:"vehicle_id,omitempty" db:"vehicle_id"`
	TimeISO                   time.Time `json:"time_iso" db:"time_iso"`
	TimeUnix                  *int64    `json:"time_unix,omitempty" db:"time_unix"`
	ElectricPowerDemand       *float64  `json:"electric_power_demand,omitempty" db:"electric_power_demand"`
	GnssAltitude              *float64  `json:"gnss_altitude,omitempty" db:"gnss_altitude"`
	GnssCourse                *float64  `json:"gnss_course,omitempty" db:"gnss_course"`
	GnssLatitude              *float64  `json:"gnss_latitude,omitempty" db:"gnss_latitude"`
	GnssLongitude             *float64  `json:"gnss_longitude,omitempty" db:"gnss_longitude"`
	ItcsBusRoute              *string   `json:"itcs_busRoute,omitempty" db:"itcs_bus_route"`
	ItcsNumberOfPassengers    *float64  `json:"itcs_number_of_passengers,omitempty" db:"itcs_number_of_passengers"`
	ItcsStopName              *string   `json:"itcs_stop_name,omitempty" db:"itcs_stop_name"`
	OdometryArticulationAngle *float64  `json:"odometry_articulation_angle,omitempty" db:"odometry_articulation_angle"`
	OdometrySteeringAngle     *float64  `json:"odometry_steering_angle,omitempty" db:"odometry_steering_angle"`
	OdometryVehicleSpeed      *float64  `json:"odometry_vehicle_speed,omitempty" db:"odometry_vehicle_speed"`
	OdometryWheelSpeedFL      *float64  `json:"odometry_wheel_speed_fl,omitempty" db:"odometry_wheel_speed_fl"`
	OdometryWheelSpeedFR      *float64  `json:"odometry_wheel_speed_fr,omitempty" db:"odometry_wheel_speed_fr"`
	OdometryWheelSpeedML      *float64  `json:"odometry_wheel_speed_ml,omitempty" db:"odometry_wheel_speed_ml"`
	OdometryWheelSpeedMR      *float64  `json:"odometry_wheel_speed_mr,omitempty" db:"odometry_wheel_speed_mr"`
	OdometryWheelSpeedRL      *float64  `json:"odometry_wheel_speed_rl,omitempty" db:"odometry_wheel_speed_rl"`
	OdometryWheelSpeedRR      *float64  `json:"odometry_wheel_speed_rr,omitempty" db:"odometry_wheel_speed_rr"`
	StatusDoorIsOpen          *int      `json:"status_door_is_open,omitempty" db:"status_door_is_open"`
	StatusGridIsAvailable     *int      `json:"status_grid_is_available,omitempty" db:"status_grid_is_available"`
	StatusHaltBrakeIsActive   *int      `json:"status_halt_brake_is_active,omitempty" db:"status_halt_brake_is_active"`
	StatusParkBrakeIsActive   *int      `json:"status_park_brake_is_active,omitempty" db:"status_park_brake_is_active"`
	TemperatureAmbient        *float64  `json:"temperature_ambient,omitempty" db:"temperature_ambient"`
	TractionBrakePressure     *float64  `json:"traction_brake_pressure,omitempty" db:"traction_brake_pressure"`
	TractionTractionForce     *float64  `json:"traction_traction_force,omitempty" db:"traction_traction_force"`
}
