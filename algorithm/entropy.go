package algorithm

import (
	"math"

	"decision-tree/models"
)

// Calculate entropy of a set of indices
func CalculateEntropy(indices []int, targetCol string) float64 {
	if len(indices) == 0 {
		return 0
	}

	// Count occurrences of each target value
	valueCount := make(map[string]int)
	for _, idx := range indices {
		value := models.Records[idx][targetCol]
		key := models.GetValueKey(value)
		valueCount[key]++
	}

	// Calculate entropy
	entropy := 0.0
	for _, count := range valueCount {
		prob := float64(count) / float64(len(indices))
		entropy -= prob * math.Log2(prob)
	}

	return entropy
}
