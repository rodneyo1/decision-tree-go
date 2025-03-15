package utils

import (
	"dt/models"
	"encoding/json"
	"fmt"
	"os"
)

func SaveModel(tree *models.TreeNode) error {
	modelData := models.ModelData{
		Tree:         tree,
		FeatureTypes: models.FeatureTypes,
		TargetColumn: *ColumnPtr,
		TargetType:   models.TargetType,
		Columns:      models.Columns,
	}

	file, err := os.Create(*OutputPtr)
	if err != nil {
		return fmt.Errorf("failed to create model file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(modelData); err != nil {
		return fmt.Errorf("failed to encode model data: %w", err)
	}

	return nil
}
