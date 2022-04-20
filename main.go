package main

import (
	"flag"
	"os"

	"github.com/fredo0522/gameOfLife/model"
)

func main() {
	// If the arguments are not specified, then take default behavior B3S23
	birth := flag.Uint("b", 3, "Indicates the number of neighbors that a cell needs to be born on the next generation.\n")
	under := flag.Uint("u", 2, "Indicates the number of neighbors that a cell needs to die due to under population.\n")
	over := flag.Uint("o", 3, "Indicates the number of neighbors that a cell needs to die due to over population.\n")

	// Parse flag with the OS's executable arguments
	flag.Parse()

	// Game instance
	game := model.GameOfLife{
		BirthCell:       *birth,
		OverPopulation:  *over,
		UnderPopulation: *under,
	}

	// Execute the view buffer of the terminal and the game itself
	view := model.View{}
	view.InitScreen(game)
	view.Run()

	// Exit the program
	os.Exit(0)
}
