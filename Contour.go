package giame

import (
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
	"image"
	"math"
)

type Contour struct {
	Scale      Scale
	PipeOption PipeOption
	Bound      image.Rectangle
	//
	Filler Filler
	Points []mgl32.Vec2
}

func pointsStringer(pts []mgl32.Vec2) string {
	const printmax = 12
	res := "["
	for i, pt := range pts {
		res += fmt.Sprintf("(%.4f, %.4f)", pt[0], pt[1])
		res += ", "
		if i+1 > printmax {
			res += "..., "
		}
	}
	res = res[:len(res)-2]
	res += fmt.Sprintf("/ length : %d]", len(pts))
	return res
}
func (s Contour) String() string {
	return fmt.Sprintf(
		"Contour(Filler : %v, Points : %s, Scale : %v, Bound : %v, PipeOption : %v)",
		s.Filler,
		pointsStringer(s.Points),
		s.Scale,
		s.Bound,
		s.PipeOption,
	)
}
func (s Contour) WorkspaceSize() (w, h int) {
	sclsqrt := s.Scale.Sqrt()
	return int(float32(s.Bound.Dx()) * sclsqrt), int(float32(s.Bound.Dy()) * sclsqrt)
}

type Scale int32

const (
	Scale1x  Scale = 1
	Scale2x  Scale = 2
	Scale4x  Scale = 4
	Scale8x  Scale = 8
	Scale16x Scale = 16
)

func (s Scale) String() string {
	switch s {
	case Scale1x:
		return "1x"
	case Scale2x:
		return "2x"
	case Scale4x:
		return "4x"
	case Scale8x:
		return "8x"
	case Scale16x:
		return "16x"

	default:
		return fmt.Sprintf("x%d(unofficial)", int(s))
	}
}
func (s Scale) Sqrt() float32 {
	return float32(math.Sqrt(float64(s)))
}

type PipeOption int32

// If PipeOptionNearestNeighbor and Scale is 1x there is no Pipe work
const (
	PipeOptionKernel          PipeOption = 0x0000000F
	PipeOptionNearestNeighbor PipeOption = 0x00000000
	PipeOptionBilinear        PipeOption = 0x00000001
	PipeOptionBicubic         PipeOption = 0x00000002
	//

	PipeOptionFilter          PipeOption = 0x000000F0
	PipeOptionNone            PipeOption = 0x00000000
	PipeOptionLaplacian       PipeOption = 0x00000010
	PipeOptionLaplacianExtend PipeOption = 0x00000020
)

func (s PipeOption) String() string {
	result := "["
	switch s & 0x0000000F {
	case PipeOptionNearestNeighbor:
		result += "NearestNeighbor"
	case PipeOptionBilinear:
		result += "Bilinear"
	case PipeOptionBicubic:
		result += "Bicubic"
	}
	result += " / "
	switch s & 0x000000F0 {
	case PipeOptionNone:
		result += "None"
	case PipeOptionLaplacian:
		result += "Laplacian"
	case PipeOptionLaplacianExtend:
		result += "LaplacianExtend"
	}
	result += "]"
	return result
}
