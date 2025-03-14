package utils

import (
	"encoding/json"
	"fmt"
	"os"

	"dt/models"
)

func LoadModels() (*models.ModelData, error) {
	jsonData, err := os.ReadFile(*ModelFilePtr)
	if err != nil {
		return nil, fmt.Errorf("failed to read model file: %w", err)
	}
	var modelData models.ModelData
	if err := json.Unmarshal(jsonData, &modelData); err != nil {
		return nil, fmt.Errorf("failed to parse model data: %w", err)
	}
	// Update global model data
	models.FeatureTypes = modelData.FeatureTypes
	models.TargetColumn = modelData.TargetColumn
	models.TargetType = modelData.TargetType
	models.Columns = modelData.Columns

	fmt.Printf("Loaded model trained for target column: %s\n", modelData.TargetColumn)
	return &modelData, nil
}
