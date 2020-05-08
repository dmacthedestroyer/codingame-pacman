package main

import (
	"bufio"
	"fmt"
	"os"
)

// Pac represents a Pac man (or woman)
type Pac struct {
	id                              int
	player                          int
	x, y                            int
	typeID                          string
	speedTurnsLeft, abilityCooldown int
}

// Pellet represents a pellet with a location and point value
type Pellet struct {
	x, y, value int
}

// Cell represents a single wall or floor of the game area
type Cell struct {
	value rune
}

// GameMap represents the walls and floors of the game area
type GameMap struct {
	width, height int
	cells         []Cell
}

// GameState represents a snapshot of the game at a point in time
type GameState struct {
	gameMap        GameMap
	scores         []int
	visiblePacs    []Pac
	visiblePellets []Pellet
}

func debug(a ...interface{}) {
	fmt.Fprintln(os.Stderr, a...)
}

/**
 * Grab the pellets as fast as you can!
 **/
func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1000000), 1000000)

	// width: size of the grid
	// height: top left corner is (x=0, y=0)
	var width, height int
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &width, &height)
	var cells []Cell

	for i := 0; i < height; i++ {
		scanner.Scan()
		row := scanner.Text() // one line of the grid: space " " is floor, pound "#" is wall
		debug(row)
		for _, cellValue := range row {
			cells = append(cells, Cell{cellValue})
		}
	}

	gameMap := GameMap{width, height, cells}

	for {
		var myScore, opponentScore int
		scanner.Scan()
		fmt.Sscan(scanner.Text(), &myScore, &opponentScore)
		// visiblePacCount: all your pacs and enemy pacs in sight
		var visiblePacCount int
		scanner.Scan()
		fmt.Sscan(scanner.Text(), &visiblePacCount)

		var visiblePacs []Pac

		for i := 0; i < visiblePacCount; i++ {
			var pacID int
			var player int
			var x, y int
			var typeID string
			var speedTurnsLeft, abilityCooldown int
			scanner.Scan()
			fmt.Sscan(scanner.Text(), &pacID, &player, &x, &y, &typeID, &speedTurnsLeft, &abilityCooldown)

			visiblePacs = append(visiblePacs, Pac{pacID, player, x, y, typeID, speedTurnsLeft, abilityCooldown})
		}
		// visiblePelletCount: all pellets in sight
		var visiblePelletCount int
		scanner.Scan()
		fmt.Sscan(scanner.Text(), &visiblePelletCount)

		var visiblePellets []Pellet

		for i := 0; i < visiblePelletCount; i++ {
			// value: amount of points this pellet is worth
			var x, y, value int
			scanner.Scan()
			fmt.Sscan(scanner.Text(), &x, &y, &value)

			visiblePellets = append(visiblePellets, Pellet{x, y, value})
		}

		gameState := GameState{gameMap, []int{myScore, opponentScore}, visiblePacs, visiblePellets}

		// fmt.Fprintln(os.Stderr, "Debug messages...")
		var cmd string
		var myPacs []Pac
		// find all of my pacs from the visible collection
		for _, pac := range gameState.visiblePacs {
			if pac.player == 1 {
				myPacs = append(myPacs, pac)
			}
		}
		// equally divide the visible pellets amongst all of my pacs (even though the order of the pellets and the position of the pacs are meaningless)
		for i, pac := range myPacs {
			pellet := gameState.visiblePellets[len(visiblePellets)/len(myPacs)*i]
			if i > 0 {
				cmd = cmd + " | "
			}
			cmd = cmd + fmt.Sprint("MOVE ", pac.id, pellet.x, pellet.y, pac.id, " ", pac.typeID) // MOVE <pacId> <x> <y>
		}

		debug(cmd)
		fmt.Println(cmd)
	}
}
