package main

import (
    "github.com/go-gl/glfw/v3.2/glfw"
    "fmt"
    "github.com/go-gl/mathgl/mgl64"
)

var keyPressedMap = make(map[glfw.Key]bool) // 键盘按下map
var cursorFirst = true              // 光标是否是第一次进入屏幕，默认是true
var cursorPos mgl64.Vec2            // 光标位置
var cursorPosLast mgl64.Vec2        // 光标上次监听事件的位置
var bufferedCursorChange mgl64.Vec2 // 光标累计变化总量

func keyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
    // timing for key events occurs differently from what the program loop requires
    // so just track what key actions occur and then access them in the program loop
    switch action {
    case glfw.Press:
        keyPressedMap[key] = true
    case glfw.Release:
        keyPressedMap[key] = false
    }
    fmt.Println("key pressed/released")
}

func cursorPosCallback(window *glfw.Window, xpos, ypos float64) {
	if cursorFirst {
		cursorPosLast[0] = xpos
		cursorPosLast[1] = ypos
		cursorFirst = false
	}

	bufferedCursorChange[0] += xpos - cursorPosLast[0]
	bufferedCursorChange[1] += ypos - cursorPosLast[1]

	cursorPosLast[0] = xpos
	cursorPosLast[1] = ypos
	fmt.Println("mouse moved")
}

func scrollCallback(w *glfw.Window, xoff float64, yoff float64) {

}

// IsKeyListPressed 返回是否按下了keyList中的所有键
func IsKeyListPressed(keyList ...glfw.Key) bool {
    for _, key := range keyList {
        if !keyPressedMap[key] {
            return false
        }
    }
    return true
}
