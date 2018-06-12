package svg

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/GUMI-golang/gumi/gcore"
)

type Vectable interface {
	MoveTo(to mgl32.Vec2)
	LineTo(to mgl32.Vec2)
	QuadTo(p0, to mgl32.Vec2)
	CubeTo(p0, p1, to mgl32.Vec2)
	CloseTo()
}

type PathDataHandler struct {
	bear *mgl32.Mat2
	last *mgl32.Vec2
	// for smooth quadto
	qcontrol0 *mgl32.Vec2
	keepq     bool
	// for smooth cubeto
	ccontrol0 *mgl32.Vec2
	keepc     bool
}

func NewPathDataHandler() *PathDataHandler {
	temp := mgl32.Rotate2D(0)
	return &PathDataHandler{
		last: new(mgl32.Vec2),
		bear: &temp,

	}
}
func (s *PathDataHandler) do() {
	if s.keepq{
		s.keepq = false
	}else {
		s.qcontrol0 = nil
	}
	if s.keepc{
		s.keepc = false
	}else {
		s.ccontrol0 = nil
	}
}

func (s *PathDataHandler)Handle(query Vectable, datas ... PathData) {
	for _, data := range datas {
		s.handle(query, data)
	}
}

func (s *PathDataHandler)handle(query Vectable, data PathData) {
	s.do()
	switch d := data.(type) {
	case AbsoluteMoveTo:
		for _, v := range d.Arg {
			query.MoveTo(v)
			*s.last = v
		}
	case RelativeMoveTo:
		for _, v := range d.Arg {
			p := s.modify(v)
			query.MoveTo(p)
			*s.last = p
		}
	case AbsoluteLineTo:
		for _, v := range d.Arg {
			query.LineTo(v)
			*s.last = v
		}
	case RelativeLineTo:
		for _, v := range d.Arg{
			p := s.modify(v)
			query.LineTo(p)
			*s.last = p
		}
	case AbsoluteHorizontalLineTo:
		for _, v := range d.Arg {
			temp := mgl32.Vec2{v, s.last.Y()}
			query.LineTo(temp)
			*s.last = temp
		}
	case RelativeHorizontalLineTo:
		for _, v := range d.Arg {
			temp := mgl32.Vec2{s.last.X() + v, s.last.Y()}
			query.LineTo(temp)
			*s.last = temp
		}
	case AbsoluteVerticalLineTo:
		for _, v := range d.Arg {
			temp := mgl32.Vec2{s.last.X(), v}
			query.LineTo(temp)
			*s.last = temp
		}
	case RelativeVerticalLineTo:
		for _, v := range d.Arg {
			temp := mgl32.Vec2{s.last.X(), s.last.Y() + v}
			query.LineTo(temp)
			*s.last = temp
		}
	case AbsoluteQuadTo:
		var a0, to mgl32.Vec2
		for _, vs := range d.Arg {
			a0 = vs[0]
			to = vs[1]
			query.QuadTo(a0, to)
		}
		*s.last = to
		*s.qcontrol0 = a0
		s.keepq = true
	case RelativeQuadTo:
		var a0, to mgl32.Vec2
		for _, vs := range d.Arg {
			a0 = s.modify(vs[0])
			to = s.modify(vs[1])
			query.QuadTo(a0, to)
		}
		*s.last = to
		*s.qcontrol0 = a0
		s.keepq = true
	case AbsoluteSmoothQuadTo:
		var p0, to mgl32.Vec2
		if s.qcontrol0 ==nil{
			p0 = *s.last
		}else {
			p0 = mirrorByPoint(*s.qcontrol0, *s.last)
		}
		to = d.Arg[0]
		//
		query.QuadTo(p0, to)
		*s.last = to
		*s.qcontrol0 = p0
		//
		for i := 1; i < len(d.Arg); i += 1 {
			var p0, to mgl32.Vec2
			p0 = mirrorByPoint(*s.qcontrol0, d.Arg[i - 1])
			to = d.Arg[i]
			//
			query.QuadTo(p0, to)
			*s.last = to
			*s.qcontrol0 = p0
		}
		s.keepq = true
	case RelativeSmoothQuadTo:
		// first case
		var p0, to mgl32.Vec2
		if s.qcontrol0 ==nil{
			p0 = *s.last
		}else {
			p0 = mirrorByPoint(*s.qcontrol0, *s.last)
		}
		prev := s.last.Add(s.bear.Mul2x1(d.Arg[0]))
		to = prev
		//
		query.QuadTo(p0, to)
		*s.last = to
		*s.qcontrol0 = p0
		//
		for i := 1; i < len(d.Arg); i += 1 {
			var p0, to mgl32.Vec2
			p0 = mirrorByPoint(*s.qcontrol0, prev)
			to = s.last.Add(s.bear.Mul2x1(d.Arg[i]))
			//
			query.QuadTo(p0, to)
			*s.last = to
			*s.qcontrol0 = p0
			prev = to
		}
		s.keepq = true
	case AbsoluteCubeTo:
		var p0, p1, to mgl32.Vec2
		for _, v := range d.Arg {
			p0, p1, to = v[0], v[1], v[2]
			query.CubeTo(p0, p1, to)
		}
		*s.ccontrol0 = p1
		*s.last = to
		s.keepc = true
	case RelativeCubeTo:
		var p0, p1, to mgl32.Vec2
		for _, v := range d.Arg {
			p0 = s.modify(v[0])
			p1 = s.modify(v[1])
			to = s.modify(v[2])
			query.CubeTo(p0, p1, to)

			*s.ccontrol0 = p1
			*s.last = to
			s.keepc = true
		}
	case AbsoluteSmoothCubeTo:
		var p0, p1, to mgl32.Vec2
		if s.ccontrol0 ==nil{
			p0 = *s.last
		}else {
			p0 = mirrorByPoint(*s.ccontrol0, *s.last)
		}
		p1 = d.Arg[0][0]
		to = d.Arg[0][1]
		//
		query.CubeTo(p0, p1, to)
		*s.last = to
		*s.ccontrol0 = p1
		//
		for _, v := range d.Arg {
			p0 = mirrorByPoint(*s.ccontrol0, p1)
			p1 = v[0]
			to = v[1]
			query.CubeTo(p0, p1, to)
			*s.last = to
			*s.ccontrol0 = p1
		}
		s.keepc = true
	case RelativeSmoothCubeTo:
		var p0, p1, to mgl32.Vec2
		if s.ccontrol0 ==nil{
			p0 = *s.last
		}else {
			p0 = mirrorByPoint(*s.ccontrol0, *s.last)
		}
		p1 = s.modify(d.Arg[0][0])
		to = s.modify(d.Arg[0][1])
		//
		query.CubeTo(p0, p1, to)
		*s.last = to
		*s.ccontrol0 = p1
		//
		for _, v := range d.Arg {
			p0 = mirrorByPoint(*s.ccontrol0, p1)
			p1 = s.modify(v[0])
			to = s.modify(v[1])
			query.CubeTo(p0, p1, to)
			*s.last = to
			*s.ccontrol0 = p1
		}
		s.keepc = true
		//
	case AbsoluteArc:
		// TODO
	case RelativeArc:
		// TODO
	case AbsoluteBearing:
		*s.bear = mgl32.Rotate2D(float32(gcore.ToRadian(float64(d.Arg[len(d.Arg) - 1]))))
	case RelativeBearing:
		*s.bear = s.bear.Mul2(mgl32.Rotate2D(float32(gcore.ToRadian(float64(d.Arg[len(d.Arg) - 1])))))
	case CloseTo:
		query.CloseTo()
	}
}

func (s *PathDataHandler) modify(vec2 mgl32.Vec2) mgl32.Vec2 {
	return s.bear.Mul2x1(s.last.Add(vec2))
}

func mirrorByPoint(a, mirror mgl32.Vec2) mgl32.Vec2 {
	return mirror.Add(mirror.Sub(a))
}