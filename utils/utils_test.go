package utils

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"testing"
	"dt/models"
	"io/ioutil"
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


func createTempCSV(content string) (string, error) {
	tmpfile, err := ioutil.TempFile("", "testdata-*.csv")
	if err != nil {
		return "", err
	}

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		tmpfile.Close()
		return "", err
	}
	if err := tmpfile.Close(); err != nil {
		return "", err
	}

	return tmpfile.Name(), nil
}

func TestLoadPredictionData(t *testing.T) {
	tests := []struct {
		name        string
		csvContent  string
		expectError bool
	}{
		{
			name: "valid CSV",
			csvContent: `col1,col2,col3
			val1,val2,val3
			val4,val5,val6`,
			expectError: false,
		},
		{
			name:        "empty CSV",
			csvContent:  ``,
			expectError: true,
		},
		{
			name: "CSV with missing values",
			csvContent: `col1,col2,col3
			val1,val2
			val4,val5,val6`,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fileName, err := createTempCSV(tt.csvContent)
			if err != nil {
				t.Fatalf("failed to create temp CSV file: %v", err)
			}
			defer os.Remove(fileName)

			InputPtr = &fileName

			err = LoadPredictionData()
			if (err != nil) != tt.expectError {
				t.Errorf("LoadPredictionData() error = %v, expectError %v", err, tt.expectError)
			}

			if !tt.expectError {
				if len(models.Columns) == 0 {
					t.Errorf("expected columns to be set, got %v", models.Columns)
				}
				if len(models.Records) == 0 {
					t.Errorf("expected records to be set, got %v", models.Records)
				}
			}
		})
	}
}
