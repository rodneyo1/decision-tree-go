package algorithm

import (
	"dt/models"
	"math"
	"reflect"
	"testing"
)

func TestCalculateEntropy(t *testing.T) {
	// Set up test data
	models.Records = []map[string]interface{}{
		{"result": "yes"},
		{"result": "no"},
		{"result": "yes"},
		{"result": "no"},
		{"result": "yes"},
		{"result": "yes"},
		{"result": "no"},
		{"result": "yes"},
	}

	type args struct {
		indices   []int
		targetCol string
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "Perfect split (all same value)",
			args: args{
				indices:   []int{0, 2, 4, 5, 7}, // all "yes"
				targetCol: "result",
			},
			want: 0.0,
		},
		{
			name: "Maximum entropy (50-50 split)",
			args: args{
				indices:   []int{0, 1}, // one "yes", one "no"
				targetCol: "result",
			},
			want: 1.0,
		},
		{
			name: "Mixed entropy",
			args: args{
				indices:   []int{0, 1, 2, 3, 4, 5, 6, 7}, // 5 yes, 3 no
				targetCol: "result",
			},
			want: 0.9544340029249649, // -((5/8)*log2(5/8) + (3/8)*log2(3/8))
		},
		{
			name: "Empty set",
			args: args{
				indices:   []int{},
				targetCol: "result",
			},
			want: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateEntropy(tt.args.indices, tt.args.targetCol)
			if math.Abs(got-tt.want) > 1e-10 {
				t.Errorf("CalculateEntropy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMostCommonTarget(t *testing.T) {

    models.Records = []map[string]interface{}{
        {"result": "yes"},        // 0
        {"result": "no"},         // 1
        {"result": "yes"},        // 2
        {"result": "no"},         // 3
        {"result": "yes"},        // 4
        {"result": "yes"},        // 5
        {"result": "maybe"},      // 6
        {"result": "maybe"},      // 7
        {"result": "yes"},        // 8
        {"result": 1},            // 9
        {"result": 1},            // 10
        {"result": 2},            // 11
    }

    type args struct {
        indices   []int
        targetCol string
    }
    tests := []struct {
        name string
        args args
        want interface{}
    }{
        {
            name: "Clear majority",
            args: args{
                indices:   []int{0, 1, 2, 4, 5, 8},  // 5 yes, 1 no
                targetCol: "result",
            },
            want: "yes",
        },
        {
            name: "Empty set",
            args: args{
                indices:   []int{},
                targetCol: "result",
            },
            want: nil,
        },
        {
            name: "Single element",
            args: args{
                indices:   []int{1},  // just "no"
                targetCol: "result",
            },
            want: "no",
        },
        {
            name: "Tie returns first encountered",
            args: args{
                indices:   []int{6, 7, 0, 1},  // 2 maybe, 1 yes, 1 no
                targetCol: "result",
            },
            want: "maybe",
        },
        {
            name: "Numeric values",
            args: args{
                indices:   []int{9, 10, 11},  // 2 ones, 1 two
                targetCol: "result",
            },
            want: 1,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := MostCommonTarget(tt.args.indices, tt.args.targetCol); !reflect.DeepEqual(got, tt.want) {
                t.Errorf("MostCommonTarget() = %v, want %v", got, tt.want)
            }
        })
    }
}

// Mock data for testing
func setupMockData() {
	models.Records = []map[string]interface{}{
		{"feature": "A", "target": "Yes"},
		{"feature": "A", "target": "No"},
		{"feature": "B", "target": "Yes"},
		{"feature": "B", "target": "Yes"},
		{"feature": "C", "target": "No"},
	}

	models.FeatureTypes = map[string]string{
		"feature": "categorical",
	}
}

// Test FindBestSplit function
func TestFindBestSplit(t *testing.T) {
	setupMockData()
	indices := []int{0, 1, 2, 3, 4}
	features := []string{"feature"}
	targetCol := "target"

	split := FindBestSplit(indices, features, targetCol)

	if split.Feature != "feature" {
		t.Errorf("Expected best split feature to be 'feature', got %v", split.Feature)
	}

	if split.GainRatio <= 0 {
		t.Errorf("Expected positive gain ratio, got %f", split.GainRatio)
	}
}