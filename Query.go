package giame

import (
	"github.com/go-gl/mathgl/mgl32"
	"image"
	"image/color"
)

type (
	Quary struct {
		builderAllocFn BuilderAllocator
		//
		c contourInfo
		// for stroke
		join  StrokeJoin
		cap   StrokeCap
		width float32
	}

	InnerQuery interface {
		MoveTo(to mgl32.Vec2)
		LineTo(to mgl32.Vec2)
		QuadTo(p0, to mgl32.Vec2)
		CubeTo(p0, p1, to mgl32.Vec2)
		CloseTo()
	}
	StrokeQuery struct {
		q *Quary
		b ContourBuilder
		p mgl32.Vec2
		s *strokewriter
	}

)

func NewQuary(bd image.Rectangle, builder BuilderAllocator) *Quary {
	return &Quary{
		//points: nil,
		builderAllocFn:builder,
		//
		c:contourInfo{
			Option:Option{},
			Filler:NewUniformFiller(color.Black),
			Bound:bd,
		},
		//
		join:  StrokeJoinBevel,
		cap:   StrokeCapButt,
		width: 1,
	}
}

//
func (s *Quary) SetOption(f Option) {
	s.c.Option = f
}
func (s *Quary) SetBound(rectangle image.Rectangle) {
	s.c.Bound = rectangle
}
func (s *Quary) SetFiller(f Filler) {
	s.c.Filler = f
}
func (s *Quary) SetStrokeCap(f StrokeCap) {
	s.cap = f
}
func (s *Quary) SetStrokeJoin(f StrokeJoin) {
	s.join = f
}
func (s *Quary) SetStrokeWidth(f float32) {
	s.width = f
}

//
func (s *Quary) GetOption() Option {
	return s.c.Option
}
func (s *Quary) GetBound() (rectangle image.Rectangle) {
	return s.c.Bound
}
func (s *Quary) GetFiller() (f Filler) {
	return s.c.Filler
}
func (s *Quary) GetStrokeCap() (f StrokeCap) {
	return s.cap
}
func (s *Quary) GetStrokeJoin() (f StrokeJoin) {
	return s.join
}
func (s *Quary) GetStrokeWidth() (f float32) {
	return s.width
}

func (s *Quary) Fill(fn func(query *FillQuery)) Contour {
	q := &FillQuery{
		q:s,
		p:Spacer(),
		b:s.builderAllocFn(s.c),
	}
	fn(q)
	return q.b.Build()
}

