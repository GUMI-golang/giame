package shaders

import (
	"github.com/go-gl/mathgl/mgl32"
	"unsafe"
	"reflect"
)

type StartDelta struct {
	Start [2]int32
	Delta mgl32.Vec2
}
func (StartDelta ) Size() int {
	return 4 * 2 + 4 * 2
}
func (s *StartDelta) Pointer() unsafe.Pointer {
	return unsafe.Pointer(reflect.ValueOf(s).Pointer())
}