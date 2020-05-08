package main

import (
	"bufio"
	"fmt"
	"os"
)

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

		var myPacIds []int

		for i := 0; i < visiblePacCount; i++ {
			// pacID: pac number (unique within a team)
			// mine: true if this pac is yours
			// x: position in the grid
			// y: position in the grid
			// typeId: unused in wood leagues
			// speedTurnsLeft: unused in wood leagues
			// abilityCooldown: unused in wood leagues
			var pacID int
			var mine bool
			var _mine int
			var x, y int
			var typeID string
			var speedTurnsLeft, abilityCooldown int
			scanner.Scan()
			fmt.Sscan(scanner.Text(), &pacID, &_mine, &x, &y, &typeID, &speedTurnsLeft, &abilityCooldown)
			mine = _mine != 0

			if mine {
				myPacIds = append(myPacIds, pacID)
			}
		}
		// visiblePelletCount: all pellets in sight
		var visiblePelletCount int
		scanner.Scan()
		fmt.Sscan(scanner.Text(), &visiblePelletCount)

		type Coord struct {
			X, Y int
		}
		var moves []Coord

		for i := 0; i < visiblePelletCount; i++ {
			// value: amount of points this pellet is worth
			var x, y, value int
			scanner.Scan()
			fmt.Sscan(scanner.Text(), &x, &y, &value)

			moves = append(moves, Coord{x, y})
		}

		// fmt.Fprintln(os.Stderr, "Debug messages...")
		var cmd string
		for i, myPacID := range myPacIds {
			move := moves[len(moves)/len(myPacIds)*i]
			if i > 0 {
				cmd = cmd + " | "
			}
			cmd = cmd + fmt.Sprintf("MOVE %v %v %v", myPacID, move.X, move.Y) // MOVE <pacId> <x> <y>
		}

		fmt.Println(cmd)
	}
}
