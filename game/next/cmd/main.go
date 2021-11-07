package main

import (
	"fmt"
	"log"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/colornames"
)

type (
	Game struct {
		count uint8
	}
)

const (
	w, h = 640, 480
)

func (g *Game) Update() error {
	g.count++

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()
	// screen.Fill(colornames.Blanchedalmond)
	// eu.DrawLine(screen, 0, 0, 100+float64(g.count), 100, colornames.Red)

	// https://github.com/oddstream/gosol
	// https://www.memotut.com/sample-to-draw-a-simple-clock-using-ebiten-961db/
	ctx := gg.NewContext(w, h)

	ctx.Push()

	ctx.SetLineWidth(2)
	ctx.SetColor(colornames.Pink)
	ctx.DrawCircle(500, 400, 50)
	ctx.Stroke()
	ctx.SetColor(colornames.Orange)
	ctx.DrawCircle(200, 200, 100)
	ctx.Stroke()

	ctx.Pop()

	image := ebiten.NewImageFromImage(ctx.Image())
	screen.DrawImage(image, &ebiten.DrawImageOptions{})

	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f (count: %d)", ebiten.CurrentTPS(), g.count))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return w, h
}

func main() {
	ebiten.SetWindowSize(w, h)
	ebiten.SetWindowTitle("Hello, World!")

	g := Game{}

	log.Println(ebiten.RunGame(&g))
}
