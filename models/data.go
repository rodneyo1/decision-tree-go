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
