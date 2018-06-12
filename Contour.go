package giame

import (
	"image"
	"github.com/go-gl/mathgl/mgl32"
)

// - *ScanlineHorizontal
//
type Contour interface {
	GetOption() Option
	GetBound() image.Rectangle
	GetFiller() Filler
}
type ContourBuilder interface {
	Write(l0, l1 mgl32.Vec2)
	Build() Contour
}
type BuilderAllocator func(c contourInfo) ContourBuilder