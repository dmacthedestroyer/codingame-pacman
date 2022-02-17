package main

import (
	"fmt"
	"testing"
)

func TestSortCoordsByProximity(t *testing.T) {
	pos := Coord{10, 10}
	tests := []struct {
		coords   []Coord
		pellets  map[Coord]int
		expected Coord
	}{
		{[]Coord{{10, 11}, {10, 12}}, map[Coord]int{}, Coord{10, 11}},
		{[]Coord{{10, 12}, {10, 11}}, map[Coord]int{}, Coord{10, 11}},
		{[]Coord{{13, 13}, {9, 9}}, map[Coord]int{}, Coord{9, 9}},
		{[]Coord{{10, 11}, {10, 5}}, map[Coord]int{{10, 11}: 1, {10, 5}: 10}, Coord{10, 5}},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			sortCoords(tt.coords, pos, tt.pellets)
			if tt.coords[0] != tt.expected {
				t.Errorf("expected first element to be %v, but was %v", tt.expected, tt.coords[0])
			}
		})
	}
}
