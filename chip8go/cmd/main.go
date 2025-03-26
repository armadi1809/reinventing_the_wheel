package main

import (
	"log"
	"os"

	"github.com/armadi1809/reinventing_the_wheel/chip8go/chip8"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	emulator *chip8.Chip8
}

const clockRate = 10

var keyboardToEmulatorMap map[ebiten.Key]int = map[ebiten.Key]int{
	ebiten.Key1: 0x1,
	ebiten.Key2: 0x2,
	ebiten.Key3: 0x3,
	ebiten.Key4: 0xC,

	ebiten.KeyQ: 0x4,
	ebiten.KeyW: 0x5,
	ebiten.KeyE: 0x6,
	ebiten.KeyR: 0xD,

	ebiten.KeyA: 0x7,
	ebiten.KeyS: 0x8,
	ebiten.KeyD: 0x9,
	ebiten.KeyF: 0xE,

	ebiten.KeyZ: 0xA,
	ebiten.KeyX: 0x0,
	ebiten.KeyC: 0xB,
	ebiten.KeyV: 0xF,
}

func (g *Game) Update() error {
	updateKeys(g.emulator)
	for range clockRate {
		g.emulator.EmulateCycle()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.emulator.DrawFlag {
		pixels := getPixelsFromEmulator(g.emulator)
		screen.WritePixels(pixels)
		g.emulator.DrawFlag = false
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 64, 32
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Invalid Arguments. Run the emulator with one rom file path")
	}
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Chip 8 In Golang")
	emulator := chip8.New()
	emulator.Initialize()
	emulator.LoadProgram(os.Args[1])
	if err := ebiten.RunGame(&Game{emulator: emulator}); err != nil {
		log.Fatal(err)
	}
}

func updateKeys(emulator *chip8.Chip8) {
	for key, val := range keyboardToEmulatorMap {
		if ebiten.IsKeyPressed(key) || inpututil.IsKeyJustPressed(key) {
			emulator.Key[val] = 1
		} else {
			emulator.Key[val] = 0
		}
	}
}

func getPixelsFromEmulator(emulator *chip8.Chip8) []byte {
	width := 64
	height := 32
	// Create a slice to hold the RGBA values
	// Length = width * height * 4 (4 values per pixel: R, G, B, A)
	rgbaArray := make([]byte, width*height*4)
	gfx := emulator.Gfx
	// Fill the array with random RGBA values
	for i := range len(gfx) {
		if gfx[i] == 0 {
			rgbaArray[i*4] = 0
			rgbaArray[i*4+1] = 0
			rgbaArray[i*4+2] = 0
			rgbaArray[i*4+3] = 255
		} else {
			rgbaArray[i*4] = 255
			rgbaArray[i*4+1] = 255
			rgbaArray[i*4+2] = 255
			rgbaArray[i*4+3] = 255
		}

	}
	return rgbaArray

}
