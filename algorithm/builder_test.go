package algorithm

import (
	"testing"

	"dt/models"
)

// Mock data setup for testing
func setupMockTreeData() {
	models.Records = []map[string]interface{}{
		{"feature": "A", "target": "Yes"},
		{"feature": "A", "target": "Yes"},
		{"feature": "B", "target": "No"},
		{"feature": "B", "target": "No"},
		{"feature": "C", "target": "Yes"},
		{"feature": "C", "target": "No"},
		{"feature": "A", "target": "Yes"},
		{"feature": "B", "target": "No"},
		{"feature": "C", "target": "Yes"},
	}

	models.Columns = []string{"feature", "target"}
	models.FeatureTypes = map[string]string{
		"feature": "categorical",
	}
}

// Test BuildTree function
func TestBuildTree(t *testing.T) {
	setupMockTreeData()

	tree, err := BuildTree("target")
	if err != nil {
		t.Fatalf("BuildTree returned an error: %v", err)
	}

	if tree == nil {
		t.Fatal("BuildTree returned a nil tree")
	}

	if !tree.IsLeaf && tree.Feature == "" {
		t.Errorf("Expected a decision node with a feature, got an empty feature")
	}
}

// Test buildTreeNode function directly
func TestBuildTreeNode(t *testing.T) {
	setupMockTreeData()

	indices := []int{0, 1, 2, 3, 4, 5, 6, 7, 8}
	features := []string{"feature"}
	targetCol := "target"

	rootNode := buildTreeNode(indices, features, targetCol, 0)

	if rootNode == nil {
		t.Fatal("buildTreeNode returned nil")
	}

	if !rootNode.IsLeaf && rootNode.Feature == "" {
		t.Errorf("Expected a split on a feature, but got an empty feature")
	}

	if rootNode.IsLeaf && rootNode.Prediction == nil {
		t.Errorf("Expected a prediction at a leaf node but got nil")
	}
}
