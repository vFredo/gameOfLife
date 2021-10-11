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

// View structure that has the game itself and the screen where everything is render
type View struct {
	screen   tcell.Screen
	game     GameOfLife
	hideMenu bool
}

// Initialize screen view and game
func (view *View) initScreen() {
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	view.screen = s

	if err := view.screen.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	view.screen.SetStyle(DefaultStyle)
	view.screen.EnableMouse()
	view.screen.Clear()

	width, height := view.screen.Size()

	// Initialize Game
	view.game.Init(height, width/2)
}

// Control input events
func (view *View) readInput() {
	for {
		// Poll event catch events on the buffer
		ev := view.screen.PollEvent()

		// Process events catch
		switch ev := ev.(type) {
		// Resize window event
		case *tcell.EventResize:
			view.screen.Sync()
			width, height := view.screen.Size()
			view.game.Resize(height, width/2)
			// Keyboard events
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC || ev.Rune() == 'q' {
				view.screen.Fini()
				os.Exit(0)
			} else if ev.Rune() == ' ' {
				view.game.Start = !view.game.Start
			} else if ev.Key() == tcell.KeyEnter {
				if !view.game.Start {
					view.game.Step()
				}
			} else if ev.Rune() == 'h' {
				view.hideMenu = !view.hideMenu
			}
			// Mouse events
		case *tcell.EventMouse:
			switch ev.Buttons() {
			case tcell.Button1:
				x, y := ev.Position()
				// If the game is in pause, let it modified, else don't
				if !view.game.Start {
					rows, cols := y, x/2
					if rows < view.game.X && cols < view.game.Y {
						if view.game.CurrentGen[rows][cols] == ALIVE {
							view.game.CurrentGen[rows][cols] = DEAD
						} else {
							view.game.CurrentGen[rows][cols] = ALIVE
						}
					}
				}
			case tcell.Button2:
				view.game.ClearGame()
			}
		}
	}
}

// How each cell is render on the terminal buffer
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

// helper that takes the string (info) and pos (x, y) to put it on the screen
func (view *View) prepareStringInfo(x int, y int, info string) {
	for i, byte := range info {
		view.screen.SetCell(x+i, y, InfoStyle, byte)
	}
}

// Information show in the menu
func (view *View) displayInfo() {
	_, height := view.screen.Size()
	view.prepareStringInfo(0, 0, "ENTER: next generation, SPC: play/pause, q/ESC/Ctrl-C: quit, h: hide menu")
	view.prepareStringInfo(0, 1, "LeftClick: switch state cell, RightClick: clear board")
	view.prepareStringInfo(0, height-1, fmt.Sprintf("Generation: %d ", view.game.Generation))
}

// Infinite loop for the terminal view buffer where the game is executed
func (view *View) Loop() {
	view.initScreen()

	framesPerSecond := 15
	sleepTime := time.Duration(1000/framesPerSecond) * time.Millisecond
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
