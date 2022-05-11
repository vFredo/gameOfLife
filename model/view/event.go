package view

const QUIT = "quit"
const INPUT = "input"

type Event struct {
	Type string
	X    uint
	Y    uint
}
