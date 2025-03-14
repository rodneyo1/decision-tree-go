package algorithm

import (
	"fmt"

	"decision-tree/models"
)

const (
	MaxDepth       = 20    // Maximum tree depth
	MinSamplesLeaf = 5     // Minimum samples required to form a leaf
	MinInfoGain    = 0.001 // Minimum information gain required to split
)

func BuildTree(targetCol string) (*models.TreeNode, error) {
	fmt.Println("Building decision tree for target:", targetCol)

	// Get available features (exclude target column)
	features := make([]string, 0)
	for _, col := range models.Columns {
		if col != targetCol {
			features = append(features, col)
		}
	}

	// Create indices for all records
	indices := make([]int, len(models.Records))
	for i := range indices {
		indices[i] = i
	}
	tree := buildTreeNode(indices, features, targetCol, 0)
	fmt.Println("Tree building complete")
	return tree, nil
}

func buildTreeNode(indices []int, features []string, targetCol string, depth int) *models.TreeNode {
}
