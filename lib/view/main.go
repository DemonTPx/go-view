package view

import (
	"fmt"
	gl "github.com/chsc/gogl/gl21"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	windowTitle = "Go View"

	windowW = 1600
	windowH = 900
)

type Main struct {
	Window   *sdl.Window
	Renderer *sdl.Renderer
	Context  sdl.GLContext

	Running bool

	Filename string

	Texture *Texture
	View    View
}

type View struct {
	X     float64
	Y     float64
	Scale float64
}

func NewMain(filename string) *Main {
	return &Main{
		Filename: filename,
		View:     View{Scale: 1},
	}
}

func (m *Main) Run() error {
	var err error

	err = sdl.Init(sdl.INIT_VIDEO | sdl.INIT_EVENTS)
	if err != nil {
		return err
	}
	defer sdl.Quit()

	_ = sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 3)
	_ = sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 3)
	_ = sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE)

	_ = sdl.GLSetAttribute(sdl.GL_MULTISAMPLESAMPLES, 4)
	_ = sdl.GLSetAttribute(sdl.GL_DOUBLEBUFFER, 1)

	m.Window, m.Renderer, err = sdl.CreateWindowAndRenderer(windowW, windowH, sdl.WINDOW_OPENGL|sdl.WINDOW_RESIZABLE)
	if err != nil {
		return err
	}
	defer m.Renderer.Destroy()
	defer m.Window.Destroy()

	m.Window.SetTitle(windowTitle)

	info, err := m.Renderer.GetInfo()
	if err != nil {
		return err
	}

	expectedFlags := uint32(sdl.RENDERER_ACCELERATED | sdl.RENDERER_TARGETTEXTURE)
	if (info.Flags & expectedFlags) != expectedFlags {
		return fmt.Errorf("failed to create opengl context")
	}

	m.Context, err = m.Window.GLCreateContext()
	if err != nil {
		return fmt.Errorf("failed to create opengl context")
	}

	err = m.InitGL()
	if err != nil {
		return err
	}

	if len(m.Filename) != 0 {
		m.Texture, err = NewTextureFromFile(m.Filename)
		if err != nil {
			return fmt.Errorf("failed to open file: %s", err)
		}

		m.View = View{
			X:     float64(windowW) / 2,
			Y:     float64(windowH) / 2,
			Scale: 1,
		}

		m.FitToWindow()
	}

	commandChannel := make(chan interface{}, 10)
	inputHandler := NewInputHandler(commandChannel)
	inputHandler.Run()

	commandHandler := NewCommandHandler(m, commandChannel)

	// Main stuff
	m.Running = true
	for m.Running {
		gl.Clear(gl.COLOR_BUFFER_BIT)

		if m.Texture != nil {
			m.Texture.DrawScale(m.View.X, m.View.Y, m.View.Scale)
		}

		m.Window.GLSwap()

		commandHandler.HandleBlocking()
	}

	return nil
}

func (m *Main) InitGL() error {
	err := gl.Init()
	if err != nil {
		return fmt.Errorf("failed to initialize opengl")
	}

	gl.ClearColor(0.2, 0.2, 0.2, 1.0)
	gl.Viewport(0, 0, gl.Sizei(windowW), gl.Sizei(windowH))

	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()

	gl.Ortho(gl.Double(0), gl.Double(windowW), gl.Double(windowH), gl.Double(0), gl.Double(-1.0), gl.Double(1.0))

	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()

	gl.Enable(gl.TEXTURE_2D)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	return nil
}

func (m *Main) FitToWindow() {
	if m.Texture == nil {
		return
	}

	if m.Texture.W > windowW || m.Texture.H > windowH {
		windowRatio := float64(windowW) / float64(windowH)
		textureRatio := float64(m.Texture.W) / float64(m.Texture.H)

		if windowRatio > textureRatio {
			m.View.Scale = windowH / float64(m.Texture.H)
		} else {
			m.View.Scale = windowW / float64(m.Texture.W)
		}
	}
}
