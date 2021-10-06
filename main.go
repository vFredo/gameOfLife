package main

import (
	"github.com/fredo0522/gameOfLife/model"
)

func main() {
	// execute the view buffer of the terminal and the game itself
	view := model.View{}
	view.Loop()
}
