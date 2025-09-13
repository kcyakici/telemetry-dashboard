package handlers

import (
	"context"
	"net/http"
	"telemetry-dashboard/my_structs"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Ingest(c *gin.Context, pool *pgxpool.Pool) {
	var data []my_structs.Telemetry
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json: " + err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db begin: " + err.Error()})
		return
	}
	defer tx.Rollback(ctx) // safe no-op if already committed

	// build rows for CopyFrom
	rows := make([][]interface{}, 0, len(data))
	for _, t := range data {
		// convert pointers to actual interface values (nil for missing)
		var vehicle interface{}
		if t.VehicleID != nil {
			vehicle = *t.VehicleID
		} else {
			vehicle = nil
		}

		var tu interface{}
		if t.TimeUnix != nil {
			tu = *t.TimeUnix
		} else {
			tu = nil
		}

		rows = append(rows, []interface{}{
			vehicle,
			t.TimeISO,
			tu,
			floatOrNil(t.ElectricPowerDemand),
			floatOrNil(t.GnssAltitude),
			floatOrNil(t.GnssCourse),
			floatOrNil(t.GnssLatitude),
			floatOrNil(t.GnssLongitude),
			stringOrNil(t.ItcsBusRoute),
			floatOrNil(t.ItcsNumberOfPassengers),
			stringOrNil(t.ItcsStopName),
			floatOrNil(t.OdometryArticulationAngle),
			floatOrNil(t.OdometrySteeringAngle),
			floatOrNil(t.OdometryVehicleSpeed),
			floatOrNil(t.OdometryWheelSpeedFL),
			floatOrNil(t.OdometryWheelSpeedFR),
			floatOrNil(t.OdometryWheelSpeedML),
			floatOrNil(t.OdometryWheelSpeedMR),
			floatOrNil(t.OdometryWheelSpeedRL),
			floatOrNil(t.OdometryWheelSpeedRR),
			intOrNil(t.StatusDoorIsOpen),
			intOrNil(t.StatusGridIsAvailable),
			intOrNil(t.StatusHaltBrakeIsActive),
			intOrNil(t.StatusParkBrakeIsActive),
			floatOrNil(t.TemperatureAmbient),
			floatOrNil(t.TractionBrakePressure),
			floatOrNil(t.TractionTractionForce),
		})
	}

	// do CopyFrom
	_, err = tx.CopyFrom(ctx, pgx.Identifier{"telemetry"}, expectedColumnsInCsv, pgx.CopyFromRows(rows))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "copy from failed: " + err.Error()})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "commit failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "inserted": len(data)})
}

// helpers
func floatOrNil(p *float64) interface{} {
	if p == nil {
		return nil
	}
	return *p
}
func stringOrNil(p *string) interface{} {
	if p == nil {
		return nil
	}
	return *p
}
func intOrNil(p *int) interface{} {
	if p == nil {
		return nil
	}
	return *p
}
