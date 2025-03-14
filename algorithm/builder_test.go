package algorithm_test

import (
	"os"
	"testing"

	"dt/models"
	"dt/algorithm"
	"dt/utils"
)

// Setup a mock dataset for testing
func setupMockData() {
	models.Columns = []string{"Feature1", "Feature2", "Target"}
	models.Records = []map[string]interface{}{
		{"Feature1": "A", "Feature2": 1.2, "Target": "Yes"},
		{"Feature1": "B", "Feature2": 2.4, "Target": "No"},
		{"Feature1": "A", "Feature2": 1.5, "Target": "Yes"},
		{"Feature1": "B", "Feature2": 2.1, "Target": "No"},
		{"Feature1": "A", "Feature2": 1.8, "Target": "Yes"},
	}
}

func TestBuildTree(t *testing.T) {
	setupMockData()

	// Ensure OutputPtr is set to a valid test file
	tempFile := "test_model.json"
	utils.OutputPtr = &tempFile // Assign a valid file path

	// Run BuildTree
	tree, err := algorithm.BuildTree("Target")
	if err != nil {
		t.Fatalf("BuildTree returned an error: %v", err)
	}

	if tree == nil {
		t.Fatal("BuildTree returned a nil tree")
	}

	if tree.IsLeaf && tree.Prediction == nil {
		t.Fatal("Leaf node has no prediction")
	}

	// Clean up the test output file
	os.Remove(tempFile)
}
