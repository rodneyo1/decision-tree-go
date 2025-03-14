package models

import (
	"testing"
	"time"
)

func TestCompareValues(t *testing.T) {
	tests := []struct {
		a, b     interface{}
		expected int
	}{
		{nil, nil, 0},
		{nil, 1, -1},
		{1, nil, 1},
		{1, 1, 0},
		{1, 2, -1},
		{2, 1, 1},
		{1, 1.0, 0},
		{1.0, 1, 0},
		{1.0, 2.0, -1},
		{2.0, 1.0, 1},
		{"a", "a", 0},
		{"a", "b", -1},
		{"b", "a", 1},
		{true, true, 0},
		{false, true, -1},
		{true, false, 1},
		{time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC), time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC), 0},
		{time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC), time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC), -1},
		{time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC), time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC), 1},
	}

	for _, test := range tests {
		result := CompareValues(test.a, test.b)
		if result != test.expected {
			t.Errorf("CompareValues(%v, %v) = %d; expected %d", test.a, test.b, result, test.expected)
		}
	}
}
