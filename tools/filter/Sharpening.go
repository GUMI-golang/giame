package filter


// 3 x 3
var (
	Laplacian = Filter3{
		Data: [3][3]float32{
			{0, -0.25, 0},
			{-0.25, 2, -0.25},
			{0, -0.25, 0},
		},
		Name: "Laplacian",
	}
	LaplacianExtend = Filter3{
		Data: [3][3]float32{
			{-.125, -.125, -.125},
			{-.125, 2, -.125},
			{-.125, -.125, -.125},
		},
		Name: "LaplacianExtend",
	}
)

