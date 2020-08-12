package view

import (
	"github.com/veandco/go-sdl2/sdl"
	"time"
)

type InputHandler struct {
	commandChannel chan interface{}
	binds          map[KeyMod]map[sdl.Keycode]interface{}
	modMap         map[uint16]KeyMod
	keyMod         KeyMod
}

type KeyMod uint16

const (
	KeyModNone    KeyMod = 0
	KeyModShift   KeyMod = 1 << 0
	KeyModControl KeyMod = 1 << 1
	KeyModAlt     KeyMod = 1 << 2
	KeyModSuper   KeyMod = 1 << 3
)

func NewInputHandler(commandChannel chan interface{}) *InputHandler {
	return &InputHandler{
		commandChannel: commandChannel,
		binds: map[KeyMod]map[sdl.Keycode]interface{}{
			KeyModNone: {
				sdl.K_ESCAPE:    QuitCommand{},
				sdl.K_PLUS:      ZoomInCommand{},
				sdl.K_KP_PLUS:   ZoomInCommand{},
				sdl.K_EQUALS:    ZoomInCommand{},
				sdl.K_KP_EQUALS: ZoomInCommand{},
				sdl.K_MINUS:     ZoomOutCommand{},
				sdl.K_KP_MINUS:  ZoomOutCommand{},
				sdl.K_1:         ZoomOriginalSizeCommand{},
				sdl.K_f:         ZoomFitToWindowCommand{},
			},
			KeyModControl: {
				sdl.K_w: QuitCommand{},
			},
		},
		modMap: map[uint16]KeyMod{
			sdl.KMOD_SHIFT: KeyModShift,
			sdl.KMOD_CTRL:  KeyModControl,
			sdl.KMOD_ALT:   KeyModAlt,
			sdl.KMOD_GUI:   KeyModSuper,
		},
		keyMod: KeyModNone,
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
				if h.keyMod&KeyModControl == KeyModControl {
					if k.Y < 0 {
						h.commandChannel <- ZoomOutCommand{}
					}
					if k.Y > 0 {
						h.commandChannel <- ZoomInCommand{}
					}
				}
			case *sdl.KeyboardEvent:
				k := e.(*sdl.KeyboardEvent)

				if k.Type == sdl.KEYDOWN {
					for sdlMod, keyMod := range h.modMap {
						if k.Keysym.Mod&sdlMod != 0 {
							h.keyMod |= keyMod
						}
					}
				}
				if k.Type == sdl.KEYUP {
					for sdlMod, keyMod := range h.modMap {
						if k.Keysym.Mod&sdlMod == 0 {
							h.keyMod &^= keyMod
						}
					}
				}

				if k.Type != sdl.KEYDOWN {
					continue
				}

				modBinds, ok := h.binds[h.keyMod]
				if !ok {
					continue
				}
				command, ok := modBinds[k.Keysym.Sym]
				if !ok {
					continue
				}
				h.commandChannel <- command
			}
		}
	}()
}
