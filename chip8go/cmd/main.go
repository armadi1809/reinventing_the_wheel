package main

import (
	"log"

	"github.com/armadi1809/reinventing_the_wheel/chip8go/chip8"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	emulator *chip8.Chip8
}

func (g *Game) Update() error {
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
