package main

import (
    "testing"
    "github.com/go-gl/mathgl/mgl32"
)

func TestNewCamera(t *testing.T) {
    camera := NewCamera(
        mgl32.Vec3{0, 0, 3},
        mgl32.Vec3{0, 0, -1},
        mgl32.Vec3{0, 1, 0},
    )
}
