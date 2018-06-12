package giamesoft

import (
	"github.com/GUMI-golang/giame"
	"image"
	"image/draw"
	"image/color"
)

var Driver = _Driver{}
type _Driver struct {}

func (s _Driver) Init() (err error) {
	return nil
}
func (s _Driver) Close() error {
	return nil
}
func (s _Driver) MakeResult(w, h int) giame.Result {
	return &Result{
		img: image.NewRGBA(image.Rect(0,0,w,h)),
	}
}

type Result struct {
	img *image.RGBA
}

func (s *Result) Size() (w, h int) {
	return s.img.Rect.Dx(), s.img.Rect.Dy()
}

func (s *Result) Image() image.Image {
	return s.img
}

func (s *Result) Clear() {
	draw.Draw(s.img, s.img.Rect, image.NewUniform(color.Transparent), image.ZP, draw.Src)
}

func (s *Result) Request(dr giame.Driver, rq giame.Contour) {
	switch v := rq.(type) {
	case *giame.ScanlineHorizontal:
		ws := NewWorkspace(v.Bound.Dx(), v.Bound.Dy())
		ScanlineHorizontal(ws, v)
		RasterFiller(s.img, ws, v.Bound.Min, v.Filler)
	}
}
