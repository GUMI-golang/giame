package giamesoft

import (
	"github.com/GUMI-golang/giame/tools/filter"
	"github.com/GUMI-golang/giame/tools"
	"github.com/go-gl/mathgl/mgl32"
)

func Filting(ws *Workspace, f filter.Filter) {
	switch af := f.(type) {
	case filter.Filter3:
		Filting3(ws, af)
	case filter.Filter5:
		Filting5(ws, af)
	case filter.Filter7:
		//Filting7(ws, af)
	default:
		data := af.DataF32()
		n := af.Radius()
		w := n * 2 + 1
		for y := 0; y < ws.Height; y++{
			for x := 0; x < ws.Width; x++{
				var sum float32
				for i, v := range data {
					pickx, picky := i % w - n, i /w - n
					sum += ws.Get(tools.Iclamp(pickx, 0, ws.Width - 1), tools.Iclamp(picky, 0, ws.Height - 1)) * float32(v)
				}
				ws.Set(x, y, mgl32.Clamp(sum, 0, 1))
			}
		}
	}
}
func Filting3(ws *Workspace, f3 filter.Filter3)  {

	for y := 0; y < ws.Height; y++{
		for x := 0; x < ws.Width; x++{
			var sum float32
			sum += ws.Get(tools.Iclamp(x - 1, 0, ws.Width - 1), tools.Iclamp(y - 1, 0, ws.Height - 1)) * float32(f3.Data[0][0])
			sum += ws.Get(tools.Iclamp(x + 0, 0, ws.Width - 1), tools.Iclamp(y - 1, 0, ws.Height - 1)) * float32(f3.Data[0][1])
			sum += ws.Get(tools.Iclamp(x + 1, 0, ws.Width - 1), tools.Iclamp(y - 1, 0, ws.Height - 1)) * float32(f3.Data[0][2])

			sum += ws.Get(tools.Iclamp(x - 1, 0, ws.Width - 1), tools.Iclamp(y + 0, 0, ws.Height - 1)) * float32(f3.Data[1][0])
			sum += ws.Get(tools.Iclamp(x + 0, 0, ws.Width - 1), tools.Iclamp(y + 0, 0, ws.Height - 1)) * float32(f3.Data[1][1])
			sum += ws.Get(tools.Iclamp(x + 1, 0, ws.Width - 1), tools.Iclamp(y + 0, 0, ws.Height - 1)) * float32(f3.Data[1][2])

			sum += ws.Get(tools.Iclamp(x - 1, 0, ws.Width - 1), tools.Iclamp(y + 1, 0, ws.Height - 1)) * float32(f3.Data[2][0])
			sum += ws.Get(tools.Iclamp(x + 0, 0, ws.Width - 1), tools.Iclamp(y + 1, 0, ws.Height - 1)) * float32(f3.Data[2][1])
			sum += ws.Get(tools.Iclamp(x + 1, 0, ws.Width - 1), tools.Iclamp(y + 1, 0, ws.Height - 1)) * float32(f3.Data[2][2])
			ws.Set(x, y, mgl32.Clamp(sum, 0, 1))
		}
	}
}
func Filting5(ws *Workspace, f5 filter.Filter5)  {

	for y := 0; y < ws.Height; y++{
		for x := 0; x < ws.Width; x++{
			var sum float32
			sum += ws.Get(tools.Iclamp(x - 2, 0, ws.Width - 1), tools.Iclamp(y - 2, 0, ws.Height - 1)) * float32(f5.Data[0][0])
			sum += ws.Get(tools.Iclamp(x - 1, 0, ws.Width - 1), tools.Iclamp(y - 2, 0, ws.Height - 1)) * float32(f5.Data[0][1])
			sum += ws.Get(tools.Iclamp(x + 0, 0, ws.Width - 1), tools.Iclamp(y - 2, 0, ws.Height - 1)) * float32(f5.Data[0][2])
			sum += ws.Get(tools.Iclamp(x + 1, 0, ws.Width - 1), tools.Iclamp(y - 2, 0, ws.Height - 1)) * float32(f5.Data[0][3])
			sum += ws.Get(tools.Iclamp(x + 2, 0, ws.Width - 1), tools.Iclamp(y - 2, 0, ws.Height - 1)) * float32(f5.Data[0][4])

			sum += ws.Get(tools.Iclamp(x - 2, 0, ws.Width - 1), tools.Iclamp(y - 1, 0, ws.Height - 1)) * float32(f5.Data[1][0])
			sum += ws.Get(tools.Iclamp(x - 1, 0, ws.Width - 1), tools.Iclamp(y - 1, 0, ws.Height - 1)) * float32(f5.Data[1][1])
			sum += ws.Get(tools.Iclamp(x + 0, 0, ws.Width - 1), tools.Iclamp(y - 1, 0, ws.Height - 1)) * float32(f5.Data[1][2])
			sum += ws.Get(tools.Iclamp(x + 1, 0, ws.Width - 1), tools.Iclamp(y - 1, 0, ws.Height - 1)) * float32(f5.Data[1][3])
			sum += ws.Get(tools.Iclamp(x + 2, 0, ws.Width - 1), tools.Iclamp(y - 1, 0, ws.Height - 1)) * float32(f5.Data[1][4])

			sum += ws.Get(tools.Iclamp(x - 2, 0, ws.Width - 1), tools.Iclamp(y + 0, 0, ws.Height - 1)) * float32(f5.Data[2][0])
			sum += ws.Get(tools.Iclamp(x - 1, 0, ws.Width - 1), tools.Iclamp(y + 0, 0, ws.Height - 1)) * float32(f5.Data[2][1])
			sum += ws.Get(tools.Iclamp(x + 0, 0, ws.Width - 1), tools.Iclamp(y + 0, 0, ws.Height - 1)) * float32(f5.Data[2][2])
			sum += ws.Get(tools.Iclamp(x + 1, 0, ws.Width - 1), tools.Iclamp(y + 0, 0, ws.Height - 1)) * float32(f5.Data[2][3])
			sum += ws.Get(tools.Iclamp(x + 2, 0, ws.Width - 1), tools.Iclamp(y + 0, 0, ws.Height - 1)) * float32(f5.Data[2][4])

			sum += ws.Get(tools.Iclamp(x - 2, 0, ws.Width - 1), tools.Iclamp(y + 1, 0, ws.Height - 1)) * float32(f5.Data[3][0])
			sum += ws.Get(tools.Iclamp(x - 1, 0, ws.Width - 1), tools.Iclamp(y + 1, 0, ws.Height - 1)) * float32(f5.Data[3][1])
			sum += ws.Get(tools.Iclamp(x + 0, 0, ws.Width - 1), tools.Iclamp(y + 1, 0, ws.Height - 1)) * float32(f5.Data[3][2])
			sum += ws.Get(tools.Iclamp(x + 1, 0, ws.Width - 1), tools.Iclamp(y + 1, 0, ws.Height - 1)) * float32(f5.Data[3][3])
			sum += ws.Get(tools.Iclamp(x + 2, 0, ws.Width - 1), tools.Iclamp(y + 1, 0, ws.Height - 1)) * float32(f5.Data[3][4])

			sum += ws.Get(tools.Iclamp(x - 2, 0, ws.Width - 1), tools.Iclamp(y + 2, 0, ws.Height - 1)) * float32(f5.Data[4][0])
			sum += ws.Get(tools.Iclamp(x - 1, 0, ws.Width - 1), tools.Iclamp(y + 2, 0, ws.Height - 1)) * float32(f5.Data[4][1])
			sum += ws.Get(tools.Iclamp(x + 0, 0, ws.Width - 1), tools.Iclamp(y + 2, 0, ws.Height - 1)) * float32(f5.Data[4][2])
			sum += ws.Get(tools.Iclamp(x + 1, 0, ws.Width - 1), tools.Iclamp(y + 2, 0, ws.Height - 1)) * float32(f5.Data[4][3])
			sum += ws.Get(tools.Iclamp(x + 2, 0, ws.Width - 1), tools.Iclamp(y + 2, 0, ws.Height - 1)) * float32(f5.Data[4][4])
			ws.Set(x, y, mgl32.Clamp(sum, 0, 1))
		}
	}
}
func Filting7(ws *Workspace, f7 filter.Filter7)  {

	for y := 0; y < ws.Height; y++{
		for x := 0; x < ws.Width; x++{
			var sum float32
			sum += ws.Get(tools.Iclamp(x - 3, 0, ws.Width - 1), tools.Iclamp(y - 3, 0, ws.Height - 1)) * float32(f7.Data[0][0])
			sum += ws.Get(tools.Iclamp(x - 2, 0, ws.Width - 1), tools.Iclamp(y - 3, 0, ws.Height - 1)) * float32(f7.Data[0][1])
			sum += ws.Get(tools.Iclamp(x - 1, 0, ws.Width - 1), tools.Iclamp(y - 3, 0, ws.Height - 1)) * float32(f7.Data[0][2])
			sum += ws.Get(tools.Iclamp(x + 0, 0, ws.Width - 1), tools.Iclamp(y - 3, 0, ws.Height - 1)) * float32(f7.Data[0][3])
			sum += ws.Get(tools.Iclamp(x + 1, 0, ws.Width - 1), tools.Iclamp(y - 3, 0, ws.Height - 1)) * float32(f7.Data[0][4])
			sum += ws.Get(tools.Iclamp(x + 2, 0, ws.Width - 1), tools.Iclamp(y - 3, 0, ws.Height - 1)) * float32(f7.Data[0][5])
			sum += ws.Get(tools.Iclamp(x + 3, 0, ws.Width - 1), tools.Iclamp(y - 3, 0, ws.Height - 1)) * float32(f7.Data[0][6])

			sum += ws.Get(tools.Iclamp(x - 3, 0, ws.Width - 1), tools.Iclamp(y - 2, 0, ws.Height - 1)) * float32(f7.Data[1][0])
			sum += ws.Get(tools.Iclamp(x - 2, 0, ws.Width - 1), tools.Iclamp(y - 2, 0, ws.Height - 1)) * float32(f7.Data[1][1])
			sum += ws.Get(tools.Iclamp(x - 1, 0, ws.Width - 1), tools.Iclamp(y - 2, 0, ws.Height - 1)) * float32(f7.Data[1][2])
			sum += ws.Get(tools.Iclamp(x + 0, 0, ws.Width - 1), tools.Iclamp(y - 2, 0, ws.Height - 1)) * float32(f7.Data[1][3])
			sum += ws.Get(tools.Iclamp(x + 1, 0, ws.Width - 1), tools.Iclamp(y - 2, 0, ws.Height - 1)) * float32(f7.Data[1][4])
			sum += ws.Get(tools.Iclamp(x + 2, 0, ws.Width - 1), tools.Iclamp(y - 2, 0, ws.Height - 1)) * float32(f7.Data[1][5])
			sum += ws.Get(tools.Iclamp(x + 3, 0, ws.Width - 1), tools.Iclamp(y - 2, 0, ws.Height - 1)) * float32(f7.Data[1][6])

			sum += ws.Get(tools.Iclamp(x - 3, 0, ws.Width - 1), tools.Iclamp(y - 1, 0, ws.Height - 1)) * float32(f7.Data[2][0])
			sum += ws.Get(tools.Iclamp(x - 2, 0, ws.Width - 1), tools.Iclamp(y - 1, 0, ws.Height - 1)) * float32(f7.Data[2][1])
			sum += ws.Get(tools.Iclamp(x - 1, 0, ws.Width - 1), tools.Iclamp(y - 1, 0, ws.Height - 1)) * float32(f7.Data[2][2])
			sum += ws.Get(tools.Iclamp(x + 0, 0, ws.Width - 1), tools.Iclamp(y - 1, 0, ws.Height - 1)) * float32(f7.Data[2][3])
			sum += ws.Get(tools.Iclamp(x + 1, 0, ws.Width - 1), tools.Iclamp(y - 1, 0, ws.Height - 1)) * float32(f7.Data[2][4])
			sum += ws.Get(tools.Iclamp(x + 2, 0, ws.Width - 1), tools.Iclamp(y - 1, 0, ws.Height - 1)) * float32(f7.Data[2][5])
			sum += ws.Get(tools.Iclamp(x + 3, 0, ws.Width - 1), tools.Iclamp(y - 1, 0, ws.Height - 1)) * float32(f7.Data[2][6])

			sum += ws.Get(tools.Iclamp(x - 3, 0, ws.Width - 1), tools.Iclamp(y + 0, 0, ws.Height - 1)) * float32(f7.Data[3][0])
			sum += ws.Get(tools.Iclamp(x - 2, 0, ws.Width - 1), tools.Iclamp(y + 0, 0, ws.Height - 1)) * float32(f7.Data[3][1])
			sum += ws.Get(tools.Iclamp(x - 1, 0, ws.Width - 1), tools.Iclamp(y + 0, 0, ws.Height - 1)) * float32(f7.Data[3][2])
			sum += ws.Get(tools.Iclamp(x + 0, 0, ws.Width - 1), tools.Iclamp(y + 0, 0, ws.Height - 1)) * float32(f7.Data[3][3])
			sum += ws.Get(tools.Iclamp(x + 1, 0, ws.Width - 1), tools.Iclamp(y + 0, 0, ws.Height - 1)) * float32(f7.Data[3][4])
			sum += ws.Get(tools.Iclamp(x + 2, 0, ws.Width - 1), tools.Iclamp(y + 0, 0, ws.Height - 1)) * float32(f7.Data[3][5])
			sum += ws.Get(tools.Iclamp(x + 3, 0, ws.Width - 1), tools.Iclamp(y + 0, 0, ws.Height - 1)) * float32(f7.Data[3][6])

			sum += ws.Get(tools.Iclamp(x - 3, 0, ws.Width - 1), tools.Iclamp(y + 1, 0, ws.Height - 1)) * float32(f7.Data[4][0])
			sum += ws.Get(tools.Iclamp(x - 2, 0, ws.Width - 1), tools.Iclamp(y + 1, 0, ws.Height - 1)) * float32(f7.Data[4][1])
			sum += ws.Get(tools.Iclamp(x - 1, 0, ws.Width - 1), tools.Iclamp(y + 1, 0, ws.Height - 1)) * float32(f7.Data[4][2])
			sum += ws.Get(tools.Iclamp(x + 0, 0, ws.Width - 1), tools.Iclamp(y + 1, 0, ws.Height - 1)) * float32(f7.Data[4][3])
			sum += ws.Get(tools.Iclamp(x + 1, 0, ws.Width - 1), tools.Iclamp(y + 1, 0, ws.Height - 1)) * float32(f7.Data[4][4])
			sum += ws.Get(tools.Iclamp(x + 2, 0, ws.Width - 1), tools.Iclamp(y + 1, 0, ws.Height - 1)) * float32(f7.Data[4][5])
			sum += ws.Get(tools.Iclamp(x + 3, 0, ws.Width - 1), tools.Iclamp(y + 1, 0, ws.Height - 1)) * float32(f7.Data[4][6])

			sum += ws.Get(tools.Iclamp(x - 3, 0, ws.Width - 1), tools.Iclamp(y + 2, 0, ws.Height - 1)) * float32(f7.Data[5][0])
			sum += ws.Get(tools.Iclamp(x - 2, 0, ws.Width - 1), tools.Iclamp(y + 2, 0, ws.Height - 1)) * float32(f7.Data[5][1])
			sum += ws.Get(tools.Iclamp(x - 1, 0, ws.Width - 1), tools.Iclamp(y + 2, 0, ws.Height - 1)) * float32(f7.Data[5][2])
			sum += ws.Get(tools.Iclamp(x + 0, 0, ws.Width - 1), tools.Iclamp(y + 2, 0, ws.Height - 1)) * float32(f7.Data[5][3])
			sum += ws.Get(tools.Iclamp(x + 1, 0, ws.Width - 1), tools.Iclamp(y + 2, 0, ws.Height - 1)) * float32(f7.Data[5][4])
			sum += ws.Get(tools.Iclamp(x + 2, 0, ws.Width - 1), tools.Iclamp(y + 2, 0, ws.Height - 1)) * float32(f7.Data[5][5])
			sum += ws.Get(tools.Iclamp(x + 3, 0, ws.Width - 1), tools.Iclamp(y + 2, 0, ws.Height - 1)) * float32(f7.Data[5][6])

			sum += ws.Get(tools.Iclamp(x - 3, 0, ws.Width - 1), tools.Iclamp(y + 3, 0, ws.Height - 1)) * float32(f7.Data[6][0])
			sum += ws.Get(tools.Iclamp(x - 2, 0, ws.Width - 1), tools.Iclamp(y + 3, 0, ws.Height - 1)) * float32(f7.Data[6][1])
			sum += ws.Get(tools.Iclamp(x - 1, 0, ws.Width - 1), tools.Iclamp(y + 3, 0, ws.Height - 1)) * float32(f7.Data[6][2])
			sum += ws.Get(tools.Iclamp(x + 0, 0, ws.Width - 1), tools.Iclamp(y + 3, 0, ws.Height - 1)) * float32(f7.Data[6][3])
			sum += ws.Get(tools.Iclamp(x + 1, 0, ws.Width - 1), tools.Iclamp(y + 3, 0, ws.Height - 1)) * float32(f7.Data[6][4])
			sum += ws.Get(tools.Iclamp(x + 2, 0, ws.Width - 1), tools.Iclamp(y + 3, 0, ws.Height - 1)) * float32(f7.Data[6][5])
			sum += ws.Get(tools.Iclamp(x + 3, 0, ws.Width - 1), tools.Iclamp(y + 3, 0, ws.Height - 1)) * float32(f7.Data[6][6])
			ws.Set(x, y, mgl32.Clamp(sum, 0, 1))
		}
	}
}
