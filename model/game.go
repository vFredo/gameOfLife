package model

// Consts that tells if a cell is dead or alive
const (
	ALIVE = 1
	DEAD  = 0
)

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

// Resize board according to the current width and height of the window (x,y)
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

// Spawn an alive cell in the [x][y] position
func (game *GameOfLife) SetCell(x int, y int) {
	game.CurrentGen[x][y] |= ALIVE

	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if (x+i < game.X && x+i >= 0 && y+j < game.Y && y+j >= 0) && (j != 0 || i != 0) {
				game.CurrentGen[x+i][y+j] += 0x02
			}
		}
	}
}

// Kill a cell in the [x][y] position
func (game *GameOfLife) ClearCell(x int, y int) {
	game.CurrentGen[x][y] &^= ALIVE

	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if (x+i < game.X && x+i >= 0 && y+j < game.Y && y+j >= 0) && (j != 0 || i != 0) {
				game.CurrentGen[x+i][y+j] -= 0x02
			}
		}
	}
}

// Returns the cell state, if it's dead (0) or alive (1)
func (game *GameOfLife) CellState(x int, y int) bool {
	return int(game.CurrentGen[x][y]&ALIVE) == ALIVE
}

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

// Go to the next generation of the game
func (game *GameOfLife) Step() {
	prevGen := game.copyGeneration()

	for i := 0; i < game.X; i++ {
		for j := 0; j < game.Y; j++ {
			neighbors := uint(prevGen[i][j] >> 1)
			if (prevGen[i][j] & ALIVE) == ALIVE {
				if neighbors < game.UnderPopulation || neighbors > game.OverPopulation {
					game.ClearCell(i, j)
				}
			} else if neighbors == game.BirthCell {
				game.SetCell(i, j)
			}
		}
	}
	game.Generation += 1
}

// Kill all the cells that are in the board
func (game *GameOfLife) ClearGame() {
	for i := 0; i < game.X; i++ {
		for j := 0; j < game.Y; j++ {
			game.CurrentGen[i][j] = DEAD
		}
	}
	game.Generation = 0
}
