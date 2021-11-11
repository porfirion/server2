package main

import (
	"fmt"
	"image/color"
	"log"
	"sync"
	"time"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/porfirion/server2/game/next"
)

type (
	Game struct {
		sync.Mutex
		count       uint
		ctx         *gg.Context
		controlChan chan next.ControlMessage
		inputChan   chan next.PlayerInput
		logic       *next.LogicImpl
		lastState   next.GameState
	}
)

const (
	w, h = 1000, 1000
)

func (g *Game) Update() error {
	g.count++

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.Lock()
	defer g.Unlock()

	// screen.Clear()

	// screen.Fill(colornames.Blanchedalmond)
	// eu.DrawLine(screen, 0, 0, 100+float64(g.count), 100, colornames.Red)

	// https://github.com/oddstream/gosol
	// https://www.memotut.com/sample-to-draw-a-simple-clock-using-ebiten-961db/
	ctx := g.ctx

	ctx.Push()

	if g.lastState != nil {
		state := g.lastState.(*next.GameStateImpl)
		n := (g.count / 10) * (0xffffff / uint(len(state.Objects)))

		for _, obj := range state.Objects {
			n += 0xffffff / uint(len(state.Objects))

			ctx.SetColor(color.RGBA{
				R: uint8(n & 0xff0000 >> 16),
				G: uint8(n & 0xff00 >> 8),
				B: uint8(n & 0xff),
				A: 255,
			})

			ctx.DrawCircle(obj.CurrentPosition.X+500, -obj.CurrentPosition.Y+500, obj.Size)
			ctx.Stroke()
		}
	}

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

	controlChan := make(chan next.ControlMessage)
	inputChan := make(chan next.PlayerInput)
	monitorChan := make(chan next.GameState)

	g := Game{
		ctx:         gg.NewContext(w, h),
		controlChan: controlChan,
		inputChan:   inputChan,
		logic:       next.NewLogic(controlChan, inputChan, next.SimulationModeContinuous, time.Second, time.Second),
	}

	g.logic.SetMonitorChan(monitorChan)
	g.logic.Start()

	go func() {
		for st := range monitorChan {
			g.Lock()
			g.lastState = st
			g.Unlock()
		}
	}()

	err := ebiten.RunGame(&g)
	log.Println(err)
}
