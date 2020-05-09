package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strings"
)

// Coord is a point in cartesian space
type Coord struct{ x, y int }

// DistanceSquared calculates the squared distance between this Coord and the given Coord
func (coord Coord) DistanceSquared(other Coord) int {
	dx, dy := coord.x-other.x, coord.y-other.y
	return dx*dx + dy*dy
}

// Pac represents a Pac man (or woman)
type Pac struct {
	id                              int
	player                          int
	pos                             Coord
	typeID                          string
	speedTurnsLeft, abilityCooldown int
}

// Pellet represents a pellet with a location and point value
type Pellet struct {
	pos   Coord
	value int
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

// GameData represents a snapshot of the game at a point in time
type GameData struct {
	round          int
	gameMap        GameMap
	scores         []int
	visiblePacs    []Pac
	visiblePellets []Pellet
}

// Agent decides on a command given game data
type Agent interface {
	makeCommand(GameData) string
}

// DansLilHeuristicBot is just a lil guy tryina eat some pellets
type DansLilHeuristicBot struct{}

// Bucketize returns the bucket (0...numBuckets-1) that value x belongs in, if evenly distributed amongst width
func Bucketize(x, numBuckets, width int) int {
	if bucket := x / (width / numBuckets); bucket < numBuckets {
		return bucket
	}
	return numBuckets - 1
}

// SortPelletsByProximity sorts the pellets by distance squared to the given Pac
func SortPelletsByProximity(pellets []Pellet, pac Pac) {
	sort.Slice(pellets, func(i, j int) bool {
		iPellet, jPellet := pellets[i], pellets[j]
		iDist, jDist := iPellet.pos.DistanceSquared(pac.pos), jPellet.pos.DistanceSquared(pac.pos)
		if iDist < jDist {
			return true
		}
		return false
	})
}

func (bot DansLilHeuristicBot) makeCommand(gameState GameData) string {
	var myPacs []Pac
	// find all of my pacs from the visible collection
	for _, pac := range gameState.visiblePacs {
		if pac.player == 1 {
			myPacs = append(myPacs, pac)
		}
	}

	pelletsByArea := make([][]Pellet, len(myPacs))
	for _, pellet := range gameState.visiblePellets {
		key := Bucketize(pellet.pos.x, len(myPacs), gameState.gameMap.width)
		pelletsByArea[key] = append(pelletsByArea[key], pellet)
	}

	var actions []string
	for iPac, pac := range myPacs {
		if pac.abilityCooldown <= 0 {
			// they're speedy lil devils, these Pacs
			actions = append(actions, fmt.Sprint("SPEED ", pac.id))
		} else {
			var pellet Pellet
			var status string
			if len(pelletsByArea[iPac]) > 0 {
				SortPelletsByProximity(pelletsByArea[iPac], pac)
				pellet = pelletsByArea[iPac][0]
				status = joinStrings("P", len(pelletsByArea[iPac]))
			} else {
				coord := func(x int) int {
					return rand.Intn(x/len(myPacs)) + (iPac * x / len(myPacs))
				}
				x, y := coord(gameState.gameMap.width), coord(gameState.gameMap.height)
				pellet = Pellet{Coord{x, y}, 1}
				status = joinStrings("S", x, y)
			}
			actions = append(actions, joinStrings("MOVE", pac.id, pellet.pos, iPac, status))
		}
	}

	return strings.Join(actions, "|")
}

func joinStrings(elems ...interface{}) string {
	elemStrings := make([]string, len(elems))
	for i, elem := range elems {
		elemStrings[i] = fmt.Sprintf("%v", elem)
	}

	return strings.Join(elemStrings, " ")
}

//-----------------------------------------------------------------------------------
// main stuff
//-----------------------------------------------------------------------------------

func debug(a ...interface{}) {
	fmt.Fprintln(os.Stderr, a...)
}

/**
 * Grab the pellets as fast as you can!
 **/
func main() {
	agent := DansLilHeuristicBot{}

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
	var gameRound = 0

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

			visiblePacs = append(visiblePacs, Pac{pacID, player, Coord{x, y}, typeID, speedTurnsLeft, abilityCooldown})
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

			visiblePellets = append(visiblePellets, Pellet{Coord{x, y}, value})
		}

		cmd := agent.makeCommand(GameData{gameRound, gameMap, []int{myScore, opponentScore}, visiblePacs, visiblePellets})
		debug(cmd)
		fmt.Println(cmd)

		gameRound++
	}
}
