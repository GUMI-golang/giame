package kernel

import (
	"fmt"
	"math"
)

const almostzero = .000001

type _NearestNeighbor struct{}

func (_NearestNeighbor) Do(x float64) float64 {
	if -.5 < x && x <= .5 {
		return 1
	}
	return 0
}
func (_NearestNeighbor) Rad() int {
	return 0
}
func (_NearestNeighbor) String() string {
	return "NearestNeighbor"
}

type _Bilinear struct{}
func (_Bilinear) Do(x float64) float64 {
	x = math.Abs(x)
	if x <= 1 {
		return 1 - x
	}
	return 0
}
func (_Bilinear) Rad() int {
	return 1
}
func (_Bilinear) String() string {
	return "Bilinear"
}

type _Bell struct{}
func (_Bell) Do(x float64) float64 {
	x = math.Abs(x)
	if x <= .5 {
		return 0.75 - x*x
	}
	if x <= 1.5 {
		return .5 * (x - 1.5) * (x - 1.5)
	}
	return 0
}
func (_Bell) Rad() int {
	return 2
}
func (_Bell) String() string {
	return "Bell"
}

type _Hermite struct{}
func (_Hermite) Do(x float64) float64 {
	x = math.Abs(x)
	if x <= 1 {
		return 2*x*x*x - 3*x*x + 1
	}
	return 0
}
func (_Hermite) Rad() int {
	return 1
}
func (_Hermite) String() string {
	return "Hermite"
}

type _Bicubic struct{ A float64 }
func (s _Bicubic) Do(x float64) float64 {
	x = math.Abs(x)
	if x <= 1 {
		return (s.A+2)*x*x*x - (s.A+3)*x*x + 1
	}
	if x < 2 {
		return (s.A * x * x * x) - (5 * s.A * x * x) + (8 * s.A * x) - (4 * s.A)
	}
	return 0
}
func (_Bicubic) Rad() int {
	return 2
}
func (_Bicubic) String() string {
	return "Bicubic"
}

type _Michell struct{ B, C float64 }
func makeMichell(b float64) _Michell {
	return _Michell{B: b, C: (1 - b)/2}
}
func (s _Michell) Do(x float64) float64 {
	x = math.Abs(x)
	if x < 1 {
		return ((12-9*s.B-6*s.C)*x*x*x + (-18+12*s.B+6*s.C)*x*x + (6 - 2*s.B)) / 6
	}
	if x < 2 {
		return ((-s.B-6*s.C)*x*x*x + (6*s.B+30*s.C)*x*x + (-12*s.B-48*s.C)*x + (8*s.B + 24*s.C)) / 6
	}
	return 0
}
func (_Michell) Rad() int {
	return 2
}
func (s _Michell) String() string {
	return fmt.Sprintf("Michell(B:%v, C:%v)", s.B, s.C)
}

type (
	_Lanczos2 struct{}
	_Lanczos3 struct{}
)

func (_Lanczos2) Do(x float64) float64 {
	const a = 2
	if x == 0 {
		return 1
	}
	if -a <= x && x <= a {
		return (a * math.Sin(math.Pi*x) * math.Sin((math.Pi*x)/a)) / (math.Pi * math.Pi * x * x)
	}
	return 0
}
func (_Lanczos2) Rad() int {
	return 2
}
func (_Lanczos2) String() string {
	return "Lanczos2"
}
func (_Lanczos3) Do(x float64) float64 {
	const a = 3
	if x == 0 {
		return 1
	}
	if -a <= x && x <= a {
		return (a * math.Sin(math.Pi*x) * math.Sin((math.Pi*x)/a)) / (math.Pi * math.Pi * x * x)
	}
	return 0
}
func (_Lanczos3) Rad() int {
	return 3
}
func (_Lanczos3) String() string {
	return "Lanczos3"
}
