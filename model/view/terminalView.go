package view

import (
	"fmt"
	"log"
	"time"

	"github.com/fredo0522/gameOfLife/model/game"
	"github.com/gdamore/tcell/v2"
)

// Styles (uses the terminal colors)
var (
	PlayStyle    = tcell.StyleDefault.Background(tcell.ColorYellow).Foreground(tcell.ColorReset)
	PauseStyle   = tcell.StyleDefault.Background(tcell.ColorBlue).Foreground(tcell.ColorReset)
	DefaultStyle = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	InfoStyle    = tcell.StyleDefault.Background(tcell.ColorGray).Foreground(tcell.ColorWhite)
)

// TermView structure that has the game itself and the screen (terminal buffer) where everything is render
type TermView struct {
	screen      tcell.Screen
	game        game.GameOfLife
	Start       bool
	hideMenu    bool
	hideMenuAll bool
	event       Event
}

// Initialize screen and game itself
func (view *TermView) InitScreen(game game.GameOfLife) {
	screenInstance, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	view.screen = screenInstance

	if err := view.screen.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	// Initialize the user event
	view.event = Event{Type: RUNNING}

	view.screen.SetStyle(DefaultStyle)
	view.screen.EnableMouse()
	view.screen.Clear()

	// Initialize Game into the view
	view.game = game
	width, height := view.screen.Size()
	view.game.Init(uint(height), uint(width/2))
}

// Helper method that render text on the screen
func (view *TermView) renderText(x int, y int, info string) {
	for i, char := range info {
		view.screen.SetContent(x+i, y, char, nil, InfoStyle)
	}
}

// How each cell is draw on the screen
func (view *TermView) displayGame() {
	view.screen.Clear()
	// Death: Background color, AlivePause: blue, AlivePlay: yellow
	for i := 0; uint(i) < view.game.X; i++ {
		for j := 0; uint(j) < view.game.Y; j++ {
			if view.game.CellState(uint(i), uint(j)) && view.Start {
				view.screen.SetContent(j*2, i, ' ', nil, PlayStyle)
				view.screen.SetContent(j*2+1, i, ' ', nil, PlayStyle)
			} else if view.game.CellState(uint(i), uint(j)) {
				view.screen.SetContent(j*2, i, ' ', nil, PauseStyle)
				view.screen.SetContent(j*2+1, i, ' ', nil, PauseStyle)
			} else {
				view.screen.SetContent(j*2, i, ' ', nil, DefaultStyle)
				view.screen.SetContent(j*2+1, i, ' ', nil, DefaultStyle)
			}
		}
	}
}

// Information shown on the menu
func (view *TermView) displayInfo() {
	width, height := view.screen.Size()

	generationText := fmt.Sprintf(" Generation: %d ", view.game.Generation)
	firstText :=
		" ENTER: next generation, SPC: play/pause, q/ESC/Ctrl-C: quit, h/H: hide menu/ ALL menu info "
	secondText :=
		" LeftClick: switch state cell, RightClick: reset board  p: create preset c: cycle presets "
	x, y := 0, 0

	if !view.hideMenu {
		if len(firstText) <= width {
			view.renderText(x, y, firstText)
			y += 1
		}

		if len(secondText) <= width {
			view.renderText(x, y, secondText)
			y += 1
		}
	}

	if !view.hideMenuAll {
		view.renderText(width-len(generationText), height-1, generationText)
	}
}

// Infinite loop for the terminal view buffer where the game is executed
func (view *TermView) Run() {
	// Get the FPS for executing the game while is on a 'start' state
	framesPerSecond := 15
	sleepTime := time.Duration(1000/framesPerSecond) * time.Millisecond

	// Read input in another routine
	go view.readInput()
	for {
		switch view.event.GetType() {
		case RUNNING:
			view.displayGame()
			if view.Start {
				view.game.Step()
				time.Sleep(sleepTime)
			}
			// Display information menu and update the screen
			view.displayInfo()
			view.screen.Show()
		case QUIT:
			view.screen.Clear()
			view.screen.Fini()
			return
		case PAUSE:
			// Update screen
			view.screen.Show()
		}
	}
}
