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

// Calculate the most common target value for a set of indices
func MostCommonTarget(indices []int, targetCol string) interface{} {
	if len(indices) == 0 {
		return nil
	}

	// Count occurrences of each target value
	valueCount := make(map[string]int)
	valueMap := make(map[string]interface{})

	for _, idx := range indices {
		value := models.Records[idx][targetCol]
		key := models.GetValueKey(value)
		valueCount[key]++
		valueMap[key] = value
	}

	// Find the most common value
	maxCount := 0
	var maxKey string
	for key, count := range valueCount {
		if count > maxCount {
			maxCount = count
			maxKey = key
		}
	}

	return valueMap[maxKey]
}

func FindBestSplit(indices []int, features []string, targetCol string) models.SplitCriteria {
	baseEntropy := CalculateEntropy(indices, targetCol)
	bestSplit := models.SplitCriteria{
		InfoGain:  -1,
		GainRatio: -1,
	}

	// If entropy is 0, no need to split
	if baseEntropy == 0 {
		return bestSplit
	}

	for _, feature := range features {
		if feature == targetCol {
			continue
		}

		featureType := models.FeatureTypes[feature]
		if featureType == "categorical" {
		
	}

	return bestSplit
}

