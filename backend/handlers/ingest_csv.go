package handlers

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"telemetry-dashboard/my_structs"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func IngestCSV(c *gin.Context, pool *pgxpool.Pool) {
	ct := c.ContentType()
	// log.Printf("Content type: " + ct)
	// for k, v := range c.Request.Header {
	// 	log.Printf("%s: %v", k, v)
	// } //TODO delete
	if ct != "multipart/form-data" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid content type, must be multipart/form-data, instead got " + ct})
		return
	}

	// Limit file size 15 MB
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 15<<20)

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file upload failed: " + err.Error()})
		return
	}
	defer file.Close()

	if !strings.HasSuffix(strings.ToLower(header.Filename), ".csv") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file extension, must be .csv"})
		return
	}

	// Extract vehicle_id from file name
	// Example: B183_2019-06-24_03-16-13_2019-06-24_18-54-06.csv
	parts := strings.Split(header.Filename, "_")
	vehicleID := parts[0]

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = len(expectedColumnsInCsv)

	headerRow, err := reader.Read()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read header row"})
		return
	}

	// Map CSV headers → DB columns
	mappedCols := make([]string, len(headerRow))
	for i, csvCol := range headerRow {
		dbCol, ok := csvHeaderToDb[csvCol]
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("unexpected column %d: got '%s'", i, csvCol),
			})
			return
		}
		mappedCols[i] = dbCol
	}

	// Collect rows
	var rows [][]interface{}
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid csv row: " + err.Error()})
			return
		}

		// Convert CSV strings → Go types
		row, convErr := parseCSVRecord(record)
		if convErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "parse error: " + convErr.Error()})
			return
		}

		fullRow := append([]interface{}{vehicleID}, row...)
		rows = append(rows, fullRow)
	}

	cols := append([]string{"vehicle_id"}, mappedCols...)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	tx, err := pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db begin: " + err.Error()})
		return
	}
	defer tx.Rollback(ctx)

	_, err = tx.CopyFrom(ctx, pgx.Identifier{"telemetry"}, cols, pgx.CopyFromRows(rows))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "copy from failed: " + err.Error()})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "commit failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "inserted": len(rows), "vehicle_id": vehicleID})
}

// ParseCSVRecord dynamically maps a CSV row into []interface{} according to Telemetry struct.
func parseCSVRecord(rec []string) ([]interface{}, error) {
	typ := reflect.TypeOf(my_structs.Telemetry{})
	numFields := typ.NumField()

	if len(rec) != numFields-1 { // -1 because VehicleID comes from filename, not CSV
		return nil, fmt.Errorf("unexpected column count: got %d, want %d", len(rec), numFields-1)
	}

	out := make([]interface{}, numFields-1)

	// Skip VehicleID (that comes from filename)
	for i := 1; i < numFields; i++ {
		field := typ.Field(i)
		val := rec[i-1] // shift because VehicleID is skipped

		switch field.Type.Kind() {
		case reflect.Struct: // time.Time
			ts, err := time.Parse(time.RFC3339, val)
			if err != nil {
				return nil, fmt.Errorf("invalid time_iso: %s", val)
			}
			out[i-1] = ts

		case reflect.Ptr:
			elemType := field.Type.Elem().Kind()
			if val == "NaN" || val == "-" {
				out[i-1] = nil
				continue
			}
			switch elemType {
			case reflect.Float64:
				f, err := strconv.ParseFloat(val, 64)
				if err != nil {
					return nil, fmt.Errorf("invalid float in %s: %s", field.Name, val)
				}
				out[i-1] = f
			case reflect.Int, reflect.Int64:
				// time_unix is int64, status_* are int
				n, err := strconv.ParseInt(val, 10, 64)
				if err != nil {
					return nil, fmt.Errorf("invalid int in %s: %s", field.Name, val)
				}
				// handle int vs int64 separately
				if field.Type.Elem().Kind() == reflect.Int {
					out[i-1] = int(n)
				} else {
					out[i-1] = n
				}
			case reflect.String:
				out[i-1] = val
			default:
				return nil, fmt.Errorf("unsupported pointer type: %s", field.Type.String())
			}

		default:
			return nil, fmt.Errorf("unsupported field kind: %s", field.Type.Kind())
		}
	}

	return out, nil
}
