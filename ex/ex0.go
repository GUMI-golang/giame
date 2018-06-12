package main

import (
	"github.com/GUMI-golang/giame/giamegl/giamegl43"
	"github.com/GUMI-golang/gumi/gcore"
	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"runtime"
	"github.com/GUMI-golang/giame"
	"image"
	"github.com/go-gl/mathgl/mgl32"
	"strings"
)

func main() {
	runtime.LockOSThread()
	// glfw, window init
	gcore.Must(glfw.Init())
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.Visible, glfw.False)
	wnd := gcore.MustValue(glfw.CreateWindow(10, 10, "", nil, nil)).(*glfw.Window)
	wnd.MakeContextCurrent()
	defer wnd.Destroy()
	// gl, giamegl init
	gcore.Must(gl.Init())
	gcore.Must(giamegl43.Driver.Init())
	defer giamegl43.Driver.Close()
	res := giamegl43.Driver.MakeResult(128, 128)
	//
	q := giame.NewQuary(image.Rect(32, 32, 96, 96), giame.SetupScanlineHorizontal)
	q.SetFiller(giame.NewUniformFiller(gcore.MustUnmarshalColor("Green")))
	giame.DefaultVFont.SetSize(16)
	c := q.Fill(func(query *giame.FillQuery) {
		//query.MoveTo(mgl32.Vec2{0, 0})
		//query.LineTo(mgl32.Vec2{32, 0})
		//query.LineTo(mgl32.Vec2{32, 32})
		//query.LineTo(mgl32.Vec2{0, 32})
		//query.CloseTo()
		query.Query(giame.SVGPath, strings.NewReader("M0,0 L32,0 L32,32 L0,32 Z"))
		giame.DefaultVFont.Text(query, "Hello", mgl32.Vec2{0,0}, 0)
	})
	res.Request(giamegl43.Driver, c)
	gcore.Capture("_out", res.Image())
}
