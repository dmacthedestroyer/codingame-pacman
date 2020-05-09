package main

import (
	"fmt"
	"testing"
)

func TestSortPelletsByProximity(t *testing.T) {
	pac := Pac{pos: Coord{10, 10}}
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
			pellets := make([]Pellet, len(tt.coords))
			for i, c := range tt.coords {
				pellets[i] = Pellet{c, 1}
			}
			SortPelletsByProximity(pellets, pac)
			if pellets[0].pos != tt.expected {
				t.Errorf("expected first element to be %v, but was %v", tt.expected, pellets[0])
			}
		})
	}
}
