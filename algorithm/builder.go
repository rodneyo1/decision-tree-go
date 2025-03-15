package algorithm

import (
	"fmt"
	"sync"

	"dt/models"
	"dt/utils"
)

const (
	MaxDepth       = 20    // Maximum tree depth
	MinSamplesLeaf = 5     // Minimum samples required to form a leaf
	MinInfoGain    = 0.001 // Minimum information gain required to split
)

// BuildTree builds a decision tree from the given dataset
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

	// Train in batches
	batchSize := 1000
	var tree *models.TreeNode
	for i := 0; i < len(indices); i += batchSize {
		end := i + batchSize
		if end > len(indices) {
			end = len(indices)
		}
		batchIndices := indices[i:end]

		if tree == nil {
			tree = buildTreeNode(batchIndices, features, targetCol, 0)
		}

		// Save the model incrementally
		if err := utils.SaveModel(tree); err != nil {
			return nil, fmt.Errorf("failed to save model: %w", err)
		}
	}

	fmt.Println("Tree building complete")
	return tree, nil
}

func buildTreeNode(indices []int, features []string, targetCol string, depth int) *models.TreeNode {
	// Create a leaf node if:
	// 1. Maximum depth reached
	// 2. Not enough samples to split
	// 3. All samples have the same target value
	if depth >= MaxDepth || len(indices) <= MinSamplesLeaf || CalculateEntropy(indices, targetCol) == 0 {
		prediction := MostCommonTarget(indices, targetCol)
		return &models.TreeNode{
			IsLeaf:     true,
			Prediction: prediction,
		}
	}

	bestSplit := FindBestSplit(indices, features, targetCol)

	// If no good split is found, create a leaf node
	if bestSplit.GainRatio < MinInfoGain {
		prediction := MostCommonTarget(indices, targetCol)
		return &models.TreeNode{
			IsLeaf:     true,
			Prediction: prediction,
		}
	}

	// Create a decision node
	node := &models.TreeNode{
		IsLeaf:     false,
		Feature:    bestSplit.Feature,
		SplitType:  bestSplit.SplitType,
		SplitValue: bestSplit.SplitValue,
	}

	// Split based on feature type
	if bestSplit.SplitType == "categorical" {
		// For categorical features, create a child for each value
		node.Children = make(map[string]*models.TreeNode)

		var wg sync.WaitGroup
		var mutex sync.Mutex

		for value, subIndices := range bestSplit.SplitIndices {
			if len(subIndices) == 0 {
				continue
			}

			wg.Add(1)
			go func(value string, subIndices []int) {
				defer wg.Done()
				childNode := buildTreeNode(subIndices, features, targetCol, depth+1)

				mutex.Lock()
				node.Children[value] = childNode
				mutex.Unlock()
			}(value, subIndices)
		}

		wg.Wait()

		// If no children were created, make it a leaf node
		if len(node.Children) == 0 {
			node.IsLeaf = true
			node.Prediction = MostCommonTarget(indices, targetCol)
			node.Children = nil
		}
	} else {
		// For numerical features, create left and right children
		if len(bestSplit.LeftIndices) > 0 {
			node.Left = buildTreeNode(bestSplit.LeftIndices, features, targetCol, depth+1)
		}

		if len(bestSplit.RightIndices) > 0 {
			node.Right = buildTreeNode(bestSplit.RightIndices, features, targetCol, depth+1)
		}

		// If both children are the same leaf, merge them
		if node.Left != nil && node.Right != nil &&
			node.Left.IsLeaf && node.Right.IsLeaf &&
			fmt.Sprintf("%v", node.Left.Prediction) == fmt.Sprintf("%v", node.Right.Prediction) {
			node.IsLeaf = true
			node.Prediction = node.Left.Prediction
			node.Left = nil
			node.Right = nil
		}
	}
	return node
}
