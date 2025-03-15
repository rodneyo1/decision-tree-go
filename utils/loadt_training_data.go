package utils

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"time"

	"dt/models"
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

	// Initialize records and target values
	models.Records = []map[string]interface{}{}
	models.TargetValues = make(map[interface{}]int)

	// Calculate mean for numeric columns and mode for non-numeric columns
	columnMeans := make(map[int]float64)
	columnModes := make(map[int]string)
	columnCounts := make(map[int]int)
	columnValueCounts := make(map[int]map[string]int)
	
	// First pass: calculate means and modes for imputation
	allRows := [][]string{}
	
	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading row: %w", err)
		}
		
		allRows = append(allRows, row)
		
		for i, value := range row {
			if value == "" {
				continue
			}
			if _, err := strconv.ParseFloat(value, 64); err == nil {
				// Numeric column
				if _, exists := columnMeans[i]; !exists {
					columnMeans[i] = 0
					columnCounts[i] = 0
				}
				val, _ := strconv.ParseFloat(value, 64)
				columnMeans[i] += val
				columnCounts[i]++
			} else {
				// Non-numeric column
				if _, exists := columnValueCounts[i]; !exists {
					columnValueCounts[i] = make(map[string]int)
				}
				columnValueCounts[i][value]++
			}
		}
	}
	
	// Calculate final means and modes
	for i := range columnMeans {
		if columnCounts[i] > 0 {
			columnMeans[i] /= float64(columnCounts[i])
		}
	}
	for i, valueCounts := range columnValueCounts {
		maxCount := 0
		for value, count := range valueCounts {
			if count > maxCount {
				maxCount = count
				columnModes[i] = value
			}
		}
	}

	// Second pass: process records with imputation
	batchSize := 1000
	batch := make([]map[string]interface{}, 0, batchSize)

	for _, row := range allRows {
		record := make(map[string]interface{})
		
		for i, val := range row {
			// Impute missing values
			if val == "" {
				if mean, exists := columnMeans[i]; exists {
					val = fmt.Sprintf("%f", mean)
				} else if mode, exists := columnModes[i]; exists {
					val = mode
				}
			}
			
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
			targetVal := record[columns[targetIndex]]
			models.TargetValues[targetVal]++
		}

		batch = append(batch, record)
		if len(batch) == batchSize {
			models.Records = append(models.Records, batch...)
			batch = make([]map[string]interface{}, 0, batchSize)
		}
	}

	// Append any remaining records
	if len(batch) > 0 {
		models.Records = append(models.Records, batch...)
	}

	// Determine target type
	if models.FeatureTypes[*ColumnPtr] == "numeric" {
		models.TargetType = "numeric" 
	} else {
		models.TargetType = "categorical"
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