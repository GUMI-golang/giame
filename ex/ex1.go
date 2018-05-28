package main

import (
	"github.com/GUMI-golang/gumi/gcore"
	"github.com/go-gl/mathgl/mgl32"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"
	"reflect"
	"runtime"
)

func main() {
	//f, err := os.Open("./example/4Pix.png")
	f, err := os.Open("./example/binarysmall.png")
	//f, err := os.Open("./example/jellybeans.png")
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

	var x = 4
	c32f := NewRGBA32F(sz.X, sz.Y)
	draw.Draw(c32f, c32f.Bounds(), img, img.Bounds().Min, draw.Src)
	quad := NewRGBA32F(sz.X*x, sz.Y*x)
	for _, fn := range []kernelfn{
		fnNearestNeighbor,
		//fnBilinear,
		//fnHermite,
	} {
		Resample2D(c32f, quad, fn, 1)
		gcore.Capture("_out_"+runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name(), quad.ToRGBA())
	}
}

var RGBA32FColorModel = color.ModelFunc(func(i color.Color) color.Color {
	if c, ok := i.(RGBA32FColor); ok {
		return c
	}
	r, g, b, a := i.RGBA()
	return RGBA32FColor{
		float32(r) / math.MaxUint16,
		float32(g) / math.MaxUint16,
		float32(b) / math.MaxUint16,
		float32(a) / math.MaxUint16,
	}
})

type RGBA32FColor mgl32.Vec4

func (s RGBA32FColor) RGBA() (r, g, b, a uint32) {
	r = uint32(exiclamp(int(s[0]*math.MaxUint16), 0, 0xFFFF))
	g = uint32(exiclamp(int(s[1]*math.MaxUint16), 0, 0xFFFF))
	b = uint32(exiclamp(int(s[2]*math.MaxUint16), 0, 0xFFFF))
	a = uint32(exiclamp(int(s[3]*math.MaxUint16), 0, 0xFFFF))
	return
}
func (s RGBA32FColor) Normalize() RGBA32FColor {
	return RGBA32FColor{
		mgl32.Clamp(s[0], 0, 1),
		mgl32.Clamp(s[1], 0, 1),
		mgl32.Clamp(s[2], 0, 1),
		mgl32.Clamp(s[3], 0, 1),
	}
}

type RGBA32F struct {
	Pix           []RGBA32FColor
	Width, Height int
}

func NewRGBA32F(w, h int) *RGBA32F {
	return &RGBA32F{
		Pix:    make([]RGBA32FColor, w*h),
		Width:  w,
		Height: h,
	}
}
func (s *RGBA32F) ColorModel() color.Model {
	return RGBA32FColorModel
}
func (s *RGBA32F) Bounds() image.Rectangle {
	return image.Rect(0, 0, s.Width, s.Height)
}
func (s *RGBA32F) At(x, y int) color.Color {
	return s.Pix[x+y*s.Width]
}
func (s *RGBA32F) Set(x, y int, c color.Color) {
	s.Pix[x+y*s.Width] = RGBA32FColorModel.Convert(c).(RGBA32FColor)
}
func (s *RGBA32F) ToRGBA() *image.RGBA {
	rgba := image.NewRGBA(image.Rect(0,0,s.Width, s.Height))
	for y := 0; y < s.Height; y++{
		for x := 0; x < s.Width; x++{
			off := rgba.PixOffset(x, y)
			r,g,b,a := color.RGBAModel.Convert(s.Pix[x + y * s.Width]).RGBA()

			rgba.Pix[off + 0] = uint8(r)
			rgba.Pix[off + 1] = uint8(g)
			rgba.Pix[off + 2] = uint8(b)
			rgba.Pix[off + 3] = uint8(a)
		}
	}
	return rgba
}
func Resample2D(from, to *RGBA32F, k kernelfn, a int) {
	dh := float32(from.Width) / float32(to.Width)
	dv := float32(from.Height) / float32(to.Height)
	sh := float32(to.Width) / float32(from.Width)
	sv := float32(to.Height) / float32(from.Height)
	ah := eximax(int(dh), a)
	av := eximax(int(dv), a)
	for x := -int(sh / 2); x < to.Width-int(sh/2); x++ {
		for y := -int(sv / 2); y < to.Height-int(sv/2); y++ {
			var temp mgl32.Vec4
			var sum float32
			sx := float32(x) * dh
			sy := float32(y) * dv
			// Linear
			//for i := eximax(int(sx) - ah, 0); i <= eximin(int(sx) + ah, from.Width - 1); i++ {
			//	kv := k(sx - float32(i) - 0.00001)
			//	temp = temp.Add(mgl32.Vec4(from.Pix[i + int(sy + .5) * from.Width]).Mul(kv))
			//	sum+= kv
			//}
			//for i := eximax(int(sy) - av, 0); i <= eximin(int(sy) + av, from.Height - 1); i++ {
			//	kv := k(sy - float32(i) - 0.00001)
			//	temp = temp.Add(mgl32.Vec4(from.Pix[int(sx + .5) + i * from.Width]).Mul(kv))
			//	sum+= kv
			//}
			// Area
			for i := eximax(int(sx)-ah, 0); i <= eximin(int(sx + .5)+ah, from.Width-1); i++ {
				for j := eximax(int(sy)-av, 0); j <= eximin(int(sy + .5)+av, from.Height-1); j++ {
					dx := float64(sx - float32(i) - 0.00001)
					dy := float64(sy - float32(j) - 0.00001)
					kv := k(float32(math.Sqrt(dx * dx + dy * dy)))
					temp = temp.Add(mgl32.Vec4(from.Pix[i+j*from.Width]).Mul(kv))
					sum += kv
				}
			}
			if sum != 0{
				to.Pix[(x+int(sh/2))+(y+int(sv/2))*to.Width] = RGBA32FColor(temp.Mul(1 / sum)).Normalize()

			}
		}
	}

}

type kernelfn func(x float32) float32

const AlmostZero = .0000001

func fnNearestNeighbor(x float32) float32 {
	if -.51 < x && x < .51 {
		return 1
	}
	return 0
}
func fnBilinear(x float32) float32 {
	absx := float32(math.Abs(float64(x)))
	if absx <= 1 {
		return 1 - absx
	}
	return 0
}
func fnHermite(x float32) float32 {
	if x < 0 {
		x = -x
	}
	if x <= 1 {
		return 2*x*x*x - 3*x*x + 1
	}
	return 0
}
func fnLanczos3(x float32) float32 {
	const a = 3
	if x == 0 {
		return 1
	}
	x64 := float64(x)
	if -a <= x && x <= a {
		return float32(a*math.Sin(math.Pi*x64)*math.Sin((math.Pi*x64)/a)) / (math.Pi * math.Pi * x * x)
	}
	return 0
}
func fnLanczos2(x float32) float32 {
	const a = 2
	if x == 0 {
		return 1
	}
	x64 := float64(x)
	if -a <= x && x <= a {
		return float32(a*math.Sin(math.Pi*x64)*math.Sin((math.Pi*x64)/a)) / (math.Pi * math.Pi * x * x)
	}
	return 0
}

func exiclamp(i, min, max int) int {
	if i < min {
		return min
	}
	if i > max {
		return max
	}
	return i
}
func eximin(i, min int) int {
	if i < min {
		return i
	}
	return min
}
func eximax(i, max int) int {
	if i < max {
		return max
	}
	return i
}