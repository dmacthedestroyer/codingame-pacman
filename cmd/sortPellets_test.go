package main

import (
	"fmt"
	"testing"
)

type coord struct{ x, y int }

func TestSortPelletsByProximity(t *testing.T) {
	pac := Pac{x: 10, y: 10}
	tests := []struct {
		coords   []coord
		expected coord
	}{
		{[]coord{{10, 11}, {10, 12}}, coord{10, 11}},
		{[]coord{{10, 12}, {10, 11}}, coord{10, 11}},
		{[]coord{{13, 13}, {9, 9}}, coord{9, 9}},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			pellets := make([]Pellet, len(tt.coords))
			for i, c := range tt.coords {
				pellets[i] = Pellet{c.x, c.y, 1}
			}
			SortPelletsByProximity(pellets, pac)
			if pellets[0].x != tt.expected.x && pellets[0].y != tt.expected.y {
				t.Errorf("expected first element to be %v, but was %v", tt.expected, pellets[0])
			}
		})
	}
}
