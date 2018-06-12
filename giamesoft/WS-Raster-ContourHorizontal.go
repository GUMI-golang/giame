package giamesoft

import (
	"github.com/GUMI-golang/giame"
	"math"
)


func ScanlineHorizontal(ws *Workspace, rq *giame.ScanlineHorizontal) {
	var (
		start int32 = 0
	)
	for y := 0; y < ws.Height; y++ {
		end := rq.Range[y]
		line := rq.Value[start:end]
		start = end
		//
		if len(line) == 0{
			continue
		}
		//
		var dir float32 = 0.
		var from, to float32
		for _, temp := range line {
			to = float32(math.Abs(float64(temp)))
			if dir == 0 {
				dir += CheckSign(temp)
				from = to
				continue
			}
			a, b := startPoint(from), endPoint(to)
			if a == b{
				ws.Set(int(a), y, to - from)
			}else {
				ws.Set(int(a), y, float32(math.Min(float64(a - from), .5)) + .5)
				ws.Set(int(b), y, float32(math.Min(float64(to - b), .5)) + .5)
				for x := a + 1; x < b; x++ {
					ws.Set(int(x), y, 1.)
				}
			}
			dir += CheckSign(temp)
			from = to
		}
	}
	return
}
func startPoint(f32 float32) float32 {
	return float32(math.Ceil(float64(f32 - .5))) + .5
}
func endPoint(f32 float32) float32 {
	return float32(math.Floor(float64(f32 + .5))) - .5
}
func CheckSign(f32 float32) float32 {
	return float32(int(math.Float32bits(f32) >> 31) * -2 + 1)
}