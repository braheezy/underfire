package main

import (
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

var palette = []color.Color{
	color.RGBA{0, 0, 0, 0},         // 0: Clear
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

type Game struct {
	framebuffer [][]int

	screenWidth  int
	screenHeight int
	fireWidth    int
	fireHeight   int
}

// Initialize the game.
func NewGame(w, h int) *Game {
	g := &Game{}
	g.resetFramebuffer(w, h)
	return g
}

func (g *Game) resetFramebuffer(width, height int) {
	g.screenWidth = width
	g.screenHeight = height
	g.fireWidth = width
	g.fireHeight = height
	g.framebuffer = make([][]int, g.fireHeight)
	for i := range g.framebuffer {
		g.framebuffer[i] = make([]int, g.fireWidth)
	}
	for x := 0; x < g.fireWidth; x++ {
		g.framebuffer[g.fireHeight-1][x] = len(palette) - 1
	}
}

func (g *Game) doFire() {
	for x := 0; x < g.fireWidth; x++ {
		for y := 1; y < g.fireHeight; y++ { // Start from 1 to leave the bottom row unchanged
			g.spreadFire(y*g.fireWidth + x)
		}
	}
}

func (g *Game) spreadFire(from int) {
	randOffset := rand.Intn(3) - 1 // Random value: -1, 0, or 1 (left, up, or right)
	to := from - g.fireWidth + randOffset

	// Ensure the target index is within bounds
	if to < 0 || to >= g.fireWidth*g.fireHeight {
		return
	}

	// Calculate source coordinates
	srcX := from % g.fireWidth
	srcY := from / g.fireWidth

	// Calculate target coordinates
	dstX := to % g.fireWidth
	dstY := to / g.fireWidth

	// Ensure propagation stays within the screen width
	if dstX < 0 || dstX >= g.fireWidth {
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
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}
	return nil
}

// Draw the framebuffer on the screen.
func (g *Game) Draw(screen *ebiten.Image) {
	// Use a single pixel buffer to optimize drawing
	pixelBuffer := make([]byte, g.fireWidth*g.fireHeight*4)
	for y := 0; y < g.fireHeight; y++ {
		for x := 0; x < g.fireWidth; x++ {
			intensity := g.framebuffer[y][x]
			idx := (y*g.fireWidth + x) * 4
			if intensity > 0 {
				c := palette[intensity].(color.RGBA)
				pixelBuffer[idx] = c.R
				pixelBuffer[idx+1] = c.G
				pixelBuffer[idx+2] = c.B
				pixelBuffer[idx+3] = c.A
			} else {
				pixelBuffer[idx] = 0
				pixelBuffer[idx+1] = 0
				pixelBuffer[idx+2] = 0
				pixelBuffer[idx+3] = 0
			}
		}
	}
	screen.WritePixels(pixelBuffer)
}

// Layout specifies the dimensions of the game screen.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	if outsideWidth != g.screenWidth || outsideHeight != g.screenHeight {
		g.resetFramebuffer(outsideWidth, outsideHeight)
	}
	return outsideWidth, outsideHeight
}

func main() {
	fullW, fullH := ebiten.Monitor().Size()

	fireHeight := int(fullH / 8)

	game := NewGame(fullW, fireHeight)
	ebiten.SetWindowFloating(true)
	ebiten.SetWindowSize(game.screenWidth, fireHeight)
	ebiten.SetWindowPosition(0, fullH-fireHeight)
	ebiten.SetWindowDecorated(false)

	if err := ebiten.RunGameWithOptions(game, &ebiten.RunGameOptions{ScreenTransparent: true}); err != nil {
		log.Fatal(err)
	}
}
