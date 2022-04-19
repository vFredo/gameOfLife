package model

// Consts that tells if a cell is dead or alive
const (
	ALIVE = 1
)

// Structure that has everything that it's need it to play the game
type GameOfLife struct {
	Width           int
	Height          int
	Length          uint
	Start           bool
	Generation      uint
	CurrentGen      []uint8
	TempGen         []uint8
	BirthCell       int
	OverPopulation  int
	UnderPopulation int
}

// Initialize a new game with zero cells alive
func (game *GameOfLife) Init(x int, y int) {
	game.Width = x
	game.Height = y
	game.Length = uint(x * y)
	field := make([]uint8, game.Length)
	game.CurrentGen = field
	game.TempGen = field
	game.Generation = 0
	game.Start = false
}

// Resize board according to the current width and height of the window (x,y)
func (game *GameOfLife) Resize(x int, y int) {
	// Create the new board
	newBoard := make([]uint8, x*y)

	// Copy cells from the previous board to the new board
	for i := 0; i < game.Width; i++ {
		for j := 0; j < game.Height; j++ {
			pos := j*game.Width + i
			if (i != 0 || j != 0) && pos < x*y {
				newBoard[pos] = game.CurrentGen[pos]
			}
		}
	}
	game.Width = x
	game.Height = y
	game.Length = uint(x * y)
	game.CurrentGen = newBoard
	game.TempGen = newBoard
}

// Spawn an alive cell in the [x][y] position
func (game *GameOfLife) SetCell(x int, y int) {
	pos := y*game.Width + x
	game.CurrentGen[pos] |= 0x01

	// Calculate the offsets to the eight neighboring cells,
	// accounting for wrapping around at the edges of the cell map
	var xLeft, xRight, yAbove, yBelow int
	if x == 0 {
		xLeft = game.Width - 1
	} else {
		xLeft = -1
	}

	if x == (game.Width - 1) {
		xRight = -(game.Width - 1)
	} else {
		xRight = 1
	}

	if y == 0 {
		yAbove = int(game.Length) - game.Width
	} else {
		yAbove = -game.Width
	}

	if y == (game.Height - 1) {
		yBelow = -(int(game.Length) - game.Width)
	} else {
		yBelow = game.Width
	}

	// Update update cell's neighbors adding 1 to the count
	// Change successive bits for neighbour counts
	game.CurrentGen[pos+yAbove+xLeft] += 0x02
	game.CurrentGen[pos+yAbove] -= 0x02
	game.CurrentGen[pos+yAbove+xRight] += 0x02
	game.CurrentGen[pos+xLeft] += 0x02
	game.CurrentGen[pos+xRight] += 0x02
	game.CurrentGen[pos+yBelow+xLeft] += 0x02
	game.CurrentGen[pos+yBelow] += 0x02
	game.CurrentGen[pos+yBelow+xRight] += 0x02
}

// Kill a cell in the [x][y] position
func (game *GameOfLife) ClearCell(x int, y int) {
	pos := y*game.Width + x
	game.CurrentGen[pos] &^= 0x01

	// Calculate the offsets to the eight neighboring cells,
	// accounting for wrapping around at the edges of the cell map
	var xLeft, xRight, yAbove, yBelow int
	if x == 0 {
		xLeft = game.Width - 1
	} else {
		xLeft = -1
	}

	if x == (game.Width - 1) {
		xRight = -(game.Width - 1)
	} else {
		xRight = 1
	}

	if y == 0 {
		yAbove = int(game.Length) - game.Width
	} else {
		yAbove = -game.Width
	}

	if y == (game.Height - 1) {
		yBelow = -(int(game.Length) - game.Width)
	} else {
		yBelow = game.Width
	}

	// Update update cell's neighbors adding 1 to the count
	// Change successive bits for neighbour counts
	game.CurrentGen[pos+yAbove+xLeft] -= 0x02
	game.CurrentGen[pos+yAbove] -= 0x02
	game.CurrentGen[pos+yAbove+xRight] -= 0x02
	game.CurrentGen[pos+xLeft] -= 0x02
	game.CurrentGen[pos+xRight] -= 0x02
	game.CurrentGen[pos+yBelow+xLeft] -= 0x02
	game.CurrentGen[pos+yBelow] -= 0x02
	game.CurrentGen[pos+yBelow+xRight] -= 0x02
}

// Returns the cell state, if it's dead (0) or alive (1)
func (game *GameOfLife) CellState(x int, y int) int {
	return int(game.CurrentGen[y*game.Width+x] & ALIVE)
}

// Go to the next generation of the game
func (game *GameOfLife) Step() {
	// Copying array
	copy(game.TempGen, game.CurrentGen)

	// Update current generation taking into account the previous generation
	for i := 0; i < game.Width; i++ {
		for j := 0; j < game.Height; j++ {
			neighbors := int(game.TempGen[j*game.Width+i] >> 0x01)
			if game.CellState(i, j) == ALIVE {
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
	for i := 0; i < game.Width; i++ {
		for j := 0; j < game.Height; j++ {
			game.ClearCell(i, j)
		}
	}
	game.Generation = 0
}
