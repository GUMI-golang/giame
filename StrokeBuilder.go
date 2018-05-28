package giame

import (
	"github.com/go-gl/mathgl/mgl32"
	"math"
	"container/list"
)

type StrokeBuilder struct {
	res90, res270 *list.List
	//
	rad float32
	join StrokeJoin
	cap StrokeCap
	//
	p0, p1, p2 mgl32.Vec2
	l01, l12 mgl32.Vec2
	vv01, vv12 mgl32.Vec2
	i int
}

func NewStrokeBuilder(rad float32, join StrokeJoin, cap StrokeCap) *StrokeBuilder {
	return &StrokeBuilder{
		res90: list.New(),
		res270: list.New(),
		//
		rad:rad,
		cap:cap,
		join:join,
		//
		i:0,
	}
}
func (s *StrokeBuilder ) Write(points ... mgl32.Vec2) {
	for _, point := range points {
		s.write(point)
	}
}

var (
	rotm90 = mgl32.Rotate2D(-math.Pi / 2)
)
func (s *StrokeBuilder ) write(points mgl32.Vec2) {
	s.i ++
	if s.i >= 2 {
		// s.i >= 2
		s.p2 = points
		s.l12 = s.p2.Sub(s.p1)
		s.vv12 = rotm90.Mul2x1(s.l12).Normalize().Mul(s.rad)
		//
		cos := float64(s.l01.Mul(-1).Dot(s.l12) / s.l01.Len() / s.l12.Len())
		dir := s.vv01.Add(s.vv12).Normalize().Mul(s.rad / float32(math.Sqrt((1 - cos)/2)))
		if s.l01.Dot(dir) < 0 {
			// b0 inside
			s.res90.PushBack(s.p1.Add(dir))
			switch s.join {
			case StrokeJoinBevel:
				s.res270.PushBack(s.p1.Sub(s.vv01))
				s.res270.PushBack(s.p1.Sub(s.vv12))
			case StrokeJoinMiter:
				s.res270.PushBack(s.p1.Sub(s.vv01))
				s.res270.PushBack(s.p1.Sub(dir))
				s.res270.PushBack(s.p1.Sub(s.vv12))
			case StrokeJoinRound:
				for _, p := range quadFromTo(s.p1.Sub(s.vv01), s.p1.Sub(dir), s.p1.Sub(s.vv12)) {
					s.res270.PushBack(p)
				}
			}

		} else {
			switch s.join {
			case StrokeJoinBevel:
				s.res90.PushBack(s.p1.Add(s.vv01))
				s.res90.PushBack(s.p1.Add(s.vv12))
			case StrokeJoinMiter:
				s.res90.PushBack(s.p1.Add(s.vv01))
				s.res90.PushBack(s.p1.Add(dir))
				s.res90.PushBack(s.p1.Add(s.vv12))
			case StrokeJoinRound:
				for _, p := range quadFromTo(s.p1.Add(s.vv01), s.p1.Add(dir), s.p1.Add(s.vv12)) {
					s.res270.PushBack(p)
				}
			}
			s.res270.PushBack(s.p1.Sub(dir))
		}
		//
		s.p0 = s.p1
		s.p1 = s.p2
		s.l01 = s.l12
		s.vv01 = s.vv12
	}else if s.i == 1{
		// s.i == 1
		// first point, second point draw function
		s.p1 = points
		s.l01 = s.p1.Sub(s.p0)
		s.vv01 = rotm90.Mul2x1(s.l01).Normalize().Mul(s.rad)
		switch s.cap {
		case StrokeCapButt:
		case StrokeCapSqaure:
			r01n := s.l01.Normalize().Mul(s.rad)
			s.res90.PushBack(s.p0.Sub(r01n).Sub(s.vv01))
			s.res90.PushBack(s.p0.Sub(r01n).Add(s.vv01))
		case StrokeCapRound:
			r01n := s.l01.Normalize().Mul(s.rad)
			for _, p := range quadFromTo(s.p0.Sub(s.vv01), s.p0.Sub(r01n).Sub(s.vv01), s.p0.Sub(r01n)) {
				s.res90.PushBack(p)
			}
			for _, p := range quadFromTo(s.p0.Sub(r01n), s.p0.Sub(r01n).Add(s.vv01), s.p0.Add(s.vv01)) {
				s.res90.PushBack(p)
			}
		}
		s.res90.PushBack(s.p0.Add(s.vv01))
		s.res270.PushBack(s.p0.Sub(s.vv01))
	}else {
		// s.i == 0
		s.p0 = points
	}
}
func (s *StrokeBuilder) Result() []mgl32.Vec2 {
	s.res90.PushBack(s.p1.Add(s.vv01))
	s.res270.PushBack(s.p1.Sub(s.vv01))
	// last cap
	switch s.cap {
	case StrokeCapButt:
	case StrokeCapSqaure:
		r01n := s.l01.Normalize().Mul(s.rad)
		s.res270.PushBack(s.p1.Add(r01n).Sub(s.vv01))
		s.res270.PushBack(s.p1.Add(r01n).Add(s.vv01))
	case StrokeCapRound:
		r01n := s.l01.Normalize().Mul(s.rad)
		for _, p := range quadFromTo(s.p1.Add(s.vv01), s.p1.Add(r01n).Add(s.vv01), s.p1.Add(r01n)) {
			s.res90.PushBack(p)
		}
		for _, p := range quadFromTo(s.p1.Add(r01n), s.p1.Add(r01n).Sub(s.vv01), s.p1.Sub(s.vv01)) {
			s.res90.PushBack(p)
		}
	}
	result := make([]mgl32.Vec2, s.res90.Len() + s.res270.Len())
	var i = 0
	for p := s.res90.Front(); p != nil; i, p = i + 1, p.Next() {
		result[i] = p.Value.(mgl32.Vec2)
	}
	for p := s.res270.Back(); p != nil; i, p = i + 1, p.Prev() {
		result[i] = p.Value.(mgl32.Vec2)
	}
	return result
}

func (s *StrokeBuilder) Reset() {
	s.res90.Init()
	s.res270.Init()
	s.i = 0
}
