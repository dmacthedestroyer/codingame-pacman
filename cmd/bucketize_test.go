package main

import (
	"fmt"
	"testing"
)

func TestBucketize(t *testing.T) {
	tests := []struct{ x, numPacs, width, expected int }{
		{0, 2, 10, 0},
		{1, 2, 10, 0},
		{2, 2, 10, 0},
		{3, 2, 10, 0},
		{4, 2, 10, 0},
		{5, 2, 10, 1},
		{6, 2, 10, 1},
		{7, 2, 10, 1},
		{8, 2, 10, 1},
		{9, 2, 10, 1},

		{0, 3, 10, 0},
		{1, 3, 10, 0},
		{2, 3, 10, 0},
		{3, 3, 10, 1},
		{4, 3, 10, 1},
		{5, 3, 10, 1},
		{6, 3, 10, 2},
		{7, 3, 10, 2},
		{8, 3, 10, 2},
		{9, 3, 10, 2},

		{33, 4, 35, 3},
	}
	for _, tt := range tests {
		testName := fmt.Sprintf("bucketize(%v, %v, %v) = %v", tt.x, tt.numPacs, tt.width, tt.expected)
		t.Run(testName, func(t *testing.T) {
			if actual := bucketize(tt.x, tt.numPacs, tt.width); actual != tt.expected {
				t.Errorf("expected %v, but got %v", tt.expected, actual)
			}
		})
	}
}
