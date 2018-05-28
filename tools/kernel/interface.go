package kernel

type Kernel interface {
	Do(x float64) float64
	Rad() int
	//Compile() Filter
}