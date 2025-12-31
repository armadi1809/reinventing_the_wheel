package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"os"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	rectX  float64
	rectY  float64
	canvas *ebiten.Image
}

func getRandomDirection() (float64, float64) {
	dir := rand.Intn(4)
	switch dir {
	case 0: // UP
		return 0, -1
	case 1: // DOWN
		return 0, 1
	case 2: // RIGHT
		return 1, 0
	default: // LEFT
		return -1, 0
	}
}

func (g *Game) Update() error {
	vx, vy := getRandomDirection()
	for range 10 {
		g.rectX += vx
		g.rectY += vy
		rect := ebiten.NewImage(2, 2)
		rect.Fill(color.White)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(g.rectX), float64(g.rectY))

		g.canvas.DrawImage(rect, op)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(g.canvas, nil)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func NewGame() *Game {
	canvas := ebiten.NewImage(640, 480)
	canvas.Fill(color.Black)

	return &Game{
		canvas: canvas,
		rectX:  320,
		rectY:  240,
	}
}

func main() {
	numAgents := 5
	args := os.Args
	if len(args) > 2 {
		fmt.Printf("Usage: %s <num-of-agents>\n", args[0])
		os.Exit(1)
	}
	if len(args) == 2 {
		customNumAgents, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Printf("Invalid number for agents: %s", args[1])
		}
		numAgents = customNumAgents
	}

	fmt.Println(numAgents)

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle(fmt.Sprintf("Hello, agents: %d", numAgents))
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}

}
