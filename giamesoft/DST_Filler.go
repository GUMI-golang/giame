package giamesoft

import (
	"image/draw"
	"github.com/GUMI-golang/giame"
	"image/color"
	"image"
)

func RasterFiller(dst draw.Image, ws *Workspace, min image.Point, f giame.Filler)  {
	f.ToBound(image.Rect(0,0,ws.Width, ws.Height))
	for x := 0; x < ws.Width; x++ {
		for y := 0; y < ws.Height; y++ {
			v := ws.Get(x, y)
			r, g, b, a := f.At(x, y).RGBA()
			dst.Set(min.X + x, min.Y + y, color.RGBA{
				R: uint8(float32(r >> 8) * v + .5),
				G: uint8(float32(g >> 8) * v + .5),
				B: uint8(float32(b >> 8) * v + .5),
				A: uint8(float32(a >> 8) * v + .5),
			})
		}
	}
}
func RasterUniform(dst draw.Image, ws *Workspace, min image.Point, f *giame.UniformFiller)  {
	for x := 0; x < ws.Width; x++ {
		for y := 0; y < ws.Height; y++ {
			v := ws.Get(x, y)
			dst.Set(min.X + x, min.Y + y, color.RGBA{
				R: uint8(float32(f.Color.R) * v + .5),
				G: uint8(float32(f.Color.G) * v + .5),
				B: uint8(float32(f.Color.B) * v + .5),
				A: uint8(float32(f.Color.A) * v + .5),
			})
		}
	}
}
func RasterFixed(dst draw.Image, ws *Workspace, min image.Point, f *giame.FixedFiller)  {
	for x := 0; x < ws.Width; x++ {
		for y := 0; y < ws.Height; y++ {
			v := ws.Get(x, y)
			r, g, b, a := f.At(x, y).RGBA()
			dst.Set(min.X + x, min.Y + y, color.RGBA{
				R: uint8(float32(r >> 8) * v + .5),
				G: uint8(float32(g >> 8) * v + .5),
				B: uint8(float32(b >> 8) * v + .5),
				A: uint8(float32(a >> 8) * v + .5),
			})
		}
	}
}
func RasterKernel(dst draw.Image, ws *Workspace, min image.Point, f *giame.KernelFiller)  {
	f.ToBound(image.Rect(0,0,ws.Width, ws.Height))
	for x := 0; x < ws.Width; x++ {
		for y := 0; y < ws.Height; y++ {
			v := ws.Get(x, y)
			r, g, b, a := f.At(x, y).RGBA()
			dst.Set(min.X + x, min.Y + y, color.RGBA{
				R: uint8(float32(r >> 8) * v + .5),
				G: uint8(float32(g >> 8) * v + .5),
				B: uint8(float32(b >> 8) * v + .5),
				A: uint8(float32(a >> 8) * v + .5),
			})
		}
	}
}
func RasterRepeat(dst draw.Image, ws *Workspace, min image.Point, f *giame.RepeatFiller)  {
	f.ToBound(image.Rect(0,0,ws.Width, ws.Height))
	for x := 0; x < ws.Width; x++ {
		for y := 0; y < ws.Height; y++ {
			v := ws.Get(x, y)
			r, g, b, a := f.At(x, y).RGBA()
			dst.Set(min.X + x, min.Y + y, color.RGBA{
				R: uint8(float32(r >> 8) * v + .5),
				G: uint8(float32(g >> 8) * v + .5),
				B: uint8(float32(b >> 8) * v + .5),
				A: uint8(float32(a >> 8) * v + .5),
			})
		}
	}
}

