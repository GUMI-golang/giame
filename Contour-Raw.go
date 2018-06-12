package giame

import "github.com/go-gl/mathgl/mgl32"

type ContourRaw struct {
	contourInfo
	Points []mgl32.Vec2
}
