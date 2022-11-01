package view

import (
	"time"

	"github.com/gdamore/tcell/v2"
)

// Control input (mouse/keyboard) events on the screen
func readInput(view *TermView) {
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
					inputFormPreset(view)
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
func inputFormPreset(view *TermView) {
	namePreset := ""

	for {
		renderInput(namePreset, view)

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

// Helper method to render the input box on the screen
func renderInput(currName string, view *TermView) {
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
