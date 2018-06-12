package giamegl43

import (
	"github.com/GUMI-golang/giame"
	"github.com/GUMI-golang/giame/giamegl/shaders"
	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"image"
	"math"
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
func (s GLResult) Request(d giame.Driver, rq giame.Contour) {
	temp, ok := d.(*GLDriver)
	if !ok {
		panic("can't convert")
	}
	switch v := rq.(type) {
	case *giame.ScanlineHorizontal:
		s.requestScanlineHorizontal(temp, v)
	default:
		panic("Unsupported Contour")
	}
	//dr := d.(*GLDriver)
	//wL, hL := rq.WorkspaceSize()
	//wR, hR := rq.Bound.Dx(), rq.Bound.Dy()
	////
	//// workspaceL init
	//// [0] : temporal workspace
	//// [1] : merging workspace
	//// [2] : stencil result
	//var workspace [2]uint32
	//var pathResult uint32
	//gl.GenTextures(2, &workspace[0])
	//defer gl.DeleteTextures(2, &workspace[0])
	//for i, v := range workspace {
	//	gl.BindTexture(gl.TEXTURE_2D, v)
	//	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	//	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	//	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	//	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	//	if i == 1{
	//		gl.TexImage2D(gl.TEXTURE_2D, 0, gl.R32I, int32(wR), int32(hR), 0, gl.RED_INTEGER, gl.INT, gl.Ptr(nil))
	//	}else {
	//		gl.TexImage2D(gl.TEXTURE_2D, 0, gl.R32I, int32(wL), int32(hL), 0, gl.RED_INTEGER, gl.INT, gl.Ptr(nil))
	//	}
	//}
	//gl.MemoryBarrier(gl.TEXTURE_UPDATE_BARRIER_BIT)
	////
	//var fillbound uint32
	//gl.GenBuffers(1, &fillbound)
	//defer gl.DeleteBuffers(1, &fillbound)
	//gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, fillbound)
	//temp := [4]int32{math.MaxInt32, math.MaxInt32, math.MinInt32, math.MinInt32}
	//gl.BufferData(gl.SHADER_STORAGE_BUFFER, 4*4, gl.Ptr(&temp[0]), gl.DYNAMIC_COPY)
	////
	////=============================================================================
	//// * path lining
	//gl.UseProgram(dr.pathLineProgram[rq.Option.Line])
	//// setup buffer
	//
	//var bufPoint uint32
	//gl.GenBuffers(1, &bufPoint)
	//defer gl.DeleteBuffers(1, &bufPoint)
	//gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, bufPoint)
	//gl.BufferData(gl.SHADER_STORAGE_BUFFER, 0, gl.Ptr(nil), gl.STATIC_READ)
	//gl.BufferData(gl.SHADER_STORAGE_BUFFER, len(rq.Points)*4*2, gl.Ptr(rq.Points), gl.STATIC_READ)
	//gl.MemoryBarrier(gl.BUFFER_UPDATE_BARRIER_BIT)
	////
	//gl.BindImageTexture(0, workspace[0], 0, false, 0, gl.READ_WRITE, gl.R32I)
	//gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 1, bufPoint)
	//gl.DispatchCompute(uint32(len(rq.Points) - 1), 1, 1)
	//gl.MemoryBarrier(gl.SHADER_STORAGE_BARRIER_BIT | gl.SHADER_IMAGE_ACCESS_BARRIER_BIT)
	//gcore.Capture("_line", pathVisualize(workspace[0]))
	////=============================================================================
	//// * path filling
	//gl.UseProgram(dr.pathFillProgram)
	//gl.BindImageTexture(0, workspace[0], 0, false, 0, gl.READ_WRITE, gl.R32I)
	//gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 2, fillbound)
	//gl.DispatchCompute(uint32(hL), 1, 1)
	//gl.MemoryBarrier(gl.SHADER_STORAGE_BARRIER_BIT | gl.SHADER_IMAGE_ACCESS_BARRIER_BIT)
	//gcore.Capture("_fill",  pathVisualize(workspace[0]))
	////=============================================================================
	//// * Pipe Downscale
	//if rq.AntiAliasing == giame.AA1x {
	//	pathResult = workspace[0]
	//}else {
	//	// TODO : giame.OptionKernel
	//	gl.UseProgram(dr.pipeDownProgram)
	//	gl.BindImageTexture(0, workspace[1], 0, false, 0, gl.READ_WRITE, gl.R32I)
	//	gl.BindImageTexture(1, workspace[0], 0, false, 0, gl.READ_WRITE, gl.R32I)
	//	gl.DispatchCompute(uint32(wR), uint32(hR), 1)
	//	gl.MemoryBarrier(gl.SHADER_IMAGE_ACCESS_BARRIER_BIT)
	//	pathResult = workspace[1]
	//}
	//
	//switch rq.Option.Filter {
	//case giame.FilterNone:
	//case giame.FilterLaplacian:
	//	var bufPoint uint32
	//	gl.GenBuffers(1, &bufPoint)
	//	defer gl.DeleteBuffers(1, &bufPoint)
	//	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, bufPoint)
	//	temp := filter.Laplacian.DataF32()
	//	gl.BufferData(gl.SHADER_STORAGE_BUFFER, len(temp)*4, gl.Ptr(&temp[0][0]), gl.STATIC_READ)
	//	gl.MemoryBarrier(gl.BUFFER_UPDATE_BARRIER_BIT)
	//	//
	//	gl.UseProgram(dr.pipeFilt33Program)
	//	gl.BindImageTexture(0, pathResult, 0, false, 0, gl.READ_WRITE, gl.R32I)
	//	gl.DispatchCompute(uint32(wR), uint32(hR), 1)
	//	gl.MemoryBarrier(gl.SHADER_IMAGE_ACCESS_BARRIER_BIT)
	//case giame.FilterLaplacianExtend:
	//	var bufPoint uint32
	//	gl.GenBuffers(1, &bufPoint)
	//	defer gl.DeleteBuffers(1, &bufPoint)
	//	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, bufPoint)
	//	temp := filter.LaplacianExtend.DataF32()
	//	gl.BufferData(gl.SHADER_STORAGE_BUFFER, len(temp)*4, gl.Ptr(&temp[0][0]), gl.STATIC_READ)
	//	gl.MemoryBarrier(gl.BUFFER_UPDATE_BARRIER_BIT)
	//	//
	//	gl.UseProgram(dr.pipeFilt33Program)
	//	gl.BindImageTexture(0, pathResult, 0, false, 0, gl.READ_WRITE, gl.R32I)
	//	gl.DispatchCompute(uint32(wR), uint32(hR), 1)
	//	gl.MemoryBarrier(gl.SHADER_IMAGE_ACCESS_BARRIER_BIT)
	//}
	////=============================================================================
	//// * Rasterize
	//switch t := rq.Filler.(type) {
	//case *giame.UniformFiller:
	//	gl.UseProgram(dr.raster[int(giame.FillerTypeUniform)])
	//	rasterUniform(
	//		s,
	//		pathResult,
	//		uint32(wR), uint32(hR),
	//		[2]int32{int32(rq.Bound.Min.X), int32(rq.Bound.Min.Y)},
	//		mgl32.Vec4{
	//			float32(t.Color.R) / math.MaxUint8,
	//			float32(t.Color.G) / math.MaxUint8,
	//			float32(t.Color.B) / math.MaxUint8,
	//			float32(t.Color.A) / math.MaxUint8,
	//		},
	//	)
	//case *giame.FixedFiller:
	//	var filler uint32
	//	gl.GenTextures(1, &filler)
	//	gl.BindTexture(gl.TEXTURE_2D, filler)
	//	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA32F, int32(t.Width), int32(t.Height), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(&t.Data[0]))
	//	gl.MemoryBarrier(gl.TEXTURE_UPDATE_BARRIER_BIT)
	//	gl.UseProgram(dr.raster[int(t.FillerType())])
	//	rasterFixed(s, pathResult, filler, fillbound, uint32(wR), uint32(hR), [2]int32{int32(rq.Bound.Min.X), int32(rq.Bound.Min.Y)},)
	//case *giame.KernelFiller:
	//	var filler uint32
	//	gl.GenTextures(1, &filler)
	//	gl.BindTexture(gl.TEXTURE_2D, filler)
	//	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA32F, int32(t.Width), int32(t.Height), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(&t.Data[0]))
	//	gl.MemoryBarrier(gl.TEXTURE_UPDATE_BARRIER_BIT)
	//	gl.UseProgram(dr.raster[int(t.FillerType())])
	//	rasterKernel(s, t.Width, t.Height,  pathResult, filler, fillbound, uint32(wR), uint32(hR), [2]int32{int32(rq.Bound.Min.X), int32(rq.Bound.Min.Y)},)
	//case *giame.RepeatFiller:
	//	var filler uint32
	//	gl.GenTextures(1, &filler)
	//	gl.BindTexture(gl.TEXTURE_2D, filler)
	//	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA32F, int32(t.Width), int32(t.Height), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(&t.Data[0]))
	//	gl.MemoryBarrier(gl.TEXTURE_UPDATE_BARRIER_BIT)
	//	gl.UseProgram(dr.raster[int(t.FillerType())])
	//	switch t.Type {
	//	case giame.FillerTypeRepeatHorizontal:
	//		rasterRepeat(s, pathResult, filler, fillbound, uint32(wR), uint32(hR),
	//			[2]int32{int32(rq.Bound.Min.X), int32(rq.Bound.Min.Y)},
	//			[2]int32{int32(t.Width), math.MaxInt32},
	//		)
	//	case giame.FillerTypeRepeatVertical:
	//		rasterRepeat(s, pathResult, filler, fillbound, uint32(wR), uint32(hR),
	//			[2]int32{int32(rq.Bound.Min.X), int32(rq.Bound.Min.Y)},
	//			[2]int32{math.MaxInt32, int32(t.Height)},
	//		)
	//	default:
	//		rasterRepeat(s, pathResult, filler, fillbound, uint32(wR), uint32(hR),
	//			[2]int32{int32(rq.Bound.Min.X), int32(rq.Bound.Min.Y)},
	//			[2]int32{int32(t.Width), int32(t.Height)},
	//		)
	//	}
	//default:
	//	panic("Unsupport FillerType")
	//}
}

func (s GLResult) requestScanlineHorizontal(d *GLDriver, rq *giame.ScanlineHorizontal) {
	var workspace uint32
	gl.GenTextures(1, &workspace)
	defer gl.DeleteTextures(1, &workspace)
	gl.BindTexture(gl.TEXTURE_2D, workspace)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.R32I, int32(rq.Bound.Dx()), int32(rq.Bound.Dy()), 0, gl.RED_INTEGER, gl.INT, gl.Ptr(nil))

	var buf [2]uint32
	gl.GenBuffers(2, &buf[0])
	defer gl.DeleteBuffers(2, &buf[0])
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, buf[0])
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, 4*len(rq.Value), gl.Ptr(&rq.Value[0]), gl.STATIC_READ)
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, buf[1])
	temp := append([]int32{0}, rq.Range...)
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, 4*len(temp), gl.Ptr(&temp[0]), gl.STATIC_READ)
	gl.MemoryBarrier(gl.BUFFER_UPDATE_BARRIER_BIT | gl.TEXTURE_UPDATE_BARRIER_BIT)

	gl.UseProgram(d.prScanlineHorizontal)
	gl.BindImageTexture(0, workspace, 0, false, 0, gl.READ_WRITE, gl.R32I)
	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 1, buf[0])
	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 2, buf[1])
	gl.DispatchCompute(1, uint32(len(rq.Range)), 1)
	gl.MemoryBarrier(gl.SHADER_IMAGE_ACCESS_BARRIER_BIT)

	//
	switch t := rq.Filler.(type) {
	case *giame.UniformFiller:
		var arg [2]uint32
		gl.GenBuffers(2, &arg[0])
		defer gl.DeleteBuffers(2, &arg[0])
		a0 := [2]int32{int32(rq.Bound.Min.X), int32(rq.Bound.Min.Y)}
		a1 := mgl32.Vec4{float32(t.Color.R) / math.MaxUint8, float32(t.Color.G) / math.MaxUint8, float32(t.Color.B) / math.MaxUint8, float32(t.Color.A) / math.MaxUint8}
		gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, arg[0])
		gl.BufferData(gl.SHADER_STORAGE_BUFFER, 4*2, gl.Ptr(&a0[0]), gl.STATIC_READ)
		gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, arg[1])
		gl.BufferData(gl.SHADER_STORAGE_BUFFER, 4*4, gl.Ptr(&a1[0]), gl.STATIC_READ)
		//
		gl.UseProgram(d.raster[int(giame.FillerTypeUniform)])
		gl.BindImageTexture(0, uint32(s), 0, false, 0, gl.READ_WRITE, gl.RGBA32F)
		gl.BindImageTexture(1, workspace, 0, false, 0, gl.READ_ONLY, gl.R32I)
		gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 2, arg[0])
		gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 3, arg[1])
		gl.DispatchCompute(uint32(rq.Bound.Dx()), uint32(rq.Bound.Dy()), 1)
		gl.MemoryBarrier(gl.SHADER_IMAGE_ACCESS_BARRIER_BIT)
	case *giame.FixedFiller:
		var filler uint32
		gl.GenTextures(1, &filler)
		gl.BindTexture(gl.TEXTURE_2D, filler)
		gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA32F, int32(t.Width), int32(t.Height), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(&t.Data[0]))
		var arg uint32
		gl.GenBuffers(1, &arg)
		defer gl.DeleteBuffers(1, &arg)
		a0 := [2]int32{int32(rq.Bound.Min.X), int32(rq.Bound.Min.Y)}
		gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, arg)
		gl.BufferData(gl.SHADER_STORAGE_BUFFER, 4*2, gl.Ptr(&a0[0]), gl.STATIC_READ)
		//
		gl.UseProgram(d.raster[int(giame.FillerTypeFixed)])
		gl.BindImageTexture(0, uint32(s), 0, false, 0, gl.READ_WRITE, gl.RGBA32F)
		gl.BindImageTexture(1, workspace, 0, false, 0, gl.READ_ONLY, gl.R32I)
		gl.BindImageTexture(2, filler, 0, false, 0, gl.READ_ONLY, gl.RGBA32F)
		gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 3, arg)
		gl.DispatchCompute(uint32(rq.Bound.Dx()), uint32(rq.Bound.Dy()), 1)
		gl.MemoryBarrier(gl.SHADER_IMAGE_ACCESS_BARRIER_BIT)
	case *giame.KernelFiller:
		var filler uint32
		gl.GenTextures(1, &filler)
		gl.BindTexture(gl.TEXTURE_2D, filler)
		gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA32F, int32(t.Width), int32(t.Height), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(&t.Data[0]))
		gl.MemoryBarrier(gl.TEXTURE_UPDATE_BARRIER_BIT)

		var arg uint32
		gl.GenBuffers(1, &arg)
		defer gl.DeleteBuffers(1, &arg)
		argraw := shaders.StartDelta{
			Start: [2]int32{int32(rq.Bound.Min.X), int32(rq.Bound.Min.Y)},
			Delta: mgl32.Vec2{
				float32(rq.Filler.Bounds().Dx()) / float32(rq.Bound.Dx()+1),
				float32(rq.Filler.Bounds().Dy()) / float32(rq.Bound.Dy()+1),
			},
		}
		gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, arg)
		gl.BufferData(gl.SHADER_STORAGE_BUFFER, argraw.Size(), argraw.Pointer(), gl.STATIC_READ)

		//

		gl.UseProgram(d.raster[int(t.FillerType())])

		gl.BindImageTexture(0, uint32(s), 0, false, 0, gl.READ_WRITE, gl.RGBA32F)
		gl.BindImageTexture(1, workspace, 0, false, 0, gl.READ_ONLY, gl.R32I)
		gl.BindImageTexture(2, filler, 0, false, 0, gl.READ_ONLY, gl.RGBA32F)
		gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 3, arg)
		gl.DispatchCompute(uint32(rq.Bound.Dx()), uint32(rq.Bound.Dy()), 1)
		gl.MemoryBarrier(gl.SHADER_IMAGE_ACCESS_BARRIER_BIT)
	case *giame.RepeatFiller:

		var filler uint32
		gl.GenTextures(1, &filler)
		gl.BindTexture(gl.TEXTURE_2D, filler)
		gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA32F, int32(t.Width), int32(t.Height), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(&t.Data[0]))
		gl.MemoryBarrier(gl.TEXTURE_UPDATE_BARRIER_BIT)

		var arg uint32
		gl.GenBuffers(1, &arg)
		defer gl.DeleteBuffers(1, &arg)
		gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, arg)
		gl.BufferData(gl.SHADER_STORAGE_BUFFER, 4*4, gl.Ptr([]int32{ int32(rq.Bound.Min.X), int32(rq.Bound.Min.Y), int32(t.RepeaterWidth), int32(t.RepeaterHeight),}), gl.STATIC_READ)
		//	//
		gl.UseProgram(d.raster[int(t.FillerType())])
		gl.BindImageTexture(0, uint32(s), 0, false, 0, gl.READ_WRITE, gl.RGBA32F)
		gl.BindImageTexture(1, workspace, 0, false, 0, gl.READ_ONLY, gl.R32I)
		gl.BindImageTexture(2, filler, 0, false, 0, gl.READ_ONLY, gl.RGBA32F)
		gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 3, arg)
		gl.DispatchCompute(uint32(rq.Bound.Dx()), uint32(rq.Bound.Dy()), 1)
		gl.MemoryBarrier(gl.SHADER_IMAGE_ACCESS_BARRIER_BIT)
	default:
		panic("Unsupport FillerType")
	}
}

func (s GLResult) GLTex() uint32 {
	return uint32(s)
}

func (s GLResult) path_line_nomod() {

}
