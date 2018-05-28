package mask

var Laplacian = Mask3{
	Data: [3][3]float64{
		{0, -0.25, 0},
		{-0.25, 2, -0.25},
		{0, -0.25, 0},
	},
	Name: "Laplacian",
}
var LaplacianExtend = Mask3{
	Data: [3][3]float64{
		{-.125, -.125, -.125},
		{-.125, 2, -.125},
		{-.125, -.125, -.125},
	},
	Name: "LaplacianExtend",
}
