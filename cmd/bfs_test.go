package main

import "testing"

func TestEnemiesWithinRange(t *testing.T) {
	gameMap := BuildGameMap(`
## ##
#   #
#####`)

	pacsByPos := map[Coord]Pac{
		{x: 1, y: 1}: {mine: true},
		{x: 3, y: 1}: {mine: false},
	}

	expected := 1
	actual := enemiesWithinRange(gameMap, pacsByPos, Coord{1, 1}, 2)
	if expected != len(actual) {
		t.Errorf("expected %v enemies, but got %v: %v", expected, len(actual), actual)
	}
}
