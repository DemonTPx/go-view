package view

import (
	"fmt"
	gl "github.com/chsc/gogl/gl21"
	"github.com/veandco/go-sdl2/sdl"
	"path/filepath"
)

const WindowTitle = "Go View"

var (
	DragThreshold   = 5.0
	DragColor       = NewColor(0.4, 0.4, 0.8, 0.5)
	DragBorderWidth = 2.0
	DragBorderColor = NewColor(0.4, 0.4, 0.8, 0.8)
)

type Main struct {
	Window   *sdl.Window
	Renderer *sdl.Renderer
	Context  sdl.GLContext

	Running bool

	Filename   string
	FileCursor *FileCursor

	Settings Settings

	Texture *Texture
	View    View
	Mouse   Mouse
}

type View struct {
	X, Y  float64
	W, H  uint32
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

	m.Settings = LoadSettings()

	_ = sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 3)
	_ = sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 3)
	_ = sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE)

	_ = sdl.GLSetAttribute(sdl.GL_MULTISAMPLESAMPLES, 4)
	_ = sdl.GLSetAttribute(sdl.GL_DOUBLEBUFFER, 1)

	m.Window, err = sdl.CreateWindow(WindowTitle, int32(m.Settings.Window.X), int32(m.Settings.Window.Y), int32(m.Settings.Window.W), int32(m.Settings.Window.H), sdl.WINDOW_OPENGL|sdl.WINDOW_RESIZABLE)
	defer m.Window.Destroy()
	if err != nil {
		return err
	}

	m.Renderer, err = sdl.CreateRenderer(m.Window, -1, 0)
	if err != nil {
		return err
	}
	defer m.Renderer.Destroy()

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
		m.FileCursor, err = NewFileCursorFromFilename(m.Filename)
		if err != nil {
			return err
		}
	} else {
		m.FileCursor, err = NewFileCursorFromWorkingDirectory()
		if err != nil {
			return err
		}
	}

	err = m.LoadFile()
	if err != nil {
		return err
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

		if m.Mouse.dragging {
			rect := m.Mouse.DragRect()
			if rect.W >= DragThreshold || rect.H >= DragThreshold {
				DrawQuadBorder(rect, DragColor, DragBorderWidth, DragBorderColor)
			}
		}

		m.Window.GLSwap()

		commandHandler.HandleBlocking()
	}

	m.SaveSettings()

	return nil
}

func (m *Main) InitGL() error {
	err := gl.Init()
	if err != nil {
		return fmt.Errorf("failed to initialize opengl")
	}

	gl.ClearColor(0.2, 0.2, 0.2, 1.0)

	gl.Enable(gl.TEXTURE_2D)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	m.ResetGLView(m.Settings.Window.W, m.Settings.Window.H)

	return nil
}

func (m *Main) ResetGLView(w, h uint32) {
	gl.Viewport(0, 0, gl.Sizei(w), gl.Sizei(h))

	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()

	gl.Ortho(gl.Double(0), gl.Double(w), gl.Double(h), gl.Double(0), gl.Double(-1.0), gl.Double(1.0))

	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()

	m.View.W = w
	m.View.H = h
	m.View.X = float64(w) / 2
	m.View.Y = float64(h) / 2
}

func (m *Main) SaveSettings() {
	x, y := m.Window.GetPosition()
	w, h := m.Window.GetSize()
	settings := Settings{
		Window: WindowSettings{
			X: uint32(x),
			Y: uint32(y),
			W: uint32(w),
			H: uint32(h),
		},
	}
	m.Settings = settings
	SaveSettings(m.Settings)
}

func (m *Main) FitToWindow() {
	if m.Texture == nil {
		return
	}

	if m.Texture.W > int32(m.View.W) || m.Texture.H > int32(m.View.H) {
		windowRatio := float64(m.View.W) / float64(m.View.H)
		textureRatio := float64(m.Texture.W) / float64(m.Texture.H)

		if windowRatio > textureRatio {
			m.View.Scale = float64(m.View.H) / float64(m.Texture.H)
		} else {
			m.View.Scale = float64(m.View.W) / float64(m.Texture.W)
		}
	}
}

func (m *Main) LoadFile() error {
	var err error

	m.Filename = m.FileCursor.GetFilename()

	if len(m.Filename) == 0 {
		return nil
	}
	fmt.Printf("loading file %s\n", m.Filename)

	if m.Texture != nil {
		m.Texture.Destroy()
	}

	m.Texture, err = NewTextureFromFile(m.Filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %s", err)
	}

	m.Window.SetTitle(fmt.Sprintf("%s - %dx%d", filepath.Base(m.Filename), m.Texture.W, m.Texture.H))

	m.View.X = float64(m.View.W) / 2
	m.View.Y = float64(m.View.H) / 2
	m.View.Scale = 1

	m.FitToWindow()

	return nil
}
