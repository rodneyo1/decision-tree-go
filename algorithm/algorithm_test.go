package algorithm

import (
	"dt/models"
	"dt/utils"
	"math"
	"os"
	"reflect"
	"testing"
)

// Setup a mock dataset for testing
func setupMockData1() {
	models.Columns = []string{"Feature1", "Feature2", "Target"}
	models.Records = []map[string]interface{}{
		{"Feature1": "A", "Feature2": 1.2, "Target": "Yes"},
		{"Feature1": "B", "Feature2": 2.4, "Target": "No"},
		{"Feature1": "A", "Feature2": 1.5, "Target": "Yes"},
		{"Feature1": "B", "Feature2": 2.1, "Target": "No"},
		{"Feature1": "A", "Feature2": 1.8, "Target": "Yes"},
	}
}

func TestBuildTree(t *testing.T) {
	setupMockData1()

	// Ensure OutputPtr is set to a valid test file
	tempFile := "test_model.json"
	utils.OutputPtr = &tempFile // Assign a valid file path

	// Run BuildTree
	tree, err := BuildTree("Target")
	if err != nil {
		t.Fatalf("BuildTree returned an error: %v", err)
	}

	if tree == nil {
		t.Fatal("BuildTree returned a nil tree")
	}

	if tree.IsLeaf && tree.Prediction == nil {
		t.Fatal("Leaf node has no prediction")
	}

	// Clean up the test output file
	os.Remove(tempFile)
}

func TestEstimateNodeSize(t *testing.T) {
	tests := []struct {
		name     string
		node     *models.TreeNode
		expected int
	}{
		{
			name:     "Nil node",
			node:     nil,
			expected: 0,
		},
		{
			name: "Leaf node",
			node: &models.TreeNode{
				IsLeaf: true,
			},
			expected: 1,
		},
		{
			name: "Categorical split with children",
			node: &models.TreeNode{
				SplitType: "categorical",
				Children: map[string]*models.TreeNode{
					"child1": {IsLeaf: true},
					"child2": {IsLeaf: true},
				},
			},
			expected: 2,
		},
		{
			name: "Numerical split with left and right children",
			node: &models.TreeNode{
				SplitType: "numerical",
				Left:      &models.TreeNode{IsLeaf: true},
				Right:     &models.TreeNode{IsLeaf: true},
			},
			expected: 2,
		},
		{
			name: "Mixed split types",
			node: &models.TreeNode{
				SplitType: "categorical",
				Children: map[string]*models.TreeNode{
					"child1": {
						SplitType: "numerical",
						Left:      &models.TreeNode{IsLeaf: true},
						Right:     &models.TreeNode{IsLeaf: true},
					},
					"child2": {IsLeaf: true},
				},
			},
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := estimateNodeSize(tt.node)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestPredict(t *testing.T) {
	tests := []struct {
		name     string
		tree     *models.TreeNode
		records  []map[string]interface{}
		expected []interface{}
	}{
		{
			name: "Single leaf node",
			tree: &models.TreeNode{
				IsLeaf:     true,
				Prediction: "A",
			},
			records: []map[string]interface{}{
				{"feature1": "value1"},
				{"feature2": "value2"},
			},
			expected: []interface{}{"A", "A"},
		},
		{
			name: "Categorical split",
			tree: &models.TreeNode{
				SplitType: "categorical",
				Feature:   "feature1",
				Children: map[string]*models.TreeNode{
					"value1": {IsLeaf: true, Prediction: "A"},
					"value2": {IsLeaf: true, Prediction: "B"},
				},
			},
			records: []map[string]interface{}{
				{"feature1": "value1"},
				{"feature1": "value2"},
			},
			expected: []interface{}{"A", "B"},
		},
		{
			name: "Numerical split",
			tree: &models.TreeNode{
				SplitType:  "numerical",
				Feature:    "feature1",
				SplitValue: 10,
				Left:       &models.TreeNode{IsLeaf: true, Prediction: "A"},
				Right:      &models.TreeNode{IsLeaf: true, Prediction: "B"},
			},
			records: []map[string]interface{}{
				{"feature1": 5},
				{"feature1": 15},
			},
			expected: []interface{}{"A", "B"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			models.Records = tt.records
			result := Predict(tt.tree)
			for i, prediction := range result {
				if prediction != tt.expected[i] {
					t.Errorf("expected %v, got %v", tt.expected[i], prediction)
				}
			}
		})
	}
}

func TestPredictRecord(t *testing.T) {
	tests := []struct {
		name     string
		record   map[string]interface{}
		node     *models.TreeNode
		expected interface{}
	}{
		{
			name: "Leaf node",
			record: map[string]interface{}{
				"feature1": "value1",
			},
			node: &models.TreeNode{
				IsLeaf:     true,
				Prediction: "A",
			},
			expected: "A",
		},
		{
			name: "Categorical split",
			record: map[string]interface{}{
				"feature1": "value1",
			},
			node: &models.TreeNode{
				SplitType: "categorical",
				Feature:   "feature1",
				Children: map[string]*models.TreeNode{
					"value1": {IsLeaf: true, Prediction: "A"},
					"value2": {IsLeaf: true, Prediction: "B"},
				},
			},
			expected: "A",
		},
		{
			name: "Numerical split",
			record: map[string]interface{}{
				"feature1": 5,
			},
			node: &models.TreeNode{
				SplitType:  "numerical",
				Feature:    "feature1",
				SplitValue: 10,
				Left:       &models.TreeNode{IsLeaf: true, Prediction: "A"},
				Right:      &models.TreeNode{IsLeaf: true, Prediction: "B"},
			},
			expected: "A",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := predictRecord(tt.record, tt.node)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

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
		{"result": "yes"},   // 0
		{"result": "no"},    // 1
		{"result": "yes"},   // 2
		{"result": "no"},    // 3
		{"result": "yes"},   // 4
		{"result": "yes"},   // 5
		{"result": "maybe"}, // 6
		{"result": "maybe"}, // 7
		{"result": "yes"},   // 8
		{"result": 1},       // 9
		{"result": 1},       // 10
		{"result": 2},       // 11
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
				indices:   []int{0, 1, 2, 4, 5, 8}, // 5 yes, 1 no
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
				indices:   []int{1}, // just "no"
				targetCol: "result",
			},
			want: "no",
		},
		{
			name: "Tie returns first encountered",
			args: args{
				indices:   []int{6, 7, 0, 1}, // 2 maybe, 1 yes, 1 no
				targetCol: "result",
			},
			want: "maybe",
		},
		{
			name: "Numeric values",
			args: args{
				indices:   []int{9, 10, 11}, // 2 ones, 1 two
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
func setupMockData2() {
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
	setupMockData2()
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

func Test_findNumericalSplit(t *testing.T) {
	// Set up test data
	models.Records = []map[string]interface{}{
		{"size": 1.0, "target": "Yes"},
		{"size": 2.0, "target": "No"},
		{"size": 3.0, "target": "Yes"},
		{"size": 4.0, "target": "No"},
		{"size": 2.0, "target": "Yes"},
	}

	models.FeatureTypes = map[string]string{
		"size": "numerical",
	}

	type args struct {
		indices     []int
		feature     string
		targetCol   string
		baseEntropy float64
	}

	tests := []struct {
		name string
		args args
		want models.SplitCriteria
	}{
		{
			name: "Single value (no split possible)",
			args: args{
				indices:     []int{0},
				feature:     "size",
				targetCol:   "target",
				baseEntropy: 0.0,
			},
			want: models.SplitCriteria{
				Feature:   "size",
				SplitType: "numerical",
				InfoGain:  -1,
				GainRatio: -1,
			},
		},
		{
			name: "Empty set",
			args: args{
				indices:     []int{},
				feature:     "size",
				targetCol:   "target",
				baseEntropy: 0.0,
			},
			want: models.SplitCriteria{
				Feature:   "size",
				SplitType: "numerical",
				InfoGain:  -1,
				GainRatio: -1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findNumericalSplit(tt.args.indices, tt.args.feature, tt.args.targetCol, tt.args.baseEntropy)

			// Check basic fields
			if got.Feature != tt.want.Feature ||
				got.SplitType != tt.want.SplitType ||
				(got.InfoGain > 0 && math.Abs(got.InfoGain-tt.want.InfoGain) > 1e-3) ||
				(got.GainRatio > 0 && math.Abs(got.GainRatio-tt.want.GainRatio) > 1e-3) {
				t.Errorf("findNumericalSplit() = %v, want %v", got, tt.want)
			}

			// Check split specific fields when a valid split is found
			if got.InfoGain > 0 {
				if got.SplitValue == nil {
					t.Error("Expected split value, got nil")
				} else if math.Abs(got.SplitValue.(float64)-tt.want.SplitValue.(float64)) > 1e-3 {
					t.Errorf("Split value = %v, want %v", got.SplitValue, tt.want.SplitValue)
				}

				// Check that indices are properly split
				if !reflect.DeepEqual(got.LeftIndices, tt.want.LeftIndices) ||
					!reflect.DeepEqual(got.RightIndices, tt.want.RightIndices) {
					t.Errorf("Indices split mismatch:\nGot: L=%v, R=%v\nWant: L=%v, R=%v",
						got.LeftIndices, got.RightIndices,
						tt.want.LeftIndices, tt.want.RightIndices)
				}
			}
		})
	}
}

func Test_findCategoricalSplit(t *testing.T) {
	// Set up test data
	models.Records = []map[string]interface{}{
		{"color": "red", "target": "yes"},   // 0
		{"color": "blue", "target": "no"},   // 1
		{"color": "red", "target": "yes"},   // 2
		{"color": "blue", "target": "no"},   // 3
		{"color": "green", "target": "yes"}, // 4
	}

	models.FeatureTypes = map[string]string{
		"color": "categorical",
	}

	type args struct {
		indices     []int
		feature     string
		targetCol   string
		baseEntropy float64
	}

	tests := []struct {
		name string
		args args
		want models.SplitCriteria
	}{
		{
			name: "Single value",
			args: args{
				indices:     []int{0},
				feature:     "color",
				targetCol:   "target",
				baseEntropy: 0.0,
			},
			want: models.SplitCriteria{
				Feature:   "color",
				SplitType: "categorical",
				InfoGain:  0.0,
				GainRatio: 0.0,
				SplitIndices: map[string][]int{
					"red": {0},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findCategoricalSplit(tt.args.indices, tt.args.feature, tt.args.targetCol, tt.args.baseEntropy)

			// Check basic fields
			if got.Feature != tt.want.Feature ||
				got.SplitType != tt.want.SplitType ||
				(got.InfoGain >= 0 && math.Abs(got.InfoGain-tt.want.InfoGain) > 1e-3) ||
				(got.GainRatio >= 0 && math.Abs(got.GainRatio-tt.want.GainRatio) > 1e-3) {
				t.Errorf("findCategoricalSplit() = %v, want %v", got, tt.want)
			}

			// Check split indices when a valid split exists
			if got.InfoGain >= 0 {
				if !reflect.DeepEqual(got.SplitIndices, tt.want.SplitIndices) {
					t.Errorf("Split indices mismatch:\nGot: %v\nWant: %v",
						got.SplitIndices, tt.want.SplitIndices)
				}
			}
		})
	}
}
