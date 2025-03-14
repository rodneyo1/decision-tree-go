package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"dt/models"
)

func SaveModel(tree *models.TreeNode) error {
	// Create the model data structure
	modelData := models.ModelData{
		Tree:         tree,
		FeatureTypes: models.FeatureTypes,
		TargetColumn: *ColumnPtr,
		TargetType:   models.TargetType,
		Columns:      models.Columns,
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(*OutputPtr)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	// Marshal to JSON
	jsonData, err := json.MarshalIndent(modelData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal model to JSON: %w", err)
	}

	// Write to file
	if err := os.WriteFile(*OutputPtr, jsonData, 0o644); err != nil {
		return fmt.Errorf("failed to write model file: %w", err)
	}

	fmt.Printf("Model saved to %s\n", *OutputPtr)
	return nil
}
