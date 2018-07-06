package main

import "github.com/go-gl/mathgl/mgl32"

type Camera struct {
    Pos   mgl32.Vec3 // 摄像机的位置vec3
    Front mgl32.Vec3 // 摄像机向前的向量
    Up    mgl32.Vec3 // 摄像机的上向量

    Yaw   float32 // 偏航角
    Pitch float32 // 俯仰角
    Roll  float32 // 滚轴角
    LastX float32
    LastY float32
    Fov   float32 // Field of view or fov defines how much we can see of the scene

    MoveSpeed float32
    CursorSensitivity float32
}

// NewCamera 创建一个新的摄像机
func NewCamera(pos, front, up mgl32.Vec3, yaw, pitch, roll, lastX, lastY, fov, moveSpeed, cursorSensitivity float32) *Camera {
    return &Camera{
        Pos : pos,
        Front:front,
        Up:up,
        Yaw:yaw,
        Pitch:pitch,
        Roll:roll,
        LastX:lastX,
        LastY:lastY,
        Fov:fov,
        MoveSpeed:moveSpeed,
        CursorSensitivity: cursorSensitivity,
    }
}


