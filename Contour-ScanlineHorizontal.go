package giame

import (
	"github.com/go-gl/mathgl/mgl32"
	"math"
)



// Contour:ScanlineHorizontal
type ScanlineHorizontal struct {
	contourInfo
	Value []float32
	Range []int32
}

func (s *ContourRaw) ToContourHorizontal() (*ScanlineHorizontal) {
	b := SetupScanlineHorizontal(s.contourInfo)
	for i := 0; i < len(s.Points)-1; i++ {
		start, end := s.Points[i], s.Points[i+1]
		if IsSpacer(start) || IsSpacer(end){
			continue
		} else {
			b.Write(start, end)
		}
	}
	return b.Build().(*ScanlineHorizontal)
}












// ContourBuilder:BuilderContourHorizontal
type BuilderContourHorizontal struct {
	data []Float32SliceIgnoreSign
	sz   mgl32.Vec2
	isEdit bool
	res *ScanlineHorizontal
}
func SetupScanlineHorizontal(c contourInfo) ContourBuilder {
	return &BuilderContourHorizontal{
		data: make([]Float32SliceIgnoreSign, c.Bound.Dy()),
		sz:   mgl32.Vec2{float32(c.Bound.Dx()), float32(c.Bound.Dy())},
		res: &ScanlineHorizontal{
			contourInfo : c,
		},
	}
}
func (s *BuilderContourHorizontal) Write(start, end mgl32.Vec2) {
	var dir float32 = 1.
	if start[1] > end[1] {
		dir = -1.
		start, end = end, start
	}

	miny := math.Floor(float64(start[1]))
	maxy := math.Ceil(float64(end[1]))

	//
	if end[1]-start[1] < 0.0001 {
		return
	}
	delta := (end[0] - start[0]) / (end[1] - start[1])
	var x = start[0] + (float32(miny)-start[1])*delta
	for y := int(miny); y <= int(maxy); y, x = y+1, x+delta {
		if !(start[1] <= float32(y) && float32(y) < end[1]) {
			continue
		}
		if x < 0 || x >= s.sz[0] || y < 0 || y >= int(s.sz[1]){
			continue
		}
		if x == 0. && dir == -1.{
			// minus zero
			s.data[y].Append(minuszero)
		}else {
			s.data[y].Append(x * dir)
		}

	}
}
func (s *BuilderContourHorizontal) Build() Contour {
	temp := s.res
	s.res = &ScanlineHorizontal{
		contourInfo : s.res.contourInfo,
	}

	var end int32
	temp.Range = make([]int32, 0, len(s.data))
	for _, line := range s.data {
		end += int32(len(line))
		temp.Value = append(temp.Value, line...)
		temp.Range = append(temp.Range, end)
	}
	return temp
}
