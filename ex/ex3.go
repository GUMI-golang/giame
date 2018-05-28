package main

import (
	"os"
	"image/png"
	"github.com/GUMI-golang/giame/tools/mask"
	"github.com/GUMI-golang/giame/tools"
	"github.com/GUMI-golang/gumi/gcore"
	"image"
	"golang.org/x/image/draw"
)

func main() {

	f, err := os.Open("./example/jellybeans.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	img, err := png.Decode(f)
	if err != nil {
		panic(err)
	}
	dst := image.NewRGBA(img.Bounds())
	for _, s := range []mask.Mask3{
		mask.Laplacian,
		mask.LaplacianExtend,
	}{
		draw.Draw(dst, dst.Rect, img, img.Bounds().Min, draw.Src)
		tools.Filting(s, dst)
		gcore.Capture("_out_" + s.Strings(), dst)
	}
}
