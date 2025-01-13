package main

import (
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 320
	screenHeight = 200
	fireWidth    = screenWidth
	fireHeight   = screenHeight
)

type Game struct {
	framebuffer [][]int
	palette     []color.Color
}

// Initialize the game.
func NewGame() *Game {
	g := &Game{
		framebuffer: make([][]int, fireHeight),
		palette:     generatePalette(),
	}
	// Initialize framebuffer rows
	for i := range g.framebuffer {
		g.framebuffer[i] = make([]int, fireWidth)
	}
	// Set the bottom row to the maximum intensity (36)
	for x := 0; x < fireWidth; x++ {
		g.framebuffer[fireHeight-1][x] = 36
	}
	return g
}

// Generate a placeholder palette with 37 shades of white for now.
func generatePalette() []color.Color {
	return []color.Color{
		color.RGBA{0, 0, 0, 255},       // 0: Black
		color.RGBA{7, 7, 7, 255},       // 1: Dark Gray
		color.RGBA{31, 7, 7, 255},      // 2: Dark Red
		color.RGBA{47, 15, 7, 255},     // 3: Red-Brown
		color.RGBA{71, 15, 7, 255},     // 4: Brown
		color.RGBA{87, 23, 7, 255},     // 5: Orange-Brown
		color.RGBA{103, 31, 7, 255},    // 6: Dark Orange
		color.RGBA{119, 31, 7, 255},    // 7: Orange
		color.RGBA{143, 39, 7, 255},    // 8: Bright Orange
		color.RGBA{159, 47, 7, 255},    // 9: Light Orange
		color.RGBA{175, 63, 7, 255},    // 10: Lighter Orange
		color.RGBA{191, 71, 7, 255},    // 11: Reddish Orange
		color.RGBA{199, 71, 7, 255},    // 12: Bright Reddish Orange
		color.RGBA{223, 79, 7, 255},    // 13: Red-Orange
		color.RGBA{223, 87, 7, 255},    // 14: Fiery Orange
		color.RGBA{223, 87, 7, 255},    // 15: Fiery Orange (repeat for gradient)
		color.RGBA{215, 95, 7, 255},    // 16: Yellowish Orange
		color.RGBA{215, 103, 15, 255},  // 17: Yellow-Orange
		color.RGBA{207, 111, 15, 255},  // 18: Yellow
		color.RGBA{207, 119, 15, 255},  // 19: Bright Yellow
		color.RGBA{207, 127, 15, 255},  // 20: Pale Yellow
		color.RGBA{207, 135, 23, 255},  // 21: Pale Yellowish White
		color.RGBA{199, 135, 23, 255},  // 22: Lighter Yellowish White
		color.RGBA{199, 143, 23, 255},  // 23: Yellow-White
		color.RGBA{199, 151, 31, 255},  // 24: Pale Yellow-White
		color.RGBA{191, 159, 31, 255},  // 25: Bright Yellow-White
		color.RGBA{191, 159, 31, 255},  // 26: Bright Yellow-White (repeat for gradient)
		color.RGBA{191, 167, 39, 255},  // 27: Very Bright Yellow
		color.RGBA{191, 167, 39, 255},  // 28: Very Bright Yellow (repeat for gradient)
		color.RGBA{255, 255, 63, 255},  // 29: Near White
		color.RGBA{255, 255, 111, 255}, // 30: Faint Yellowish White
		color.RGBA{255, 255, 159, 255}, // 31: Faint Pale White
		color.RGBA{255, 255, 191, 255}, // 32: Pale White
		color.RGBA{255, 255, 223, 255}, // 33: Near Full White
		color.RGBA{255, 255, 239, 255}, // 34: Almost Full White
		color.RGBA{255, 255, 247, 255}, // 35: Subtle White
		color.RGBA{255, 255, 255, 255}, // 36: Full White
	}
}

func (g *Game) doFire() {
	for x := 0; x < fireWidth; x++ {
		for y := 1; y < fireHeight; y++ { // Start from 1 to leave the bottom row unchanged
			g.spreadFire(y*fireWidth + x)
		}
	}
}

func (g *Game) spreadFire(from int) {
	randOffset := rand.Intn(3) - 1 // Random value: -1, 0, or 1 (left, up, or right)
	to := from - fireWidth + randOffset

	// Ensure the target index is within bounds
	if to < 0 || to >= fireWidth*fireHeight {
		return
	}

	// Calculate source coordinates
	srcX := from % fireWidth
	srcY := from / fireWidth

	// Calculate target coordinates
	dstX := to % fireWidth
	dstY := to / fireWidth

	// Ensure propagation stays within the screen width
	if dstX < 0 || dstX >= fireWidth {
		return
	}

	// Reduce the heat, ensuring it doesn't go below zero
	newValue := g.framebuffer[srcY][srcX] - rand.Intn(2) // Random decay: 0 or 1
	if newValue < 0 {
		newValue = 0
	}
	g.framebuffer[dstY][dstX] = newValue
}

// Update is called once every frame to update the game logic.
func (g *Game) Update() error {
	g.doFire()
	return nil
}

// Draw the framebuffer on the screen.
func (g *Game) Draw(screen *ebiten.Image) {
	// Draw each pixel based on the framebuffer and palette
	for y := 0; y < fireHeight; y++ {
		for x := 0; x < fireWidth; x++ {
			intensity := g.framebuffer[y][x]
			screen.Set(x, y, g.palette[intensity])
		}
	}
}

// Layout specifies the dimensions of the game screen.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight // Game resolution
}

func main() {
	game := NewGame()
	ebiten.SetWindowTitle("Doom Fire Animation")
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2) // Scale window for visibility
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
