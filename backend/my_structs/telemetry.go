package my_structs

import "time"

type Telemetry struct {
	VehicleID                 *string   `json:"vehicle_id,omitempty"`
	TimeISO                   time.Time `json:"time_iso"` // RFC3339 string in JSON
	TimeUnix                  *int64    `json:"time_unix,omitempty"`
	ElectricPowerDemand       *float64  `json:"electric_power_demand,omitempty"`
	GnssAltitude              *float64  `json:"gnss_altitude,omitempty"`
	GnssCourse                *float64  `json:"gnss_course,omitempty"`
	GnssLatitude              *float64  `json:"gnss_latitude,omitempty"`
	GnssLongitude             *float64  `json:"gnss_longitude,omitempty"`
	ItcsBusRoute              *string   `json:"itcs_busRoute,omitempty"`
	ItcsNumberOfPassengers    *float64  `json:"itcs_number_of_passengers,omitempty"`
	ItcsStopName              *string   `json:"itcs_stop_name,omitempty"`
	OdometryArticulationAngle *float64  `json:"odometry_articulation_angle,omitempty"`
	OdometrySteeringAngle     *float64  `json:"odometry_steering_angle,omitempty"`
	OdometryVehicleSpeed      *float64  `json:"odometry_vehicle_speed,omitempty"`
	OdometryWheelSpeedFL      *float64  `json:"odometry_wheel_speed_fl,omitempty"`
	OdometryWheelSpeedFR      *float64  `json:"odometry_wheel_speed_fr,omitempty"`
	OdometryWheelSpeedML      *float64  `json:"odometry_wheel_speed_ml,omitempty"`
	OdometryWheelSpeedMR      *float64  `json:"odometry_wheel_speed_mr,omitempty"`
	OdometryWheelSpeedRL      *float64  `json:"odometry_wheel_speed_rl,omitempty"`
	OdometryWheelSpeedRR      *float64  `json:"odometry_wheel_speed_rr,omitempty"`
	StatusDoorIsOpen          *int      `json:"status_door_is_open,omitempty"`
	StatusGridIsAvailable     *int      `json:"status_grid_is_available,omitempty"`
	StatusHaltBrakeIsActive   *int      `json:"status_halt_brake_is_active,omitempty"`
	StatusParkBrakeIsActive   *int      `json:"status_park_brake_is_active,omitempty"`
	TemperatureAmbient        *float64  `json:"temperature_ambient,omitempty"`
	TractionBrakePressure     *float64  `json:"traction_brake_pressure,omitempty"`
	TractionTractionForce     *float64  `json:"traction_traction_force,omitempty"`
}
