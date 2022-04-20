# Game of Life

This is yet another Conway's Game of Life interpretation made in Golang on the terminal.
It's not the best but I just did it to learn a little bit more of the Go programming
language.
More about the game [Game of Life](https://en.wikipedia.org/wiki/Conway's_Game_of_Life).

## Key binds

   |          Key bind          |         Description                  |
   | -------------------------- | ------------------------------------ |
   | `Space`                    | Play/Pause the board                 |
   | `Enter`                    | Move to the next generation on pause |
   | `LeftClick`                | Toggle cell state spawn/kill on pause|
   | `RightClick`               | Clear board                          |
   | `q` or `Esc` or `Ctrl + c` | Exit the game                        |
   | `h`                        | Hide information menu                |

## Execution
Create the executable with:
```bash
go build
```
There are flags so you can especify the behavior of the game. You can see all the
option by doing.

```bash
# On the executable
./gameOfLife -h

# On runtime
go run ./main.go -h
```

## TODO
- [ ] Option to wrap edges of the window (snake-like game edges)
- [ ] Make presets like the glider, glider gun, pulsar etc...
- [ ] Save custom presets

## Inspired
- [go-life](https://github.com/sachaos/go-life)
- [Michael Abrash's Graphics Programming Black Book (1997)](http://www.jagregory.com/abrash-black-book/)
