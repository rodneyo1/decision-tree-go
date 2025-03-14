package utils

import (
	"io/ioutil"
	"os"
	"testing"

	"dt/models"
)

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