package main

import (
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Point struct {
	x, y, z float64
}

type Edge struct {
	from, to int
}

const WIDTH = 640
const HEIGHT = 480

var FOREGROUND_COLOR = color.RGBA{0, 128, 0, 1}

var cubeEdges = []Edge{
	// Front face
	{0, 1}, {1, 3}, {3, 2}, {2, 0},
	// Back face
	{4, 5}, {5, 7}, {7, 6}, {6, 4},
	// Connecting edges
	{0, 4}, {1, 5}, {2, 6}, {3, 7},
}

type Game struct {
	canvas *ebiten.Image
	points []Point
	dz     float64
	angle  float64
}

func screenPosition(x, y float64) (float64, float64) {
	return (x + 1) * 0.5 * WIDTH, (1 - (y+1)*0.5) * HEIGHT
}

func rotate_xz(p Point, angle float64) Point {
	c := math.Cos(angle)
	s := math.Sin(angle)
	return Point{
		x: p.x*c - p.z*s,
		y: p.y,
		z: p.x*s + p.z*c,
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(g.canvas, nil)
}

func (g *Game) updateAndDrawPoints() {
	screenPoints := make([][2]float64, len(g.points))
	for i, p := range g.points {
		np := rotate_xz(p, g.angle)
		newZ := np.z + g.dz
		newX, newY := np.x/newZ, np.y/newZ
		x, y := screenPosition(newX, newY)
		screenPoints[i] = [2]float64{x, y}
	}

	for _, edge := range cubeEdges {
		p1 := screenPoints[edge.from]
		p2 := screenPoints[edge.to]
		vector.StrokeLine(g.canvas,
			float32(p1[0]), float32(p1[1]),
			float32(p2[0]), float32(p2[1]),
			2, // stroke width
			FOREGROUND_COLOR,
			false, // anti-alias
		)
	}

	// Draw vertices as circles
	for _, sp := range screenPoints {
		vector.FillCircle(g.canvas,
			float32(sp[0]), float32(sp[1]),
			5, // radius
			FOREGROUND_COLOR,
			false, // anti-alias
		)
	}
}

func (g *Game) Update() error {
	dt := 1 / float64(ebiten.TPS())
	// g.dz += dt
	g.angle += math.Pi * dt
	g.canvas.Clear()
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
			{0.25, 0.25, 0.25},
			{-0.25, 0.25, 0.25},
			{0.25, -0.25, 0.25},
			{-0.25, -0.25, 0.25},

			{0.25, 0.25, -0.25},
			{-0.25, 0.25, -0.25},
			{0.25, -0.25, -0.25},
			{-0.25, -0.25, -0.25},
		},
		dz:    1,
		angle: 0,
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
