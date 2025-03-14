package algorithm

import (
	"testing"

	"dt/models"
)

func TestEstimateNodeSize(t *testing.T) {
	tests := []struct {
		name     string
		node     *models.TreeNode
		expected int
	}{
		{
			name:     "Nil node",
			node:     nil,
			expected: 0,
		},
		{
			name: "Leaf node",
			node: &models.TreeNode{
				IsLeaf: true,
			},
			expected: 1,
		},
		{
			name: "Categorical split with children",
			node: &models.TreeNode{
				SplitType: "categorical",
				Children: map[string]*models.TreeNode{
					"child1": {IsLeaf: true},
					"child2": {IsLeaf: true},
				},
			},
			expected: 2,
		},
		{
			name: "Numerical split with left and right children",
			node: &models.TreeNode{
				SplitType: "numerical",
				Left:      &models.TreeNode{IsLeaf: true},
				Right:     &models.TreeNode{IsLeaf: true},
			},
			expected: 2,
		},
		{
			name: "Mixed split types",
			node: &models.TreeNode{
				SplitType: "categorical",
				Children: map[string]*models.TreeNode{
					"child1": {
						SplitType: "numerical",
						Left:      &models.TreeNode{IsLeaf: true},
						Right:     &models.TreeNode{IsLeaf: true},
					},
					"child2": {IsLeaf: true},
				},
			},
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := estimateNodeSize(tt.node)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestPredict(t *testing.T) {
	tests := []struct {
		name     string
		tree     *models.TreeNode
		records  []map[string]interface{}
		expected []interface{}
	}{
		{
			name: "Single leaf node",
			tree: &models.TreeNode{
				IsLeaf:     true,
				Prediction: "A",
			},
			records: []map[string]interface{}{
				{"feature1": "value1"},
				{"feature2": "value2"},
			},
			expected: []interface{}{"A", "A"},
		},
		{
			name: "Categorical split",
			tree: &models.TreeNode{
				SplitType: "categorical",
				Feature:   "feature1",
				Children: map[string]*models.TreeNode{
					"value1": {IsLeaf: true, Prediction: "A"},
					"value2": {IsLeaf: true, Prediction: "B"},
				},
			},
			records: []map[string]interface{}{
				{"feature1": "value1"},
				{"feature1": "value2"},
			},
			expected: []interface{}{"A", "B"},
		},
		{
			name: "Numerical split",
			tree: &models.TreeNode{
				SplitType:  "numerical",
				Feature:    "feature1",
				SplitValue: 10,
				Left:       &models.TreeNode{IsLeaf: true, Prediction: "A"},
				Right:      &models.TreeNode{IsLeaf: true, Prediction: "B"},
			},
			records: []map[string]interface{}{
				{"feature1": 5},
				{"feature1": 15},
			},
			expected: []interface{}{"A", "B"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			models.Records = tt.records
			result := Predict(tt.tree)
			for i, prediction := range result {
				if prediction != tt.expected[i] {
					t.Errorf("expected %v, got %v", tt.expected[i], prediction)
				}
			}
		})
	}
}

func TestPredictRecord(t *testing.T) {
	tests := []struct {
		name     string
		record   map[string]interface{}
		node     *models.TreeNode
		expected interface{}
	}{
		{
			name: "Leaf node",
			record: map[string]interface{}{
				"feature1": "value1",
			},
			node: &models.TreeNode{
				IsLeaf:     true,
				Prediction: "A",
			},
			expected: "A",
		},
		{
			name: "Categorical split",
			record: map[string]interface{}{
				"feature1": "value1",
			},
			node: &models.TreeNode{
				SplitType: "categorical",
				Feature:   "feature1",
				Children: map[string]*models.TreeNode{
					"value1": {IsLeaf: true, Prediction: "A"},
					"value2": {IsLeaf: true, Prediction: "B"},
				},
			},
			expected: "A",
		},
		{
			name: "Numerical split",
			record: map[string]interface{}{
				"feature1": 5,
			},
			node: &models.TreeNode{
				SplitType:  "numerical",
				Feature:    "feature1",
				SplitValue: 10,
				Left:       &models.TreeNode{IsLeaf: true, Prediction: "A"},
				Right:      &models.TreeNode{IsLeaf: true, Prediction: "B"},
			},
			expected: "A",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := predictRecord(tt.record, tt.node)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
