package main

import (
	"testing"

	"github.com/fogleman/gg"
	"golang.org/x/image/colornames"
)

func BenchmarkGame_Draw(b *testing.B) {
	ctx := gg.NewContext(w, h)
	for i := 0; i < b.N; i++ {
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

func BenchmarkGame_DrawEvery(b *testing.B) {
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

func BenchmarkGame_DrawHorizontal(b *testing.B) {
	ctx := gg.NewContext(w, h)
	for i := 0; i < b.N; i++ {
		ctx.SetLineWidth(2)
		ctx.SetColor(colornames.Pink)
		ctx.DrawLine(0, 0, 100, 0)
		ctx.Stroke()
	}
}

func BenchmarkGame_DrawVertical(b *testing.B) {
	ctx := gg.NewContext(w, h)
	for i := 0; i < b.N; i++ {
		ctx.SetLineWidth(2)
		ctx.SetColor(colornames.Pink)
		ctx.DrawLine(0, 0, 0, 100)
		ctx.Stroke()
	}
}

func BenchmarkGame_DrawRect(b *testing.B) {
	ctx := gg.NewContext(w, h)
	for i := 0; i < b.N; i++ {
		ctx.SetLineWidth(2)
		ctx.SetColor(colornames.Pink)
		ctx.DrawRectangle(0, 0, 100, 100)
		ctx.Stroke()
	}
}
