package main

import (
	"flag"
	"os"

	"github.com/fredo0522/gameOfLife/model"
)

func main() {
	// If the arguments are not specified, then take default behavior B3S23
	birth := flag.Uint("b", 3, "Indicates the exact number of neighbors that a cell needs in order to be spawn in the next generation.")
	under := flag.Uint("u", 2, "Indicates the number of neighbors a cell needs to be killed due to under population.")
	over := flag.Uint("o", 3, "Indicates the number of neighbors a cell needs to be killed due to over population.")
	wrapEdge := flag.Bool("wrap", false, "Use wrap edges on the board.")

	// Parse flag with the OS's executable arguments
	flag.Parse()

	// Game instance
	game := model.GameOfLife{
		BirthCell:       *birth,
		OverPopulation:  *over,
		UnderPopulation: *under,
		WrapEdges:       *wrapEdge,
	}

	// Execute the view buffer on the terminal and initialize game
	view := model.View{}
	view.InitScreen(game)
	view.Run()

	// Exit the program
	os.Exit(0)
}
