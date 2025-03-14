package algorithm

import (
	"math"
	"sort"

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
	}

	return bestSplit
}

// Find the best split for a categorical feature
func findCategoricalSplit(indices []int, feature string, targetCol string, baseEntropy float64) models.SplitCriteria {
	// Group indices by feature value
	valueIndices := make(map[string][]int)
	for _, idx := range indices {
		value := models.Records[idx][feature]
		key := models.GetValueKey(value)
		valueIndices[key] = append(valueIndices[key], idx)
	}

	// Calculate weighted entropy
	weightedEntropy := 0.0
	splitInfo := 0.0

	for _, subIndices := range valueIndices {
		prob := float64(len(subIndices)) / float64(len(indices))
		weightedEntropy += prob * CalculateEntropy(subIndices, targetCol)
		splitInfo -= prob * math.Log2(prob)
	}

	// Calculate information gain and gain ratio
	infoGain := baseEntropy - weightedEntropy
	gainRatio := 0.0
	if splitInfo > 0 {
		gainRatio = infoGain / splitInfo
	}

	return models.SplitCriteria{
		Feature:      feature,
		SplitType:    "categorical",
		InfoGain:     infoGain,
		GainRatio:    gainRatio,
		SplitIndices: valueIndices,
	}
}

// Find the best split for a numerical feature
func findNumericalSplit(indices []int, feature string, targetCol string, baseEntropy float64) models.SplitCriteria {
	// Collect unique values
	values := make([]interface{}, 0)
	valuesMap := make(map[string]bool)

	for _, idx := range indices {
		value := models.Records[idx][feature]
		if value != nil {
			key := models.GetValueKey(value)
			if !valuesMap[key] {
				values = append(values, value)
				valuesMap[key] = true
			}
		}
	}

	// Sort values
	sort.Slice(values, func(i, j int) bool {
		return models.CompareValues(values[i], values[j]) < 0
	})

	// If no valid split points, return empty criteria
	if len(values) <= 1 {
		return models.SplitCriteria{
			Feature:   feature,
			SplitType: "numerical",
			InfoGain:  -1,
			GainRatio: -1,
		}
	}

	// Try each split point and find the best
	bestSplit := models.SplitCriteria{
		Feature:   feature,
		SplitType: "numerical",
		InfoGain:  -1,
		GainRatio: -1,
	}

	for i := 0; i < len(values)-1; i++ {
		a := values[i]
		b := values[i+1]

		// Calculate midpoint for threshold
		var threshold interface{}
		switch va := a.(type) {
		case int:
			switch vb := b.(type) {
			case int:
				threshold = (va + vb) / 2
			case float64:
				threshold = (float64(va) + vb) / 2
			default:
				continue
			}
		case float64:
			switch vb := b.(type) {
			case int:
				threshold = (va + float64(vb)) / 2
			case float64:
				threshold = (va + vb) / 2
			default:
				continue
			}
		default:
			continue
		}

		// Evaluate the split
		leftIndices := make([]int, 0)
		rightIndices := make([]int, 0)

		for _, idx := range indices {
			value := models.Records[idx][feature]
			if value == nil || models.CompareValues(value, threshold) < 0 {
				leftIndices = append(leftIndices, idx)
			} else {
				rightIndices = append(rightIndices, idx)
			}
		}

		// Skip if all records end up in one branch
		if len(leftIndices) == 0 || len(rightIndices) == 0 {
			continue
		}

		// Calculate entropies and gain
		leftProb := float64(len(leftIndices)) / float64(len(indices))
		rightProb := float64(len(rightIndices)) / float64(len(indices))

		leftEntropy := CalculateEntropy(leftIndices, targetCol)
		rightEntropy := CalculateEntropy(rightIndices, targetCol)

		weightedEntropy := leftProb*leftEntropy + rightProb*rightEntropy
		infoGain := baseEntropy - weightedEntropy

		// Calculate split info for gain ratio
		splitInfo := -leftProb*math.Log2(leftProb) - rightProb*math.Log2(rightProb)
		gainRatio := 0.0
		if splitInfo > 0 {
			gainRatio = infoGain / splitInfo
		}

		// Update best split if this is better
		if gainRatio > bestSplit.GainRatio {
			bestSplit.SplitValue = threshold
			bestSplit.InfoGain = infoGain
			bestSplit.GainRatio = gainRatio
			bestSplit.LeftIndices = leftIndices
			bestSplit.RightIndices = rightIndices
		}
	}

	return bestSplit
}
