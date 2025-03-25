package main

import (
	"log"

	"github.com/armadi1809/reinventing_the_wheel/chip8go/chip8"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	emulator *chip8.Chip8
}

var emulatorToKeyboardKeyMap map[int]ebiten.Key = map[int]ebiten.Key{
	0x1: ebiten.Key1,
	0x2: ebiten.Key2,
	0x3: ebiten.Key3,
	0xC: ebiten.Key4,

	0x4: ebiten.KeyQ,
	0x5: ebiten.KeyW,
	0x6: ebiten.KeyE,
	0xD: ebiten.KeyR,

	0x7: ebiten.KeyA,
	0x9: ebiten.KeyD,
	0xE: ebiten.KeyF,
	0x8: ebiten.KeyS,

	0xA: ebiten.KeyZ,
	0x0: ebiten.KeyX,
	0xB: ebiten.KeyC,
}

func (g *Game) Update() error {
	updateKeys(g.emulator)
	g.emulator.EmulateCycle()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	pixels := getPixelsFromEmulator(g.emulator)
	screen.WritePixels(pixels)
	g.emulator.DrawFlag = false
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 64, 32
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Chip 8 In Golang")
	emulator := chip8.New()
	emulator.Initialize()
	emulator.LoadProgram("./pong2.c8")
	if err := ebiten.RunGame(&Game{emulator: emulator}); err != nil {
		log.Fatal(err)
	}
}

func updateKeys(emulator *chip8.Chip8) {
	for key, val := range emulatorToKeyboardKeyMap {
		if inpututil.IsKeyJustPressed(val) {
			emulator.Key[key] = 1
		} else if inpututil.IsKeyJustReleased(val) {
			emulator.Key[key] = 0
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
