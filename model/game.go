package model

import (
	"fmt"
	"strconv"
)

// Consts that tells if a cell is dead or alive
const (
	ALIVE = 1
	DEAD  = 0
)

type GameOfLife struct {
	X          int
	Y          int
	Start      bool
	Generation uint
	CurrentGen [][]int
}

// Initialize a new game with cero cells alive
func (game *GameOfLife) Init(x int, y int) {
	field := make([][]int, x)
	for i := range field {
		field[i] = make([]int, y)
	}
	game.X = x
	game.Y = y
	game.CurrentGen = field
	game.Generation = 0
	game.Start = false
}

// Count how many neighbors are alive in the board
func countAlive(x int, y int, rows int, cols int, board [][]int) int {
	totalNeighbors := 0

	// Getting the total of neighbors that has the cell on field[x][y]
	for i := -1; i < 2; i++ {
		for j := -1; j < 2; j++ {
			if x+i < rows && x+i >= 0 && y+j < cols && y+j >= 0 && board[x+i][y+j] == ALIVE {
				totalNeighbors += board[x+i][y+j]
			}
		}
	}
	// Since we are counting all the neighbors on the for loop we're counting the cell itself
	// So we have to deleted it, this happens beacuse i and j go from -1 to 1
	// So there's a case where i = 0 and j = 0, in that case we are adding the cell perse
	totalNeighbors -= board[x][y]

	return totalNeighbors
}

// Go to the next generation of the game
func (game *GameOfLife) Step() {
	previousGen := make([][]int, game.X)
	for i := 0; i < game.X; i++ {
		previousGen[i] = make([]int, game.Y)
	}

	for i := 0; i < game.X; i++ {
		for j := 0; j < game.Y; j++ {
			previousGen[i][j] = game.CurrentGen[i][j]
		}
	}

	for i := 0; i < game.X; i++ {
		for j := 0; j < game.Y; j++ {
			aliveNeighbors := countAlive(i, j, game.X, game.Y, previousGen)
			if previousGen[i][j] == ALIVE {
				if aliveNeighbors < 2 || aliveNeighbors > 3 {
					game.CurrentGen[i][j] = DEAD
				}
			} else if aliveNeighbors == 3 {
				game.CurrentGen[i][j] = ALIVE
			}
		}
	}
	game.Generation += 1
}

// Kill all the cells that are in the field/board
func (game *GameOfLife) ClearGame() {
	for i := 0; i < game.X; i++ {
		for j := 0; j < game.Y; j++ {
			game.CurrentGen[i][j] = DEAD
		}
	}
	game.Generation = 0
}

func stringBoard(board [][]int, x int, y int) {
	stringBoard := ""
	for i := 0; i < x; i++ {
		for j := 0; j < y; j++ {
			stringBoard += strconv.Itoa(board[i][j])
		}
		stringBoard += "\n"
	}
	stringBoard += "\n"
	fmt.Println(stringBoard)
}
