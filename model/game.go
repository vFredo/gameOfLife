package model

// Structure that has everything that it's need it to play the game
type GameOfLife struct {
	X               int
	Y               int
	Start           bool
	Generation      uint
	CurrentGen      [][]uint8
	BirthCell       uint
	OverPopulation  uint
	UnderPopulation uint
}

// Initialize a new game with zero cells alive
func (game *GameOfLife) Init(x int, y int) {
	field := make([][]uint8, x)
	for i := 0; i < x; i++ {
		field[i] = make([]uint8, y)
	}
	game.X = x
	game.Y = y
	game.CurrentGen = field
	game.Generation = 0
	game.Start = false
}

// Update the count of each adjacent neighbor taking into account the new state of the cell
// If state == true add 1, else delete 1 to the count
func (game *GameOfLife) updateNeighbors(x int, y int, state bool) {
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			// don't include the cell itself and cells that are beyond the the matrix indexes
			if (x+i < game.X && x+i >= 0 && y+j < game.Y && y+j >= 0) && (j != 0 || i != 0) {
				if state {
					game.CurrentGen[x+i][y+j] += 0x02
				} else {
					game.CurrentGen[x+i][y+j] -= 0x02
				}
			}
		}
	}
}

// Copy the data of the current generation into a new matrix
func (game *GameOfLife) copyGeneration() [][]uint8 {
	previous := make([][]uint8, game.X)
	for i := 0; i < game.X; i++ {
		previous[i] = make([]uint8, game.Y)
	}

	for i := 0; i < game.X; i++ {
		for j := 0; j < game.Y; j++ {
			previous[i][j] = game.CurrentGen[i][j]
		}
	}
	return previous
}

// Spawn an alive cell in the [x][y] position
func (game *GameOfLife) SpawnCell(x int, y int) {
	// Spawning the cell
	game.CurrentGen[x][y] |= 0x01
	game.updateNeighbors(x, y, true)
}

// Kill a cell in the [x][y] position
func (game *GameOfLife) KillCell(x int, y int) {
	// Killing the cell
	game.CurrentGen[x][y] &^= 0x01
	game.updateNeighbors(x, y, false)
}

// Returns the cell state, if it's dead (false) or alive (true)
func (game *GameOfLife) CellState(x int, y int) bool {
	return game.CurrentGen[x][y]&0x01 == 0x01
}

// Go to the next generation of cells
func (game *GameOfLife) Step() {
	prevGen := game.copyGeneration()

	for i := 0; i < game.X; i++ {
		for j := 0; j < game.Y; j++ {
			// Jump quickly through as many dead cells with no neighbors
			for prevGen[i][j] == 0x00 {
				j += 1
				if j >= game.Y {
					goto RowDone
				}
			}

			// Since the neighbor count start at the second less significant bit
			neighbors := uint(prevGen[i][j] >> 1)
			if (prevGen[i][j] & 0x01) == 0x01 { // prevGen Cell it's alive?
				if neighbors < game.UnderPopulation || neighbors > game.OverPopulation {
					game.KillCell(i, j)
				}
			} else if neighbors == game.BirthCell {
				game.SpawnCell(i, j)
			}
		}
	RowDone:
	}
	game.Generation += 1
}

// Resize the board according to the new width and height of the window [x][y]
func (game *GameOfLife) Resize(x int, y int) {
	// Create the new board
	newBoard := make([][]uint8, x)
	for i := 0; i < x; i++ {
		newBoard[i] = make([]uint8, y)
	}

	// Copy cells from the previous board to the new board
	for i := 0; i < game.X; i++ {
		for j := 0; j < game.Y; j++ {
			if i < x && j < y {
				newBoard[i][j] = game.CurrentGen[i][j]
			}
		}
	}
	game.X = x
	game.Y = y
	game.CurrentGen = newBoard
}

// Kill all the cells that are in the board and reset neighbor's counters
func (game *GameOfLife) ClearGame() {
	for i := 0; i < game.X; i++ {
		for j := 0; j < game.Y; j++ {
			game.CurrentGen[i][j] = 0x00
		}
	}
	game.Generation = 0
}
