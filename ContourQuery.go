package giame

import (
	"container/list"
	"github.com/go-gl/mathgl/mgl32"
	"image"
	"image/color"
	"io"
)

type ContourQuary struct {
	points *list.List
	//
	filler Filler
	scale Scale
	pipeoption PipeOption
	bound  image.Rectangle
	// for stroke
	join  StrokeJoin
	cap   StrokeCap
	width float32
}

func NewContourQuary(bd image.Rectangle) *ContourQuary {
	return &ContourQuary{
		points: nil,
		filler: NewUniformFiller(color.Black),
		bound:  bd,
		scale:  Scale1x,
		//
		join:  StrokeJoinBevel,
		cap:   StrokeCapButt,
		width: 1,
	}
}

// common Path Raster funcs
// panic if not moveTo called
//
func (s *ContourQuary) MoveTo(to mgl32.Vec2) {
	s.RawMoveTo(to.Mul(float32(s.scale)))
}
func (s *ContourQuary) LineTo(to mgl32.Vec2) {
	s.RawLineTo(to.Mul(float32(s.scale)))
}
func (s *ContourQuary) QuadTo(p0, to mgl32.Vec2) {
	s.QuadTo(p0.Mul(float32(s.scale)), to.Mul(float32(s.scale)))
}
func (s *ContourQuary) CubeTo(p0, p1, to mgl32.Vec2) {
	s.RawCubeTo(p0.Mul(float32(s.scale)), p1.Mul(float32(s.scale)), to.Mul(float32(s.scale)))
}
// Raw* Path Raster funcs
// it pass value without scale
// care to use this
func (s *ContourQuary) RawMoveTo(to mgl32.Vec2) {
	if s.points == nil {
		s.points = list.New()
	}
	s.points.PushBack(to)
}
func (s *ContourQuary) RawLineTo(to mgl32.Vec2) {
	s.points.PushBack(to)
}
func (s *ContourQuary) RawQuadTo(p0, to mgl32.Vec2) {
	from := s.points.Back().Value.(mgl32.Vec2)
	for _, to := range quadFromTo(from, p0, to) {
		s.RawLineTo(to)
	}
}
func (s *ContourQuary) RawCubeTo(p0, p1, to mgl32.Vec2) {
	from := s.points.Back().Value.(mgl32.Vec2)
	for _, to := range cubeFromTo(from, p0, p1, to) {
		s.RawLineTo(to)
	}
}
// CloseTo can use both Raw*, common
func (s *ContourQuary) CloseTo() {
	var p *list.Element
	for p = s.points.Back(); p != nil; p = p.Prev() {
		if p.Value == nil {
			break
		}
	}
	if p == nil {
		p = s.points.Front()
	} else {
		p = p.Next()
	}
	s.points.PushBack(p.Value)
	s.points.PushBack(nil)
}
//
func (s *ContourQuary) SetScale(f Scale) {
	s.scale = f
}
func (s *ContourQuary) SetPipeOption(f PipeOption) {
	s.pipeoption = f
}
func (s *ContourQuary) SetBound(rectangle image.Rectangle) {
	s.bound = rectangle
}
func (s *ContourQuary) SetFiller(f Filler) {
	s.filler = f
}
func (s *ContourQuary) SetStrokeCap(f StrokeCap) {
	s.cap = f
}
func (s *ContourQuary) SetStrokeJoin(f StrokeJoin) {
	s.join = f
}
func (s *ContourQuary) SetStrokeWidth(f float32) {
	s.width = f
}
//
func (s *ContourQuary) GetScale() (f Scale) {
	return s.scale
}
func (s *ContourQuary) GetPipeOption() PipeOption {
	return s.pipeoption
}
func (s *ContourQuary) GetBound() (rectangle image.Rectangle) {
	return s.bound
}
func (s *ContourQuary) GetFiller()(f Filler) {
	return s.filler
}
func (s *ContourQuary) GetStrokeCap()(f StrokeCap) {
	return s.cap
}
func (s *ContourQuary) GetStrokeJoin()(f StrokeJoin) {
	return s.join
}
func (s *ContourQuary) GetStrokeWidth()(f float32) {
	return s.width
}
//
func (s *ContourQuary) Stroke() (res *Contour) {
	res = &Contour{
		Points: nil,
		Filler: s.filler,
		Bound:  s.bound,
		Scale:  s.scale,
	}
	builder := NewStrokeBuilder(s.width/2, s.join, s.cap)
	for p := s.points.Front(); p != nil; p = p.Next() {
		if v, ok := p.Value.(mgl32.Vec2); ok {
			builder.write(v)
		} else {
			res.Points = append(res.Points, builder.Result()...)
			res.Points = append(res.Points, spacer())
			builder.Reset()
		}
	}
	return res
}
func (s *ContourQuary) Fill() (res *Contour) {
	res = &Contour{
		Points: nil,
		Filler: s.filler,
		PipeOption:s.pipeoption,
		Bound:  s.bound,
		Scale:  s.scale,
	}
	for p := s.points.Front(); p != nil; p = p.Next() {
		if v, ok := p.Value.(mgl32.Vec2); ok {
			res.Points = append(res.Points, v)
		} else {
			res.Points = append(res.Points, spacer())
		}
	}
	return res
}
func (s *ContourQuary) Reset() {
	s.points = nil
	s.filler = NewUniformFiller(color.Black)
	//
	s.join = StrokeJoinBevel
	s.cap = StrokeCapButt
	s.width = 1
}
//
type ContourQueryKind uint8
const (
	SVGPath ContourQueryKind = iota
)
func (s *ContourQuary) Query(q ContourQueryKind, reader io.Reader) error {
	switch q {
	case SVGPath:
		return s.svgpath(reader)
	default:
		panic("Undefined Query Type")
	}

}