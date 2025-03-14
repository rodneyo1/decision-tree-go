package algorithm

import (
	"testing"

	"decision-tree/models"
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