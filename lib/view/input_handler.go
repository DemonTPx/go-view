package view

import (
	"github.com/veandco/go-sdl2/sdl"
	"time"
)

type InputHandler struct {
	commandChannel chan<- interface{}

	keyBinds  map[KeyMod]map[sdl.Keycode]interface{}
	keyModMap map[uint16]KeyMod

	mouseWheelBinds map[KeyMod]map[MouseWheel]interface{}

	currentKeyMod KeyMod
}

type KeyMod uint16
type MouseWheel uint8

const (
	KeyModNone    KeyMod = 0
	KeyModShift   KeyMod = 1 << 0
	KeyModControl KeyMod = 1 << 1
	KeyModAlt     KeyMod = 1 << 2
	KeyModSuper   KeyMod = 1 << 3

	MouseWheelUp    = 1
	MouseWheelDown  = 2
	MouseWheelLeft  = 3
	MouseWheelRight = 4
)

func NewInputHandler(commandChannel chan<- interface{}) *InputHandler {
	return &InputHandler{
		commandChannel: commandChannel,
		keyBinds: map[KeyMod]map[sdl.Keycode]interface{}{
			KeyModNone: {
				sdl.K_ESCAPE:    QuitCommand{},
				sdl.K_PLUS:      ZoomCommand{Scale: 1.25},
				sdl.K_KP_PLUS:   ZoomCommand{Scale: 1.25},
				sdl.K_EQUALS:    ZoomCommand{Scale: 1.25},
				sdl.K_KP_EQUALS: ZoomCommand{Scale: 1.25},
				sdl.K_UP:        ZoomCommand{Scale: 1.25},
				sdl.K_MINUS:     ZoomCommand{Scale: 0.8},
				sdl.K_KP_MINUS:  ZoomCommand{Scale: 0.8},
				sdl.K_DOWN:      ZoomCommand{Scale: 0.8},
				sdl.K_1:         ZoomOriginalSizeCommand{},
				sdl.K_f:         ZoomFitToWindowCommand{},
				sdl.K_PAGEDOWN:  NextFileCommand{},
				sdl.K_RIGHT:     NextFileCommand{},
				sdl.K_PAGEUP:    PreviousFileCommand{},
				sdl.K_LEFT:      PreviousFileCommand{},
				sdl.K_HOME:      FirstFileCommand{},
				sdl.K_END:       LastFileCommand{},
			},
			KeyModControl: {
				sdl.K_w:     QuitCommand{},
				sdl.K_LEFT:  MoveViewCommand{X: -10},
				sdl.K_RIGHT: MoveViewCommand{X: 10},
				sdl.K_UP:    MoveViewCommand{Y: -10},
				sdl.K_DOWN:  MoveViewCommand{Y: 10},
			},
		},
		keyModMap: map[uint16]KeyMod{
			sdl.KMOD_SHIFT: KeyModShift,
			sdl.KMOD_CTRL:  KeyModControl,
			sdl.KMOD_ALT:   KeyModAlt,
			sdl.KMOD_GUI:   KeyModSuper,
		},

		mouseWheelBinds: map[KeyMod]map[MouseWheel]interface{}{
			KeyModNone: {
				MouseWheelUp:   PreviousFileCommand{},
				MouseWheelDown: NextFileCommand{},
			},
			KeyModControl: {
				MouseWheelUp:   ZoomToMouseCursorCommand{Scale: 0.8},
				MouseWheelDown: ZoomToMouseCursorCommand{Scale: 1.25},
			},
		},

		currentKeyMod: KeyModNone,
	}
}

func (h *InputHandler) Run() {

	go func() {
		for {
			e := sdl.PollEvent()

			if e == nil {
				time.Sleep(20 * time.Millisecond)
				continue
			}

			switch e.(type) {
			case *sdl.QuitEvent:
				h.commandChannel <- QuitCommand{}

			case *sdl.MouseWheelEvent:
				m := e.(*sdl.MouseWheelEvent)
				var direction MouseWheel
				if m.X < 0 {
					direction = MouseWheelLeft
				}
				if m.X > 0 {
					direction = MouseWheelRight
				}
				if m.Y < 0 {
					direction = MouseWheelUp
				}
				if m.Y > 0 {
					direction = MouseWheelDown
				}

				modBinds, ok := h.mouseWheelBinds[h.currentKeyMod]
				if !ok {
					continue
				}
				command, ok := modBinds[direction]
				if !ok {
					continue
				}
				h.commandChannel <- command

			case *sdl.MouseMotionEvent:
				m := e.(*sdl.MouseMotionEvent)
				h.commandChannel <- MouseCursorPositionCommand{
					X: uint32(m.X),
					Y: uint32(m.Y),
				}

			case *sdl.MouseButtonEvent:
				m := e.(*sdl.MouseButtonEvent)
				if m.Button == sdl.BUTTON_LEFT {
					if m.State == sdl.PRESSED {
						h.commandChannel <- StartDragCommand{}
					} else {
						h.commandChannel <- StopDragCommand{}
					}
				}

			case *sdl.KeyboardEvent:
				k := e.(*sdl.KeyboardEvent)

				if k.Type == sdl.KEYDOWN {
					for sdlMod, keyMod := range h.keyModMap {
						if k.Keysym.Mod&sdlMod != 0 {
							h.currentKeyMod |= keyMod
						}
					}
				}
				if k.Type == sdl.KEYUP {
					for sdlMod, keyMod := range h.keyModMap {
						if k.Keysym.Mod&sdlMod == 0 {
							h.currentKeyMod &^= keyMod
						}
					}
				}

				if k.Type != sdl.KEYDOWN {
					continue
				}

				modBinds, ok := h.keyBinds[h.currentKeyMod]
				if !ok {
					continue
				}
				command, ok := modBinds[k.Keysym.Sym]
				if !ok {
					continue
				}
				h.commandChannel <- command

			case *sdl.WindowEvent:
				w := e.(*sdl.WindowEvent)
				if w.Event == sdl.WINDOWEVENT_SIZE_CHANGED {
					h.commandChannel <- UpdateWindowSizeCommand{W: uint32(w.Data1), H: uint32(w.Data2)}
					h.commandChannel <- SaveSettingsCommand{}
				}
				if w.Event == sdl.WINDOWEVENT_MOVED {
					h.commandChannel <- SaveSettingsCommand{}
				}
			}
		}
	}()
}
