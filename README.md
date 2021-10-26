# Game of Life

This is yet another Conway's Game of Life interpretation made in Golang on the terminal.
It's not the best but I just did it to learn a little bit more of the Go programming
language.
More about the game [Game of Life](https://en.wikipedia.org/wiki/Conway's_Game_of_Life).

## Key binds

   |          Key bind          |         Description                  |
   | -------------------------- | ------------------------------------ |
   | `Space`                    | Play/Pause the execution             |
   | `Enter`                    | Move to the next generation on pause |
   | `LeftClick`                | Toggle cell dead/alive on pause      |
   | `RightClick`               | Clear board                          |
   | `q` or `Esc` or `Ctrl + c` | Exit the game                        |
   | `h`                        | Hide menu                            |

## TODO
- [ ] More efficient way to check if a cell is alive or dead
- [ ] Make presets like the glider, glider gun, pulsar etc...
- [ ] Make available other rules, not only B3S23 (Birth: 3, Survival: 2 or 3)
- [ ] Make menu information responsive

## Inspired
- [go-life](https://github.com/sachaos/go-life)
