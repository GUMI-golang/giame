package svg

import (
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
	"reflect"
)

type PathData interface {
	fmt.Stringer
	empty_PathData()
}

// https://www.w3.org/TR/SVG2/paths.html#PathDataMovetoCommands
type (
	AbsoluteMoveTo struct {
		Arg []mgl32.Vec2
	}
	RelativeMoveTo struct {
		Arg []mgl32.Vec2
	}
)

func (AbsoluteMoveTo) empty_PathData() {
	panic("implement me")
}
func (s AbsoluteMoveTo) String() string {
	return fmt.Sprintf("%s : %v", reflect.TypeOf(s).Name(), s.Arg)
}
func (RelativeMoveTo) empty_PathData() {
	panic("implement me")
}
func (s RelativeMoveTo) String() string {
	return fmt.Sprintf("%s : %v", reflect.TypeOf(s).Name(), s.Arg)
}

// https://www.w3.org/TR/SVG2/paths.html#PathDataLinetoCommands
type (
	AbsoluteLineTo struct {
		Arg []mgl32.Vec2
	}
	RelativeLineTo struct {
		Arg []mgl32.Vec2
	}
	AbsoluteVerticalLineTo struct {
		Arg []float32
	}
	RelativeVerticalLineTo struct {
		Arg []float32
	}
	AbsoluteHorizontalLineTo struct {
		Arg []float32
	}
	RelativeHorizontalLineTo struct {
		Arg []float32
	}
)

func (AbsoluteLineTo) empty_PathData() {
	panic("implement me")
}
func (s AbsoluteLineTo) String() string {
	return fmt.Sprintf("%s : %v", reflect.TypeOf(s).Name(), s.Arg)
}
func (RelativeLineTo) empty_PathData() {
	panic("implement me")
}
func (s RelativeLineTo) String() string {
	return fmt.Sprintf("%s : %v", reflect.TypeOf(s).Name(), s.Arg)
}
func (AbsoluteVerticalLineTo) empty_PathData() {
	panic("implement me")
}
func (s AbsoluteVerticalLineTo) String() string {
	return fmt.Sprintf("%s : %v", reflect.TypeOf(s).Name(), s.Arg)
}
func (RelativeVerticalLineTo) empty_PathData() {
	panic("implement me")
}
func (s RelativeVerticalLineTo) String() string {
	return fmt.Sprintf("%s : %v", reflect.TypeOf(s).Name(), s.Arg)
}
func (AbsoluteHorizontalLineTo) empty_PathData() {
	panic("implement me")
}
func (s AbsoluteHorizontalLineTo) String() string {
	return fmt.Sprintf("%s : %v", reflect.TypeOf(s).Name(), s.Arg)
}
func (RelativeHorizontalLineTo) empty_PathData() {
	panic("implement me")
}
func (s RelativeHorizontalLineTo) String() string {
	return fmt.Sprintf("%s : %v", reflect.TypeOf(s).Name(), s.Arg)
}

// https://www.w3.org/TR/SVG2/paths.html#PathDataClosePathCommand
type (
	CloseTo struct{}
)

func (CloseTo) empty_PathData() {
	panic("implement me")
}
func (s CloseTo) String() string {
	return fmt.Sprintf("%s", reflect.TypeOf(s).Name())
}

// https://www.w3.org/TR/SVG2/paths.html#PathDataQuadraticBezierCommands
type (
	AbsoluteQuadTo struct {
		Arg [][2]mgl32.Vec2
	}
	RelativeQuadTo struct {
		Arg [][2]mgl32.Vec2
	}
	AbsoluteSmoothQuadTo struct {
		Arg []mgl32.Vec2
	}
	RelativeSmoothQuadTo struct {
		Arg []mgl32.Vec2
	}
)

func (AbsoluteQuadTo) empty_PathData() {
	panic("implement me")
}
func (s AbsoluteQuadTo) String() string {
	return fmt.Sprintf("%s : %v", reflect.TypeOf(s).Name(), s.Arg)
}
func (RelativeQuadTo) empty_PathData() {
	panic("implement me")
}
func (s RelativeQuadTo) String() string {
	return fmt.Sprintf("%s : %v", reflect.TypeOf(s).Name(), s.Arg)
}
func (AbsoluteSmoothQuadTo) empty_PathData() {
	panic("implement me")
}
func (s AbsoluteSmoothQuadTo) String() string {
	return fmt.Sprintf("%s : %v", reflect.TypeOf(s).Name(), s.Arg)
}
func (RelativeSmoothQuadTo) empty_PathData() {
	panic("implement me")
}
func (s RelativeSmoothQuadTo) String() string {
	return fmt.Sprintf("%s : %v", reflect.TypeOf(s).Name(), s.Arg)
}

// https://www.w3.org/TR/SVG2/paths.html#PathDataCubicBezierCommands
type (
	AbsoluteCubeTo struct {
		Arg [][3]mgl32.Vec2
	}
	RelativeCubeTo struct {
		Arg [][3]mgl32.Vec2
	}
	AbsoluteSmoothCubeTo struct {
		Arg [][2]mgl32.Vec2
	}
	RelativeSmoothCubeTo struct {
		Arg [][2]mgl32.Vec2
	}
)

func (AbsoluteCubeTo) empty_PathData() {
	panic("implement me")
}
func (s AbsoluteCubeTo) String() string {
	return fmt.Sprintf("%s : %v", reflect.TypeOf(s).Name(), s.Arg)
}
func (RelativeCubeTo) empty_PathData() {
	panic("implement me")
}
func (s RelativeCubeTo) String() string {
	return fmt.Sprintf("%s : %v", reflect.TypeOf(s).Name(), s.Arg)
}
func (AbsoluteSmoothCubeTo) empty_PathData() {
	panic("implement me")
}
func (s AbsoluteSmoothCubeTo) String() string {
	return fmt.Sprintf("%s : %v", reflect.TypeOf(s).Name(), s.Arg)
}
func (RelativeSmoothCubeTo) empty_PathData() {
	panic("implement me")
}
func (s RelativeSmoothCubeTo) String() string {
	return fmt.Sprintf("%s : %v", reflect.TypeOf(s).Name(), s.Arg)
}

// TODO
// https://www.w3.org/TR/SVG2/paths.html#PathDataEllipticalArcCommands
type (
	AbsoluteArc struct {
		Arg []ArcArguments
	}
	RelativeArc struct {
		Arg []ArcArguments
	}
	ArcArguments struct {
		Radius mgl32.Vec2
		Rotation float32
		LargeArc bool
		Sweep bool
		To mgl32.Vec2
	}
)

func (s ArcArguments) String() string {
	flags := ""
	if s.LargeArc{
		flags += "LargeArc, "
	}
	if s.Sweep{
		flags += "Sweep, "
	}
	return fmt.Sprintf("{Radius : %v, Rotation : %f, %sTo : %v}", s.Radius, s.Rotation, flags, s.To)
}

func (AbsoluteArc) empty_PathData() {
	panic("implement me")
}
func (s AbsoluteArc) String() string {
	return fmt.Sprintf("%s : %v", reflect.TypeOf(s).Name(), s.Arg)
}
func (RelativeArc) empty_PathData() {
	panic("implement me")
}
func (s RelativeArc) String() string {
	return fmt.Sprintf("%s : %v", reflect.TypeOf(s).Name(), s.Arg)
}

// https://www.w3.org/TR/SVG2/paths.html#PathDataBearingCommands
type (
	AbsoluteBearing struct {
		Arg []float32
	}
	RelativeBearing struct {
		Arg []float32
	}
)

func (AbsoluteBearing) empty_PathData() {
	panic("implement me")
}
func (s AbsoluteBearing) String() string {
	return fmt.Sprintf("%s : %v", reflect.TypeOf(s).Name(), s.Arg)
}
func (RelativeBearing) empty_PathData() {
	panic("implement me")
}
func (s RelativeBearing) String() string {
	return fmt.Sprintf("%s : %v", reflect.TypeOf(s).Name(), s.Arg)
}
