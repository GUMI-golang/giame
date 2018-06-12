package giame

import (
	"container/list"
	"github.com/go-gl/mathgl/mgl32"
	"math"
)


var (
	rotm90 = mgl32.Rotate2D(-math.Pi / 2)
)

type strokewriter struct {
	res *list.List
	//
	rad  float32
	join StrokeJoin
	cap  StrokeCap
	//
	p0, p1, p2 mgl32.Vec2
	l01, l12   mgl32.Vec2
	vv01, vv12 mgl32.Vec2
	i          int
}

func newstrokewriter(rad float32, join StrokeJoin, cap StrokeCap) *strokewriter {
	return &strokewriter{
		res: list.New(),
		//
		rad:  rad,
		cap:  cap,
		join: join,
		//
		i: 0,
	}
}
func (s *strokewriter) Write(points ...mgl32.Vec2) {
	for _, point := range points {
		s.write(point)
	}
}
func (s *strokewriter) Pad() {
	s.res.PushBack(nil)
}
func (s *strokewriter) write(points mgl32.Vec2) {
	if s.i >= 2 {
		s.p2 = points
		s.l12 = s.p2.Sub(s.p1)
		s.vv12 = rotm90.Mul2x1(s.l12).Normalize().Mul(s.rad)
		//
		cos := float64(s.l01.Mul(-1).Dot(s.l12) / s.l01.Len() / s.l12.Len())
		dir := s.vv01.Add(s.vv12).Normalize().Mul(s.rad / float32(math.Sqrt((1-cos)/2)))
		if s.l01.Dot(dir) < 0 {
			// b0 insideb
			switch s.join {
			case StrokeJoinBevel:
				s.res.PushBack([]mgl32.Vec2{
					s.p1,
					s.p1.Sub(s.vv12),
					s.p1.Sub(s.vv01),
					s.p1,
					Spacer(),
				})
			case StrokeJoinMiter:
				s.res.PushBack([]mgl32.Vec2{
					s.p1,
					s.p1.Sub(s.vv01),
					s.p1.Sub(dir),
					s.p2.Sub(s.vv12),
					s.p1,
					Spacer(),
				})
			case StrokeJoinRound:
				var temp []mgl32.Vec2
				temp = append(temp, s.p1)
				temp = append(temp, quadFromTo(s.p1.Sub(s.vv01), s.p1.Sub(dir), s.p1.Sub(s.vv12))...)
				temp = append(temp, s.p1, Spacer())
			}

		} else {
			switch s.join {
			case StrokeJoinBevel:
				s.res.PushBack([]mgl32.Vec2{
					s.p1,
					s.p1.Add(s.vv01),
					s.p1.Add(s.vv12),
					s.p1,
					Spacer(),
				})
			case StrokeJoinMiter:
				s.res.PushBack([]mgl32.Vec2{
					s.p1,
					s.p1.Add(s.vv01),
					s.p1.Add(dir),
					s.p2.Add(s.vv12),
					s.p1,
					Spacer(),
				})
			case StrokeJoinRound:
				var temp []mgl32.Vec2
				temp = append(temp, s.p1)
				temp = append(temp, quadFromTo(s.p1.Add(s.vv01), s.p1.Add(dir), s.p1.Add(s.vv12))...)
				temp = append(temp, s.p1, Spacer())
			}
		}
		//

		s.res.PushBack([]mgl32.Vec2{
			s.p1.Add(s.vv12),
			s.p2.Add(s.vv12),
			s.p2.Sub(s.vv12),
			s.p1.Sub(s.vv12),
			s.p1.Add(s.vv12),
			Spacer(),
		})
		//
		s.p0 = s.p1
		s.p1 = s.p2
		s.l01 = s.l12
		s.vv01 = s.vv12
	} else if s.i == 1 {
		// s.i == 1
		// first point, second point draw function
		s.p1 = points
		s.l01 = s.p1.Sub(s.p0)
		s.vv01 = rotm90.Mul2x1(s.l01).Normalize().Mul(s.rad)
		//
		switch s.cap {
		case StrokeCapButt:
		case StrokeCapSqaure:
			r01n := s.l01.Normalize().Mul(s.rad)
			s.res.PushBack([]mgl32.Vec2{
				s.p0.Sub(r01n).Add(s.vv01),
				s.p0.Add(s.vv01),
				s.p0.Sub(s.vv01),
				s.p0.Sub(r01n).Sub(s.vv01),
				s.p0.Sub(r01n).Add(s.vv01),
				Spacer(),
			})
		case StrokeCapRound:
			r01n := s.l01.Normalize().Mul(s.rad)
			var temp []mgl32.Vec2
			temp = append(temp, quadFromTo(s.p0.Sub(s.vv01), s.p0.Sub(r01n).Sub(s.vv01), s.p0.Sub(r01n))...)
			temp = append(temp, quadFromTo(s.p0.Sub(r01n), s.p0.Sub(r01n).Add(s.vv01), s.p0.Add(s.vv01))...)
			temp = append(temp, s.p0.Sub(s.vv01), Spacer())
			s.res.PushBack(temp)
		}
		//
		s.res.PushBack([]mgl32.Vec2{
			s.p0.Add(s.vv01),
			s.p1.Add(s.vv01),
			s.p1.Sub(s.vv01),
			s.p0.Sub(s.vv01),
			s.p0.Add(s.vv01),
			Spacer(),
		})

	} else {
		// s.i == 0
		s.p0 = points
	}

	s.i++
}
func (s *strokewriter) Result() (res []mgl32.Vec2) {
	//switch s.cap {
	//case StrokeCapButt:
	//case StrokeCapSqaure:
	//	r12n := s.l12.Normalize().Mul(s.rad)
	//	s.res.PushBack([]mgl32.Vec2{
	//		s.p2.Add(r12n).Add(s.vv12),
	//		s.p2.Add(s.vv12),
	//		s.p2.Sub(s.vv12),
	//		s.p2.Add(r12n).Sub(s.vv12),
	//		s.p2.Add(r12n).Add(s.vv12),
	//		Spacer(),
	//	})
	//case StrokeCapRound:
	//	r12n := s.l12.Normalize().Mul(s.rad)
	//	var temp []mgl32.Vec2
	//	fmt.Println(
	//		s.p2.Add(s.vv12),
	//		s.p2.Add(r12n).Add(s.vv12),
	//		s.p2.Add(r12n),
	//		s.p2.Add(r12n),
	//		s.p2.Add(r12n).Sub(s.vv12),
	//		s.p2.Sub(s.vv12),
	//	)
	//	temp = Append(temp, quadFromTo(
	//		s.p2.Add(s.vv12),
	//		s.p2.Add(r12n).Add(s.vv12),
	//		s.p2.Add(r12n),
	//	)...)
	//	temp = Append(temp, quadFromTo(
	//		s.p2.Add(r12n),
	//		s.p2.Add(r12n).Sub(s.vv12),
	//		s.p2.Sub(s.vv12),
	//	)...)
	//	temp = Append(temp, temp[0], Spacer())
	//	s.res.PushBack(temp)
	//}

	for p := s.res.Front(); p != nil; p = p.Next() {
		if v, ok := p.Value.([]mgl32.Vec2); ok {
			res = append(res, v...)
		}
	}
	return
}
