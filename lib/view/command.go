package view

import (
	"log"
	"time"
)

type QuitCommand struct{}
type ZoomCommand struct {
	Scale float64
}
type ZoomToMouseCursorCommand struct {
	Scale float64
}
type ZoomOriginalSizeCommand struct{}
type ZoomFitToWindowCommand struct{}
type FirstFileCommand struct{}
type LastFileCommand struct{}
type NextFileCommand struct{}
type PreviousFileCommand struct{}
type UpdateWindowSizeCommand struct {
	W, H uint32
}
type SaveSettingsCommand struct{}
type MouseCursorPositionCommand struct {
	X, Y uint32
}
type MoveViewCommand struct {
	X, Y float64
}

type CommandHandler struct {
	main           *Main
	commandChannel <-chan interface{}

	mouseX, mouseY uint32
}

func NewCommandHandler(main *Main, commandChannel <-chan interface{}) *CommandHandler {
	return &CommandHandler{main: main, commandChannel: commandChannel}
}
func (h *CommandHandler) HandleCommand(command interface{}) {
	switch command.(type) {
	case QuitCommand:
		h.main.Running = false
	case ZoomCommand:
		c := command.(ZoomCommand)
		h.main.View.Scale *= c.Scale
	case ZoomToMouseCursorCommand:
		c := command.(ZoomToMouseCursorCommand)
		if c.Scale < 1 {
			h.main.View.X += (float64(h.main.View.W)/2 - h.main.View.X) * (1 - c.Scale)
			h.main.View.Y += (float64(h.main.View.H)/2 - h.main.View.Y) * (1 - c.Scale)
		} else {
			h.main.View.X += (float64(h.mouseX) - h.main.View.X) * (1 - c.Scale)
			h.main.View.Y += (float64(h.mouseY) - h.main.View.Y) * (1 - c.Scale)
		}
		h.main.View.Scale *= c.Scale
	case ZoomOriginalSizeCommand:
		h.main.View.Scale = 1
	case ZoomFitToWindowCommand:
		h.main.FitToWindow()
	case FirstFileCommand:
		h.main.FileCursor.First()
		h.main.LoadFile()
	case LastFileCommand:
		h.main.FileCursor.Last()
		h.main.LoadFile()
	case NextFileCommand:
		h.main.FileCursor.Next()
		h.main.LoadFile()
	case PreviousFileCommand:
		h.main.FileCursor.Previous()
		h.main.LoadFile()
	case UpdateWindowSizeCommand:
		c := command.(UpdateWindowSizeCommand)
		h.main.ResetGLView(c.W, c.H)
	case SaveSettingsCommand:
		h.main.SaveSettings()
	case MouseCursorPositionCommand:
		c := command.(MouseCursorPositionCommand)
		h.mouseX = c.X
		h.mouseY = c.Y
	case MoveViewCommand:
		c := command.(MoveViewCommand)
		h.main.View.X += c.X
		h.main.View.Y += c.Y
	default:
		log.Printf("unexpected command: %#v", command)
	}
}

func (h *CommandHandler) HandleBlocking() {
	select {
	case command := <-h.commandChannel:
		h.HandleCommand(command)
	}
}

func (h *CommandHandler) HandleTimeout(timeout time.Duration) {
	select {
	case command := <-h.commandChannel:
		log.Printf("received command: %#v", command)
		h.HandleCommand(command)
	case <-time.After(timeout):
		log.Printf("timeout reached")
	}
}
