package main

import (
	"testing"

	"github.com/fogleman/gg"
	"golang.org/x/image/colornames"
)

func BenchmarkGame_Draw(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ctx := gg.NewContext(w, h)

		ctx.Push()

		ctx.SetLineWidth(2)
		ctx.SetColor(colornames.Pink)
		ctx.DrawCircle(500, 400, 50)
		ctx.Stroke()
		ctx.SetColor(colornames.Orange)
		ctx.DrawCircle(200, 200, 100)
		ctx.Stroke()
		ctx.DrawCircle(100, 100, 100)
		ctx.Stroke()

		ctx.Pop()

		img := ctx.Image()
		_ = img.At(0, 0)
	}
}
