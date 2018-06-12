package giame

import (
	"github.com/go-gl/mathgl/mgl32"
	"math"
)



// Contour:ContourVertical
type ScanlineVertical struct {
	contourInfo
	Value []float32
	Range []int32
}

func (s *ContourRaw) ToContourVertical() (*ScanlineVertical) {
	b := SetupScanlineVertical(s.contourInfo)
	for i := 0; i < len(s.Points)-1; i++ {
		start, end := s.Points[i], s.Points[i+1]
		if IsSpacer(start) || IsSpacer(end){
			continue
		} else {
			b.Write(start, end)
		}
	}
	return b.Build().(*ScanlineVertical)
}



// ContourBuilder:BuilderContourVertical
type BuilderContourVertical struct {
	data []Float32SliceIgnoreSign
	sz   mgl32.Vec2
	isEdit bool
	res *ScanlineVertical
}
func SetupScanlineVertical(c contourInfo) ContourBuilder {
	return &BuilderContourVertical{
		data: make([]Float32SliceIgnoreSign, c.Bound.Dy()),
		sz:   mgl32.Vec2{float32(c.Bound.Dx()), float32(c.Bound.Dy())},
		res: &ScanlineVertical{
			contourInfo : c,
		},
	}
}
func (s *BuilderContourVertical) Write(start, end mgl32.Vec2) {
	var dir float32 = 1.
	if start[0] > end[0] {
		dir = -1.
		start, end = end, start
	}
	minx := math.Floor(float64(start[0]))
	maxx := math.Ceil(float64(end[0]))

	//
	if end[0]-start[0] < 0.0001 {
		return
	}
	delta := (end[1] - start[1]) / (end[0] - start[0])
	var y = start[1] + (float32(minx)-start[0])*delta
	for x := int(minx); x <= int(maxx); x, y = x+1, y+delta {
		if !(start[0] <= float32(x) && float32(x) < end[0]) {
			continue
		}
		if y < 0 || y >= s.sz[1] || x < 0 || x >= int(s.sz[0]){
			continue
		}
		if y == 0. && dir == -1.{
			// minus zero
			s.data[x].Append(minuszero)
		}else {
			s.data[x].Append(y * dir)
		}

	}
}
func (s *BuilderContourVertical) Build() Contour {
	temp := s.res
	s.res = &ScanlineVertical{
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
