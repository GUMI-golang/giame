package giame

import "image"

type (
	Driver interface {
		MakeResult(w, h int) Result
		Init() (err error)
		Close() error
	}

	Result interface {
		Size() (w, h int)
		Image() image.Image
		Clear()
		Request(dr Driver, rq Contour)
	}
)
