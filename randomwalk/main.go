package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Agent struct {
	posX float64
	posY float64
	col  color.Color
}
type Game struct {
	canvas *ebiten.Image
	agents []Agent
}

const SCALE = 10 // defines the velocity scale
var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

func getRandomDirection() (float64, float64) {
	dir := rng.Intn(4)
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

func (g *Game) UpdateAgent(idx int) {
	agent := &g.agents[idx]
	vx, vy := getRandomDirection()
	for range SCALE {
		agent.posX += vx
		agent.posY += vy
		rect := ebiten.NewImage(2, 2)
		rect.Fill(agent.col)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(agent.posX, agent.posY)

		g.canvas.DrawImage(rect, op)
	}
}

func (g *Game) Update() error {
	for i := range g.agents {
		g.UpdateAgent(i)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(g.canvas, nil)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func generateRandomColor() color.Color {
	return color.RGBA{
		R: uint8(rng.Intn(256)),
		G: uint8(rng.Intn(256)),
		B: uint8(rng.Intn(256)),
		A: 255,
	}
}

func NewGame(numAgents int) *Game {
	canvas := ebiten.NewImage(640, 480)
	canvas.Fill(color.Black)

	agents := []Agent{}

	for range numAgents {
		agent := Agent{
			posX: 320,
			posY: 240,
			col:  generateRandomColor(),
		}
		agents = append(agents, agent)

	}

	return &Game{
		canvas: canvas,
		agents: agents,
	}
}

func main() {
	numAgents := 5
	args := os.Args
	if len(args) > 2 {
		fmt.Printf("Usage: %s <num-of-agents>\n", args[0])
		os.Exit(-1)
	}
	if len(args) == 2 {
		customNumAgents, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Printf("Invalid number for agents: %s", args[1])
			os.Exit(-1)
		}
		numAgents = customNumAgents
	}

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle(fmt.Sprintf("Random Walk: %d agents", numAgents))
	if err := ebiten.RunGame(NewGame(numAgents)); err != nil {
		log.Fatal(err)
	}

}
