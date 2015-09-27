package grid

import (
	"fmt"
	"log"
	"runtime"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
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
	Setup()
	Draw(dt float32)
	Update(dt float32)
}

func Run(g *Game) error {
	if g.Scene == nil {
		return fmt.Errorf("Scene property of given Game struct is empty")
	}

	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
	if err := glfw.Init(); err != nil {
		return fmt.Errorf("glfw.Init failed: %s", err)
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
		return fmt.Errorf("glfw.CreateWindow failed: %s", err)
	}
	window.MakeContextCurrent()

	// Initialize Glow
	if err := gl.Init(); err != nil {
		return fmt.Errorf("gl.Init failed:", err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("gl.Init successful, OpenGL version:", version)

	var previousTime, deltaTime, time float32
	previousTime = float32(glfw.GetTime()) - 1.0/60.0

	g.Scene.Setup()

	for !window.ShouldClose() {
		glfw.PollEvents()

		time = float32(glfw.GetTime())
		deltaTime = time - previousTime

		gl.ClearColor(0.2, 0.3, 0.3, 0.5)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		g.Scene.Update(deltaTime)
		g.Scene.Draw(deltaTime)

		previousTime = time

		window.SwapBuffers()
	}

	return nil
}
