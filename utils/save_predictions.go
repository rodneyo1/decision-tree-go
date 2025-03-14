package utils

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
)

func SavePredictions(predictions []interface{}) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(*OutputPtr)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	// Create CSV file
	file, err := os.Create(*OutputPtr)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	// Create CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	if err := writer.Write([]string{"prediction"}); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Write predictions
	for _, pred := range predictions {
		var strPred string
		if pred == nil {
			strPred = "unknown"
		} else {
			strPred = fmt.Sprintf("%v", pred)
		}
		if err := writer.Write([]string{strPred}); err != nil {
			return fmt.Errorf("failed to write prediction: %w", err)
		}
	}

	fmt.Printf("Predictions saved to %s\n", *OutputPtr)
	return nil
}
