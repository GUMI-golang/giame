package filter

import "fmt"

// Filter ONLY can be Filter3, Filter5, Filter7, FilterN
//
type (
	Filter interface {
		DataF32() ([]float32)
		Size() (w, h int)
		Radius() int
		Length() int
		fmt.Stringer
	}
	Filter3 struct {
		Data [3][3]float32
		Name string
	}
	Filter5 struct {
		Data [5][5]float32
		Name string
	}
	Filter7 struct {
		Data [7][7]float32
		Name string
	}
	FilterN struct {
		Data [][]float32
		n int
		Name string
	}
)


func (s Filter3) DataF32() ([]float32) {
	return []float32{
		float32(s.Data[0][0]), float32(s.Data[0][1]), float32(s.Data[0][2]),
		float32(s.Data[1][0]), float32(s.Data[1][1]), float32(s.Data[1][2]),
		float32(s.Data[2][0]), float32(s.Data[2][1]), float32(s.Data[2][2]),
	}
}
func (s Filter3) Size() (w, h int) {
	return 3,3
}
func (s Filter3) Radius() int {
	return 1
}
func (s Filter3) Length() int {
	return 9
}
func (s Filter3) String() (string) {
	return s.Name
}

func (s Filter5) DataF32() ([]float32) {
	return []float32{
		float32(s.Data[0][0]), float32(s.Data[0][1]), float32(s.Data[0][2]), float32(s.Data[0][3]), float32(s.Data[0][4]),
		float32(s.Data[1][0]), float32(s.Data[1][1]), float32(s.Data[1][2]), float32(s.Data[1][3]), float32(s.Data[1][4]),
		float32(s.Data[2][0]), float32(s.Data[2][1]), float32(s.Data[2][2]), float32(s.Data[2][3]), float32(s.Data[2][4]),
		float32(s.Data[3][0]), float32(s.Data[3][1]), float32(s.Data[3][2]), float32(s.Data[3][3]), float32(s.Data[3][4]),
		float32(s.Data[4][0]), float32(s.Data[4][1]), float32(s.Data[4][2]), float32(s.Data[4][3]), float32(s.Data[4][4]),
	}
}
func (s Filter5) Size() (w, h int) {
	return 5,5
}
func (s Filter5) Radius() int {
	return 2
}
func (s Filter5) Length() int {
	return 25
}
func (s Filter5) String() (string) {
	return s.Name
}

func (s Filter7) DataF32() ([]float32) {
	return []float32{
		float32(s.Data[0][0]), float32(s.Data[0][1]), float32(s.Data[0][2]), float32(s.Data[0][3]), float32(s.Data[0][4]), float32(s.Data[0][5]), float32(s.Data[0][6]),
		float32(s.Data[1][0]), float32(s.Data[1][1]), float32(s.Data[1][2]), float32(s.Data[1][3]), float32(s.Data[1][4]), float32(s.Data[1][5]), float32(s.Data[1][6]),
		float32(s.Data[2][0]), float32(s.Data[2][1]), float32(s.Data[2][2]), float32(s.Data[2][3]), float32(s.Data[2][4]), float32(s.Data[2][5]), float32(s.Data[2][6]),
		float32(s.Data[3][0]), float32(s.Data[3][1]), float32(s.Data[3][2]), float32(s.Data[3][3]), float32(s.Data[3][4]), float32(s.Data[3][5]), float32(s.Data[3][6]),
		float32(s.Data[4][0]), float32(s.Data[4][1]), float32(s.Data[4][2]), float32(s.Data[4][3]), float32(s.Data[4][4]), float32(s.Data[4][5]), float32(s.Data[4][6]),
		float32(s.Data[5][0]), float32(s.Data[5][1]), float32(s.Data[5][2]), float32(s.Data[5][3]), float32(s.Data[5][4]), float32(s.Data[5][5]), float32(s.Data[5][6]),
		float32(s.Data[6][0]), float32(s.Data[6][1]), float32(s.Data[6][2]), float32(s.Data[6][3]), float32(s.Data[6][4]), float32(s.Data[6][5]), float32(s.Data[6][6]),
	}
}
func (s Filter7) Size() (w, h int) {
	return 7,7
}
func (s Filter7) Radius() int {
	return 3
}
func (s Filter7) Length() int {
	return 49
}
func (s Filter7) String() (string) {
	return s.Name
}

func (s FilterN) DataF32() ([]float32) {
	w, h := s.Size()
	temp := make([]float32, s.Length())
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			temp[x + y * w] = float32(s.Data[y][x])
		}
	}
	return []float32{
		float32(s.Data[0][0]), float32(s.Data[0][1]), float32(s.Data[0][2]), float32(s.Data[0][3]), float32(s.Data[0][4]), float32(s.Data[0][5]), float32(s.Data[0][6]),
		float32(s.Data[1][0]), float32(s.Data[1][1]), float32(s.Data[1][2]), float32(s.Data[1][3]), float32(s.Data[1][4]), float32(s.Data[1][5]), float32(s.Data[1][6]),
		float32(s.Data[2][0]), float32(s.Data[2][1]), float32(s.Data[2][2]), float32(s.Data[2][3]), float32(s.Data[2][4]), float32(s.Data[2][5]), float32(s.Data[2][6]),
		float32(s.Data[3][0]), float32(s.Data[3][1]), float32(s.Data[3][2]), float32(s.Data[3][3]), float32(s.Data[3][4]), float32(s.Data[3][5]), float32(s.Data[3][6]),
		float32(s.Data[4][0]), float32(s.Data[4][1]), float32(s.Data[4][2]), float32(s.Data[4][3]), float32(s.Data[4][4]), float32(s.Data[4][5]), float32(s.Data[4][6]),
		float32(s.Data[5][0]), float32(s.Data[5][1]), float32(s.Data[5][2]), float32(s.Data[5][3]), float32(s.Data[5][4]), float32(s.Data[5][5]), float32(s.Data[5][6]),
		float32(s.Data[6][0]), float32(s.Data[6][1]), float32(s.Data[6][2]), float32(s.Data[6][3]), float32(s.Data[6][4]), float32(s.Data[6][5]), float32(s.Data[6][6]),
	}
}
func (s FilterN) Size() (w, h int) {
	return s.n * 2 + 1, s.n * 2 + 1
}
func (s FilterN) Radius() int {
	return s.n
}
func (s FilterN) Length() int {
	w, h := s.Size()
	return w * h
}
func (s FilterN) String() (string) {
	return s.Name
}
