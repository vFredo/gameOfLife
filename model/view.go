package model

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

// Styles (uses the terminal colors)
var (
	AliveStylePlay  = tcell.StyleDefault.Background(tcell.ColorYellow).Foreground(tcell.ColorReset)
	AliveStylePause = tcell.StyleDefault.Background(tcell.ColorBlue).Foreground(tcell.ColorReset)
	DefaultStyle    = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	InfoStyle       = tcell.StyleDefault.Background(tcell.ColorGray).Foreground(tcell.ColorWhite)
)

// View structure that has the game itself and the screen (terminal buffer) where everything is render
type View struct {
	screen   tcell.Screen
	game     GameOfLife
	hideMenu bool
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

	view.screen.SetStyle(DefaultStyle)
	view.screen.EnableMouse()
	view.screen.Clear()

	// Adding the game to the view
	view.game = game

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
				view.screen.Clear()
				view.screen.Fini()
				os.Exit(0)
			} else if ev.Rune() == ' ' { // space
				view.game.Start = !view.game.Start
			} else if ev.Key() == tcell.KeyEnter && !view.game.Start {
				view.game.Step()
			} else if ev.Rune() == 'h' {
				view.hideMenu = !view.hideMenu
			}
		case *tcell.EventMouse:
			switch ev.Buttons() {
			case tcell.Button1: // left click
				x, y := ev.Position()
				row, col := y, x/2
				// If the game is in pause, let it modified, else don't
				if row < view.game.X && col < view.game.Y && !view.game.Start {
					if view.game.CurrentGen[row][col] == ALIVE {
						view.game.CurrentGen[row][col] = DEAD
					} else {
						view.game.CurrentGen[row][col] = ALIVE
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
			if view.game.CurrentGen[i][j] == ALIVE && view.game.Start {
				view.screen.SetContent(j*2, i, ' ', nil, AliveStylePlay)
				view.screen.SetContent(j*2+1, i, ' ', nil, AliveStylePlay)
			} else if view.game.CurrentGen[i][j] == ALIVE {
				view.screen.SetContent(j*2, i, ' ', nil, AliveStylePause)
				view.screen.SetContent(j*2+1, i, ' ', nil, AliveStylePause)
			} else {
				view.screen.SetContent(j*2, i, ' ', nil, DefaultStyle)
				view.screen.SetContent(j*2+1, i, ' ', nil, DefaultStyle)
			}
		}
	}
}

// Helper method that takes the string (info) and pos (x, y) to put it on the screen
func (view *View) renderInfo(x int, y int, info string) {
	for i, byte := range info {
		view.screen.SetCell(x+i, y, InfoStyle, byte)
	}
}

// Information shown on the menu
func (view *View) displayInfo() {
	width, height := view.screen.Size()
	view.renderInfo(0, 0, " ENTER: next generation, SPC: play/pause, q/ESC/Ctrl-C: quit, h: hide menu ")
	view.renderInfo(0, 1, " LeftClick: switch state cell, RightClick: clear board ")
	genString := fmt.Sprintf(" Generation: %d ", view.game.Generation)
	view.renderInfo(width-len(genString), height-1, genString)
}

// Infinite loop for the terminal view buffer where the game is executed
func (view *View) Run() {
	// Get the FPS for executing the game while is on 'play' taking into account the refresh rate of the screen
	framesPerSecond := 15
	sleepTime := time.Duration(1000/framesPerSecond) * time.Millisecond

	// Read input in another routine, not the exact same as where the game is executed
	go view.readInput()

	for {
		view.displayGame()
		if view.game.Start {
			view.game.Step()
			time.Sleep(sleepTime)
		} else if !view.hideMenu {
			view.displayInfo()
		}
		// Update screen
		view.screen.Show()
	}
}
