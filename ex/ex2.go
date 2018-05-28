package main

import (
	"github.com/GUMI-golang/giame/tools/kernel"
	"github.com/GUMI-golang/gumi/gcore"
	"image"
	"image/draw"
	"image/png"
	"os"
	"reflect"
	"github.com/GUMI-golang/giame/tools"
)

func main() {
	//f, err := os.Open("./example/4Pix.png")
	//f, err := os.Open("./example/binarysmall.png")
	f, err := os.Open("./example/jellybeans.png")
	//f, err := os.Open("./example/text.png")
	//f, err := os.Open("./example/largetext.png")
	//f, err := os.Open("./example/binary.png")
	//f, err := os.Open("./example/Pix16.png")
	//f, err := os.Open("./example/_out_fill.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	img, err := png.Decode(f)
	if err != nil {
		panic(err)
	}
	sz := img.Bounds().Size()
	//

	var scale = 4.
	src := image.NewRGBA(image.Rect(0, 0, sz.X, sz.Y))
	draw.Draw(src, src.Bounds(), img, img.Bounds().Min, draw.Src)
	quad := image.NewRGBA(image.Rect(0, 0, int(float64(sz.X)*scale), int(float64(sz.Y)*scale)))
	for _, k := range []kernel.Kernel{
		kernel.NearestNeighbor,
		kernel.Bilinear,
		kernel.Bell,
		kernel.Hermite,
		kernel.BicubicHalf,
		kernel.MitchellOneThird,
		kernel.Lanczos2,
		kernel.Lanczos3,
	} {
		tools.Resampling(k, quad, quad.Rect, src, src.Rect.Min)
		gcore.Capture("_out_"+reflect.ValueOf(k).Type().Name(), quad)
	}
}
//func Resample(src image.Image, dst draw.Image, k kernel.Kernel) {
//	cvsrc, ok0 := src.(*image.RGBA)
//	cvdst, ok1 := dst.(*image.RGBA)
//	if ok0 && ok1 {
//		resampleRGBA(cvsrc, cvdst, k)
//		return
//	}
//	resampleImage(src, dst, k)
//}
//func resampleImage(src image.Image, dst draw.Image, k kernel.Kernel) {
//	srcbd := src.Bounds()
//	dstbd := dst.Bounds()
//	srcsz := srcbd.Size()
//	dstsz := dstbd.Size()
//	//
//	delta := mgl64.Vec2{
//		float64(srcsz.X) / float64(dstsz.X),
//		float64(srcsz.Y) / float64(dstsz.Y),
//	}
//
//	a := float64(k.Rad())
//	for x := 0; x < dstsz.X; x++ {
//		for y := 0; y < dstsz.Y; y++ {
//			sx := delta[0] * float64(x) - .5
//			sy := delta[1] * float64(y) - .5
//			//
//			var temp [5]float64
//			var hori = [2]int{int(sx - a), int(sx + a + .5)}
//			var vert = [2]int{int(sy - a), int(sy + a + .5)}
//			for tx := hori[0]; tx <= hori[1]; tx ++ {
//				for ty := vert[0]; ty <= vert[1]; ty ++ {
//					rx := iclamp(tx, 0, srcsz.X - 1)
//					ry := iclamp(ty, 0, srcsz.Y - 1)
//					r, g, b, a := src.At(rx, ry).RGBA()
//					kr := k.Do(float64(tx) - sx) * k.Do(float64(ty) - sy)
//					temp[0] += float64(r) / math.MaxUint16 * kr
//					temp[1] += float64(g) / math.MaxUint16 * kr
//					temp[2] += float64(b) / math.MaxUint16 * kr
//					temp[3] += float64(a) / math.MaxUint16 * kr
//					temp[4] += kr
//				}
//			}
//			alpha := iclamp(int(temp[3] / temp[4] * math.MaxUint8), 0, math.MaxUint8)
//			green := iclamp(int(temp[2] / temp[4] * math.MaxUint8), 0, alpha)
//			blue := iclamp(int(temp[2] / temp[4] * math.MaxUint8), 0, alpha)
//			red := iclamp(int(temp[2] / temp[4] * math.MaxUint8), 0, alpha)
//			dst.Set(x, y, color.RGBA{
//				R:uint8(red),
//				G:uint8(green),
//				B:uint8(blue),
//				A:uint8(alpha),
//			})
//		}
//	}
//}
//func resampleRGBA(src *image.RGBA, dst *image.RGBA, k kernel.Kernel) {
//	srcbd := src.Bounds()
//	dstbd := dst.Bounds()
//	srcsz := srcbd.Size()
//	dstsz := dstbd.Size()
//	//
//	delta := mgl64.Vec2{
//		float64(srcsz.X) / float64(dstsz.X),
//		float64(srcsz.Y) / float64(dstsz.Y),
//	}
//
//	a := float64(k.Rad())
//	for x := 0; x < dstsz.X; x++ {
//		for y := 0; y < dstsz.Y; y++ {
//			sx := delta[0] * float64(x) - .5
//			sy := delta[1] * float64(y) - .5
//			//
//			var temp [5]float64
//			var hori = [2]int{int(sx - a), int(sx + a + .5)}
//			var vert = [2]int{int(sy - a), int(sy + a + .5)}
//			for tx := hori[0]; tx <= hori[1]; tx ++ {
//				for ty := vert[0]; ty <= vert[1]; ty ++ {
//					rx := iclamp(tx, 0, srcsz.X - 1)
//					ry := iclamp(ty, 0, srcsz.Y - 1)
//					off := src.PixOffset(rx, ry)
//					kr := k.Do(float64(tx) - sx) * k.Do(float64(ty) - sy)
//					temp[0] += float64(src.Pix[off + 0]) / math.MaxUint8 * kr
//					temp[1] += float64(src.Pix[off + 1]) / math.MaxUint8 * kr
//					temp[2] += float64(src.Pix[off + 2]) / math.MaxUint8 * kr
//					temp[3] += float64(src.Pix[off + 3]) / math.MaxUint8 * kr
//					temp[4] += kr
//				}
//			}
//			off := dst.PixOffset(x,y)
//			dst.Pix[off + 3] = uint8(iclamp(int(temp[3] / temp[4] * math.MaxUint8), 0, math.MaxUint8))
//			alpha := int(dst.Pix[off + 3])
//			dst.Pix[off + 2] = uint8(iclamp(int(temp[2] / temp[4] * math.MaxUint8), 0, alpha))
//			dst.Pix[off + 1] = uint8(iclamp(int(temp[1] / temp[4] * math.MaxUint8), 0, alpha))
//			dst.Pix[off + 0] = uint8(iclamp(int(temp[0] / temp[4] * math.MaxUint8), 0, alpha))
//		}
//	}
//}
