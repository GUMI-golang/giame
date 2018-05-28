package tools

import (
	"github.com/GUMI-golang/giame/tools/mask"
	"image/draw"
	"image"
	"image/color"
	"math"
)

type Filter struct {
	s mask.Mask3
}

func NewFilter(s mask.Mask3) *Filter {
	return &Filter{
		s:s,
	}
}
func Filting(s mask.Mask3, dst draw.Image) {
	NewFilter(s).Filt(dst)
}

func (s Filter) Filt(dst draw.Image) {
	if v, ok := allRGBA(dst); ok {
		// both RGBA
		filterRGBA(v[0], s.s)
	} else {
		filterImage(dst, s.s)
	}
}
func filterImage(dst draw.Image, s mask.Mask3){
	bd := dst.Bounds()
	w, h := s.Size()
	hw, hh := w/2, h/2
	for x := bd.Min.X; x < bd.Max.X; x++ {
		for y := bd.Min.Y; y < bd.Max.Y; y++ {
			var r, g, b, a float64
			//
			for i := 0; i < w; i++{
				for j := 0; j < h; j++{
					fromx := Iclamp(x + i - hw, bd.Min.X, bd.Max.X - 1)
					fromy := Iclamp(y + j - hh, bd.Min.Y, bd.Max.Y - 1)
					tempr,tempg,tempb,tempa := dst.At(fromx, fromy).RGBA()
					val := s.Data[j][i]
					r += float64(tempr) * val
					g += float64(tempg) * val
					b += float64(tempb) * val
					a += float64(tempa) * val
				}
			}
			//
			dst.Set(x, y, color.RGBA64{
				R:uint16(Iclamp(int(r), 0, math.MaxUint16)),
				G:uint16(Iclamp(int(g), 0, math.MaxUint16)),
				B:uint16(Iclamp(int(b), 0, math.MaxUint16)),
				A:uint16(Iclamp(int(a), 0, math.MaxUint16)),
			})
		}
	}
}
func filterRGBA(dst *image.RGBA, s mask.Mask3){
	bd := dst.Bounds()
	w, h := s.Size()
	hw, hh := w/2, h/2
	for x := bd.Min.X; x < bd.Max.X; x++ {
		for y := bd.Min.Y; y < bd.Max.Y; y++ {
			var r, g, b, a float64
			//
			for i := 0; i < w; i++{
				for j := 0; j < h; j++{
					fromx := Iclamp(x + i - hw, bd.Min.X, bd.Max.X - 1)
					fromy := Iclamp(y + j - hh, bd.Min.Y, bd.Max.Y - 1)
					//
					off := dst.PixOffset(fromx, fromy)
					val := s.Data[j][i]
					r += float64(dst.Pix[off + 0]) * val
					g += float64(dst.Pix[off + 1]) * val
					b += float64(dst.Pix[off + 2]) * val
					a += float64(dst.Pix[off + 3]) * val
				}
			}
			//
			off := dst.PixOffset(x, y)
			dst.Pix[off + 0] = uint8(Iclamp(int(r), 0, math.MaxUint8))
			dst.Pix[off + 1] = uint8(Iclamp(int(g), 0, math.MaxUint8))
			dst.Pix[off + 2] = uint8(Iclamp(int(b), 0, math.MaxUint8))
			dst.Pix[off + 3] = uint8(Iclamp(int(a), 0, math.MaxUint8))
		}
	}
}