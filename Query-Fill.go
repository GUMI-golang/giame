package giame

import (
	"github.com/go-gl/mathgl/mgl32"
	"io"
	"github.com/GUMI-golang/giame/tools/svg"
)

type FillQuery struct {
	q *Quary
	b ContourBuilder
	start mgl32.Vec2
	p mgl32.Vec2
}


// FillQuery
func (s *FillQuery) MoveTo(to mgl32.Vec2) {
	s.start = to
	s.p = to
}
func (s *FillQuery) LineTo(to mgl32.Vec2) {
	if IsSpacer(s.p){
		// TODO Err
		panic("Invalid")
	}
	s.b.Write(s.p, to)
	s.p = to
}
func (s *FillQuery) QuadTo(p0, to mgl32.Vec2) {
	if IsSpacer(s.p){
		// TODO Err
		panic("Invalid")
	}
	for _, to := range quadFromTo(s.p, p0, to) {
		s.b.Write(s.p, to)
		s.p = to
	}
}
func (s *FillQuery) CubeTo(p0, p1, to mgl32.Vec2) {
	if IsSpacer(s.p){
		// TODO Err
		panic("Invalid")
	}
	for _, to := range cubeFromTo(s.p, p0, p1, to) {
		s.b.Write(s.p, to)
		s.p = to
	}
}
func (s *FillQuery) CloseTo() {
	if !IsSpacer(s.start){
		s.LineTo(s.start)
		s.start = Spacer()
	}
	s.p = Spacer()
}
//
func (s *FillQuery) Query(q ContourQueryKind, reader io.Reader) error {
	switch q {
	case SVGPath:
		//return svgpath(s, reader)
		parser := svg.NewSVGPathParser(reader)
		handle := svg.NewPathDataHandler()
		for p, err := parser.Next(); true; p, err = parser.Next() {
			if err != nil{
				if err == io.EOF{
					break
				}else {
					return err
				}
			}
			handle.Handle(s, p)
		}
		return nil
	default:
		panic("Undefined Query Type")
	}

}
