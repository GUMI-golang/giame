package main

import (
	"os"
	"image/png"
	"image"
	"image/color"
)

func main() {
	f, err := os.OpenFile("Pix16.png", os.O_CREATE | os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	temp := image.NewRGBA(image.Rect(0,0,4,4))
	temp.Set(0,0, color.RGBA{255, 0, 0, 255})
	temp.Set(0,1, color.RGBA{0, 255, 0, 255})
	temp.Set(0,2, color.RGBA{0, 0, 255, 255})
	temp.Set(0,3, color.RGBA{255, 0, 0, 255})

	temp.Set(1,0, color.RGBA{0, 255, 0, 255})
	temp.Set(1,1, color.RGBA{0, 0, 255, 255})
	temp.Set(1,2, color.RGBA{255, 0, 0, 255})
	temp.Set(1,3, color.RGBA{0, 255, 0, 255})

	temp.Set(2,0, color.RGBA{0, 0, 255, 255})
	temp.Set(2,1, color.RGBA{255, 0, 0, 255})
	temp.Set(2,2, color.RGBA{0, 255, 0, 255})
	temp.Set(2,3, color.RGBA{0, 0, 255, 255})

	temp.Set(3,0, color.RGBA{255, 0, 0, 255})
	temp.Set(3,1, color.RGBA{0, 255, 0, 255})
	temp.Set(3,2, color.RGBA{0, 0, 255, 255})
	temp.Set(3,3, color.RGBA{255, 0, 0, 255})

	png.Encode(f, temp)
}