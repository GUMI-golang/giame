package main

import (
	"github.com/GUMI-golang/giame"
	"image"
	"github.com/go-gl/glfw/v3.2/glfw"
	"runtime"
	"github.com/GUMI-golang/gumi/gcore"
	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/GUMI-golang/giame/giamegl/giamegl43"
	"os"
	"fmt"
	"strings"
)

func main()  {
	runtime.LockOSThread()
	gcore.Must(glfw.Init())
	fmt.Println("GLFW Init Complete")
	//
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.Visible, glfw.False)
	wnd := gcore.MustValue(glfw.CreateWindow(10, 10, "", nil,  nil)).(*glfw.Window)
	wnd.MakeContextCurrent()
	defer wnd.Destroy()
	//
	gcore.Must(gl.Init())
	fmt.Println("GL Init Complete")
	//
	dr := giamegl43.NewDriver()
	gcore.Must(dr.Init())
	fmt.Println("Driver Init Complete")
	defer dr.Close()
	res := dr.MakeResult(128, 128)
	//
	//res.Request(dr, Color())
	res.Request(dr, Fixed())
	//res.Request(dr, FontX(giame.Scale(4), "GLFW Init"))
	//res.Request(dr, SVGPath(
	//	"M0,0 L128,128 0,128 Z",
	//))
	//
	gcore.Capture("_out", res.Image())
}

func Color() *giame.Contour {
	cb := giame.NewContourQuary(image.Rect(0,0,128, 128))
	cb.MoveTo(mgl32.Vec2{0,0})
	cb.LineTo(mgl32.Vec2{128, 128})
	cb.LineTo(mgl32.Vec2{0,128})
	cb.CloseTo()
	//
	return cb.Fill()
}
func Fixed() *giame.Contour {
	f, err := os.Open("./example/cubes_64.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}
	//
	cb := giame.NewContourQuary(image.Rect(0,0,128, 128))
	cb.MoveTo(mgl32.Vec2{0,0})
	cb.LineTo(mgl32.Vec2{128, 0})
	cb.LineTo(mgl32.Vec2{128, 128})
	cb.LineTo(mgl32.Vec2{0,128})
	cb.CloseTo()
	//
	//cb.SetFiller(giame.NewFiller(giame.FillerTypeFixed, img))
	//cb.SetFiller(giame.NewFiller(giame.FillerTypeNearestNeighbor, img))
	//cb.SetFiller(giame.NewFiller(giame.FillerTypeBilinear, img))
	//cb.SetFiller(giame.NewFiller(giame.FillerTypeRepeat, img))
	cb.SetFiller(giame.NewFiller(giame.FillerTypeRepeatHorizontal, img))
	return cb.Fill()
}
func FontX(scale giame.Scale, txt string) *giame.Contour {
	//
	cq := giame.NewContourQuary(image.Rect(0,0,128, 128))
	cq.SetScale(scale)
	giame.DefaultVFont.SetSize(16)
	fmt.Println(giame.DefaultVFont.MeasureText(cq, txt))
	return giame.DefaultVFont.TextInRect(cq, txt, cq.GetBound(), gcore.AlignRight | gcore.AlignBottom)
}
func SVGPath(script string) *giame.Contour {
	cb := giame.NewContourQuary(image.Rect(0,0,128, 128))
	cb.Query(giame.SVGPath, strings.NewReader(script))
	//
	return cb.Fill()
}