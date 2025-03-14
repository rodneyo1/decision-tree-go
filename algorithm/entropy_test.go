package algorithm

import (
	//"math"
	"testing"

	"dt/models"
)

// Mock data for testing
func setupMockData() {
	models.Records = []map[string]interface{}{
		{"feature": "A", "target": "Yes"},
		{"feature": "A", "target": "No"},
		{"feature": "B", "target": "Yes"},
		{"feature": "B", "target": "Yes"},
		{"feature": "C", "target": "No"},
	}

	models.FeatureTypes = map[string]string{
		"feature": "categorical",
	}
}



// Test FindBestSplit function
func TestFindBestSplit(t *testing.T) {
	setupMockData()
	indices := []int{0, 1, 2, 3, 4}
	features := []string{"feature"}
	targetCol := "target"

	split := FindBestSplit(indices, features, targetCol)

	if split.Feature != "feature" {
		t.Errorf("Expected best split feature to be 'feature', got %v", split.Feature)
	}

	if split.GainRatio <= 0 {
		t.Errorf("Expected positive gain ratio, got %f", split.GainRatio)
	}
}
