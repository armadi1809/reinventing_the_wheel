package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Point struct {
	x, y, z float64
}

const WIDTH = 640
const HEIGHT = 480

var FOREGROUND_COLOR = color.RGBA{0, 128, 0, 1}

type Game struct {
	canvas *ebiten.Image
	points []Point
	dz     float64
}

func screenPosition(x, y float64) (float64, float64) {
	return (x + 1) * 0.5 * WIDTH, (1 - (y+1)*0.5) * HEIGHT
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(g.canvas, nil)
}

func (g *Game) updateAndDrawPoints() {
	for _, p := range g.points {
		newZ := p.z + g.dz
		newX, newY := p.x/newZ, p.y/newZ

		x, y := screenPosition(newX, newY)
		size := 10.0
		rect := ebiten.NewImage(int(size), int(size))
		rect.Fill(FOREGROUND_COLOR)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(x-size/2, y-size/2)
		g.canvas.DrawImage(rect, op)
	}
}

func (g *Game) Update() error {
	g.canvas.Clear()
	g.dz += 1 / float64(ebiten.TPS())
	g.updateAndDrawPoints()

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
		points: []Point{
			{0.5, 0.5, 0},
			{-0.5, 0.5, 0},
			{0.5, -0.5, 0},
			{-0.5, -0.5, 0},
		},
		dz: 0,
	}
}

func main() {

	ebiten.SetWindowSize(640, 480)
	ebiten.SetTPS(60)
	ebiten.SetWindowTitle("3D Projection Demo")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}

}
