# Game of Life

This is yet another Conway's Game of Life interpretation made in Golang on the terminal.
It's not the best but I just did it to learn a little bit more of the Go programing language.

## Keybinds

   |          Keybind          |         Description                  |
   | ------------------------- | ------------------------------------ |
   | `Space`                   | Play/Pause the execution             |
   | `l`                       | Move to the next generation on pause |
   | `left click`              | Toggle cell dead/alive on pause      |
   | `right click`             | Clear board                          |
   | `q` or `Esc` or `Ctrl + c`| Exit the game                        |

## TODO
- [] When all the cells die restart to a new instance of the game
- [] use go rutines for the keyboard/mouse events and the loop of the game
- [] More efficient way to check if a cell is alive or dead
- [] Fix when the game is play the cells pause for some frames and then update fast

