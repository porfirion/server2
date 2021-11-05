package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	eu "github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
)

type Game struct {
	initialized bool

	vertices      []ebiten.Vertex
	indeces       []uint16
	shader        *ebiten.Shader
	shaderOptions *ebiten.DrawTrianglesShaderOptions
	count         uint8
	image         *ebiten.Image
}

var (
	emptyImage = ebiten.NewImage(3, 3)
)

func init() {
	emptyImage.Fill(color.White)
}

func (g *Game) Update() error {
	// log.Println("update")

	if !g.initialized {
		p := vector.Path{}
		p.Arc(200, 200, 100, 0, math.Pi*1, vector.Clockwise)
		g.vertices, g.indeces = p.AppendVerticesAndIndicesForFilling(nil, nil)

		// https://github.com/oddstream/gosol
		// https://www.memotut.com/sample-to-draw-a-simple-clock-using-ebiten-961db/
		ctx := gg.NewContext(640, 480)
		ctx.SetLineWidth(2)
		ctx.SetColor(colornames.Pink)
		ctx.DrawCircle(500, 400, 50)
		ctx.Stroke()

		g.image = ebiten.NewImageFromImage(ctx.Image())
	}

	g.count++

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()

	eu.DrawLine(screen, 0, 0, 100+float64(g.count), 100, colornames.Red)

	src := emptyImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)

	cf := float64(g.count)
	v, i := line(100, 100, 300, 100, colornames.Lime)
	screen.DrawTriangles(v, i, src, nil)
	v, i = line(50, 150, 50, 350, colornames.Blue)
	screen.DrawTriangles(v, i, src, nil)
	v, i = line(50+float32(cf), 100+float32(cf), 200+float32(cf), 250, colornames.Purple)
	screen.DrawTriangles(v, i, src, nil)

	opts := ebiten.DrawTrianglesOptions{
		FillRule: ebiten.EvenOdd,
	}

	screen.DrawTriangles(g.vertices, g.indeces, src, &opts)

	screen.DrawImage(g.image, &ebiten.DrawImageOptions{})

	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f (count: %d)", ebiten.CurrentTPS(), g.count))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")
	err := ebiten.RunGame(&Game{})
	if err != nil {
		log.Fatal(err)
	}
}

func line(x0, y0, x1, y1 float32, clr color.RGBA) ([]ebiten.Vertex, []uint16) {
	const width = 1

	theta := math.Atan2(float64(y1-y0), float64(x1-x0))
	theta += math.Pi / 2
	dx := float32(math.Cos(theta))
	dy := float32(math.Sin(theta))

	r := float32(clr.R) / 0xff
	g := float32(clr.G) / 0xff
	b := float32(clr.B) / 0xff
	a := float32(clr.A) / 0xff

	return []ebiten.Vertex{
		{
			DstX:   x0 - width*dx/2,
			DstY:   y0 - width*dy/2,
			SrcX:   1,
			SrcY:   1,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
		{
			DstX:   x0 + width*dx/2,
			DstY:   y0 + width*dy/2,
			SrcX:   1,
			SrcY:   1,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
		{
			DstX:   x1 - width*dx/2,
			DstY:   y1 - width*dy/2,
			SrcX:   1,
			SrcY:   1,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
		{
			DstX:   x1 + width*dx/2,
			DstY:   y1 + width*dy/2,
			SrcX:   1,
			SrcY:   1,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
	}, []uint16{0, 1, 2, 1, 2, 3}
}
