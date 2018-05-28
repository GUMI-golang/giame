package giame

import (
	"image"
	"image/color"
	"image/draw"
	"fmt"
	"strings"
	"github.com/pkg/errors"
	"github.com/GUMI-golang/giame/tools/kernel"
)

type Filler interface {
	image.Image
	fmt.Stringer
	FillerType()FillerType
	ToBound(rectangle image.Rectangle)
	//
	Serial() *RawFiller
	Deserialize(filler *RawFiller) error
}

func NewFiller(tp FillerType, data image.Image) Filler {
	switch tp {
	case FillerTypeUniform:
		return NewUniformFiller(data.At(0,0))
	case FillerTypeFixed:
		return NewFixedFiller(data)
	case FillerTypeNearestNeighbor:
		return NewKernelFiller(tp, data)
	case FillerTypeBilinear:
		return NewKernelFiller(tp, data)
	case FillerTypeBell:
		return NewKernelFiller(tp, data)
	case FillerTypeHermite:
		return NewKernelFiller(tp, data)
	case FillerTypeBicubicHalf:
		return NewKernelFiller(tp, data)
	case FillerTypeMitchellOneThird:
		return NewKernelFiller(tp, data)
	case FillerTypeLanczos2:
		return NewKernelFiller(tp, data)
	case FillerTypeLanczos3:
		return NewKernelFiller(tp, data)
	case FillerTypeRepeat:
		return NewRepeatFiller(data)
	case FillerTypeRepeatHorizontal:
		return NewRepeatHorizontalFiller(data)
	case FillerTypeRepeatVertical:
		return NewRepeatVerticalFiller(data)
	default:
		panic("Unknwon FillerType")
	}
	return nil
}
func NewUniformFiller(c color.Color) Filler {
	return &UniformFiller{
		Color: color.RGBAModel.Convert(c).(color.RGBA),
	}
}
func NewFixedFiller(img image.Image) Filler {
	res := new(FixedFiller)
	//
	bd := img.Bounds()
	bd = bd.Sub(bd.Min)
	temp := image.NewRGBA(bd)
	draw.Draw(temp, temp.Rect, img, img.Bounds().Min, draw.Src)
	//
	res.Width = temp.Rect.Dx()
	res.Height = temp.Rect.Dy()
	res.Data = temp.Pix
	return res
}
//

func NewKernelFiller(t FillerType, img image.Image) Filler {
	res := new(KernelFiller)
	//
	switch t {
	case FillerTypeNearestNeighbor:
		res.Kernel = kernel.NearestNeighbor
	case FillerTypeBilinear:
		res.Kernel = kernel.Bilinear
	case FillerTypeBell:
		res.Kernel = kernel.Bell
	case FillerTypeHermite:
		res.Kernel = kernel.Hermite
	case FillerTypeBicubicHalf:
		res.Kernel = kernel.BicubicHalf
	case FillerTypeMitchellOneThird:
		res.Kernel = kernel.MitchellOneThird
	case FillerTypeLanczos2:
		res.Kernel = kernel.Lanczos2
	case FillerTypeLanczos3:
		res.Kernel = kernel.Lanczos3
	default:
		panic("Unknwon FillerType")
	}
	bd := img.Bounds()
	bd = bd.Sub(bd.Min)
	temp := image.NewRGBA(bd)
	draw.Draw(temp, temp.Rect, img, img.Bounds().Min, draw.Src)
	//
	res.Width = temp.Rect.Dx()
	res.Height = temp.Rect.Dy()
	res.Data = temp.Pix
	return res
}
//
func NewRepeatFiller(img image.Image) Filler {
	res := new(RepeatFiller)
	//
	bd := img.Bounds()
	bd = bd.Sub(bd.Min)
	temp := image.NewRGBA(bd)
	draw.Draw(temp, temp.Rect, img, img.Bounds().Min, draw.Src)
	//
	res.Type = FillerTypeRepeat
	res.Width = temp.Rect.Dx()
	res.Height = temp.Rect.Dy()
	res.Data = temp.Pix
	return res
}
func NewRepeatHorizontalFiller(img image.Image) Filler {
	res := new(RepeatFiller)
	//
	bd := img.Bounds()
	bd = bd.Sub(bd.Min)
	temp := image.NewRGBA(bd)
	draw.Draw(temp, temp.Rect, img, img.Bounds().Min, draw.Src)
	//
	res.Type = FillerTypeRepeatHorizontal
	res.Width = temp.Rect.Dx()
	res.Height = temp.Rect.Dy()
	res.Data = temp.Pix
	return res
}
func NewRepeatVerticalFiller(img image.Image) Filler {
	res := new(RepeatFiller)
	//
	bd := img.Bounds()
	bd = bd.Sub(bd.Min)
	temp := image.NewRGBA(bd)
	draw.Draw(temp, temp.Rect, img, img.Bounds().Min, draw.Src)
	//
	res.Type = FillerTypeRepeatVertical
	res.Width = temp.Rect.Dx()
	res.Height = temp.Rect.Dy()
	res.Data = temp.Pix
	return res
}

type FillerType uint8

const (
	FillerTypeUniform FillerType = iota
	//
	FillerTypeFixed FillerType = iota
	//
	FillerTypeNearestNeighbor FillerType = iota
	FillerTypeBilinear FillerType = iota
	FillerTypeBell FillerType = iota
	FillerTypeHermite FillerType = iota
	FillerTypeBicubicHalf FillerType = iota
	FillerTypeMitchellOneThird FillerType = iota
	FillerTypeLanczos2 FillerType = iota
	FillerTypeLanczos3 FillerType = iota
	//
	FillerTypeRepeat FillerType = iota
	FillerTypeRepeatHorizontal FillerType = iota
	FillerTypeRepeatVertical FillerType = iota
	// Special define
	FillerTypeLength FillerType = iota
	FillerTypeInvalid FillerType = iota
)

func ParseFillertype(s string) FillerType {
	switch strings.ToLower(s) {
	case "uniform":
		return FillerTypeUniform
	case "fixed":
		return FillerTypeFixed
	//case "gaussian":
	//	return FillerTypeGaussian
	//case "nearest":
	//	return FillerTypeNearest
	//case "nearest-neighbor":
	//	return FillerTypeNearestNeighbor
	case "repeat":
		return FillerTypeRepeat
	case "repeat-horizontal":
		return FillerTypeRepeatHorizontal
	case "repeat-vertical":
		return FillerTypeRepeatVertical
	}
	return FillerTypeInvalid
}
func (s FillerType ) Strings() string {
	switch s {
	case FillerTypeUniform:
		return "uniform"
	case FillerTypeFixed:
		return "fixed"
	//case FillerTypeGaussian:
	//	return "gaussian"
	//case FillerTypeNearest:
	//	return "nearest"
	case FillerTypeNearestNeighbor:
		return "nearest-neighbor"
	case FillerTypeRepeat:
		return "repeat"
	case FillerTypeRepeatHorizontal:
		return "repeat-horizontal"
	case FillerTypeRepeatVertical:
		return "repeat-vertical"
	}
	return "unknown"
}

var UnmatchingFillerType = errors.New("Can't use this type to here")