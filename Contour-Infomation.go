package giame

import "image"

type contourInfo struct {
	Option       Option
	Bound        image.Rectangle
	Filler       Filler
}

func (s *contourInfo) GetOption() Option {
	return s.Option
}

func (s *contourInfo) GetBound() image.Rectangle {
	return s.Bound
}

func (s *contourInfo) GetFiller() Filler {
	return s.Filler
}

type Option struct {
	Filter OptionFilter
}

func (s Option) String() string {
	result := "["
	switch s.Filter {
	case FilterNone:
		result += "FilterNone"
	case FilterLaplacian:
		result += "FilterLaplacian"
	case FilterLaplacianExtend:
		result += "FilterLaplacianExtend"
	}
	result += "]"
	return result
}

type OptionFilter int32

const (
	FilterNone            OptionFilter = iota
	FilterLaplacian       OptionFilter = iota
	FilterLaplacianExtend OptionFilter = iota
)
