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

// GameState represents a snapshot of the game at a point in time
type GameState struct {
	scores  []int
	pacs    []Pac
	pellets []Pellet
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

	for i := 0; i < height; i++ {
		scanner.Scan()
		//row := scanner.Text() // one line of the grid: space " " is floor, pound "#" is wall
	}
	for {
		var myScore, opponentScore int
		scanner.Scan()
		fmt.Sscan(scanner.Text(), &myScore, &opponentScore)
		// visiblePacCount: all your pacs and enemy pacs in sight
		var visiblePacCount int
		scanner.Scan()
		fmt.Sscan(scanner.Text(), &visiblePacCount)

		var pacs []Pac

		for i := 0; i < visiblePacCount; i++ {
			// pacID: pac number (unique within a team)
			// player: 0 if this pac is yours
			// x: position in the grid
			// y: position in the grid
			// typeId: unused in wood leagues
			// speedTurnsLeft: unused in wood leagues
			// abilityCooldown: unused in wood leagues
			var pacID int
			var player int
			var x, y int
			var typeID string
			var speedTurnsLeft, abilityCooldown int
			scanner.Scan()
			fmt.Sscan(scanner.Text(), &pacID, &player, &x, &y, &typeID, &speedTurnsLeft, &abilityCooldown)

			pacs = append(pacs, Pac{pacID, player, x, y, typeID, speedTurnsLeft, abilityCooldown})
		}
		// visiblePelletCount: all pellets in sight
		var visiblePelletCount int
		scanner.Scan()
		fmt.Sscan(scanner.Text(), &visiblePelletCount)

		var pellets []Pellet

		for i := 0; i < visiblePelletCount; i++ {
			// value: amount of points this pellet is worth
			var x, y, value int
			scanner.Scan()
			fmt.Sscan(scanner.Text(), &x, &y, &value)

			pellets = append(pellets, Pellet{x, y, value})
		}

		gameState := GameState{[]int{myScore, opponentScore}, pacs, pellets}

		// fmt.Fprintln(os.Stderr, "Debug messages...")
		var cmd string
		var myPacs []Pac
		for _, pac := range gameState.pacs {
			if pac.player == 1 {
				fmt.Fprintln(os.Stderr, pac)
				myPacs = append(myPacs, pac)
			}
		}
		for i, pac := range myPacs {
			pellet := gameState.pellets[len(pellets)/len(myPacs)*i]
			if i > 0 {
				cmd = cmd + " | "
			}
			cmd = cmd + fmt.Sprint("MOVE ", pac.id, pellet.x, pellet.y, pac.id, " ", pac.typeID) // MOVE <pacId> <x> <y>
		}

		fmt.Fprintln(os.Stderr, cmd)
		fmt.Println(cmd)
	}
}
