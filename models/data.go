package models

import (
	"fmt"
	"time"
)

var (
	Records      []map[string]interface{}
	Columns      []string
	FeatureTypes map[string]string
	TargetValues map[interface{}]int
	TargetType   string
	TargetColumn string
)

type TreeNode struct {
	IsLeaf     bool                 `json:"is_leaf"`
	Prediction interface{}          `json:"prediction,omitempty"`
	Feature    string               `json:"feature,omitempty"`
	SplitValue interface{}          `json:"split_value,omitempty"`
	SplitType  string               `json:"split_type,omitempty"` // "categorical" or "numerical"
	Children   map[string]*TreeNode `json:"children,omitempty"`   // For categorical features
	Left       *TreeNode            `json:"left,omitempty"`       // For numerical features (< threshold)
	Right      *TreeNode            `json:"right,omitempty"`      // For numerical features (>= threshold)
}

func GetValueKey(val interface{}) string {
	if val == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", val)
}

type SplitCriteria struct {
	Feature      string
	SplitValue   interface{}
	SplitType    string
	InfoGain     float64
	GainRatio    float64
	SplitIndices map[string][]int // For categorical splits
	LeftIndices  []int            // For numerical splits (<)
	RightIndices []int            // For numerical splits (>=)
}

// ModelData represents the serializable model structure
type ModelData struct {
	Tree         *TreeNode         `json:"tree"`
	FeatureTypes map[string]string `json:"feature_types"`
	TargetColumn string            `json:"target_column"`
	TargetType   string            `json:"target_type"`
	Columns      []string          `json:"columns"`
}

// CompareValues compares two values based on their types
func CompareValues(a, b interface{}) int {
	// Handle nil values
	if a == nil && b == nil {
		return 0
	}
	if a == nil {
		return -1
	}
	if b == nil {
		return 1
	}

	// Compare based on types
	switch va := a.(type) {
	case int:
		switch vb := b.(type) {
		case int:
			if va < vb {
				return -1
			} else if va > vb {
				return 1
			}
			return 0
		case float64:
			fa := float64(va)
			if fa < vb {
				return -1
			} else if fa > vb {
				return 1
			}
			return 0
		}
	case float64:
		switch vb := b.(type) {
		case int:
			fb := float64(vb)
			if va < fb {
				return -1
			} else if va > fb {
				return 1
			}
			return 0
		case float64:
			if va < vb {
				return -1
			} else if va > vb {
				return 1
			}
			return 0
		}
	case string:
		if s, ok := b.(string); ok {
			if va < s {
				return -1
			} else if va > s {
				return 1
			}
			return 0
		}
	case bool:
		if bo, ok := b.(bool); ok {
			if !va && bo {
				return -1
			} else if va && !bo {
				return 1
			}
			return 0
		}
	case time.Time:
		if t, ok := b.(time.Time); ok {
			if va.Before(t) {
				return -1
			} else if va.After(t) {
				return 1
			}
			return 0
		}
	}

	// Fallback: convert to strings and compare
	aStr := fmt.Sprintf("%v", a)
	bStr := fmt.Sprintf("%v", b)
	if aStr < bStr {
		return -1
	} else if aStr > bStr {
		return 1
	}
	return 0
}

// BinEdges stores bin edges for each numerical feature
var BinEdges map[string][]float64

// GetBinnedValue returns the binned value for a numerical feature
func GetBinnedValue(feature string, value interface{}) interface{} {
	if value == nil {
		return nil
	}
	
	// Skip binning for non-numerical features
	if FeatureTypes[feature] != "numerical" {
		return value
	}
	
	// Check if binning is enabled for this feature
	edges := BinEdges[feature]
	if edges == nil {
		return value // No binning for this feature
	}
	
	var floatValue float64
	switch v := value.(type) {
	case int:
		floatValue = float64(v)
	case float64:
		floatValue = v
	default:
		return value // Cannot bin this value type
	}
	
	// Find the bin
	for i := 0; i < len(edges)-1; i++ {
		if floatValue >= edges[i] && floatValue < edges[i+1] {
			// Return the bin index or the bin midpoint
			// Option 1: Return bin index as string
			// return fmt.Sprintf("bin_%d", i)
			
			// Option 2: Return midpoint of the bin
			return (edges[i] + edges[i+1]) / 2
		}
	}
	
	// Handle edge case: value equals the maximum edge
	if floatValue == edges[len(edges)-1] {
		return (edges[len(edges)-2] + edges[len(edges)-1]) / 2
	}
	
	return value // Value outside bin ranges
}