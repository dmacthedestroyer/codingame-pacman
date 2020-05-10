package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestBuildGameMap(t *testing.T) {
	gm := BuildGameMap(`
#####
# # #
# # #
# # #
#   #
#####`)
	if expected, actual := 5, gm.width; actual != expected {
		t.Errorf("expected width %v, but got %v", expected, actual)
	}
	if expected, actual := 6, gm.height; actual != expected {
		t.Errorf("expected height %v, but got %v", expected, actual)
	}

	tests := []struct {
		x, y     int
		expected rune
	}{
		{0, 0, '#'},
		{0, 5, '#'},
		{1, 1, ' '},
		{3, 1, ' '},
		{4, 5, '#'},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("(%v,%v) = '%v'", tt.x, tt.y, tt.expected), func(t *testing.T) {
			if expected, actual := tt.expected, gm.GetCell(Coord{tt.x, tt.y}).value; expected != actual {
				t.Errorf("expected '%c', but got '%c'", expected, actual)
			}
		})
	}
	if expected, actual := ' ', gm.GetCell(Coord{1, 1}).value; actual != expected {
		t.Errorf("")
	}
}

func TestGetCoord(t *testing.T) {
	gm := GameMap{width: 10, height: 5}

	tests := []struct {
		pos      int
		expected Coord
	}{
		{0, Coord{0, 0}},
		{1, Coord{1, 0}},
		{10, Coord{0, 1}},
		{25, Coord{5, 2}},
		{49, Coord{9, 4}},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.pos), func(t *testing.T) {
			if expected, actual := tt.expected, gm.GetCoord(tt.pos); expected != actual {
				t.Errorf("expected %v, but got %v", expected, actual)
			}
		})
	}
}

func TestCoordAndAbsolutePositionAreIsomorphic(t *testing.T) {
	gm := GameMap{width: 12, height: 14}
	for pos := 0; pos < gm.width*gm.height; pos++ {
		t.Run(fmt.Sprint(pos), func(t *testing.T) {
			coord := gm.GetCoord(pos)
			backtoPos := gm.GetAbsolutePosition(coord)
			if pos != backtoPos {
				t.Errorf("%v -> %v was %v, but expected %v", pos, coord, backtoPos, pos)
			}
		})
	}
}

func TestVisibleCells(t *testing.T) {
	gm := BuildGameMap(`
### #
# ###
   # 
#   #
### #`)

	tests := []struct {
		x, y     int
		expected []Coord
	}{
		{1, 1, []Coord{{1, 1}, {1, 2}, {1, 3}}},
		{1, 2, []Coord{{1, 2}, {1, 1}, {1, 3}, {0, 2}, {4, 2}, {2, 2}}},
		{3, 3, []Coord{{3, 3}, {3, 4}, {3, 0}, {2, 3}, {1, 3}}},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("(%v,%v)", tt.x, tt.y), func(t *testing.T) {
			if expected, actual := tt.expected, gm.VisibleCells(Coord{tt.x, tt.y}); !reflect.DeepEqual(expected, actual) {
				t.Errorf("expected %v, but got %v", expected, actual)
			}
		})
	}
}
