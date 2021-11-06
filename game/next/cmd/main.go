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

type (
	Shape struct {
		v []ebiten.Vertex
		i []uint16
		m *ebiten.DrawTrianglesOptions
	}

	Game struct {
		shapes []Shape

		shader        *ebiten.Shader
		shaderOptions *ebiten.DrawTrianglesShaderOptions
		count         uint8
		image         *ebiten.Image
		src           *ebiten.Image

		modes    []ebiten.CompositeMode
		modesInd int
	}
)

func (g *Game) Update() error {
	g.count++

	if g.count%20 == 0 {
		g.modesInd++
		g.modesInd = g.modesInd % len(g.modes)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()
	screen.Fill(colornames.Blanchedalmond)
	eu.DrawLine(screen, 0, 0, 100+float64(g.count), 100, colornames.Red)

	var (
		v []ebiten.Vertex
		i []uint16
	)

	cf := float64(g.count)
	// v, i = line(100, 100, 300, 100, colornames.Lime)
	// screen.DrawTriangles(v, i, g.src, nil)
	v, i = line(50, 150, 50, 350, colornames.Blue)
	screen.DrawTriangles(v, i, g.src, nil)

	v, i = line(50+float32(cf), 100+float32(cf), 200+float32(cf), 250, colornames.Purple)
	screen.DrawTriangles(v, i, g.src, nil)

	var opts ebiten.DrawTrianglesOptions

	for _, shape := range g.shapes {
		if shape.m != nil {
			opts = *shape.m
		} else {
			opts = ebiten.DrawTrianglesOptions{
				CompositeMode: g.modes[g.modesInd],
				FillRule:      ebiten.EvenOdd,
			}
		}

		screen.DrawTriangles(shape.v, shape.i, g.src, &opts)
	}

	screen.DrawImage(g.image, &ebiten.DrawImageOptions{})

	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f (count: %d, mode: %d)", ebiten.CurrentTPS(), g.count, g.modesInd))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")
	g := Game{}

	emptyImage := ebiten.NewImage(3, 3)
	emptyImage.Fill(color.White)
	g.src = emptyImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)

	var (
		v []ebiten.Vertex
		i []uint16
	)

	p := vector.Path{}
	p.Arc(200, 200, 100, 0, math.Pi*5/8, vector.CounterClockwise)

	v, i = p.AppendVerticesAndIndicesForFilling(nil, nil)
	g.shapes = append(g.shapes, Shape{v, i, nil})

	p = vector.Path{}
	p.Arc(400, 300, 100, 0, math.Pi*1, vector.Clockwise)
	v, i = p.AppendVerticesAndIndicesForFilling(nil, nil)
	g.shapes = append(g.shapes, Shape{v, i, nil})

	// v, i := line(100, 100, 300, 100, colornames.Lime)
	p = vector.Path{}
	p.MoveTo(100, 100)
	p.LineTo(300, 100)
	v, i = p.AppendVerticesAndIndicesForFilling(nil, nil)
	g.shapes = append(g.shapes, Shape{v, i, nil})

	// https://github.com/oddstream/gosol
	// https://www.memotut.com/sample-to-draw-a-simple-clock-using-ebiten-961db/
	ctx := gg.NewContext(640, 480)
	ctx.SetLineWidth(2)
	ctx.SetColor(colornames.Pink)
	ctx.DrawCircle(500, 400, 50)
	ctx.Stroke()
	ctx.SetColor(colornames.Orange)
	ctx.DrawCircle(200, 200, 100)
	ctx.Stroke()

	g.image = ebiten.NewImageFromImage(ctx.Image())

	g.modes = []ebiten.CompositeMode{
		// Regular alpha blending
		// c_out = c_src + c_dst × (1 - α_src)
		ebiten.CompositeModeSourceOver,
		// c_out = 0
		ebiten.CompositeModeClear,
		// c_out = c_src
		ebiten.CompositeModeCopy,
		// c_out = c_dst
		ebiten.CompositeModeDestination,
		// c_out = c_src × (1 - α_dst) + c_dst
		ebiten.CompositeModeDestinationOver,
		// c_out = c_src × α_dst
		ebiten.CompositeModeSourceIn,
		// c_out = c_dst × α_src
		ebiten.CompositeModeDestinationIn,
		// c_out = c_src × (1 - α_dst)
		ebiten.CompositeModeSourceOut,
		// c_out = c_dst × (1 - α_src)
		ebiten.CompositeModeDestinationOut,
		// c_out = c_src × α_dst + c_dst × (1 - α_src)
		ebiten.CompositeModeSourceAtop,
		// c_out = c_src × (1 - α_dst) + c_dst × α_src
		ebiten.CompositeModeDestinationAtop,
		// c_out = c_src × (1 - α_dst) + c_dst × (1 - α_src)
		ebiten.CompositeModeXor,
		// Sum of source and destination (a.k.a. 'plus' or 'additive')
		// c_out = c_src + c_dst
		ebiten.CompositeModeLighter,
		// The product of source and destination (a.k.a 'multiply blend mode')
		// c_out = c_src * c_dst
		ebiten.CompositeModeMultiply,
	}

	var (
		x, y float32 = 50, 20
	)
	const (
		w = 20
	)
	for ind := range g.modes {
		l := float32(ind)
		p = vector.Path{}
		p.MoveTo(x+w*l, y+0)
		p.LineTo(x+w*l+w, y+0)
		p.LineTo(x+w*l+w, y+w)
		p.LineTo(x+w*l, y+w)
		v, i := p.AppendVerticesAndIndicesForFilling(nil, nil)
		opts := ebiten.DrawTrianglesOptions{
			CompositeMode: g.modes[ind],
			ColorM:        ebiten.ColorM{},
		}
		// opts.ColorM.Scale(1, 0.5, 0.5, 1)
		// c := opts.ColorM.Apply(color.RGBA{R: 128, G: 0, B: 128, A: 150})
		c := color.RGBA{R: 128, G: 0, B: 128, A: 150}
		r, gr, b, a := c.RGBA()
		opts.ColorM.Scale(
			float64(r),
			float64(gr),
			float64(b),
			float64(a),
		)

		g.shapes = append(g.shapes, Shape{v, i, &opts})

	}
	y += 10
	for ind := range g.modes {
		l := float32(ind)
		p = vector.Path{}
		p.MoveTo(x+w*l, y+0)
		p.LineTo(x+w*l+w, y+0)
		p.LineTo(x+w*l+w, y+w)
		p.LineTo(x+w*l, y+w)
		v, i := p.AppendVerticesAndIndicesForFilling(nil, nil)
		opts := ebiten.DrawTrianglesOptions{
			CompositeMode: g.modes[ind],
			ColorM:        ebiten.ColorM{},
		}
		// opts.ColorM.Scale(1, 0.5, 0.5, 1)
		c := opts.ColorM.Apply(colornames.Purple)
		r, gr, b, a := c.RGBA()
		opts.ColorM.Scale(
			float64(r),
			float64(gr),
			float64(b),
			float64(a),
		)

		g.shapes = append(g.shapes, Shape{v, i, &opts})
	}

	log.Println(ebiten.RunGame(&g))
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
