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
				sdl.K_PLUS:      ZoomInCommand{},
				sdl.K_KP_PLUS:   ZoomInCommand{},
				sdl.K_EQUALS:    ZoomInCommand{},
				sdl.K_KP_EQUALS: ZoomInCommand{},
				sdl.K_UP:        ZoomInCommand{},
				sdl.K_MINUS:     ZoomOutCommand{},
				sdl.K_KP_MINUS:  ZoomOutCommand{},
				sdl.K_DOWN:      ZoomOutCommand{},
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
				sdl.K_w: QuitCommand{},
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
				MouseWheelUp:   ZoomOutCommand{},
				MouseWheelDown: ZoomInCommand{},
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
				k := e.(*sdl.MouseWheelEvent)
				var direction MouseWheel
				if k.X < 0 {
					direction = MouseWheelLeft
				}
				if k.X > 0 {
					direction = MouseWheelRight
				}
				if k.Y < 0 {
					direction = MouseWheelUp
				}
				if k.Y > 0 {
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
