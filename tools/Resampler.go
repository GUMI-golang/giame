package tools

import (
	"github.com/GUMI-golang/giame/tools/kernel"
	"github.com/go-gl/mathgl/mgl64"
	"image"
	"image/color"
	"image/draw"
	"math"
)

type Resampler struct {
	k kernel.Kernel
}

func NewResampler(k kernel.Kernel) *Resampler {
	return &Resampler{
		k: k,
	}
}
func Resampling(k kernel.Kernel, dst draw.Image, r image.Rectangle, src image.Image, sp image.Point) {
	NewResampler(k).Draw(dst, r, src, sp)
}
func (s *Resampler) Draw(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point) {
	if v, ok := allRGBA(src, dst); ok {
		// both RGBA
		resampleRGBA(v[0], v[1], s.k, r, sp)
	} else {
		resampleImage(src, dst, s.k, r, sp)
	}
}

func resampleImage(src image.Image, dst draw.Image, k kernel.Kernel, r image.Rectangle, sp image.Point) {
	srcbd := src.Bounds()
	dstbd := dst.Bounds()
	dstbd = dstbd.Intersect(r)
	srcbd = srcbd.Sub(sp)
	//
	delta := mgl64.Vec2{
		float64(srcbd.Dx()) / float64(dstbd.Dx()),
		float64(srcbd.Dy()) / float64(dstbd.Dy()),
	}

	a := float64(k.Rad())
	for x := dstbd.Min.X; x < dstbd.Max.X; x++ {
		for y := dstbd.Min.Y; y < dstbd.Max.Y; y++ {
			sx := delta[0]*float64(x) - .5
			sy := delta[1]*float64(y) - .5
			//
			var temp [5]float64
			var hori = [2]int{int(sx - a), int(sx + a + .5)}
			var vert = [2]int{int(sy - a), int(sy + a + .5)}
			for tx := hori[0]; tx <= hori[1]; tx++ {
				for ty := vert[0]; ty <= vert[1]; ty++ {
					rx := Iclamp(tx, srcbd.Min.X, srcbd.Max.X-1)
					ry := Iclamp(ty, srcbd.Min.Y, srcbd.Max.Y-1)
					r, g, b, a := src.At(rx, ry).RGBA()
					kr := k.Do(float64(tx)-sx) * k.Do(float64(ty)-sy)
					temp[0] += float64(r) / math.MaxUint16 * kr
					temp[1] += float64(g) / math.MaxUint16 * kr
					temp[2] += float64(b) / math.MaxUint16 * kr
					temp[3] += float64(a) / math.MaxUint16 * kr
					temp[4] += kr
				}
			}
			alpha := Iclamp(int(temp[3]/temp[4]*math.MaxUint8), 0, math.MaxUint8)
			green := Iclamp(int(temp[2]/temp[4]*math.MaxUint8), 0, alpha)
			blue := Iclamp(int(temp[2]/temp[4]*math.MaxUint8), 0, alpha)
			red := Iclamp(int(temp[2]/temp[4]*math.MaxUint8), 0, alpha)
			dst.Set(x, y, color.RGBA{
				R: uint8(red),
				G: uint8(green),
				B: uint8(blue),
				A: uint8(alpha),
			})
		}
	}
}
func resampleRGBA(src *image.RGBA, dst *image.RGBA, k kernel.Kernel, r image.Rectangle, sp image.Point) {
	srcbd := src.Bounds()
	dstbd := dst.Bounds()
	dstbd = dstbd.Intersect(r)
	srcbd = srcbd.Sub(sp)
	//
	delta := mgl64.Vec2{
		float64(srcbd.Dx()) / float64(dstbd.Dx()),
		float64(srcbd.Dy()) / float64(dstbd.Dy()),
	}

	a := float64(k.Rad())
	for x := dstbd.Min.X; x < dstbd.Max.X; x++ {
		for y := dstbd.Min.Y; y < dstbd.Max.Y; y++ {
			sx := delta[0]*float64(x) - .5
			sy := delta[1]*float64(y) - .5
			//
			var temp [5]float64
			var hori = [2]int{int(sx - a), int(sx + a + .5)}
			var vert = [2]int{int(sy - a), int(sy + a + .5)}
			for tx := hori[0]; tx <= hori[1]; tx++ {
				for ty := vert[0]; ty <= vert[1]; ty++ {
					rx := Iclamp(tx, srcbd.Min.X, srcbd.Max.X-1)
					ry := Iclamp(ty, srcbd.Min.Y, srcbd.Max.Y-1)
					off := src.PixOffset(rx, ry)
					kr := k.Do(float64(tx)-sx) * k.Do(float64(ty)-sy)
					temp[0] += float64(src.Pix[off+0]) / math.MaxUint8 * kr
					temp[1] += float64(src.Pix[off+1]) / math.MaxUint8 * kr
					temp[2] += float64(src.Pix[off+2]) / math.MaxUint8 * kr
					temp[3] += float64(src.Pix[off+3]) / math.MaxUint8 * kr
					temp[4] += kr
				}
			}
			off := dst.PixOffset(x, y)
			dst.Pix[off+3] = uint8(Iclamp(int(temp[3]/temp[4]*math.MaxUint8), 0, math.MaxUint8))
			alpha := int(dst.Pix[off+3])
			dst.Pix[off+2] = uint8(Iclamp(int(temp[2]/temp[4]*math.MaxUint8), 0, alpha))
			dst.Pix[off+1] = uint8(Iclamp(int(temp[1]/temp[4]*math.MaxUint8), 0, alpha))
			dst.Pix[off+0] = uint8(Iclamp(int(temp[0]/temp[4]*math.MaxUint8), 0, alpha))
		}
	}
}
