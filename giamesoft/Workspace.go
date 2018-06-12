package giamesoft

import (
	"github.com/GUMI-golang/giame/tools"
	"image"
	"github.com/go-gl/mathgl/mgl32"
	"math"
)

type Workspace struct {
	// DO NOT manipulate this values!
	// readonly recommanded
	Data          []float32
	Width, Height int
	limitW, limitH int
}

func NewWorkspace(w, h int) *Workspace {
	s := &Workspace{
		Data:   make([]float32, w * h),
		Width:  w,
		Height: h,

	}

	s.limitW, s.limitH = s.Width, s.Height
	return s
}
func (s Workspace) Get(x, y int) float32 {
	x = tools.Iclamp(x, 0, s.limitW - 1)
	y = tools.Iclamp(y, 0, s.limitH - 1)
	return s.Data[x + s.limitW * y]
}
func (s Workspace) Set(x, y int, v float32) {
	x = tools.Iclamp(x, 0, s.limitW - 1)
	y = tools.Iclamp(y, 0, s.limitH - 1)
	s.Data[x + s.limitW * y] = v
}
func (s Workspace) Add(x, y int, v float32) {
	x = tools.Iclamp(x, 0, s.limitW - 1)
	y = tools.Iclamp(y, 0, s.limitH - 1)
	s.Data[x + s.limitW * y] += v
}
func (s Workspace) Visualize() image.Image{
	temp := image.NewRGBA(image.Rect(0,0,s.limitW, s.limitH))
	for x := 0; x < s.limitW; x++ {
		for y := 0; y < s.limitH; y++ {
			off := temp.PixOffset(x, y)
			v := s.Get(x, y)
			if v > 0{
				temp.Pix[off + 0] = uint8(mgl32.Clamp(v* math.MaxUint8, 0, math.MaxUint8))
				if v > 1{
					temp.Pix[off + 2] = temp.Pix[off + 0]
				}
				temp.Pix[off + 3] = temp.Pix[off + 0]
			}else if v < 0{

				temp.Pix[off + 1] = uint8(mgl32.Clamp(-v* math.MaxUint8, 0, math.MaxUint8))
				if -v > 1{
					temp.Pix[off + 2] = temp.Pix[off + 1]
				}
				temp.Pix[off + 3] = temp.Pix[off + 1]
			}
		}
	}
	return temp
}