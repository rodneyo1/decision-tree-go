package utils

import (
	"dt/models"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"time"
)

func LoadTrainingData() error {
	csvFile, err := os.Open(*InputPtr)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer csvFile.Close()
	csvReader := csv.NewReader(csvFile)

	// Read Header Row
	columns, err := csvReader.Read()
	if err != nil {
		return fmt.Errorf("failed to read header row: %w", err)
	}
	models.Columns = columns

	// Verify target column exists
	if !slices.Contains(columns, *ColumnPtr) {
		return fmt.Errorf("target column '%s' not found in dataset", *ColumnPtr)
	}

	targetIndex := -1
	models.FeatureTypes = make(map[string]string)

	for i, col := range columns {
		models.FeatureTypes[col] = "unknown"
		if col == *ColumnPtr {
			targetIndex = i
		}
	}

	// Read and process data
	models.Records = []map[string]interface{}{}
	models.TargetValues = make(map[interface{}]int)

	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading row: %w", err)
		}

		record := make(map[string]interface{})
		for i, val := range row {
			parsedVal := parseValue(val)
			record[columns[i]] = parsedVal

			// Detect feature type if not yet determined
			if models.FeatureTypes[columns[i]] == "unknown" {
				switch parsedVal.(type) {
				case int:
					models.FeatureTypes[columns[i]] = "numeric"
				case float64:
					models.FeatureTypes[columns[i]] = "numeric"
				case time.Time:
					models.FeatureTypes[columns[i]] = "date"
				default:
					models.FeatureTypes[columns[i]] = "categorical"
				}
			}
		}

		// Track unique target values for classification
		if targetIndex >= 0 && targetIndex < len(row) {
			targetVal := parseValue(row[targetIndex])
			models.TargetValues[targetVal]++
		}

		models.Records = append(models.Records, record)
	}

	if len(models.TargetValues) <= 10 || models.FeatureTypes[*ColumnPtr] == "categorical" {
		models.TargetType = "categorical"
	} else {
		models.TargetType = "numeric"
	}

	fmt.Printf("Loaded %d records with %d columns\n", len(models.Records), len(columns))
	return nil
}

func parseValue(value string) interface{} {
	// Handle empty values
	if value == "" {
		return nil
	}

	// Try parsing as int
	if intVal, err := strconv.Atoi(value); err == nil {
		return intVal
	}

	// Try parsing as float
	if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
		return floatVal
	}

	// Try parsing as boolean
	if value == "true" || value == "false" {
		return value == "true"
	}

	// Try parsing as date
	if dateVal, err := time.Parse("2006-01-02", value); err == nil {
		return dateVal
	}

	// Default to string
	return value
}
