package tools

import "image"

func allRGBA(imgs ... image.Image) (rgbas []*image.RGBA, ok bool) {
	rgbas = make([]*image.RGBA, len(imgs))
	for i, img := range imgs {
		temp, check := img.(*image.RGBA)
		if !check{
			return nil, ok
		}
		rgbas[i] = temp
	}
	ok = true
	return
}

func Iclamp(i, min, max int) int {
	if i <= min {
		return min
	}
	if i >= max {
		return max
	}
	return i
}