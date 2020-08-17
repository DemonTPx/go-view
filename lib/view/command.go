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
type StartDragCommand struct{}
type StopDragCommand struct{}

type CommandHandler struct {
	main           *Main
	commandChannel <-chan interface{}
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
			h.main.View.X += (float64(h.main.Mouse.X) - h.main.View.X) * (1 - c.Scale)
			h.main.View.Y += (float64(h.main.Mouse.Y) - h.main.View.Y) * (1 - c.Scale)
		}
		h.main.View.Scale *= c.Scale
	case ZoomOriginalSizeCommand:
		h.main.View.X = float64(h.main.View.W) / 2
		h.main.View.Y = float64(h.main.View.H) / 2
		h.main.View.Scale = 1
	case ZoomFitToWindowCommand:
		h.main.View.X = float64(h.main.View.W) / 2
		h.main.View.Y = float64(h.main.View.H) / 2
		h.main.FitToWindow()
	case FirstFileCommand:
		h.main.FileCursor.First()
		_ = h.main.LoadFile()
	case LastFileCommand:
		h.main.FileCursor.Last()
		_ = h.main.LoadFile()
	case NextFileCommand:
		h.main.FileCursor.Next()
		_ = h.main.LoadFile()
	case PreviousFileCommand:
		h.main.FileCursor.Previous()
		_ = h.main.LoadFile()
	case UpdateWindowSizeCommand:
		c := command.(UpdateWindowSizeCommand)
		h.main.ResetGLView(c.W, c.H)
	case SaveSettingsCommand:
		h.main.SaveSettings()
	case MouseCursorPositionCommand:
		c := command.(MouseCursorPositionCommand)
		h.main.Mouse.X = c.X
		h.main.Mouse.Y = c.Y
	case StartDragCommand:
		h.main.Mouse.dragX = h.main.Mouse.X
		h.main.Mouse.dragY = h.main.Mouse.Y
		h.main.Mouse.dragging = true
	case StopDragCommand:
		h.main.Mouse.dragging = false

		dragRect := h.main.Mouse.DragRect()

		if dragRect.W < DragThreshold || dragRect.H < DragThreshold {
			break
		}

		dragRatio := dragRect.W / dragRect.H
		windowRatio := float64(h.main.View.W) / float64(h.main.View.H)

		var scale float64
		if dragRatio > windowRatio {
			scale = float64(h.main.View.W) / dragRect.W
		} else {
			scale = float64(h.main.View.H) / dragRect.H
		}

		h.main.View.X += ((dragRect.X + dragRect.W/2) - h.main.View.X) * (1 - scale)
		h.main.View.Y += ((dragRect.Y + dragRect.H/2) - h.main.View.Y) * (1 - scale)
		h.main.View.Scale *= scale

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
