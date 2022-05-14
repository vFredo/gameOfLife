# Game of Life

Yet another Conway's Game of Life implementation made in Go on the terminal.
It's not the best but I just did it to learn a little bit more of the Go programming
language. Read more about the [Game of Life](https://en.wikipedia.org/wiki/Conway's_Game_of_Life).

## Key binds

   |          Key bind          |         Description                           |
   | -------------------------- | ----------------------------------------------|
   | `Space`                    | Play/Pause the board                          |
   | `Enter`                    | Move to the next generation on pause          |
   | `LeftClick`                | Toggle cell state spawn/kill on pause         |
   | `RightClick`               | Clear board                                   |
   | `q` or `Esc` or `Ctrl + c` | Exit the game                                 |
   | `h` or `H`                 | Hide information menu or all of it            |
   | `p`                        | Create a preset on pause of the current board |
   | `c`                        | Cycle through the presets available           |

## Execution
Create the executable with:
```bash
go build
```
There are some flags as parameter for the executable so you can specify some details on
the behavior of the game. You can see all the options by doing:
```bash
# On the executable
./gameOfLife -h

# On runtime
go run ./main.go -h
```

## TODO
- [ ] The user when saving a preset can give a name to it on the TUI
- [ ] Show on TUI the current preset

## Inspired
- [go-life](https://github.com/sachaos/go-life)
- [Michael Abrash's Graphics Programming Black Book (1997)](http://www.jagregory.com/abrash-black-book/)
