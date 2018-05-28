package giamegl43

import (
	"github.com/GUMI-golang/giame"
	"github.com/GUMI-golang/giame/giamegl/shaders"
	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/pkg/errors"
)

type GLDriver struct {
	pathLineProgram uint32
	pathFillProgram uint32
	//
	pipeDownProgram   uint32
	pipeFilt33Program uint32
	pipeFilt55Program uint32
	//
	raster [giame.FillerTypeLength]uint32
	//
	utilMixing   uint32
	//
	isInit bool
}


func NewDriver() *GLDriver {
	return new(GLDriver)
}
func (GLDriver) MakeResult(w, h int) giame.Result {
	var temp uint32
	gl.GenTextures(1, &temp)
	gl.BindTexture(gl.TEXTURE_2D, temp)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA32F, int32(w), int32(h), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(nil))
	return GLResult(temp)
}
func (s *GLDriver) Init() (err error) {
	if s.isInit {
		return errors.New("Already init")
	}
	// path init
	s.pathLineProgram, err = compile(string(shaders.MustAsset("Path0_Line.cs.glsl")))
	defer func() {
		if err != nil && s.pathLineProgram != 0{
			gl.DeleteProgram(s.pathLineProgram)
		}
	}()
	if err != nil {
		return err
	}
	s.pathFillProgram, err = compile(string(shaders.MustAsset("Path1_Fill.cs.glsl")))
	defer func() {
		if err != nil && s.pathFillProgram != 0{
			gl.DeleteProgram(s.pathFillProgram)
		}
	}()
	if err != nil {
		return err
	}
	// pipe init
	s.pipeDownProgram, err = compile(string(shaders.MustAsset("Pipe0_Downscale.cs.glsl")))
	defer func() {
		if err != nil && s.pipeDownProgram != 0{
			gl.DeleteProgram(s.pipeDownProgram)
		}
	}()
	if err != nil {
		return err
	}
	s.pipeFilt33Program, err = compile(string(shaders.MustAsset("Pipe1_Filter3x3.cs.glsl")))
	defer func() {
		if err != nil && s.pipeFilt33Program != 0{
			gl.DeleteProgram(s.pipeFilt33Program)
		}
	}()
	if err != nil {
		return err
	}
	s.pipeFilt55Program, err = compile(string(shaders.MustAsset("Pipe1_Filter5x5.cs.glsl")))
	defer func() {
		if err != nil && s.pipeFilt55Program != 0{
			gl.DeleteProgram(s.pipeFilt55Program)
		}
	}()
	if err != nil {
		return err
	}
	// raster init
	s.raster[giame.FillerTypeUniform], err = compile(string(shaders.MustAsset("Raster_Uniform.cs.glsl")))
	defer func() {
		if err != nil && s.raster[giame.FillerTypeUniform] != 0{
			gl.DeleteProgram(s.raster[giame.FillerTypeUniform])
		}
	}()
	if err != nil {
		return err
	}

	s.raster[giame.FillerTypeFixed], err = compile(string(shaders.MustAsset("Raster_Fixed.cs.glsl")))
	defer func() {
		if err != nil && s.raster[giame.FillerTypeFixed] != 0{
			gl.DeleteProgram(s.raster[giame.FillerTypeFixed])
		}
	}()
	if err != nil {
		return err
	}

	// Kernel Using
	s.raster[giame.FillerTypeNearestNeighbor], err = compile(string(shaders.MustAsset("Raster_Kernel_NearestNeighbor.cs.glsl")))
	defer func() {
		if err != nil && s.raster[giame.FillerTypeNearestNeighbor] != 0{
			gl.DeleteProgram(s.raster[giame.FillerTypeNearestNeighbor])
		}
	}()
	if err != nil {
		return err
	}
	s.raster[giame.FillerTypeBilinear], err = compile(string(shaders.MustAsset("Raster_Kernel_Bilinear.cs.glsl")))
	defer func() {
		if err != nil && s.raster[giame.FillerTypeBilinear] != 0{
			gl.DeleteProgram(s.raster[giame.FillerTypeBilinear])
		}
	}()
	if err != nil {
		return err
	}
	s.raster[giame.FillerTypeBell], err = compile(string(shaders.MustAsset("Raster_Kernel_Bell.cs.glsl")))
	defer func() {
		if err != nil && s.raster[giame.FillerTypeBell] != 0{
			gl.DeleteProgram(s.raster[giame.FillerTypeBell])
		}
	}()
	if err != nil {
		return err
	}
	s.raster[giame.FillerTypeHermite], err = compile(string(shaders.MustAsset("Raster_Kernel_Hermite.cs.glsl")))
	defer func() {
		if err != nil && s.raster[giame.FillerTypeHermite] != 0{
			gl.DeleteProgram(s.raster[giame.FillerTypeHermite])
		}
	}()
	if err != nil {
		return err
	}
	s.raster[giame.FillerTypeBicubicHalf], err = compile(string(shaders.MustAsset("Raster_Kernel_BicubicHalf.cs.glsl")))
	defer func() {
		if err != nil && s.raster[giame.FillerTypeBicubicHalf] != 0{
			gl.DeleteProgram(s.raster[giame.FillerTypeBicubicHalf])
		}
	}()
	if err != nil {
		return err
	}
	s.raster[giame.FillerTypeMitchellOneThird], err = compile(string(shaders.MustAsset("Raster_Kernel_MitchellOneThird.cs.glsl")))
	defer func() {
		if err != nil && s.raster[giame.FillerTypeMitchellOneThird] != 0{
			gl.DeleteProgram(s.raster[giame.FillerTypeMitchellOneThird])
		}
	}()
	if err != nil {
		return err
	}
	s.raster[giame.FillerTypeLanczos2], err = compile(string(shaders.MustAsset("Raster_Kernel_Lanczos2.cs.glsl")))
	defer func() {
		if err != nil && s.raster[giame.FillerTypeLanczos2] != 0{
			gl.DeleteProgram(s.raster[giame.FillerTypeLanczos2])
		}
	}()
	if err != nil {
		return err
	}
	s.raster[giame.FillerTypeLanczos3], err = compile(string(shaders.MustAsset("Raster_Kernel_Lanczos3.cs.glsl")))
	defer func() {
		if err != nil && s.raster[giame.FillerTypeLanczos3] != 0{
			gl.DeleteProgram(s.raster[giame.FillerTypeLanczos3])
		}
	}()
	if err != nil {
		return err
	}

	s.raster[giame.FillerTypeRepeat], err = compile(string(shaders.MustAsset("Raster_Repeat.cs.glsl")))
	defer func() {
		if err != nil && s.raster[giame.FillerTypeRepeat] != 0{
			gl.DeleteProgram(s.raster[giame.FillerTypeRepeat])
		}
	}()
	if err != nil {
		return err
	}
	s.raster[giame.FillerTypeRepeatHorizontal] = s.raster[giame.FillerTypeRepeat]
	s.raster[giame.FillerTypeRepeatVertical] = s.raster[giame.FillerTypeRepeat]
	// util init

	s.utilMixing, err = compile(string(shaders.MustAsset("util-Mixing.cs.glsl")))
	defer func() {
		if err != nil && s.utilMixing != 0{
			gl.DeleteProgram(s.utilMixing)
		}
	}()
	if err != nil {
		return err
	}


	// Set flag isInit
	s.isInit = true
	return nil
}
func (s *GLDriver) Close() error {
	if !s.isInit {
		return errors.New("Not Inited")
	}
	gl.DeleteProgram(s.pathLineProgram)
	gl.DeleteProgram(s.pathFillProgram)
	for j := giame.FillerType(0); j < giame.FillerTypeLength; j++ {
		gl.DeleteProgram(s.raster[j])
	}
	return nil
}
