package main

import (
	"github.com/GUMI-golang/giame/giamesoft"
	"github.com/GUMI-golang/giame"
	"image"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/GUMI-golang/gumi/gcore"
)

func main() {
	giamesoft.Driver.Init()
	defer giamesoft.Driver.Close()
	//
	f := giame.NewFiller(giame.FillerTypeBilinear, gcore.MustValue(gcore.Load("example/jellybeans.png")).(image.Image))
	//
	r := giamesoft.Driver.MakeResult(128, 128)
	q := giame.NewQuary(image.Rect(32, 0, 96, 64), giame.SetupScanlineHorizontal)
	q.SetFiller(f)
	c := q.Fill(func(query *giame.FillQuery) {
		query.MoveTo(mgl32.Vec2{0,0})
		query.LineTo(mgl32.Vec2{64,0})
		query.LineTo(mgl32.Vec2{64,64})
		query.LineTo(mgl32.Vec2{0,64})
		query.CloseTo()
	})
	r.Request(giamesoft.Driver, c)
	gcore.Capture("cpu-fill", r.Image())
}
