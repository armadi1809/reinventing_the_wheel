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
	g.emulator.EmulateCycle()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.emulator.DrawFlag {
		// TODO get pixels from emulator graphics
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 160, 120
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Chip 8 In Golang")
	emulator := chip8.New()
	emulator.Initialize()
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}

func updateKeys(emulator *chip8.Chip8) {
	pressedKeys := inpututil.AppendJustPressedKeys(nil)
	for _, key := range pressedKeys {
		emulator.Key[keyboardToEmulatorMap[key]] = 1
	}

	releasedKeys := inpututil.AppendJustReleasedKeys(nil)
	for _, key := range releasedKeys {
		emulator.Key[keyboardToEmulatorMap[key]] = 1
	}
}
