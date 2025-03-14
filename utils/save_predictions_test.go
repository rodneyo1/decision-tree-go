package utils

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"testing"
)

func TestSavePredictions(t *testing.T) {
	tempDir := t.TempDir()
	outputFile := filepath.Join(tempDir, "predictions.csv")
	OutputPtr = &outputFile

	tests := []struct {
		name        string
		predictions []interface{}
		wantErr     bool
	}{
		{
			name:        "valid predictions",
			predictions: []interface{}{"A", "B", "C"},
			wantErr:     false,
		},
		{
			name:        "nil predictions",
			predictions: []interface{}{nil, "B", nil},
			wantErr:     false,
		},
		{
			name:        "empty predictions",
			predictions: []interface{}{},
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SavePredictions(tt.predictions)
			if (err != nil) != tt.wantErr {
				t.Errorf("SavePredictions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Verify the file content
			file, err := os.Open(outputFile)
			if err != nil {
				t.Fatalf("failed to open output file: %v", err)
			}
			defer file.Close()

			reader := csv.NewReader(file)
			records, err := reader.ReadAll()
			if err != nil {
				t.Fatalf("failed to read CSV file: %v", err)
			}

			if len(records) != len(tt.predictions)+1 {
				t.Errorf("expected %d records, got %d", len(tt.predictions)+1, len(records))
			}

			for i, pred := range tt.predictions {
				expected := "unknown"
				if pred != nil {
					expected = pred.(string)
				}
				if records[i+1][0] != expected {
					t.Errorf("expected prediction %v, got %v", expected, records[i+1][0])
				}
			}
		})
	}
}