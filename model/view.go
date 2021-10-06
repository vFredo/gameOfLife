package model

import (
	"log"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

// Styles (terminal color theme)
var (
	AliveStylePlay  = tcell.StyleDefault.Background(tcell.ColorYellow).Foreground(tcell.ColorReset)
	AliveStylePause = tcell.StyleDefault.Background(tcell.ColorBlue).Foreground(tcell.ColorReset)
	DefaultStyle    = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
)

type View struct {
	screen tcell.Screen
	game   GameOfLife
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

	//Initialize Game
	view.game.Init(height, width/2)
}

// Control key strokes and mouse cliks events
func (view *View) readInput() {
	quit := func() {
		view.screen.Fini()
		os.Exit(0)
	}
	// Poll event catch events on the buffer
	ev := view.screen.PollEvent()

	// Process events catch
	switch ev := ev.(type) {
	case *tcell.EventResize:
		view.screen.Sync()
	case *tcell.EventKey:
		if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC || ev.Rune() == 'q' {
			quit()
		} else if ev.Rune() == ' ' {
			view.game.Start = !view.game.Start
		} else if ev.Rune() == 'l' {
			if !view.game.Start {
				view.game.Step()
			}
		}
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

// How each cell is render on the terminal buffer
func (view *View) displayGame() {
	view.screen.Clear()

	// Death: Background color, AlivePause: Blue, AlivePlay: Yellow
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

// Infinite loop for the terminal view buffer where the game is executed
func (view *View) Loop() {
	view.initScreen()

	framesPerSecond := 15
	sleepTime := time.Duration(1000/framesPerSecond) * time.Millisecond
	for {
		view.readInput()

		if view.game.Start {
			view.game.Step()
			time.Sleep(sleepTime)
		}
		view.displayGame()
		// Update screen
		view.screen.Show()
	}
}
