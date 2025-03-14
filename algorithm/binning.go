package algorithm

import (
	"math"
	"sort"
	"dt/models"
)

// BinningOptions defines how binning should be performed
type BinningOptions struct {
	Method      string // "equal_width", "equal_frequency", or "custom"
	NumBins     int    // Number of bins to create
	CustomEdges map[string][]float64 // Custom bin edges for specific features
}

// Binning applies binning to numerical features in the dataset
func ApplyBinning(options BinningOptions) error {
	// Create a map to store bin edges for each feature
	binEdges := make(map[string][]float64)
	
	// Process each feature
	for feature, featureType := range models.FeatureTypes {
		// Only bin numerical features
		if featureType != "numerical" {
			continue
		}
		
		// Check if custom edges are provided for this feature
		if edges, ok := options.CustomEdges[feature]; ok {
			binEdges[feature] = edges
			continue
		}
		
		// Calculate bin edges based on the method
		switch options.Method {
		case "equal_width":
			binEdges[feature] = calculateEqualWidthBins(feature, options.NumBins)
		case "equal_frequency":
			binEdges[feature] = calculateEqualFrequencyBins(feature, options.NumBins)
		default:
			// Default to equal width
			binEdges[feature] = calculateEqualWidthBins(feature, options.NumBins)
		}
	}
	
	// Store bin edges in models package for use during tree building and prediction
	models.BinEdges = binEdges
	
	return nil
}

// Calculate bin edges using equal width method
func calculateEqualWidthBins(feature string, numBins int) []float64 {
	// Collect all values for the feature
	values := make([]float64, 0)
	for _, record := range models.Records {
		if val := record[feature]; val != nil {
			switch v := val.(type) {
			case int:
				values = append(values, float64(v))
			case float64:
				values = append(values, v)
			}
		}
	}
	
	if len(values) == 0 {
		return nil
	}
	
	// Find min and max
	min, max := values[0], values[0]
	for _, v := range values {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	
	// Create bin edges
	edges := make([]float64, numBins+1)
	width := (max - min) / float64(numBins)
	
	for i := 0; i <= numBins; i++ {
		edges[i] = min + float64(i)*width
	}
	
	return edges
}

// Calculate bin edges using equal frequency method
func calculateEqualFrequencyBins(feature string, numBins int) []float64 {
	// Collect and sort all values
	values := make([]float64, 0)
	for _, record := range models.Records {
		if val := record[feature]; val != nil {
			switch v := val.(type) {
			case int:
				values = append(values, float64(v))
			case float64:
				values = append(values, v)
			}
		}
	}
	
	if len(values) == 0 {
		return nil
	}
	
	sort.Float64s(values)
	
	// Create bin edges
	edges := make([]float64, numBins+1)
	edges[0] = values[0] - 0.0001 // Slightly below minimum
	edges[numBins] = values[len(values)-1] + 0.0001 // Slightly above maximum
	
	// Calculate quantile positions
	for i := 1; i < numBins; i++ {
		position := int(math.Round(float64(i) * float64(len(values)) / float64(numBins)))
		if position >= len(values) {
			position = len(values) - 1
		}
		edges[i] = values[position]
	}
	
	return edges
}

// GetBin returns the bin index for a value
func GetBin(feature string, value interface{}) int {
	edges := models.BinEdges[feature]
	if edges == nil {
		return -1 // No binning for this feature
	}
	
	var floatValue float64
	switch v := value.(type) {
	case int:
		floatValue = float64(v)
	case float64:
		floatValue = v
	default:
		return -1 // Cannot bin this value type
	}
	
	// Find the bin
	for i := 0; i < len(edges)-1; i++ {
		if floatValue >= edges[i] && floatValue < edges[i+1] {
			return i
		}
	}
	
	// Handle edge case: value equals the maximum edge
	if floatValue == edges[len(edges)-1] {
		return len(edges) - 2
	}
	
	return -1 // Value outside bin ranges
}