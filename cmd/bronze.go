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

// Pac represents a Pac man (or woman)
type Pac struct {
	// id is the pac's id (unique for a given player)
	id int
	// mine is true if this pac belongs to the player
	mine bool
	// pos is the pac's positoin
	pos Coord
	// typeID is the pac's type (ROCK or PAPER or SCISSORS). In the next league, a pac that has died will be of type DEAD.
	typeID string
	// speedTurnsLeft is the number of remaining turns before the speed effect fades
	speedTurnsLeft int
	// abilityCooldown is the number of turns until you can request a new ability for this pac (SWITCH and SPEED)
	abilityCooldown int
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
	return gm.cells[gm.GetAbsolutePosition(pos)]
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

// Wrap normalizes a coordinate that may be outside of the bounds of the map by wrapping it around to the other side
func (gm GameMap) Wrap(coord Coord) Coord {
	_wrap := func(d, m int) int {
		var res int = d % m
		if (res < 0 && m > 0) || (res > 0 && m < 0) {
			return res + m
		}
		return res
	}

	return Coord{_wrap(coord.x, gm.width), _wrap(coord.y, gm.height)}
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

//-----------------------------------------------------------------------------------
// general utility stuff
//-----------------------------------------------------------------------------------

func adjacentCoords(pos Coord) []Coord {
	adjacents := []Coord{}

	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			if dx != 0 || dy != 0 {
				adjacents = append(adjacents, Coord{pos.x + dx, pos.y + dy})
			}
		}
	}
	return adjacents
}

// enemiesWithinRange returns all enemies within the given distance, sorted by distance
// TODO: more tests
func enemiesWithinRange(gameMap GameMap, pacsByPosition map[Coord]Pac, pos Coord, distance int) []Pac {
	type SearchNode struct {
		pos   Coord
		depth int
	}

	queue := []SearchNode{{pos, 0}}
	// list of nodes already visited to avoid cycles
	visited := map[Coord]bool{pos: true}

	enemies := []Pac{}

	for len(queue) > 0 {
		// queue.peek
		node := queue[0]
		// queue.dequeue
		queue = queue[1:]

		// if this position contains an enemy pac, add it to our list of enemies
		pac, isPacPresent := pacsByPosition[node.pos]
		if isPacPresent && !pac.mine {
			enemies = append(enemies, pac)
		}

		// walk each orthogonal (x, y) coordinate
		for _, adjacent := range adjacentCoords(node.pos) {
			dPos := gameMap.Wrap(adjacent)
			alreadyVisited := visited[dPos]
			// if we haven't visited this coordinate before in this search
			if !alreadyVisited {
				visited[dPos] = true
				cell := gameMap.GetCell(dPos)
				// if the cell is a floor (instead of a wall)
				if cell.value == ' ' {
					// don't process nodes that are farther away than our max distance
					if node.depth+1 <= distance {
						queue = append(queue, SearchNode{dPos, node.depth + 1})
					}
				}
			}
		}
	}

	return enemies
}

func getWinningTypeId(typeId string) string {
	if typeId == "ROCK" {
		return "PAPER"
	} else if typeId == "PAPER" {
		return "SCISSORS"
	} else if typeId == "SCISSORS" {
		return "ROCK"
	} else {
		panic(fmt.Sprintf("unknown typeId: %v", typeId))
	}
}

// awayFrom returns a traversable coordinate one move away from me in the opposite direction of them
func awayFrom(me Coord, them Coord, gameMap GameMap) Coord {
	dx := []int{0}
	dy := []int{0}

	// TODO: doesn't account for wrapping around the map
	if them.x >= me.x {
		dx = append(dx, -1)
	}
	if them.x <= me.x {
		dx = append(dx, +1)
	}
	if them.y >= me.y {
		dy = append(dy, -1)
	}
	if them.y <= me.y {
		dy = append(dy, +1)
	}

	for _, newX := range dx {
		for _, newY := range dy {
			if newX != 0 || newY != 0 {
				newCoord := gameMap.Wrap(Coord{me.x + newX, me.y + newY})
				if gameMap.GetCell(newCoord).value == ' ' {
					return newCoord
				}
			}
		}
	}

	// ¯\_(ツ)_/¯
	return me
}

// sortCoords sorts (in place) area by value, then distance
func sortCoords(area []Coord, pos Coord, pelletValuesByPos map[Coord]int) {
	sort.Slice(area, func(i, j int) bool {
		iCoord, jCoord := area[i], area[j]
		iValue, iExists := pelletValuesByPos[iCoord]
		jValue, jExists := pelletValuesByPos[jCoord]
		if !iExists && jExists {
			return true
		} else if iExists && !jExists {
			return false
		} else if iValue < jValue {
			return false
		} else if iValue > jValue {
			return true
		} else {
			iDistance, jDistance := iCoord.distanceSquared(pos), jCoord.distanceSquared(pos)
			return iDistance < jDistance
		}
	})
}

// distanceSquared calculates the squared distance between this Coord and the given Coord
func (coord Coord) distanceSquared(other Coord) int {
	dx, dy := coord.x-other.x, coord.y-other.y
	return dx*dx + dy*dy
}

// bucketize returns the bucket (0...numBuckets-1) that value x belongs in, if evenly distributed amongst width
func bucketize(x, numBuckets, width int) int {
	if bucket := x / (width / numBuckets); bucket < numBuckets {
		return bucket
	}
	return numBuckets - 1
}

func joinStrings(elems ...interface{}) string {
	elemStrings := make([]string, len(elems))
	for i, elem := range elems {
		elemStrings[i] = fmt.Sprintf("%v", elem)
	}

	return strings.Join(elemStrings, " ")
}

//-----------------------------------------------------------------------------------
// Bot implementation
//-----------------------------------------------------------------------------------

// DansLilHeuristicBot is just a lil guy tryina eat some pellets
type DansLilHeuristicBot struct {
	// pelletValuesByPos keeps track of each pellet value based on its absolute position in the grid
	pelletValuesByPos []int
	// pelletValuesByCoord keeps track of each pellet value based on its coordinate position. I made this because I regretted storing the info in an array in pelletValuesByPos
	pelletValuesByCoord map[Coord]int
	pacsByPos           map[Coord]Pac
}

func (bot *DansLilHeuristicBot) init(gameMap GameMap) {
	bot.pelletValuesByPos = make([]int, gameMap.width*gameMap.height)
	bot.pelletValuesByCoord = make(map[Coord]int, gameMap.width*gameMap.height)
	for pos, cell := range gameMap.cells {
		coord := gameMap.GetCoord(pos)
		var value int
		if cell.value == ' ' {
			value = 1
		} else {
			value = 0
		}
		bot.pelletValuesByPos[pos] = value
		bot.pelletValuesByCoord[coord] = value
	}
	bot.pacsByPos = make(map[Coord]Pac)
}

func (bot *DansLilHeuristicBot) update(gameData GameData) {
	// clear all known pellet values for all currently visible cells -- we'll replace the existing values based on observed data next
	for _, pac := range gameData.visiblePacs {
		// only clear cells for my pacs
		if pac.mine {
			for _, coord := range gameData.gameMap.VisibleCells(pac.pos) {
				bot.pelletValuesByPos[gameData.gameMap.GetAbsolutePosition(coord)] = 0
			}
		}
	}
	// update pellet values for all visible pellets
	for _, pellet := range gameData.visiblePellets {
		bot.pelletValuesByPos[gameData.gameMap.GetAbsolutePosition(pellet.pos)] = pellet.value
		bot.pelletValuesByCoord[pellet.pos] = pellet.value
	}

	// update pacs by position
	bot.pacsByPos = make(map[Coord]Pac)
	for _, pac := range gameData.visiblePacs {
		bot.pacsByPos[pac.pos] = pac
	}
}

func (bot DansLilHeuristicBot) makeCommand(gameData GameData) string {
	bot.update(gameData)

	var myPacs []Pac
	// find all of my pacs from the visible collection
	for _, pac := range gameData.visiblePacs {
		if pac.mine {
			myPacs = append(myPacs, pac)
		}
	}

	// partition pellets into mutually exclusive zones for each pac
	pelletsByArea := make([][]Coord, len(myPacs))
	for pos, pelletValue := range bot.pelletValuesByPos {
		if pelletValue > 0 {
			coord := gameData.gameMap.GetCoord(pos)
			key := bucketize(coord.x, len(myPacs), gameData.gameMap.width)
			pelletsByArea[key] = append(pelletsByArea[key], coord)
		}
	}

	var actions []string
	for iPac, pac := range myPacs {
		speed := func(status string) string { return joinStrings("SPEED ", pac.id, status) }
		move := func(pos Coord, status string) string { return joinStrings("MOVE", pac.id, pos.x, pos.y, iPac, status) }
		switchType := func(typeId string) string { return joinStrings("SWITCH", pac.id, typeId) }
		var action string

		// find any enemies within "striking distance"
		enemies := enemiesWithinRange(gameData.gameMap, bot.pacsByPos, pac.pos, 4)
		if len(enemies) > 0 {
			nearest := enemies[0]
			winningTypeId := getWinningTypeId(nearest.typeID)
			if winningTypeId == pac.typeID {
				if pac.abilityCooldown <= 0 {
					action = speed("ZOOM")
				} else {
					action = move(nearest.pos, "NOM")
				}
			} else if pac.abilityCooldown <= 0 {
				action = switchType(winningTypeId)
			} else {
				action = move(awayFrom(pac.pos, nearest.pos, gameData.gameMap), "EEK!")
			}
		} else {
			// choose closest pellet TODO: fix locking conditions
			myArea := pelletsByArea[iPac]
			if len(myArea) > 0 {
				sortCoords(myArea, pac.pos, bot.pelletValuesByCoord)
				action = move(pelletsByArea[iPac][0], joinStrings("P", len(pelletsByArea[iPac])))
			} else {
				// wander aimlessly, hoping to find more delicious pellets
				coord := func(x int) int {
					return rand.Intn(x/len(myPacs)) + (iPac * x / len(myPacs))
				}
				x, y := coord(gameData.gameMap.width), coord(gameData.gameMap.height)
				action = move(Coord{x, y}, joinStrings("S", x, y))
			}
		}

		if len(action) > 0 {
			actions = append(actions, action)
		}
	}

	return strings.Join(actions, "|")
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

			visiblePacs = append(visiblePacs, Pac{pacID, player == 1, Coord{x, y}, typeID, speedTurnsLeft, abilityCooldown})
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
