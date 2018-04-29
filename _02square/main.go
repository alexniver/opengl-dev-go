package main

import (
	"fmt"
	"log"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

var vertexShaderSource1 = `
#version 330
in vec3 vp;
void main() {
	gl_Position = vec4(vp, 1.0);
}
` + "\x00"

// this vertexShader makes point pos's x=-x
var vertexShaderSource2 = `
#version 330
in vec3 vp;
void main() {
	gl_Position = vec4(-vp.x, vp.y, vp.z, 1.0);
}
` + "\x00"

var fragmentShaderSource1 = `
#version 330
out vec4 frag_colour;
void main() {
	frag_colour = vec4(0, 1, 1, 1);
}
` + "\x00"

var fragmentShaderSource2 = `
#version 330
out vec4 frag_colour;
void main() {
	frag_colour = vec4(0, 0, 1, 1);
}
` + "\x00"

func init() {
	runtime.LockOSThread()
}

func initGlfw() *glfw.Window {
	// init glfw
	if err := glfw.Init(); nil != err {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(800, 600, "Trinagle", nil, nil)
	if nil != err {
		panic(err)
	}
	window.MakeContextCurrent()

	return window
}

func initOpenGL() {
	// init gl and get a program
	if err := gl.Init(); nil != err {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("Opengl version", version)
}

func makeVbo(vertices []float32) uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)

	return vbo
}

func makeVao(vbo uint32) uint32 {
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	return vao
}

func makeEbo(indices []uint32) uint32 {
	var ebo uint32
	gl.GenBuffers(1, &ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(indices), gl.Ptr(indices), gl.STATIC_DRAW)
	return ebo
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)
	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if gl.FALSE == status {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))
		return 0, fmt.Errorf("failed to compile %v : %v", source, log)
	}

	return shader, nil
}

func main() {
	window := initGlfw()
	defer glfw.Terminate()

	initOpenGL()
	prog1 := gl.CreateProgram()
	prog2 := gl.CreateProgram()

	// attach shader
	vertexShader1, err := compileShader(vertexShaderSource1, gl.VERTEX_SHADER)
	if nil != err {
		log.Panic(err)
	}
	fragmentShader1, err := compileShader(fragmentShaderSource1, gl.FRAGMENT_SHADER)
	if nil != err {
		log.Panic(err)
	}
	vertexShader2, err := compileShader(vertexShaderSource2, gl.VERTEX_SHADER)
	if nil != err {
		log.Panic(err)
	}
	fragmentShader2, err := compileShader(fragmentShaderSource2, gl.FRAGMENT_SHADER)
	if nil != err {
		log.Panic(err)
	}
	gl.AttachShader(prog1, vertexShader1)
	gl.AttachShader(prog1, fragmentShader1)
	gl.LinkProgram(prog1)
	gl.AttachShader(prog2, vertexShader2)
	gl.AttachShader(prog2, fragmentShader2)
	gl.LinkProgram(prog2)

	vertices := []float32{
		-1, 0.5, 0,
		-0.5, 0.5, 0,
		-0.5, -0.5, 0,
		-1, -0.5, 0,
	}

	indices := []uint32{
		0, 1, 3,
		1, 2, 3,
	}

	vbo := makeVbo(vertices)
	vao := makeVao(vbo)
	ebo := makeEbo(indices)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	defer gl.DeleteBuffers(1, &vbo)
	defer gl.DeleteBuffers(1, &vao)
	defer gl.DeleteBuffers(1, &ebo)

	// gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)

	for !window.ShouldClose() {
		gl.ClearColor(0.5, 0.5, 1, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(prog1)

		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))

		gl.UseProgram(prog2)
		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))

		glfw.PollEvents()
		window.SwapBuffers()
	}
}
