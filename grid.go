package grid

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"log"
	"runtime"
)

const (
	OpenGlVerMajor = 3
	OpenGLVerMinor = 3
)

type Game struct {
	Width      int32
	Height     int32
	Fullscreen bool
	Resizable  bool
	Title      string
	Scene      Scene
}

type Scene interface {
	Draw(dt float32)
	Update(dt float32)
}

func Run(g *Game) {
	if g.Scene == nil {
		log.Fatalln("Game has a empty Scene property ")
	}

	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
	if err := glfw.Init(); err != nil {
		log.Fatalln("glfw.Init failed:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, OpenGLVerMinor)
	glfw.WindowHint(glfw.ContextVersionMinor, OpenGLVerMinor)

	glfw.WindowHint(glfw.Resizable, glfw.False)
	if g.Resizable {
		glfw.WindowHint(glfw.Resizable, glfw.True)
	}

	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(int(g.Width), int(g.Height), g.Title, nil, nil)
	if err != nil {
		log.Fatalln("glfw.CreateWindow failed:", err)
	}
	window.MakeContextCurrent()

	// Initialize Glow
	if err := gl.Init(); err != nil {
		log.Fatalln("gl.Init failed:", err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("gl.Init successful, OpenGL version:", version)

	var previousTime, deltaTime float32

	deltaTime = 1.0 / 60.0

	for !window.ShouldClose() {
		previousTime = float32(glfw.GetTime())

		g.Scene.Update(deltaTime)
		g.Scene.Draw(deltaTime)

		window.SwapBuffers()
		glfw.PollEvents()
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		deltaTime = float32(glfw.GetTime()) - previousTime
	}

}
