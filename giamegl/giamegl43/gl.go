package giamegl43

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/gl/v4.3-core/gl"
	"strings"
	"fmt"
	"image"
	"math"
	"github.com/GUMI-golang/giame/giamegl/shaders"
)

func compile(src string) (uint32, error){
	prog := gl.CreateProgram()
	shader := gl.CreateShader(gl.COMPUTE_SHADER)
	csources, free := gl.Strs(src + "\x00")
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	// check compile success
	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))
		gl.DeleteProgram(prog)
		return 0, fmt.Errorf("failed to compile %v: %v", src, log)
	}
	gl.AttachShader(prog, shader)
	gl.DeleteShader(shader)
	gl.LinkProgram(prog)
	gl.GetProgramiv(prog, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(prog, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(prog, logLength, nil, gl.Str(log))
		gl.DeleteProgram(prog)
		return 0, fmt.Errorf("failed to link program: %v", log)
	}
	return prog, nil
}
func pathVisualize(workspace uint32, w, h int) *image.RGBA{
	mem := make([]int32, w * h)
	gl.BindTexture(gl.TEXTURE_2D, workspace)
	gl.GetTexImage(gl.TEXTURE_2D, 0, gl.RED_INTEGER, gl.INT, gl.Ptr(mem))
	rgba := image.NewRGBA(image.Rect(0,0,w,h))

	for i, value := range mem {
		x, y := i % w, i / w
		offset := rgba.PixOffset(x, y)
		if value > 0{
			if value > math.MaxUint16{
				value = math.MaxUint16
			}
			rgba.Pix[offset + 0] = uint8(value >> 8)
			rgba.Pix[offset + 3] = uint8(value >> 8)
 		}else {
			value = -value
			if value > math.MaxUint16{
				value = math.MaxUint16
			}
			rgba.Pix[offset + 1] = uint8(value >> 8)
			rgba.Pix[offset + 3] = uint8(value >> 8)
		}
	}
	return rgba
}
func rasterUniform(res GLResult, workspace uint32, w, h uint32, min [2]int32, c mgl32.Vec4) {
	// setup buffer
	// [0] : min ivec2
	// [1] : color vec4
	var bufPoint [2]uint32
	gl.GenBuffers(2, &bufPoint[0])
	defer gl.DeleteBuffers(2, &bufPoint[0])
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, bufPoint[0])
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, 4*2, gl.Ptr(&min[0]), gl.STATIC_READ)
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, bufPoint[1])
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, 4*4, gl.Ptr(&c[0]), gl.STATIC_READ)
	gl.MemoryBarrier(gl.BUFFER_UPDATE_BARRIER_BIT)
	//
	gl.BindImageTexture(0, uint32(res), 0, false, 0, gl.READ_WRITE, gl.RGBA32F)
	gl.BindImageTexture(1, workspace, 0, false, 0, gl.READ_ONLY, gl.R32I)
	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 2, bufPoint[0])
	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 3, bufPoint[1])
	gl.DispatchCompute(w, h, 1)

	gl.MemoryBarrier(gl.SHADER_STORAGE_BARRIER_BIT | gl.SHADER_IMAGE_ACCESS_BARRIER_BIT)
}
func rasterFixed(res GLResult, workspace, filler uint32, fillbound uint32, w, h uint32, min [2]int32) {
	// setup buffer
	// start point
	var bufPoint uint32
	gl.GenBuffers(1, &bufPoint)
	defer gl.DeleteBuffers(1, &bufPoint)
	//
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, bufPoint)
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, 4*2, gl.Ptr(&min[0]), gl.DYNAMIC_COPY)
	gl.MemoryBarrier(gl.BUFFER_UPDATE_BARRIER_BIT)
	//
	gl.BindImageTexture(0, uint32(res), 0, false, 0, gl.READ_WRITE, gl.RGBA32F)
	gl.BindImageTexture(1, workspace, 0, false, 0, gl.READ_ONLY, gl.R32I)
	gl.BindImageTexture(2, filler, 0, false, 0, gl.READ_ONLY, gl.RGBA32F)
	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 3, bufPoint)
	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 4, fillbound)
	gl.DispatchCompute(w, h, 1)
	gl.MemoryBarrier(gl.SHADER_STORAGE_BARRIER_BIT | gl.SHADER_IMAGE_ACCESS_BARRIER_BIT)
}
func rasterKernel(res GLResult, fillerw, fillerh int, workspace, filler uint32, fillbound uint32, w, h uint32, min [2]int32) {
	// setup buffer
	// start point
	var bufPoint uint32
	gl.GenBuffers(1, &bufPoint)
	defer gl.DeleteBuffers(1, &bufPoint)
	//

	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, fillbound)
	var temp [4]int32
	gl.GetBufferSubData(gl.SHADER_STORAGE_BUFFER, 0, 4 * 4, gl.Ptr(&temp[0]))
	data := shaders.StartDelta{
		Start: min,
		Delta:mgl32.Vec2{
			float32(fillerw)/float32(temp[2] - temp[0] + 1),
			float32(fillerh)/float32(temp[3] - temp[1] + 1),
		},
	}
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, bufPoint)
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, data.Size(), data.Pointer(), gl.DYNAMIC_COPY)
	gl.MemoryBarrier(gl.BUFFER_UPDATE_BARRIER_BIT)
	//
	gl.BindImageTexture(0, uint32(res), 0, false, 0, gl.READ_WRITE, gl.RGBA32F)
	gl.BindImageTexture(1, workspace, 0, false, 0, gl.READ_ONLY, gl.R32I)
	gl.BindImageTexture(2, filler, 0, false, 0, gl.READ_ONLY, gl.RGBA32F)
	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 3, bufPoint)
	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 4, fillbound)
	gl.DispatchCompute(w, h, 1)
	gl.MemoryBarrier(gl.SHADER_STORAGE_BARRIER_BIT | gl.SHADER_IMAGE_ACCESS_BARRIER_BIT)
}
func rasterRepeat(res GLResult, workspace, filler uint32, fillbound uint32, w, h uint32, min [2]int32, repeator [2]int32) {
	// setup buffer
	// start point
	var bufPoint uint32
	gl.GenBuffers(1, &bufPoint)
	defer gl.DeleteBuffers(1, &bufPoint)
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, bufPoint)
	//gl.BufferData(gl.SHADER_STORAGE_BUFFER, data.Size(), gl.Ptr(data.Data()), gl.DYNAMIC_COPY)
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, 4 * 4, gl.Ptr([]int32{
		min[0], min[1], repeator[0], repeator[1],
	}), gl.DYNAMIC_COPY)
	gl.MemoryBarrier(gl.BUFFER_UPDATE_BARRIER_BIT)
	//
	gl.BindImageTexture(0, uint32(res), 0, false, 0, gl.READ_WRITE, gl.RGBA32F)
	gl.BindImageTexture(1, workspace, 0, false, 0, gl.READ_ONLY, gl.R32I)
	gl.BindImageTexture(2, filler, 0, false, 0, gl.READ_ONLY, gl.RGBA32F)
	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 3, bufPoint)
	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 4, fillbound)
	gl.DispatchCompute(w, h, 1)
	gl.MemoryBarrier(gl.SHADER_STORAGE_BARRIER_BIT | gl.SHADER_IMAGE_ACCESS_BARRIER_BIT)
}