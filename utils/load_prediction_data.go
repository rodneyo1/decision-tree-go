package utils

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"decision-tree/models"
)

func LoadPredictionData() error {
	csvFile, err := os.Open(*InputPtr)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer csvFile.Close()

	csvReader := csv.NewReader(csvFile)

	columns, err := csvReader.Read()
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

		models.Records = append(models.Records, record)
	}

	fmt.Printf("Loaded %d records with %d columns for prediction\n", len(models.Records), len(columns))
	return nil
}
