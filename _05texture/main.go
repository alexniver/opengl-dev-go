package main

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const windowWidth = 800
const windowHeight = 600

var rate = float32(0.5)

func init() {
	runtime.LockOSThread()
}

func main() {
	window := initGlfw()
	defer glfw.Terminate()

	initOpenGL()

	window.SetKeyCallback(keyCallback)

	shaderProgram := gl.CreateProgram()

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
	gl.AttachShader(shaderProgram, verticsShaderHandle)
	gl.AttachShader(shaderProgram, fragmentShaderHandle)
	gl.LinkProgram(shaderProgram)

	// vertices and indices
	vertices := []float32{
		-0.5, 0.5, 0.0, // pos
		1.0, 0.0, 0.0, // color
		0, 1, // uv

		0.5, 0.5, 0.0,
		0.0, 1.0, 0.0,
		1, 1,

		0.5, -0.5, 0.0,
		0.0, 0.0, 1.0,
		1, 0,

		-0.5, -0.5, 0.0,
		0.5, 0.5, 0.0,
		0, 0,
	}

	indices := []uint32{
		0, 1, 3,
		1, 2, 3,
	}

	VBO := makeVbo(vertices)
	makeVao(VBO)
	// pos
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 8*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	// color
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 8*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	// uv
	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, 8*4, gl.PtrOffset(6*4))
	gl.EnableVertexAttribArray(2)

	makeEbo(indices)

	texture0, err := newTexture("texture/funny.jpg")
	texture1, err := newTexture("texture/wall.jpeg")

	gl.UseProgram(shaderProgram)

	gl.Uniform1i(gl.GetUniformLocation(shaderProgram, gl.Str("texture0"+"\x00")), 0)
	gl.Uniform1i(gl.GetUniformLocation(shaderProgram, gl.Str("texture1"+"\x00")), 1)

	if nil != err {
		log.Fatal(err)
	}

	for !window.ShouldClose() {
		gl.ClearColor(0.5, 0.5, 1, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// set rate
		gl.Uniform1f(gl.GetUniformLocation(shaderProgram, gl.Str("rate"+"\x00")), rate)

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture0)

		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, texture1)

		gl.UseProgram(shaderProgram)
		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))

		glfw.PollEvents()
		window.SwapBuffers()
	}

}

func initGlfw() *glfw.Window {
	// init glfw
	if err := glfw.Init(); nil != err {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.True)
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

func readShaderFromFile(file string, sType uint32) (uint32, error) {
	// read and compile
	src, err := ioutil.ReadFile(file)
	if nil != err {
		return 0, err
	}

	return compileShader(string(src)+"\x00", sType)
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

func newTexture(filePath string) (uint32, error) {
	imgFile, err := os.Open(filePath)
	if nil != err {
		return 0, fmt.Errorf("texture %q not found on disk: %v", filePath, err)
	}

	img, _, err := image.Decode(imgFile)
	img = imaging.FlipV(img) // flip the image
	if nil != err {
		return 0, fmt.Errorf("texture %q decode error: %v", filePath, err)
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return 0, fmt.Errorf("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.MIRRORED_REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix),
	)
	return texture, nil
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

	if action == glfw.Press {
		switch key {
		case glfw.KeyUp:
			rate = rate + 0.1
		case glfw.KeyDown:
			rate = rate - 0.1
		}
	}
}
