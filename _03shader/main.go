package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"runtime"
	"strings"
	"unsafe"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const windowWidth = 800
const windowHeight = 600

func init() {
	runtime.LockOSThread()
}

func main() {
	if err := glfw.Init(); nil != err {
		log.Fatal("failed to inifitialize glfw:", err)
	}

	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(windowWidth, windowWidth, "basic shaders", nil, nil)
	if nil != err {
		log.Panic(err)
	}

	window.MakeContextCurrent()

	if err := gl.Init(); nil != err {
		log.Panic(err)
	}

	window.SetKeyCallback(keyCallback)

	// read shader from files
	verticsShaderHandle, err := readShaderFromFile("./shaders/vertices.vert", gl.VERTEX_SHADER)
	if nil != err {
		log.Panic(err)
	}
	fragmentShaderHandle, err := readShaderFromFile("./shaders/fragment.frag", gl.FRAGMENT_SHADER)
	if nil != err {
		log.Panic(err)
	}

	// gen program
	shaderProgram := gl.CreateProgram()
	gl.AttachShader(shaderProgram, verticsShaderHandle)
	gl.AttachShader(shaderProgram, fragmentShaderHandle)
	gl.LinkProgram(shaderProgram)

	// vertices and indices
	vertices := []float32{
		0.0, 0.5, 0.0, // top
		1.0, 0.0, 0.0, // color red

		0.5, -0.5, 0.0, // right
		0.0, 1.0, 0.0, // color green

		-0.5, -0.5, 0.0, // left
		0.0, 0.0, 1.0, // color blue
	}

	indices := []float32{
		0, 1, 2,
	}

	// set vbo vao ebo
	var VAO uint32
	var VBO uint32
	var EBO uint32

	gl.GenVertexArrays(1, &VAO)
	gl.GenBuffers(1, &VBO)
	gl.GenBuffers(1, &EBO)

	gl.BindVertexArray(VAO)

	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, EBO)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	// position
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	// color
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	for !window.ShouldClose() {

		gl.ClearColor(0.2, 0.5, 0.5, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		gl.UseProgram(shaderProgram)
		gl.BindVertexArray(VAO)
		gl.DrawElements(gl.TRIANGLES, 3, gl.UNSIGNED_INT, unsafe.Pointer(nil))

		window.SwapBuffers()
		glfw.PollEvents()
	}

}

func readShaderFromFile(file string, sType uint32) (uint32, error) {
	// read and compile
	src, err := ioutil.ReadFile(file)
	if nil != err {
		return 0, err
	}

	fmt.Printf("string(src) = %+v\n", string(src))
	shader := gl.CreateShader(sType)
	csources, free := gl.Strs(string(src) + "\x00")
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	// check error
	var success int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &success)
	if success == gl.FALSE {
		var logLen int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLen)
		log := gl.Str(strings.Repeat("\x00", int(logLen)))
		gl.GetShaderInfoLog(shader, logLen, nil, log)

		return 0, fmt.Errorf("%s: %s", "ReadShaderFromFile error", gl.GoStr(log))
	}

	return shader, nil
}

func keyCallback(
	window *glfw.Window,
	key glfw.Key,
	scancode int,
	action glfw.Action,
	mods glfw.ModifierKey,
) {
	if key == glfw.KeyEscape && action == glfw.Press {
		window.SetShouldClose(true)
	}
}
