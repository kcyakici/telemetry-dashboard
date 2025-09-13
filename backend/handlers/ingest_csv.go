package handlers

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func IngestCSV(c *gin.Context, pool *pgxpool.Pool) {
	ct := c.ContentType()
	if ct != "multipart/form-data" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid content type, must be multipart/form-data"})
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

// Helper: convert CSV strings into proper types (NULL for NaN / "-")
func parseCSVRecord(rec []string) ([]interface{}, error) {
	out := make([]interface{}, len(expectedColumnsInCsv))

	// 0: time_iso
	ts, err := time.Parse(time.RFC3339, rec[0])
	if err != nil {
		return nil, fmt.Errorf("invalid time_iso: %s", rec[0])
	}
	out[0] = ts

	// 1: time_unix
	if rec[1] == "NaN" || rec[1] == "-" {
		out[1] = nil
	} else {
		val, _ := strconv.ParseInt(rec[1], 10, 64)
		out[1] = val
	}

	// Generic numeric / string parsing for the rest
	for i := 2; i < len(rec); i++ {
		v := rec[i]
		if v == "NaN" || v == "-" {
			out[i] = nil
			continue
		}
		// try float first
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			out[i] = f
		} else {
			out[i] = v
		}
	}
	return out, nil
}
