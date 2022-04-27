package model

import (
	"fmt"
	"log"
	"time"

	"github.com/gdamore/tcell/v2"
)

// Styles (uses the terminal colors)
var (
	PlayStyle    = tcell.StyleDefault.Background(tcell.ColorYellow).Foreground(tcell.ColorReset)
	PauseStyle   = tcell.StyleDefault.Background(tcell.ColorBlue).Foreground(tcell.ColorReset)
	DefaultStyle = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	InfoStyle    = tcell.StyleDefault.Background(tcell.ColorGray).Foreground(tcell.ColorWhite)
)

// View structure that has the game itself and the screen (terminal buffer) where everything is render
type View struct {
	screen   tcell.Screen
	game     GameOfLife
	Start    bool
	hideMenu bool
	quit     chan struct{}
}

// Initialize screen and game itself
func (view *View) InitScreen(game GameOfLife) {
	screenInstance, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	view.screen = screenInstance

	if err := view.screen.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	// Initialize the signal to quit
	view.quit = make(chan struct{})

	view.screen.SetStyle(DefaultStyle)
	view.screen.EnableMouse()
	view.screen.Clear()

	// Adding the game to the view
	view.game = game
	view.Start = false

	// Initialize Game
	width, height := view.screen.Size()
	view.game.Init(height, width/2)
}

// Control input (mouse/keyboard) events on the screen
func (view *View) readInput() {
	for {
		// Catch events that are triggered on the buffer
		ev := view.screen.PollEvent()
		// Process each event
		switch ev := ev.(type) {
		case *tcell.EventResize:
			view.screen.Sync()
			width, height := view.screen.Size()
			view.game.Resize(height, width/2)
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC || ev.Rune() == 'q' {
				view.quit <- struct{}{}
				return
			} else if ev.Rune() == ' ' { // space
				view.Start = !view.Start
			} else if ev.Key() == tcell.KeyEnter && !view.Start {
				view.game.Step()
			} else if ev.Rune() == 'h' {
				view.hideMenu = !view.hideMenu
			} else if ev.Rune() == 'c' {
        if !view.game.Start {
          view.game.SaveBoard(time.Now().Format("01-06-2006_15:04:05"))
        }
			}
		case *tcell.EventMouse:
			switch ev.Buttons() {
			case tcell.Button1: // left click
				x, y := ev.Position()
				// x it's divided by 2 because each cell it's represent with 2 pixels wide
				row, col := y, x/2
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

// How each cell is draw on the screen
func (view *View) displayGame() {
	view.screen.Clear()
	// Death: Background color, AlivePause: blue, AlivePlay: yellow
	for i := 0; i < view.game.X; i++ {
		for j := 0; j < view.game.Y; j++ {
			if view.game.CellState(i, j) && view.Start {
				view.screen.SetContent(j*2, i, ' ', nil, PlayStyle)
				view.screen.SetContent(j*2+1, i, ' ', nil, PlayStyle)
			} else if view.game.CellState(i, j) {
				view.screen.SetContent(j*2, i, ' ', nil, PauseStyle)
				view.screen.SetContent(j*2+1, i, ' ', nil, PauseStyle)
			} else {
				view.screen.SetContent(j*2, i, ' ', nil, DefaultStyle)
				view.screen.SetContent(j*2+1, i, ' ', nil, DefaultStyle)
			}
		}
	}
}

// Helper method that takes the string (info) and pos [x][y] and put it on the screen
func (view *View) renderInfo(x int, y int, info string) {
	for i, byte := range info {
		view.screen.SetCell(x+i, y, InfoStyle, byte)
	}
}

// Information shown on the menu
func (view *View) displayInfo() {
	width, height := view.screen.Size()

	generationText := fmt.Sprintf(" Generation: %d ", view.game.Generation)
	firstText := " ENTER: next generation, SPC: play/pause, q/ESC/Ctrl-C: quit, h: hide menu "
	secondText := " LeftClick: switch state cell, RightClick: reset board  c: create preset "
	x, y := 0, 0

	if len(firstText) <= width {
		view.renderInfo(x, y, firstText)
		y += 1
	}

	if len(secondText) <= width {
		view.renderInfo(x, y, secondText)
		y += 1
	}

	view.renderInfo(width-len(generationText), height-1, generationText)
}

// Infinite loop for the terminal view buffer where the game is executed
func (view *View) Run() {
	// Get the FPS for executing the game while is on a 'start' state
	framesPerSecond := 15
	sleepTime := time.Duration(1000/framesPerSecond) * time.Millisecond

	// Read input in another routine
	go view.readInput()

	for {
		// Keep running until the user wants to quit the game
		select {
		case <-view.quit:
			// Close the screen
			view.screen.Clear()
			view.screen.Fini()
			return
		default:
			view.displayGame()

			if view.Start {
				view.game.Step()
				time.Sleep(sleepTime)
			}

			if !view.hideMenu {
				view.displayInfo()
			}
			// Update screen
			view.screen.Show()
		}
	}
}
