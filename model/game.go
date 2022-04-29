package model

import (
	"log"
	"math"
)

// Structure that has everything that it's need it to play the game
type GameOfLife struct {
	X               int
	Y               int
	Lenght          uint
	Start           bool
	WrapEdges       bool
	Generation      uint
	CurrentGen      []uint8
	BirthCell       uint
	OverPopulation  uint
	UnderPopulation uint
	PresetManager   PresetManager
}

// Initialize a new game with zero cells alive
func (game *GameOfLife) Init(x int, y int) {
	game.X = x
	game.Y = y
	game.Lenght = uint(x * y)
	game.CurrentGen = make([]uint8, game.Lenght)
	game.Generation = 0
	game.PresetManager = PresetManager{}
	game.PresetManager.FetchPresets()
}

// Update the count of each adjacent neighbor taking into account the new state of the cell
func (game *GameOfLife) updateNeighbors(x int, y int, state bool) {
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			// Don't include the cell itself
			if i == 0 && j == 0 {
				j += 1
			}

			// Useful temp variables when WrapEdges is enabled
			aboveBelow, leftRight := i, j

			if game.WrapEdges {
				if x == 0 && i == -1 {
					aboveBelow = game.X - 1
				} else if x == (game.X-1) && i == 1 {
					aboveBelow = -(game.X - 1)
				}

				if y == 0 && j == -1 {
					leftRight = game.Y - 1
				} else if y == (game.Y-1) && j == 1 {
					leftRight = -(game.Y - 1)
				}
			}

			// Don't include cells that are beyond the array's indexes
			if x+aboveBelow < game.X && x+aboveBelow >= 0 && y+leftRight < game.Y && y+leftRight >= 0 {
				pos := ((x + aboveBelow) * game.Y) + (y + leftRight)
				if state {
					game.CurrentGen[pos] += 0x02
				} else {
					game.CurrentGen[pos] -= 0x02
				}
			}
		}
	}
}

// Spawn an alive cell in the [x][y] position
func (game *GameOfLife) SpawnCell(x int, y int) {
	// Spawning the cell
	game.CurrentGen[x*game.Y+y] |= 0x01
	game.updateNeighbors(x, y, true)
}

// Kill a cell in the [x][y] position
func (game *GameOfLife) KillCell(x int, y int) {
	// Killing the cell
	game.CurrentGen[x*game.Y+y] &^= 0x01
	game.updateNeighbors(x, y, false)
}

// Returns the cell state, if it's dead (false) or alive (true)
func (game *GameOfLife) CellState(x int, y int) bool {
	return game.CurrentGen[x*game.Y+y]&0x01 == 0x01
}

// Go to the next generation of the cells
func (game *GameOfLife) Step() {
	// Copy the data of the current generation into a new array
	prevGen := make([]uint8, game.Lenght)
	copy(prevGen, game.CurrentGen)

	for i := 0; i < game.X; i++ {
		for j := 0; j < game.Y; j++ {
			pos := i*game.Y + j
			// Skip quickly through as many dead cells with no neighbors
			for prevGen[pos] == 0x00 {
				pos += 1
				j += 1
				if j >= game.Y {
					goto JumpRow
				}
			}

			// Since the neighbor count start at the second less significant bit
			neighbors := uint(prevGen[pos] >> 1)
			if (prevGen[pos] & 0x01) == 0x01 { // prevGen Cell it's alive?
				if neighbors < game.UnderPopulation || neighbors > game.OverPopulation {
					game.KillCell(i, j)
				}
			} else if neighbors == game.BirthCell {
				game.SpawnCell(i, j)
			}
		}
	JumpRow:
	}
	game.Generation += 1
}

// Resize the board according to the new width and height of the window [x][y]
func (game *GameOfLife) Resize(x int, y int) {
	// Create the new board
	newBoard := make([]uint8, x*y)

	// Copy cells from the previous board to the new board
	for i := 0; i < game.X; i++ {
		for j := 0; j < game.Y; j++ {
			// Don't copy the cells that are beyond the size of the new board
			// it's the case where the new board it's smaller than the old one
			if i < x && j < y {
				newPos := i*y + j
				pos := i*game.Y + j
				newBoard[newPos] = game.CurrentGen[pos]
			}
		}
	}
	game.Lenght = uint(x * y)
	game.X = x
	game.Y = y
	game.CurrentGen = newBoard
}

// Kill all the cells that are in the board and reset neighbor's counters
func (game *GameOfLife) ClearGame() {
	for i := 0; i < int(game.Lenght); i++ {
		game.CurrentGen[i] = 0x00
	}
	game.Generation = 0
}

// Open preset that it's on the list of PresetManager with the name given
func (game *GameOfLife) OpenPreset(name string) {
	game.ClearGame()

	currPreset, err := game.PresetManager.GetPreset(name)
	if err != nil {
		log.Fatalf("Error while getting preset: %s", err)
	}

	// TODO: Center preset on currentGen
	for k := 0; k < len(currPreset.AliveCells); k++ {
		aliveCell := currPreset.AliveCells[k]
		if int(aliveCell[0]) < game.X && int(aliveCell[1]) < game.Y {
			game.SpawnCell(int(aliveCell[0]), int(aliveCell[1]))
		}
	}
}

// Save the current board as a json file with the name given
func (game *GameOfLife) SaveBoard(name string) {
	max_width, max_height := math.Inf(-1), math.Inf(-1)
	min_width, min_height := math.Inf(1), math.Inf(1)
	var alive [][]uint

	// Save the position of the alive cells and found the max and min
	// of the width and height
	for i := game.X - 1; i > 0; i-- {
		for j := game.Y - 1; j > 0; j-- {
			pos := (i * game.Y) + j
			if (game.CurrentGen[pos] & 0x01) == 0x01 {
				max_width = math.Max(max_width, float64(i))
				min_width = math.Min(min_width, float64(i))

				max_height = math.Max(max_height, float64(j))
				min_height = math.Min(min_height, float64(j))
				alive = append(alive, []uint{uint(i), uint(j)})
			}
		}
	}

	// FIX: This doesn't work at how it's suposed to, but nevertheless.. it works in most cases xD
	// To make a smaller preset we have to know the dimentions of the section
	// that have alive cells on the board; that's why we need the min and max
	width := uint(max_width - min_width)
	height := uint(max_height - min_height)

	// Translate the position of the alive cells to the new smaller board of 'width' 'height'
	for i := 0; i < len(alive); i++ {
		alive[i][0] -= width
		alive[i][1] -= height
	}

	game.PresetManager.CreatePreset(name, alive, width, height)
}
