package main

import (
	"fmt"
	"testing"
)

func TestSortCoordsByProximity(t *testing.T) {
	pos := Coord{10, 10}
	tests := []struct {
		coords   []Coord
		expected Coord
	}{
		{[]Coord{{10, 11}, {10, 12}}, Coord{10, 11}},
		{[]Coord{{10, 12}, {10, 11}}, Coord{10, 11}},
		{[]Coord{{13, 13}, {9, 9}}, Coord{9, 9}},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			SortCoordsByProximity(tt.coords, pos)
			if tt.coords[0] != tt.expected {
				t.Errorf("expected first element to be %v, but was %v", tt.expected, tt.coords[0])
			}
		})
	}
}
