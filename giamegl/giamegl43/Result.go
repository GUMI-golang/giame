package giamegl43

import (
	"github.com/GUMI-golang/giame"
	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"image"
	"math"
	"github.com/GUMI-golang/giame/tools/mask"
)

type GLResult uint32

func (s GLResult) Clear() {
	w, h := s.Size()
	gl.BindTexture(gl.TEXTURE_2D, uint32(s))
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA32F, int32(w), int32(h), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(nil))
	return
}
func (s GLResult) Size() (w, h int) {
	var wi32, hi32 int32
	gl.BindTexture(gl.TEXTURE_2D, uint32(s))
	gl.GetTexLevelParameteriv(gl.TEXTURE_2D, 0, gl.TEXTURE_WIDTH, &wi32)
	gl.GetTexLevelParameteriv(gl.TEXTURE_2D, 0, gl.TEXTURE_HEIGHT, &hi32)
	return int(wi32), int(hi32)
}
func (s GLResult) Image() image.Image {
	w, h := s.Size()
	rgba := image.NewRGBA(image.Rect(0, 0, w, h))
	var temp = make([]uint8, w*h*4)
	gl.BindTexture(gl.TEXTURE_2D, uint32(s))
	gl.GetTexImage(gl.TEXTURE_2D, 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(temp))
	for y := 0; y < h; y++ {
		copy(rgba.Pix[rgba.Stride*y:rgba.Stride*(y+1)], temp[rgba.Stride*(h-y-1):rgba.Stride*(h-y)])
	}
	return rgba
}
func (s GLResult) Request(d giame.Driver, rq *giame.Contour) {
	dr := d.(*GLDriver)
	wL, hL := rq.WorkspaceSize()
	wR, hR := rq.Bound.Dx(), rq.Bound.Dy()
	//
	// workspaceL init
	var workspace [2]uint32
	var pathResult uint32
	gl.GenTextures(2, &workspace[0])
	defer gl.DeleteTextures(2, &workspace[0])
	for i, v := range workspace {
		gl.BindTexture(gl.TEXTURE_2D, v)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
		if i == 0{
			gl.TexImage2D(gl.TEXTURE_2D, 0, gl.R32I, int32(wL), int32(hL), 0, gl.RED_INTEGER, gl.INT, gl.Ptr(nil))
		}else {
			gl.TexImage2D(gl.TEXTURE_2D, 0, gl.R32I, int32(wR), int32(hR), 0, gl.RED_INTEGER, gl.INT, gl.Ptr(nil))
		}
	}
	gl.MemoryBarrier(gl.TEXTURE_UPDATE_BARRIER_BIT)
	//
	var fillbound uint32
	gl.GenBuffers(1, &fillbound)
	defer gl.DeleteBuffers(1, &fillbound)
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, fillbound)
	temp := [4]int32{math.MaxInt32, math.MaxInt32, math.MinInt32, math.MinInt32}
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, 4*4, gl.Ptr(&temp[0]), gl.DYNAMIC_COPY)
	//=============================================================================
	// * path lining
	gl.UseProgram(dr.pathLineProgram)
	// setup buffer
	var bufPoint uint32
	gl.GenBuffers(1, &bufPoint)
	defer gl.DeleteBuffers(1, &bufPoint)
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, bufPoint)
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, len(rq.Points)*4*2, gl.Ptr(rq.Points), gl.STATIC_READ)
	gl.MemoryBarrier(gl.BUFFER_UPDATE_BARRIER_BIT)
	//
	gl.BindImageTexture(0, workspace[0], 0, false, 0, gl.READ_WRITE, gl.R32I)
	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 1, bufPoint)
	gl.DispatchCompute(uint32(len(rq.Points) - 1), 1, 1)
	gl.MemoryBarrier(gl.SHADER_STORAGE_BARRIER_BIT | gl.SHADER_IMAGE_ACCESS_BARRIER_BIT)

	//gcore.Capture("_out_line", pathVisualize(workspaceL, wL, hL))
	//=============================================================================
	// * path filling
	gl.UseProgram(dr.pathFillProgram)
	gl.BindImageTexture(0, workspace[0], 0, false, 0, gl.READ_WRITE, gl.R32I)
	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 1, fillbound)
	gl.DispatchCompute(uint32(hL), 1, 1)
	gl.MemoryBarrier(gl.SHADER_STORAGE_BARRIER_BIT | gl.SHADER_IMAGE_ACCESS_BARRIER_BIT)
	//gcore.Capture("_out_fill", pathVisualize(workspaceL, wL, hL))
	//=============================================================================

	// * Pipe Downscale
	if rq.Scale == giame.Scale1x{
		pathResult = workspace[0]
	}else {
		// TODO : giame.PipeOptionKernel
		gl.UseProgram(dr.pipeDownProgram)
		gl.BindImageTexture(0, workspace[1], 0, false, 0, gl.READ_WRITE, gl.R32I)
		gl.BindImageTexture(1, workspace[0], 0, false, 0, gl.READ_WRITE, gl.R32I)
		gl.DispatchCompute(uint32(wR), uint32(hR), 1)
		gl.MemoryBarrier(gl.SHADER_IMAGE_ACCESS_BARRIER_BIT)
		pathResult = workspace[1]
	}

	switch rq.PipeOption & giame.PipeOptionFilter {
	case giame.PipeOptionLaplacian:
		var bufPoint uint32
		gl.GenBuffers(1, &bufPoint)
		defer gl.DeleteBuffers(1, &bufPoint)
		gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, bufPoint)
		temp := mask.Laplacian.DataF32()
		gl.BufferData(gl.SHADER_STORAGE_BUFFER, len(temp)*4, gl.Ptr(&temp[0][0]), gl.STATIC_READ)
		gl.MemoryBarrier(gl.BUFFER_UPDATE_BARRIER_BIT)
		//
		gl.UseProgram(dr.pipeFilt33Program)
		gl.BindImageTexture(0, pathResult, 0, false, 0, gl.READ_WRITE, gl.R32I)
		gl.DispatchCompute(uint32(wR), uint32(hR), 1)
		gl.MemoryBarrier(gl.SHADER_IMAGE_ACCESS_BARRIER_BIT)
	case giame.PipeOptionLaplacianExtend:
		var bufPoint uint32
		gl.GenBuffers(1, &bufPoint)
		defer gl.DeleteBuffers(1, &bufPoint)
		gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, bufPoint)
		temp := mask.LaplacianExtend.DataF32()
		gl.BufferData(gl.SHADER_STORAGE_BUFFER, len(temp)*4, gl.Ptr(&temp[0][0]), gl.STATIC_READ)
		gl.MemoryBarrier(gl.BUFFER_UPDATE_BARRIER_BIT)
		//
		gl.UseProgram(dr.pipeFilt33Program)
		gl.BindImageTexture(0, pathResult, 0, false, 0, gl.READ_WRITE, gl.R32I)
		gl.DispatchCompute(uint32(wR), uint32(hR), 1)
		gl.MemoryBarrier(gl.SHADER_IMAGE_ACCESS_BARRIER_BIT)
	}

	//gcore.Capture("_out_down", pathVisualize(workspaceR, wR, hR))
	//=============================================================================
	// * Rasterize
	switch t := rq.Filler.(type) {
	case *giame.UniformFiller:
		gl.UseProgram(dr.raster[int(giame.FillerTypeUniform)])
		rasterUniform(
			s,
			pathResult,
			uint32(wR), uint32(hR),
			[2]int32{int32(rq.Bound.Min.X), int32(rq.Bound.Min.Y)},
			mgl32.Vec4{
				float32(t.Color.R) / math.MaxUint8,
				float32(t.Color.G) / math.MaxUint8,
				float32(t.Color.B) / math.MaxUint8,
				float32(t.Color.A) / math.MaxUint8,
			},
		)
	case *giame.FixedFiller:
		var filler uint32
		gl.GenTextures(1, &filler)
		gl.BindTexture(gl.TEXTURE_2D, filler)
		gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA32F, int32(t.Width), int32(t.Height), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(&t.Data[0]))
		gl.MemoryBarrier(gl.TEXTURE_UPDATE_BARRIER_BIT)
		gl.UseProgram(dr.raster[int(t.FillerType())])
		rasterFixed(s, pathResult, filler, fillbound, uint32(wR), uint32(hR), [2]int32{int32(rq.Bound.Min.X), int32(rq.Bound.Min.Y)},)
	case *giame.KernelFiller:
		var filler uint32
		gl.GenTextures(1, &filler)
		gl.BindTexture(gl.TEXTURE_2D, filler)
		gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA32F, int32(t.Width), int32(t.Height), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(&t.Data[0]))
		gl.MemoryBarrier(gl.TEXTURE_UPDATE_BARRIER_BIT)
		gl.UseProgram(dr.raster[int(t.FillerType())])
		rasterKernel(s, t.Width, t.Height,  pathResult, filler, fillbound, uint32(wR), uint32(hR), [2]int32{int32(rq.Bound.Min.X), int32(rq.Bound.Min.Y)},)
	case *giame.RepeatFiller:
		var filler uint32
		gl.GenTextures(1, &filler)
		gl.BindTexture(gl.TEXTURE_2D, filler)
		gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA32F, int32(t.Width), int32(t.Height), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(&t.Data[0]))
		gl.MemoryBarrier(gl.TEXTURE_UPDATE_BARRIER_BIT)
		gl.UseProgram(dr.raster[int(t.FillerType())])
		switch t.Type {
		case giame.FillerTypeRepeatHorizontal:
			rasterRepeat(s, pathResult, filler, fillbound, uint32(wR), uint32(hR),
				[2]int32{int32(rq.Bound.Min.X), int32(rq.Bound.Min.Y)},
				[2]int32{int32(t.Width), math.MaxInt32},
			)
		case giame.FillerTypeRepeatVertical:
			rasterRepeat(s, pathResult, filler, fillbound, uint32(wR), uint32(hR),
				[2]int32{int32(rq.Bound.Min.X), int32(rq.Bound.Min.Y)},
				[2]int32{math.MaxInt32, int32(t.Height)},
			)
		default:
			rasterRepeat(s, pathResult, filler, fillbound, uint32(wR), uint32(hR),
				[2]int32{int32(rq.Bound.Min.X), int32(rq.Bound.Min.Y)},
				[2]int32{int32(t.Width), int32(t.Height)},
			)
		}
	default:
		panic("Unsupport FillerType")
	}
}
func (s GLResult) GLTex() uint32 {
	return uint32(s)
}
