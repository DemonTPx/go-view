package view

import (
	"github.com/veandco/go-sdl2/sdl"
	"time"
)

type InputHandler struct {
	commandChannel chan interface{}
	binds          map[KeyMod]map[sdl.Keycode]interface{}
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
		},
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
			case *sdl.KeyboardEvent:
				k := e.(*sdl.KeyboardEvent)
				if k.Type != sdl.KEYDOWN {
					continue
				}

				mod := h.ResolveMod(k.Keysym.Mod)
				modBinds, ok := h.binds[mod]
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

func (h *InputHandler) ResolveMod(mod uint16) KeyMod {
	keyMod := KeyModNone
	if mod&sdl.KMOD_SHIFT != 0 {
		keyMod |= KeyModShift
	}
	if mod&sdl.KMOD_CTRL != 0 {
		keyMod |= KeyModControl
	}
	if mod&sdl.KMOD_ALT != 0 {
		keyMod |= KeyModAlt
	}
	if mod&sdl.KMOD_GUI != 0 {
		keyMod |= KeyModSuper
	}
	return keyMod
}
