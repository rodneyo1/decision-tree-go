package models

import "fmt"

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
