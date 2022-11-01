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

// Helper method to render the input box on the screen
func (view *TermView) renderInput(currName string) {
	label := " Preset name: "
	border := ""
	width, height := view.screen.Size()

	// Initialize borders
	sizeBox := 3*(width/4) - width/4
	for i := 0; i < sizeBox; i++ {
		border += "-"
	}

	// Top border
	view.renderText(width/4, (height/2)-1, border)
	// Center input
	middleText := "|" + label + currName
	for i := len(middleText); i < sizeBox-1; i++ {
		middleText += " "
	}
	middleText += "|"
	if len(middleText) <= sizeBox {
		view.renderText(width/4, height/2, middleText)
	}
	// Bottom Border
	view.renderText(width/4, (height/2)+1, border)
}

// Control input (mouse/keyboard) events on the screen
func (view *TermView) readInput() {
	for {
		// Catch events that are triggered on the buffer
		ev := view.screen.PollEvent()
		// Process each event
		switch ev := ev.(type) {
		case *tcell.EventResize:
			view.screen.Sync()
			width, height := view.screen.Size()
			view.game.Resize(uint(height), uint(width/2))
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC || ev.Rune() == 'q' {
				view.event.SetType(QUIT)
				return
			} else if ev.Rune() == ' ' { // space
				view.Start = !view.Start
			} else if ev.Key() == tcell.KeyEnter && !view.Start {
				view.game.Step()
			} else if ev.Rune() == 'h' {
				view.hideMenu = !view.hideMenu
				if !view.hideMenu {
					view.hideMenuAll = false
				}
			} else if ev.Rune() == 'H' {
				view.hideMenuAll = !view.hideMenuAll
				if view.hideMenuAll {
					view.hideMenu = true
				}
			} else if ev.Rune() == 'p' {
				if !view.game.Start {
					view.event.SetType(PAUSE)
					view.inputFormPreset()
				}
			} else if ev.Rune() == 'c' {
				view.game.CyclePresets()
			}
		case *tcell.EventMouse:
			switch ev.Buttons() {
			case tcell.Button1: // left click
				x, y := ev.Position()
				// x it's divided by 2 because each cell it's represent with 2 pixels wide
				row, col := uint(y), uint(x/2)
				// If the game is in pause, let it modified
				if row < view.game.X && col < view.game.Y && !view.Start {
					if view.game.CellState(row, col) {
						view.game.KillCell(row, col)
					} else {
						view.game.SpawnCell(row, col)
					}
				}
			case tcell.Button2: // right click
				view.game.ClearGame()
			}
		}
	}
}

// Manage events of the input form when the user wants to save a preset
func (view *TermView) inputFormPreset() {
	namePreset := ""

	for {
		view.renderInput(namePreset)

		ev := view.screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape {
				view.event.SetType(RUNNING)
				return
			} else if ev.Key() == tcell.KeyBackspace || ev.Key() == tcell.KeyBackspace2 {
				if len(namePreset) > 0 {
					namePreset = namePreset[:len(namePreset)-1]
				}
			} else if ev.Key() == tcell.KeyEnter {
				// If the preset doesn't have a name, put the current date and time
				if namePreset == "" {
					namePreset = time.Now().Format("01-06-2006_15:04:05")
				}
				view.game.SaveBoard(namePreset)
				view.event.SetType(RUNNING)
				return
			} else {
				namePreset += string(ev.Rune())
			}
		}
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
