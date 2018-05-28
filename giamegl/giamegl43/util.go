package giamegl43

import "github.com/go-gl/gl/v4.3-core/gl"

// mix `a`, `b` image to `to`
// It is ok `a` equare `to` or `b` equare `to`
func Mixing(dr *GLDriver, to, a, b GLResult) {
	w, h := to.Size()
	gl.UseProgram(dr.utilMixing)
	gl.BindImageTexture(0, uint32(to), 0, false, 0, gl.WRITE_ONLY, gl.RGBA32F)
	gl.BindImageTexture(1, uint32(a), 0, false, 0, gl.READ_ONLY, gl.RGBA32F)
	gl.BindImageTexture(2, uint32(b), 0, false, 0, gl.READ_ONLY, gl.RGBA32F)
	gl.DispatchCompute(uint32(w), uint32(h), 1)
}
