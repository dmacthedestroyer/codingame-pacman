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
	// space " " is floor, pound "#" is wall
	value rune
}

// GameMap represents the walls and floors of the game area
type GameMap struct {
	width, height int
	cells         []Cell
}

// GetCell gets the Cell value at the given Coord, or panics if not in range
func (gm GameMap) GetCell(pos Coord) Cell {
	return gm.cells[pos.x+(pos.y*gm.width)]
}

// GetCoord returns the (x,y) translation for the provided 1-dimensional position, based on the width and height values for this GameMap
func (gm GameMap) GetCoord(pos int) Coord {
	return Coord{pos % gm.width, pos / gm.width}
}

// GetAbsolutePosition returns the absolute position transation for the provided Coord, based on the width and height values for this GameMap
func (gm GameMap) GetAbsolutePosition(coord Coord) int {
	return coord.x + coord.y*gm.width
}

// VisibleCells returns all cells with line of sight to the given Coord or panics if pos is outside of the map's size
func (gm GameMap) VisibleCells(pos Coord) (visibleCoords []Coord) {
	walk := func(dx, dy int) {
		var walkPos = Coord{pos.x, pos.y}
		step := func() {
			walkPos.x += dx
			if walkPos.x < 0 {
				walkPos.x += gm.width
			}
			if walkPos.x >= gm.width {
				walkPos.x %= gm.width
			}
			walkPos.y += dy
			if walkPos.y < 0 {
				walkPos.y += gm.height
			}
			if walkPos.y >= gm.height {
				walkPos.y %= gm.height
			}
		}
		step()
		for gm.GetCell(walkPos).value != '#' && walkPos != pos {
			visibleCoords = append(visibleCoords, walkPos)
			step()
		}
	}

	visibleCoords = append(visibleCoords, pos)
	walk(0, -1)
	walk(0, 1)
	walk(-1, 0)
	walk(1, 0)

	return
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
type DansLilHeuristicBot struct {
	pelletValuesByPos []int
}

func (bot *DansLilHeuristicBot) init(gameMap GameMap) {
	bot.pelletValuesByPos = make([]int, gameMap.width*gameMap.height)
	for pos, cell := range gameMap.cells {
		var value int
		if cell.value == ' ' {
			value = 1
		} else {
			value = 0
		}
		bot.pelletValuesByPos[pos] = value
	}
}

func (bot *DansLilHeuristicBot) update(gameData GameData) {
	// clear all known pellet values for all currently visible cells -- we'll replace the existing values based on observed data next
	for _, pac := range gameData.visiblePacs {
		// don't clear cells for opponents, since we don't have their line of sight
		if pac.player == 1 {
			for _, coord := range gameData.gameMap.VisibleCells(pac.pos) {
				bot.pelletValuesByPos[gameData.gameMap.GetAbsolutePosition(coord)] = 0
			}
		}
	}
	// update pellet values for all visible pellets
	for _, pellet := range gameData.visiblePellets {
		bot.pelletValuesByPos[gameData.gameMap.GetAbsolutePosition(pellet.pos)] = pellet.value
	}
}

// Bucketize returns the bucket (0...numBuckets-1) that value x belongs in, if evenly distributed amongst width
func Bucketize(x, numBuckets, width int) int {
	if bucket := x / (width / numBuckets); bucket < numBuckets {
		return bucket
	}
	return numBuckets - 1
}

// SortCoordsByProximity sorts the pellets by distance squared to the given Pac
func SortCoordsByProximity(coords []Coord, pos Coord) {
	sort.Slice(coords, func(i, j int) bool {
		iCoord, jCoord := coords[i], coords[j]
		iDist, jDist := iCoord.DistanceSquared(pos), jCoord.DistanceSquared(pos)
		if iDist < jDist {
			return true
		}
		return false
	})
}

func (bot DansLilHeuristicBot) makeCommand(gameData GameData) string {
	bot.update(gameData)

	var myPacs []Pac
	// find all of my pacs from the visible collection
	for _, pac := range gameData.visiblePacs {
		if pac.player == 1 {
			myPacs = append(myPacs, pac)
		}
	}

	pelletsByArea := make([][]Coord, len(myPacs))
	for pos, pelletValue := range bot.pelletValuesByPos {
		if pelletValue > 0 {
			coord := gameData.gameMap.GetCoord(pos)
			key := Bucketize(coord.x, len(myPacs), gameData.gameMap.width)
			pelletsByArea[key] = append(pelletsByArea[key], coord)
		}
	}

	var actions []string
	for iPac, pac := range myPacs {
		if pac.abilityCooldown <= 0 {
			// they're speedy lil devils, these Pacs
			actions = append(actions, fmt.Sprint("SPEED ", pac.id))
		} else {
			var pos Coord
			var status string
			// choose closest pellet
			if len(pelletsByArea[iPac]) > 0 {
				SortCoordsByProximity(pelletsByArea[iPac], pac.pos)
				pos = pelletsByArea[iPac][0]
				status = joinStrings("P", len(pelletsByArea[iPac]))
			} else {
				// wander aimlessly, hoping to find more delicious pellets
				coord := func(x int) int {
					return rand.Intn(x/len(myPacs)) + (iPac * x / len(myPacs))
				}
				x, y := coord(gameData.gameMap.width), coord(gameData.gameMap.height)
				pos = Coord{x, y}
				status = joinStrings("S", x, y)
			}
			actions = append(actions, joinStrings("MOVE", pac.id, pos.x, pos.y, iPac, status))
		}
	}

	return strings.Join(actions, "|")
}

//-----------------------------------------------------------------------------------
// general utility stuff
//-----------------------------------------------------------------------------------

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

func debugf(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
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
	agent.init(gameMap)
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
		gameData := GameData{gameRound, gameMap, []int{myScore, opponentScore}, visiblePacs, visiblePellets}
		cmd := agent.makeCommand(gameData)
		debug(cmd)
		fmt.Println(cmd)

		gameRound++
	}
}
