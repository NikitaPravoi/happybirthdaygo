package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	targetText = "happy birthday go!"
	width      = 50
	height     = 20
)

type Drop struct {
	x, y    int
	char    rune
	fixed   bool
	target  bool
	targetX int
}

func main() {
	// Clear screen and hide cursor
	fmt.Print("\033[2J\033[?25l")
	defer fmt.Print("\033[?25h") // Show cursor when done

	// Initialize target positions
	targetPositions := make(map[int]rune)
	startX := (width - len(targetText)) / 2
	for i, char := range targetText {
		targetPositions[startX+i] = char
	}

	// Initialize drops
	var drops []Drop
	active := make(map[int]bool)

	// Random characters to use in animation
	chars := "abcdefghijklmnopqrstuvwxyz!/\\1234567890+- "

	// Main animation loop
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for {
		// Add new drops randomly
		if rand.Float64() < 0.3 {
			x := rand.Intn(width)
			if !active[x] {
				targetChar, isTarget := targetPositions[x]
				drops = append(drops, Drop{
					x:       x,
					y:       0,
					char:    rune(chars[rand.Intn(len(chars))]),
					target:  isTarget,
					targetX: x,
					fixed:   false,
				})
				if isTarget {
					drops[len(drops)-1].char = targetChar
				}
				active[x] = true
			}
		}

		// Create screen buffer
		screen := make([][]rune, height)
		for i := range screen {
			screen[i] = make([]rune, width)
			for j := range screen[i] {
				screen[i][j] = ' '
			}
		}

		// Update drops
		newDrops := []Drop{}
		for _, drop := range drops {
			if !drop.fixed {
				if drop.y < height-1 {
					if !drop.target {
						drop.char = rune(chars[rand.Intn(len(chars))])
					}
					drop.y++
					newDrops = append(newDrops, drop)
				} else {
					if drop.target {
						drop.fixed = true
						newDrops = append(newDrops, drop)
					} else {
						active[drop.x] = false
					}
				}
			} else {
				newDrops = append(newDrops, drop)
			}
			screen[drop.y][drop.x] = drop.char
		}
		drops = newDrops

		// Render screen
		fmt.Print("\033[H") // Move cursor to top-left
		for _, row := range screen {
			fmt.Println(string(row))
		}

		// Check if animation is complete
		complete := true
		for x, targetChar := range targetPositions {
			if screen[height-1][x] != targetChar {
				complete = false
				break
			}
		}
		if complete {
			time.Sleep(500 * time.Millisecond)
			// gopher ASCII art
			gopher := "         ,_---~~~~~----._         \n  _,,_,*^____      _____``*g*\\\"*, \n / __/ /'     ^.  /      \\ ^@q   f \n[  @f | @))    |  | @))   l  0 _/  \n \\`/   \\~____ / __ \\_____/    \\   \n  |           _l__l_           I   \n  }          [______]           I  \n  ]            | | |            |  \n  ]             ~ ~             |  \n  |                            |   \n   |                           |   \n"
			fmt.Println(gopher)
			break
		}

		<-ticker.C
	}
}
