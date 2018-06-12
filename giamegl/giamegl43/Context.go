package giamegl43

import (
	"github.com/GUMI-golang/giame"
	"github.com/GUMI-golang/giame/giamegl/shaders"
	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/pkg/errors"
)

var Driver = new(GLDriver)

type GLDriver struct {
	prScanlineHorizontal uint32
	//
	raster [giame.FillerTypeLength]uint32
	//
	utilMixing   uint32
	//
	isInit bool
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
	s.prScanlineHorizontal, err = compile(string(shaders.MustAsset("PR-Scanline-horizontal.glsl")))
	defer func() {
		if err != nil && s.prScanlineHorizontal != 0{
			gl.DeleteProgram(s.prScanlineHorizontal)
		}
	}()
	if err != nil {
		return err
	}
	// raster init
	for k, v := range ass {
		if vstr, ok := v.(string);ok{
			s.raster[k], err = compile(string(shaders.MustAsset(vstr)))
			defer func(to uint32) {
				if err != nil && to != 0{
					gl.DeleteProgram(to)
				}
			}(s.raster[k])
			if err != nil {
				return err
			}
		}else {
			s.raster[k] = s.raster[v.(giame.FillerType)]
		}
	}
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
	gl.DeleteProgram(s.prScanlineHorizontal)
	for k, v := range ass {
		if _, ok := v.(string); ok{
			gl.DeleteProgram(s.raster[k])
		}
	}
	gl.DeleteProgram(s.utilMixing)
	return nil
}


var ass = map[giame.FillerType]interface{}{
	giame.FillerTypeUniform : "RS_Uniform.cs.glsl",
	//
	giame.FillerTypeFixed : "RS_Fixed.cs.glsl",
	//
	giame.FillerTypeNearestNeighbor : "RS_Kernel_NearestNeighbor.cs.glsl",
	giame.FillerTypeBilinear : "RS_Kernel_Bilinear.cs.glsl",
	giame.FillerTypeBell : "RS_Kernel_Bell.cs.glsl",
	giame.FillerTypeBicubicHalf : "RS_Kernel_BicubicHalf.cs.glsl",
	giame.FillerTypeHermite : "RS_Kernel_Hermite.cs.glsl",
	giame.FillerTypeMitchellOneThird : "RS_Kernel_MitchellOneThird.cs.glsl",
	giame.FillerTypeLanczos2 : "RS_Kernel_Lanczos2.cs.glsl",
	giame.FillerTypeLanczos3 : "RS_Kernel_Lanczos3.cs.glsl",
	//
	giame.FillerTypeRepeat: "RS_Repeat.cs.glsl",
	giame.FillerTypeRepeatVertical: giame.FillerTypeRepeat,
	giame.FillerTypeRepeatHorizontal: giame.FillerTypeRepeat,
}