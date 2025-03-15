package utils

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"

	"dt/models"
)

func LoadPredictionData() error {
	csvFile, err := os.Open(*InputPtr)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer csvFile.Close()

	csvReader := csv.NewReader(csvFile)
	csvReader.FieldsPerRecord = -1

	// Read Header Row
	columns, err := csvReader.Read()
	if err == io.EOF {
		return fmt.Errorf("input file is empty")
	}
	if err != nil {
		return fmt.Errorf("failed to read header row: %w", err)
	}
	models.Columns = columns

	// Calculate mean for numeric columns and mode for non-numeric columns
	columnMeans := make(map[int]float64)
	columnModes := make(map[int]string)
	columnCounts := make(map[int]int)
	columnValueCounts := make(map[int]map[string]int)
	
	// First pass: read all rows and calculate means and modes for imputation
	allRows := [][]string{}
	
	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading row: %w", err)
		}
		
		// Ensure row has enough elements
		if len(row) < len(columns) {
			extendedRow := make([]string, len(columns))
			copy(extendedRow, row)
			row = extendedRow
		}
		
		allRows = append(allRows, row)
		
		for i, value := range row {
			if i >= len(columns) {
				break
			}
			
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
	models.Records = []map[string]interface{}{}
	
	for _, row := range allRows {
		record := make(map[string]interface{})
		
		for i, col := range columns {
			val := ""
			if i < len(row) {
				val = row[i]
			}
			
			// Impute missing values
			if val == "" {
				if mean, exists := columnMeans[i]; exists {
					val = fmt.Sprintf("%f", mean)
				} else if mode, exists := columnModes[i]; exists {
					val = mode
				}
			}
			
			record[col] = parseValue(val)
		}
		
		models.Records = append(models.Records, record)
	}

	fmt.Printf("Loaded %d records with %d columns for prediction\n", len(models.Records), len(columns))
	return nil
}