package kernel

var(
	NearestNeighbor = _NearestNeighbor{}
	Bilinear = _Bilinear{}
	Bell = _Bell{}
	Hermite = _Hermite{}
	BicubicHalf = _Bicubic{A:-0.5}
	MitchellOneThird = makeMichell(1./3.)
	Lanczos2 = _Lanczos2{}
	Lanczos3 = _Lanczos3{}
)
