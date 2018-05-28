package example

import (
	"github.com/GUMI-golang/giame"
	"image"
	"github.com/GUMI-golang/gumi/gcore"
)

func Font8x(bound image.Rectangle) *giame.Contour {
	cq := giame.NewContourQuary(bound)
	cq.SetScale(giame.Scale8x)
	return giame.DefaultVFont.TextInRect(cq, "Hello?", bound, gcore.AlignCenter)
}
