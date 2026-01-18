package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

var FOREGROUND_COLOR = color.RGBA{0, 128, 0, 1}

type Game struct {
	canvas *ebiten.Image
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(g.canvas, nil)
}

func (g *Game) Update() error {

	rect := ebiten.NewImage(10, 10)
	rect.Fill(FOREGROUND_COLOR)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(100, 100)

	g.canvas.DrawImage(rect, op)

	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 480
}

func NewGame() *Game {
	canvas := ebiten.NewImage(640, 480)
	canvas.Fill(color.Black)
	return &Game{
		canvas: canvas,
	}
}

func main() {

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("3D Projection Demo")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}

}
