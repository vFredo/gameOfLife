package main

import (
	"flag"

	"github.com/fredo0522/gameOfLife/model"
)

func main() {

	// Getting info from flag arguments
	// If the arguments are not specified, then take default behaviour B3S23
	birth := flag.Int("b", 3, "Indicates the number of neighbors that a cell needs to be born on the next generation.")
	under := flag.Int("u", 2, "Indicates the number of neighbors that a cell needs to die due to under population.")
	over := flag.Int("o", 3, "Indicates the number of neighbors that a cell needs to die due to over population.")

	// parse flag with the operating system arguments
	flag.Parse()

	// execute the view buffer of the terminal and the game itself
	view := model.View{}
	view.InitScreen(*birth, *under, *over)
	view.StartLoop()
}
