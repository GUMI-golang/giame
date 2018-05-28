package example

import (
	"github.com/GUMI-golang/giame"
	"image"
	"github.com/go-gl/mathgl/mgl32"
)

func Triangle(bound image.Rectangle) *giame.Contour {
	cq := giame.NewContourQuary(bound)
	cq.MoveTo(mgl32.Vec2{float32(bound.Dx()) / 2, 0})
	cq.LineTo(mgl32.Vec2{0, float32(bound.Dy())})
	cq.LineTo(mgl32.Vec2{float32(bound.Dx()), float32(bound.Dy())})
	return cq.Fill()
}
