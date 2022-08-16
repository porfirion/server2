package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/fogleman/gg"
	e "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/porfirion/server2/game/next"
	"github.com/porfirion/server2/world"
	"golang.org/x/image/colornames"
)

type (
	Game struct {
		sync.Mutex
		count       uint
		ctx         *gg.Context
		controlChan chan next.ControlMessage
		inputChan   chan next.PlayerInput
		Logic       *next.Logic
		lastState   *next.GameState
		Obj         *world.MapObject
	}
)

const (
	w, h = 1000, 1000
)

var (
	colors []color.RGBA
)

func init() {
	for i := range colornames.Names {
		if colornames.Names[i] == "black" {
			continue
		}

		colors = append(colors, colornames.Map[colornames.Names[i]])
	}
	rand.Shuffle(len(colors), func(i, j int) {
		colors[i], colors[j] = colors[j], colors[i]
	})
}

func (g *Game) Update() error {
	g.count++

	if e.IsMouseButtonPressed(e.MouseButtonLeft) {
		g.Obj.StartMoveTo(g.toPointInt(e.CursorPosition()))
	}

	if inpututil.IsMouseButtonJustReleased(e.MouseButtonLeft) {
		g.Obj.StartMoveTo(g.toPointInt(e.CursorPosition()))
	}

	return nil
}

func (g *Game) Draw(screen *e.Image) {
	g.Lock()
	defer g.Unlock()

	// screen.Clear()

	// screen.Fill(colornames.Blanchedalmond)
	// eu.DrawLine(screen, 0, 0, 100+float64(g.count), 100, colornames.Red)

	// https://github.com/oddstream/gosol
	// https://www.memotut.com/sample-to-draw-a-simple-clock-using-ebiten-961db/
	ctx := g.ctx

	ctx.Clear()
	ctx.Push()

	if g.lastState != nil {
		state := g.lastState
		// n := (g.count / 10) * (0xffffff / uint(len(state.Objects)))

		for _, obj := range state.Objects {

			// n += 0xffffff / uint(len(state.Objects))
			// ctx.SetColor(color.RGBA{
			// 	R: uint8(n & 0xff0000 >> 16),
			// 	G: uint8(n & 0xff00 >> 8),
			// 	B: uint8(n & 0xff),
			// 	A: 255,
			// })

			ctx.SetColor(colors[int(obj.Id)%len(colors)])

			ctx.DrawCircle(obj.CurrentPosition.X+500, 500-obj.CurrentPosition.Y, obj.Size)
			ctx.MoveTo(g.fromPoint(obj.CurrentPosition))
			ctx.LineTo(g.fromPoint(obj.DestinationPosition))
			ctx.Stroke()
		}
	}

	ctx.Pop()

	image := e.NewImageFromImage(ctx.Image())
	defer image.Dispose()
	screen.DrawImage(image, &e.DrawImageOptions{})

	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f (count: %d)", e.CurrentTPS(), g.count))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return w, h
}

func (g *Game) fromPoint(pos world.Point2D) (float64, float64) {
	return pos.X + 500, 500 - pos.Y
}

func (g *Game) toPoint(x, y float64) world.Point2D {
	return world.Point2D{X: x - 500, Y: 500 - y}
}
func (g *Game) toPointInt(x, y int) world.Point2D {
	return world.Point2D{X: float64(x - 500), Y: float64(500 - y)}
}

func main() {
	controlChan := make(chan next.ControlMessage)
	inputChan := make(chan next.PlayerInput)
	monitorChan := make(chan *next.GameState)

	g := Game{
		ctx:         gg.NewContext(w, h),
		controlChan: controlChan,
		inputChan:   inputChan,
		Logic:       next.NewLogic(controlChan, inputChan, next.SimulationModeContinuous, time.Second, time.Second),
	}

	var rad float64 = 100
	for i := 0; i < 12; i++ {
		obj := g.Logic.State.NewObject(world.Point2D{0, 0})
		obj.Size = 10
		angle := math.Pi * 2 * float64(i) / 12
		obj.CurrentPosition = world.Point2D{X: math.Cos(angle) * rad, Y: math.Sin(angle) * rad}
		obj.DestinationPosition = world.Point2D{X: math.Cos(angle) * (rad + 100), Y: math.Sin(angle) * (rad + 100)}
	}

	rad = 200
	for i := 0; i < 24; i++ {
		obj := g.Logic.State.NewObject(world.Point2D{0, 0})
		obj.Size = 10
		angle := math.Pi * 2 * (float64(i) + 0.5) / 24
		obj.CurrentPosition = world.Point2D{X: math.Cos(angle) * rad, Y: math.Sin(angle) * rad}
		obj.DestinationPosition = world.Point2D{X: math.Cos(angle) * (rad + 100), Y: math.Sin(angle) * (rad + 100)}
	}

	rad = 300
	for i := 0; i < 48; i++ {
		obj := g.Logic.State.NewObject(world.Point2D{0, 0})
		obj.Size = 10
		angle := math.Pi * 2 * (float64(i) + 0.5) / 48
		obj.CurrentPosition = world.Point2D{X: math.Cos(angle) * rad, Y: math.Sin(angle) * rad}
		obj.DestinationPosition = world.Point2D{X: math.Cos(angle) * (rad + 100), Y: math.Sin(angle) * (rad + 100)}
	}

	obj := g.Logic.State.WorldMap.NewObject(world.Point2D{0, 0})
	obj.DestinationPosition = world.Point2D{0, 0}
	g.Obj = obj
	g.Logic.SetMonitorChan(monitorChan)
	g.Logic.Start()

	go func() {
		for st := range monitorChan {
			g.Lock()
			g.lastState = st
			g.Unlock()
		}
	}()

	e.SetWindowSize(w, h)
	e.SetWindowTitle("Hello, World!")

	err := e.RunGame(&g)
	log.Println(err)
}
