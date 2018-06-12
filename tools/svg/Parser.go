package svg

import (
	"io"
	"unicode"
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
)

const bufsize = 1024

type SVGPathParser struct {
	rd       io.Reader
	buf      []byte
	from int
}
func NewSVGPathParser(Reader io.Reader) *SVGPathParser {
	return &SVGPathParser{
		rd:    Reader,
		buf:   nil,
		from:  0,
	}
}

func (s *SVGPathParser) Raw() ([]byte, int, error) {
	if len(s.buf) == 0{
		s.buf = make([]byte, bufsize)
		n, err := s.rd.Read(s.buf)
		s.buf = s.buf[:n]
		if err != nil {
			return nil, 0, err
		}
	}
	var (
		v byte
		to int
	)
	if !unicode.Is(unicode.Letter, rune(s.buf[0])){
		return nil, 0, &SVGError{
			What: "Expected Command, but no command",
			Where:[2]int{s.from, s.from + 1},
			Why: string(s.buf[0]),
		}
	}
	for to, v = range s.buf[1:] {
		if unicode.Is(unicode.Letter, rune(v)){
			break
		}
	}
	//
	res := s.buf[:to + 1]
	s.buf = s.buf[to + 1:]
	from := s.from
	s.from += len(res)
	return res, from, nil
}
func (s *SVGPathParser) Next() (PathData, error) {
	test, from, err := s.Raw()
	if err != nil {
		return nil, err
	}
	if !isCommand(test[0]){
		return nil, &SVGError{
			What:fmt.Sprintf("Command '%s' is not valid", string(test[0])),
			Where:[2]int{from, from + len(test)},
			Why: string(test),
		}
	}
	//
	switch test[0] {
	case 'M':
		vs, err := vectors(test[1:])
		if err != nil {
			return nil, &SVGError{
				What: fmt.Sprintf("Argument Invalid : %s", err.Error()),
				Where:[2]int{from, from + len(test)},
				Why: string(test),
			}
		}
		return AbsoluteMoveTo{
			Arg:vs,
		},  nil
	case 'm':
		vs, err := vectors(test[1:])
		if err != nil {
			return nil, &SVGError{
				What: fmt.Sprintf("Argument Invalid : %s", err.Error()),
				Where:[2]int{from, from + len(test)},
				Why: string(test),
			}
		}
		return RelativeMoveTo{
			Arg:vs,
		},  nil

	case 'L':
		vs, err := vectors(test[1:])
		if err != nil {
			return nil, &SVGError{
				What: fmt.Sprintf("Argument Invalid : %s", err.Error()),
				Where:[2]int{from, from + len(test)},
				Why: string(test),
			}
		}
		return AbsoluteLineTo{
			Arg:vs,
		},  nil
	case 'l':
		vs, err := vectors(test[1:])
		if err != nil {
			return nil, &SVGError{
				What: fmt.Sprintf("Argument Invalid : %s", err.Error()),
				Where:[2]int{from, from + len(test)},
				Why: string(test),
			}
		}
		return RelativeLineTo{
			Arg:vs,
		},  nil
	case 'H':
		vs, err := floats(test[1:])
		if err != nil {
			return nil, &SVGError{
				What: fmt.Sprintf("Argument Invalid : %s", err.Error()),
				Where:[2]int{from, from + len(test)},
				Why: string(test),
			}
		}
		return AbsoluteHorizontalLineTo{
			Arg:vs,
		},  nil
	case 'h':
		vs, err := floats(test[1:])
		if err != nil {
			return nil, &SVGError{
				What: fmt.Sprintf("Argument Invalid : %s", err.Error()),
				Where:[2]int{from, from + len(test)},
				Why: string(test),
			}
		}
		return RelativeHorizontalLineTo{
			Arg:vs,
		},  nil
	case 'V':
		vs, err := floats(test[1:])
		if err != nil {
			return nil, &SVGError{
				What: fmt.Sprintf("Argument Invalid : %s", err.Error()),
				Where:[2]int{from, from + len(test)},
				Why: string(test),
			}
		}
		return AbsoluteVerticalLineTo{
			Arg:vs,
		},  nil
	case 'v':
		vs, err := floats(test[1:])
		if err != nil {
			return nil, &SVGError{
				What: fmt.Sprintf("Argument Invalid : %s", err.Error()),
				Where:[2]int{from, from + len(test)},
				Why: string(test),
			}
		}
		return RelativeVerticalLineTo{
			Arg:vs,
		},  nil

	case 'Q':
		vs, err := vectors(test[1:])
		if err != nil {
			return nil, &SVGError{
				What: fmt.Sprintf("Argument Invalid : %s", err.Error()),
				Where:[2]int{from, from + len(test)},
				Why: string(test),
			}
		}
		if len(vs) % 2 != 0{
			return nil, &SVGError{
				What: fmt.Sprintf("Not enough vector, Q must have 2n vector"),
				Where:[2]int{from, from + len(test)},
				Why: string(test),
			}
		}
		temp := make([][2]mgl32.Vec2, len(vs) /2)
		for i, v := range vs {
			temp[i/2][i%2] = v
		}
		return AbsoluteQuadTo{
			Arg:temp,
		},  nil
	case 'q':
		vs, err := vectors(test[1:])
		if err != nil {
			return nil, &SVGError{
				What: fmt.Sprintf("Argument Invalid : %s", err.Error()),
				Where:[2]int{from, from + len(test)},
				Why: string(test),
			}
		}
		if len(vs) % 2 != 0{
			return nil, &SVGError{
				What: fmt.Sprintf("Not enough vector, q must have 2n vector"),
				Where:[2]int{from, from + len(test)},
				Why: string(test),
			}
		}
		temp := make([][2]mgl32.Vec2, len(vs) /2)
		for i, v := range vs {
			temp[i/2][i%2] = v
		}
		return RelativeQuadTo{
			Arg:temp,
		},  nil
	case 'T':
		vs, err := vectors(test[1:])
		if err != nil {
			return nil, &SVGError{
				What: fmt.Sprintf("Argument Invalid : %s", err.Error()),
				Where:[2]int{from, from + len(test)},
				Why: string(test),
			}
		}
		return AbsoluteSmoothQuadTo{
			Arg:vs,
		},  nil
	case 't':
		vs, err := vectors(test[1:])
		if err != nil {
			return nil, &SVGError{
				What: fmt.Sprintf("Argument Invalid : %s", err.Error()),
				Where:[2]int{from, from + len(test)},
				Why: string(test),
			}
		}
		return RelativeSmoothQuadTo{
			Arg:vs,
		},  nil

	case 'C':
		vs, err := vectors(test[1:])
		if err != nil {
			return nil, &SVGError{
				What: fmt.Sprintf("Argument Invalid : %s", err.Error()),
				Where:[2]int{from, from + len(test)},
				Why: string(test),
			}
		}
		if len(vs) % 3 != 0{
			return nil, &SVGError{
				What: fmt.Sprintf("Not enough vector, C must have en vector"),
				Where:[2]int{from, from + len(test)},
				Why: string(test),
			}
		}
		temp := make([][3]mgl32.Vec2, len(vs) /3)
		for i, v := range vs {
			temp[i/3][i%3] = v
		}
		return AbsoluteCubeTo{
			Arg:temp,
		},  nil
	case 'c':
		vs, err := vectors(test[1:])
		if err != nil {
			return nil, &SVGError{
				What: fmt.Sprintf("Argument Invalid : %s", err.Error()),
				Where:[2]int{from, from + len(test)},
				Why: string(test),
			}
		}
		if len(vs) % 3 != 0{
			return nil, &SVGError{
				What: fmt.Sprintf("Not enough vector, c must have en vector"),
				Where:[2]int{from, from + len(test)},
				Why: string(test),
			}
		}
		temp := make([][3]mgl32.Vec2, len(vs) /3)
		for i, v := range vs {
			temp[i/3][i%3] = v
		}
		return RelativeCubeTo{
			Arg:temp,
		},  nil
	case 'S':
		vs, err := vectors(test[1:])
		if err != nil {
			return nil, &SVGError{
				What: fmt.Sprintf("Argument Invalid : %s", err.Error()),
				Where:[2]int{from, from + len(test)},
				Why: string(test),
			}
		}
		if len(vs) % 2 != 0{
			return nil, &SVGError{
				What: fmt.Sprintf("Not enough vector, S must have 2n vector"),
				Where:[2]int{from, from + len(test)},
				Why: string(test),
			}
		}
		temp := make([][2]mgl32.Vec2, len(vs) /2)
		for i, v := range vs {
			temp[i/2][i%2] = v
		}
		return AbsoluteSmoothCubeTo{
			Arg:temp,
		},  nil
	case 's':
		vs, err := vectors(test[1:])
		if err != nil {
			return nil, &SVGError{
				What: fmt.Sprintf("Argument Invalid : %s", err.Error()),
				Where:[2]int{from, from + len(test)},
				Why: string(test),
			}
		}
		if len(vs) % 2 != 0{
			return nil, &SVGError{
				What: fmt.Sprintf("Not enough vector, s must have 2n vector"),
				Where:[2]int{from, from + len(test)},
				Why: string(test),
			}
		}
		temp := make([][2]mgl32.Vec2, len(vs) /2)
		for i, v := range vs {
			temp[i/2][i%2] = v
		}
		return RelativeSmoothCubeTo{
			Arg:temp,
		},  nil

	case 'A':
		vs, err := arcArgs(test[1:])
		if err != nil {
			return nil, &SVGError{
				What: fmt.Sprintf("Argument Invalid : %s", err.Error()),
				Where:[2]int{from, from + len(test)},
				Why: string(test),
			}
		}
		return AbsoluteArc{
			Arg:vs,
		}, nil
	case 'a':
		vs, err := arcArgs(test[1:])
		if err != nil {
			return nil, &SVGError{
				What: fmt.Sprintf("Argument Invalid : %s", err.Error()),
				Where:[2]int{from, from + len(test)},
				Why: string(test),
			}
		}
		return RelativeArc{
			Arg:vs,
		}, nil

	case 'B':
		vs, err := floats(test[1:])
		if err != nil {
			return nil, &SVGError{
				What: fmt.Sprintf("Argument Invalid : %s", err.Error()),
				Where:[2]int{from, from + len(test)},
				Why: string(test),
			}
		}
		return AbsoluteBearing{
			Arg:vs,
		}, nil
	case 'b':
		vs, err := floats(test[1:])
		if err != nil {
			return nil, &SVGError{
				What: fmt.Sprintf("Argument Invalid : %s", err.Error()),
				Where:[2]int{from, from + len(test)},
				Why: string(test),
			}
		}
		return RelativeBearing{
			Arg:vs,
		}, nil

	case 'Z':
		fallthrough
	case 'z':
		return CloseTo{}, nil
	}
	return nil, nil
}