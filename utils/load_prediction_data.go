package utils

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

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

	columns, err := csvReader.Read()
	if err == io.EOF {
		return fmt.Errorf("input file is empty")
	}
	if err != nil {
		return fmt.Errorf("failed to read header row: %w", err)
	}

	models.Columns = columns
	models.Records = []map[string]interface{}{}

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
			if i < len(columns) {
				record[columns[i]] = parseValue(val)
			}
		}

		// Fill in missing values with nil
		for _, col := range columns {
			if _, exists := record[col]; !exists {
				record[col] = nil
			}
		}

		models.Records = append(models.Records, record)
	}

	fmt.Printf("Loaded %d records with %d columns for prediction\n", len(models.Records), len(columns))
	return nil
}
