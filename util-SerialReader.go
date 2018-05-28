package giame

import (
	"io"
)

type serialReader struct {
	rd       io.Reader
	buf      []byte
	from, absFrom int

	cache    []byte
}

const bufsize = 1024

func newSVGReader(Reader io.Reader) *serialReader {
	return &serialReader{
		rd:    Reader,
		buf:   nil,
		from:  0,
		absFrom:    0,
		cache: nil,
	}
}
func (s *serialReader) NextFromTo() (data []byte, from, to int, err error) {
	if s.buf == nil || len(s.buf) <= s.from{
		s.buf = make([]byte, bufsize)
		n, err := s.rd.Read(s.buf)
		if err == io.EOF  && s.cache != nil{
			temp := s.cache
			s.cache = nil
			s.buf = nil
			return temp, s.absFrom - len(temp), s.absFrom, nil
		}
		if err != nil {
			return nil, 0,0, err
		}
		s.from = 0
		s.buf = s.buf[:n]
	}
	for i := s.from + 1; i < len(s.buf); i++ {
		s.absFrom += 1
		if IsCommand(s.buf[i]) {
			dts := append(s.cache, s.buf[s.from:i]...)
			s.from = i
			return dts, s.absFrom - len(dts), s.absFrom, nil
		}
	}
	s.absFrom += 1
	s.cache = append(s.cache, s.buf[s.from:]...)
	s.from = len(s.buf)
	return s.NextFromTo()
}

var command = []byte{
	'M', 'm', // https://www.w3.org/TR/SVG2/paths.html#PathDataMovetoCommands
	'A', 'a', // https://www.w3.org/TR/SVG2/paths.html#PathDataEllipticalArcCommands
	'B', 'b', // https://www.w3.org/TR/SVG2/paths.html#PathDataBearingCommands
	'Z', 'z', // https://www.w3.org/TR/SVG2/paths.html#PathDataClosePathCommand
	'Q', 'q', 'T', 't', // https://www.w3.org/TR/SVG2/paths.html#PathDataQuadraticBezierCommands
	'C', 'c', 'S', 's', // https://www.w3.org/TR/SVG2/paths.html#PathDataCubicBezierCommands
	'L', 'l', 'H', 'h', 'V', 'v', // https://www.w3.org/TR/SVG2/paths.html#PathDataLinetoCommands
}
func IsCommand(b byte) bool {
	for _, v := range command {
		if b == v{
			return true
		}
	}
	return false
}

