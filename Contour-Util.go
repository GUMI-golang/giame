package giame

import (
	"math"
	"github.com/go-gl/mathgl/mgl32"
)

type Float32SliceIgnoreSign []float32
func (s *Float32SliceIgnoreSign) Append(x float32) {
	p := float32(math.Abs(float64(x)))
	to := -1
	for i, v := range *s {
		v = float32(math.Abs(float64(v)))
		if p > v{
			to = i
		}else {
			break
		}
	}
	*s = append((*s)[:to + 1], append([]float32{x}, (*s)[to + 1:]...)...)
}


var minuszero = math.Float32frombits(0x80000000)


func devSquared(a, b, c mgl32.Vec2) float32 {
	devx := a[0] - 2*b[0] + c[0]
	devy := a[1] - 2*b[1] + c[1]
	return devx*devx + devy*devy
}
func lerp(t float32, p, q mgl32.Vec2) mgl32.Vec2 {
	return [2]float32{p[0] + t*(q[0]-p[0]), p[1] + t*(q[1]-p[1])}
}
func quadFromTo(from, pivot, to mgl32.Vec2) (res []mgl32.Vec2) {
	devsq := devSquared(from, pivot, to)
	if devsq >= 0.333 {
		const tol = 3
		n := 1 + int(math.Sqrt(math.Sqrt(tol*float64(devsq))))
		t, nInv := float32(0), 1/float32(n)
		for i := 0; i < n-1; i++ {
			t += nInv
			fromPivot := lerp(t, from, pivot)
			pivotTo := lerp(t, pivot, to)
			res = append(res, lerp(t, fromPivot, pivotTo))
		}
	}
	res = append(res, to)
	return res
}
func cubeFromTo(from, pivot1, pivot2, to mgl32.Vec2) (res []mgl32.Vec2) {
	devsq := devSquared(from, pivot1, to)
	if devsqAlt := devSquared(from, pivot2, to); devsq < devsqAlt {
		devsq = devsqAlt
	}
	if devsq >= 0.333 {
		const tol = 3
		n := 1 + int(math.Sqrt(math.Sqrt(tol*float64(devsq))))
		t, nInv := float32(0), 1/float32(n)
		for i := 0; i < n-1; i++ {
			t += nInv
			ab := lerp(t, from, pivot1)
			bc := lerp(t, pivot1, pivot2)
			cd := lerp(t, pivot2, to)
			abc := lerp(t, ab, bc)
			bcd := lerp(t, bc, cd)
			res = append(res, lerp(t, abc, bcd))
		}
	}
	res = append(res, to)
	return res
}

func Spacer() mgl32.Vec2 {
	return mgl32.Vec2{
		float32(math.NaN()),
		float32(math.NaN()),
	}
}
func IsSpacer(v mgl32.Vec2) bool {
	return math.IsNaN(float64(v.X()))
}