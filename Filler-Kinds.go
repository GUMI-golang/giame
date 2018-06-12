package giame

import (
	"fmt"
	"github.com/GUMI-golang/gumi/gcore"
	"image"
	"image/color"
	"github.com/GUMI-golang/giame/tools/kernel"
	"github.com/go-gl/mathgl/mgl64"
	"math"
	"github.com/GUMI-golang/giame/tools"
)

const (
	R         = 0
	G         = 1
	B         = 2
	A         = 3
	PIXLENGTH = 4
)

type (
	UniformFiller struct {
		Color color.RGBA
	}
	FixedFiller struct {
		Data          []uint8
		Width, Height int
	}
	RepeatFiller struct {
		Data          []uint8
		Width, Height int
		RepeaterWidth, RepeaterHeight int
	}
	KernelFiller struct {
		Data          []uint8
		Width, Height int
		To image.Rectangle
		Kernel kernel.Kernel
	}
)

func offset(x, y, w int) int {
	return x*PIXLENGTH + w*y*PIXLENGTH
}

func (s *UniformFiller) FillerType() FillerType {
	return FillerTypeUniform
}
func (s *UniformFiller) ColorModel() color.Model {
	return color.RGBAModel
}
func (s *UniformFiller) Bounds() image.Rectangle {
	return image.Rectangle{Min: image.Point{X: -1e9, Y: -1e9}, Max: image.Point{X: 1e9, Y: 1e9}}
}
func (s *UniformFiller) At(x, y int) color.Color {
	return s.Color
}
func (s *UniformFiller) ToBound(rectangle image.Rectangle) {
}
func (s *UniformFiller) String() string {
	return fmt.Sprintf("UniformFiller(%s)", gcore.MarshalColor(s.Color))
}
func (s *UniformFiller) Serial() *RawFiller {
	return &RawFiller{
		Type: s.FillerType(),
		Data: []uint8{s.Color.R, s.Color.G, s.Color.B, s.Color.A},
	}
}
func (s *UniformFiller) Deserialize(filler *RawFiller) error {
	if s.FillerType() != filler.Type {
		return UnmatchingFillerType
	}
	s.Color.R = filler.Data[0]
	s.Color.G = filler.Data[1]
	s.Color.B = filler.Data[2]
	s.Color.A = filler.Data[3]
	return nil
}

func (s *FixedFiller) FillerType() FillerType {
	return FillerTypeFixed
}
func (s *FixedFiller) ColorModel() color.Model {
	return color.RGBAModel
}
func (s *FixedFiller) Bounds() image.Rectangle {
	return image.Rect(0, 0, s.Width, s.Height)
}
func (s *FixedFiller) At(x, y int) color.Color {
	o := offset(x, y, s.Width)
	if o < 0 || o >= len(s.Data) {
		return color.Transparent
	}
	return color.RGBA{
		R: s.Data[o+R],
		G: s.Data[o+G],
		B: s.Data[o+B],
		A: s.Data[o+A],
	}
}
func (s *FixedFiller) ToBound(rectangle image.Rectangle) {
}
func (s *FixedFiller) String() string {
	return fmt.Sprintf("Fixed(size : [%d/%d], memory : %d)", s.Width, s.Height, len(s.Data))
}
func (s *FixedFiller) Serial() *RawFiller {
	temp := new(RawFiller)
	temp.Type = s.FillerType()
	temp.Width = s.Width
	temp.Height = s.Height
	temp.Data = make([]uint8, len(s.Data))
	copy(temp.Data, s.Data)
	return temp
}
func (s *FixedFiller) Deserialize(filler *RawFiller) error {
	if s.FillerType() != filler.Type {
		return UnmatchingFillerType
	}
	s.Width = filler.Width
	s.Height = filler.Height
	s.Data = make([]uint8, len(filler.Data))
	copy(s.Data, filler.Data)
	return nil
}

func (s *RepeatFiller) ColorModel() color.Model {
	return color.RGBAModel
}
func (s *RepeatFiller) Bounds() image.Rectangle {
	return image.Rect(0, 0, s.Width, s.Height)
}
func (s *RepeatFiller) At(x, y int) color.Color {
	var o = offset(x%s.RepeaterWidth, y%s.RepeaterHeight, s.Width)

	return color.RGBA{
		R: s.Data[o+R],
		G: s.Data[o+G],
		B: s.Data[o+B],
		A: s.Data[o+A],
	}
}
func (s *RepeatFiller) FillerType() FillerType {
	return FillerTypeRepeat
}
func (s *RepeatFiller) ToBound(rectangle image.Rectangle) {
}
func (s *RepeatFiller) String() string {
	return fmt.Sprintf("Repeat(size : [%d/%d], memory : %d)", s.Width, s.Height, len(s.Data))
}
func (s *RepeatFiller) Serial() *RawFiller {
	temp := new(RawFiller)
	temp.Type = s.FillerType()
	temp.Width = s.Width
	temp.Height = s.Height
	temp.Data = make([]uint8, len(s.Data))
	copy(temp.Data, s.Data)
	return temp
}
func (s *RepeatFiller) Deserialize(filler *RawFiller) error {
	if s.FillerType() != filler.Type {
		return UnmatchingFillerType
	}
	s.Width = filler.Width
	s.Height = filler.Height
	s.Data = make([]uint8, len(filler.Data))
	copy(s.Data, filler.Data)
	return nil
}


func (s *KernelFiller) ColorModel() color.Model {
	return color.RGBAModel
}
func (s *KernelFiller) Bounds() image.Rectangle {
	return image.Rect(0, 0, s.Width, s.Height)
}
func (s *KernelFiller) At(x, y int) color.Color {
	srcbd := s.Bounds()
	dstbd := s.To
	//
	delta := mgl64.Vec2{
		float64(srcbd.Dx()) / float64(dstbd.Dx()),
		float64(srcbd.Dy()) / float64(dstbd.Dy()),
	}
	a := float64(s.Kernel.Rad())
	sx := delta[0]*float64(x) - .5
	sy := delta[1]*float64(y) - .5
	//
	var temp [5]float64
	var hori = [2]int{int(sx - a), int(sx + a + .5)}
	var vert = [2]int{int(sy - a), int(sy + a + .5)}
	for tx := hori[0]; tx <= hori[1]; tx++ {
		for ty := vert[0]; ty <= vert[1]; ty++ {
			rx := tools.Iclamp(tx, srcbd.Min.X, srcbd.Max.X-1)
			ry := tools.Iclamp(ty, srcbd.Min.Y, srcbd.Max.Y-1)
			off := offset(rx, ry, s.Width)
			kr := s.Kernel.Do(float64(tx)-sx) * s.Kernel.Do(float64(ty)-sy)
			temp[0] += float64(s.Data[off+0]) / math.MaxUint8 * kr
			temp[1] += float64(s.Data[off+1]) / math.MaxUint8 * kr
			temp[2] += float64(s.Data[off+2]) / math.MaxUint8 * kr
			temp[3] += float64(s.Data[off+3]) / math.MaxUint8 * kr
			temp[4] += kr
		}
	}
	alpha := tools.Iclamp(int(temp[3]/temp[4]*math.MaxUint8), 0, math.MaxUint8)
	return color.RGBA{
		R: uint8(tools.Iclamp(int(temp[0]/temp[4]*math.MaxUint8), 0, alpha)),
		G: uint8(tools.Iclamp(int(temp[1]/temp[4]*math.MaxUint8), 0, alpha)),
		B: uint8(tools.Iclamp(int(temp[2]/temp[4]*math.MaxUint8), 0, alpha)),
		A: uint8(alpha),
	}
}
func (s *KernelFiller) FillerType() FillerType {
	switch s.Kernel {
	case kernel.NearestNeighbor:
		return FillerTypeNearestNeighbor
	case kernel.Bilinear:
		return FillerTypeBilinear
	case kernel.Bell:
		return FillerTypeBell
	case kernel.Hermite:
		return FillerTypeHermite
	case kernel.BicubicHalf:
		return FillerTypeBicubicHalf
	case kernel.MitchellOneThird:
		return FillerTypeMitchellOneThird
	case kernel.Lanczos2:
		return FillerTypeLanczos2
	case kernel.Lanczos3:
		return FillerTypeLanczos3
	}
	panic("UndefinedType")
}
func (s *KernelFiller) ToBound(rectangle image.Rectangle) {
	s.To = rectangle
}
func (s *KernelFiller) String() string {
	return fmt.Sprintf("Kernel(kernel %v, size : [%d/%d], memory : %d)", s.Kernel,s.Width, s.Height, len(s.Data))
}
func (s *KernelFiller) Serial() *RawFiller {
	temp := new(RawFiller)
	temp.Type = s.FillerType()
	temp.Width = s.Width
	temp.Height = s.Height
	temp.Data = make([]uint8, len(s.Data))
	copy(temp.Data, s.Data)
	return temp
}
func (s *KernelFiller) Deserialize(filler *RawFiller) error {
	switch filler.Type {
	case FillerTypeNearestNeighbor:
		s.Kernel = kernel.NearestNeighbor
	case FillerTypeBilinear:
		s.Kernel = kernel.Bilinear
	case FillerTypeBell:
		s.Kernel = kernel.Bell
	case FillerTypeHermite:
		s.Kernel = kernel.Hermite
	case FillerTypeBicubicHalf:
		s.Kernel = kernel.BicubicHalf
	case FillerTypeMitchellOneThird:
		s.Kernel = kernel.MitchellOneThird
	case FillerTypeLanczos2:
		s.Kernel = kernel.Lanczos2
	case FillerTypeLanczos3:
		s.Kernel = kernel.Lanczos3
	default:
		return UnmatchingFillerType
	}
	s.Width = filler.Width
	s.Height = filler.Height
	s.Data = make([]uint8, len(filler.Data))
	copy(s.Data, filler.Data)
	return nil
}