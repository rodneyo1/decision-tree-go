package utils

import (
	"os"
	"testing"

	"dt/models"
)

func TestLoadTrainingData(t *testing.T) {
	tests := []struct {
		name           string
		inputCSV       string
		columnPtr      string
		expectedErr    bool
		expectedTarget string
	}{
		{
			name: "Valid CSV with numeric target",
			inputCSV: `feature1,feature2,target
			1,2,3
			4,5,6
			7,8,9`,
			columnPtr:      "target",
			expectedErr:    false,
			expectedTarget: "numeric",
		},
		{
			name: "Valid CSV with categorical target",
			inputCSV: `feature1,feature2,target
			a,b,c
			d,e,f
			g,h,i`,
			columnPtr:      "target",
			expectedErr:    false,
			expectedTarget: "categorical",
		},
		{
			name: "Missing target column",
			inputCSV: `feature1,feature2
			1,2
			4,5
			7,8`,
			columnPtr:   "target",
			expectedErr: true,
		},
		{
			name:        "Empty CSV",
			inputCSV:    ``,
			columnPtr:   "target",
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary CSV file
			tmpFile, err := os.CreateTemp("", "test*.csv")
			if err != nil {
				t.Fatalf("failed to create temp file: %v", err)
			}
			defer os.Remove(tmpFile.Name())

			// Write the input CSV content to the temporary file
			if _, err := tmpFile.WriteString(tt.inputCSV); err != nil {
				t.Fatalf("failed to write to temp file: %v", err)
			}
			tmpFile.Close()

			inputFileName := tmpFile.Name()
			columnName := tt.columnPtr

			// Set the input and column pointers
			InputPtr = &inputFileName
			ColumnPtr = &columnName

			err = LoadTrainingData()

			if (err != nil) != tt.expectedErr {
				t.Errorf("expected error: %v, got: %v", tt.expectedErr, err)
			}

			if !tt.expectedErr && models.TargetType != tt.expectedTarget {
				t.Errorf("expected target type: %v, got: %v", tt.expectedTarget, models.TargetType)
			}
		})
	}
}
