package algorithm

import (
	"sync"

	"decision-tree/models"
)

// Predict makes predictions for all records in the dataset
func Predict(tree *models.TreeNode) []interface{} {
	predictions := make([]interface{}, len(models.Records))

	// Use goroutines for parallel prediction
	var wg sync.WaitGroup
	workers := 4 // Number of worker goroutines
	batchSize := (len(models.Records) + workers - 1) / workers

	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			start := workerID * batchSize
			end := (workerID + 1) * batchSize
			if end > len(models.Records) {
				end = len(models.Records)
			}

			for i := start; i < end; i++ {
				predictions[i] = predictRecord(models.Records[i], tree)
			}
		}(w)
	}

	wg.Wait()
	return predictions
}

// predictRecord makes a prediction for a single record
func predictRecord(record map[string]interface{}, node *models.TreeNode) interface{} {
	// If it's a leaf node, return the prediction
	if node.IsLeaf {
		return node.Prediction
	}
	// Get the feature value
	featureValue := record[node.Feature]

	// Handle missing values (null) by going to the majority branch
	if featureValue == nil {
		if node.SplitType == "categorical" {
			// Find the child with the most examples (first key as fallback)
			var bestChild *models.TreeNode
			maxSamples := -1

			for _, child := range node.Children {
				if bestChild == nil || estimateNodeSize(child) > maxSamples {
					bestChild = child
					maxSamples = estimateNodeSize(child)
				}
			}

			if bestChild != nil {
				return predictRecord(record, bestChild)
			}
			return node.Prediction // Fallback to current node's prediction
		} else {
			// For numerical splits, go to the side with more samples
			leftSize := estimateNodeSize(node.Left)
			rightSize := estimateNodeSize(node.Right)

			if leftSize >= rightSize && node.Left != nil {
				return predictRecord(record, node.Left)
			} else if node.Right != nil {
				return predictRecord(record, node.Right)
			}
			return node.Prediction
		}
	}
	// Split based on feature type
	if node.SplitType == "categorical" {
		// For categorical features, find the matching child
		valueKey := models.GetValueKey(featureValue)
		if child, ok := node.Children[valueKey]; ok {
			return predictRecord(record, child)
		}

		// If no matching child, use the most common child
		var bestChild *models.TreeNode
		maxSamples := -1

		for _, child := range node.Children {
			if bestChild == nil || estimateNodeSize(child) > maxSamples {
				bestChild = child
				maxSamples = estimateNodeSize(child)
			}
		}

		if bestChild != nil {
			return predictRecord(record, bestChild)
		}
		return node.Prediction
	} else {
		// For numerical features, compare with the threshold
		if models.CompareValues(featureValue, node.SplitValue) < 0 {
			if node.Left != nil {
				return predictRecord(record, node.Left)
			}
		} else {
			if node.Right != nil {
				return predictRecord(record, node.Right)
			}
		}
		return node.Prediction
	}
}

// estimateNodeSize estimates the number of samples in a node
func estimateNodeSize(node *models.TreeNode) int {
	if node == nil {
		return 0
	}

	if node.IsLeaf {
		return 1
	}

	size := 0
	if node.SplitType == "categorical" {
		for _, child := range node.Children {
			size += estimateNodeSize(child)
		}
	} else {
		// For numerical splits
		size += estimateNodeSize(node.Left)
		size += estimateNodeSize(node.Right)
	}

	return size
}
