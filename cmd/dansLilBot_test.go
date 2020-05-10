package main

import (
	"fmt"
	"testing"
)

func TestInit(t *testing.T) {
	bot := DansLilHeuristicBot{}

	bot.init(BuildGameMap(`
## ##
#   #
#####`))

	tests := []int{0, 0, 1, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0}
	if expected, actual := len(tests), len(bot.pelletValuesByPos); expected != actual {
		t.Errorf("expected length %v, but got %v", expected, actual)
	} else {
		for pos, tt := range tests {
			t.Run(fmt.Sprint(pos), func(t *testing.T) {
				if expected, actual := tt, bot.pelletValuesByPos[pos]; expected != actual {
					t.Errorf("expected position %v to have pellet value %v, but got %v", pos, expected, actual)
				}
			})
		}
	}
}

func TestUpdate(t *testing.T) {
	gameMap := BuildGameMap(`
#####
# ###
# ###
# ###
#####`)
	bot := DansLilHeuristicBot{}
	bot.init(gameMap)

	preCheckPositions := []Coord{{1, 2}, {1, 3}}
	for _, preCheckPosition := range preCheckPositions {
		if actual := bot.pelletValuesByPos[gameMap.GetAbsolutePosition(preCheckPosition)]; actual <= 0 {
			t.Errorf("prerequisite failed: expected pellet value >0 at position %v, but got value %v", preCheckPosition, actual)
			t.FailNow()
		}
	}

	pacCoord := Coord{1, 1}
	visiblePellets := []Pellet{{Coord{1, 3}, 1}}
	bot.update(GameData{gameMap: gameMap, visiblePacs: []Pac{{player: 1, pos: pacCoord}}, visiblePellets: visiblePellets})

	tests := []struct {
		pos           Coord
		expectedValue int
	}{
		{Coord{1, 2}, 0},
		{Coord{1, 3}, 1},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprint(tt), func(t *testing.T) {
			ttPos := gameMap.GetAbsolutePosition(tt.pos)
			if expected, actual := tt.expectedValue, bot.pelletValuesByPos[ttPos]; expected != actual {
				t.Errorf("expected %v, but got %v", expected, actual)
			}
		})
	}
}
