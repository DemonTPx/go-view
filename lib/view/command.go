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
	W, H float64
}
type SaveSettingsCommand struct{}
type MouseCursorPositionCommand struct {
	X, Y float64
}
type MoveViewCommand struct {
	X, Y float64
}
type StartDragLeftCommand struct{}
type StopDragLeftCommand struct{}
type StartDragRightCommand struct{}
type StopDragRightCommand struct{}

type CommandHandler struct {
	main           *Main
	commandChannel <-chan interface{}
}

func NewCommandHandler(main *Main, commandChannel <-chan interface{}) *CommandHandler {
	return &CommandHandler{main: main, commandChannel: commandChannel}
}

func (h *CommandHandler) HandleCommand(command interface{}) (waitForCommand bool) {
	switch command.(type) {
	case QuitCommand:
		h.main.Running = false

	case ZoomCommand:
		c := command.(ZoomCommand)
		h.main.View.Scale *= c.Scale

	case ZoomToMouseCursorCommand:
		c := command.(ZoomToMouseCursorCommand)
		if c.Scale < 1 {
			h.main.View.X += (h.main.View.W/2 - h.main.View.X) * (1 - c.Scale)
			h.main.View.Y += (h.main.View.H/2 - h.main.View.Y) * (1 - c.Scale)
		} else {
			h.main.View.X += (h.main.Mouse.X - h.main.View.X) * (1 - c.Scale)
			h.main.View.Y += (h.main.Mouse.Y - h.main.View.Y) * (1 - c.Scale)
		}
		h.main.View.Scale *= c.Scale

	case ZoomOriginalSizeCommand:
		h.main.View.X = h.main.View.W / 2
		h.main.View.Y = h.main.View.H / 2
		h.main.View.Scale = 1

	case ZoomFitToWindowCommand:
		h.main.View.X = h.main.View.W / 2
		h.main.View.Y = h.main.View.H / 2
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
		waitForCommand = true

	case MouseCursorPositionCommand:
		c := command.(MouseCursorPositionCommand)

		if h.main.Mouse.DragRight.Dragging {
			h.main.View.X += c.X - h.main.Mouse.X
			h.main.View.Y += c.Y - h.main.Mouse.Y
		}

		h.main.Mouse.X = c.X
		h.main.Mouse.Y = c.Y

		waitForCommand = !h.main.Mouse.DragLeft.Dragging && !h.main.Mouse.DragRight.Dragging

	case StartDragLeftCommand:
		h.main.Mouse.DragLeft = MouseDrag{
			Dragging: true,
			X:        h.main.Mouse.X,
			Y:        h.main.Mouse.Y,
		}
		waitForCommand = true

	case StopDragLeftCommand:
		h.main.Mouse.DragLeft.Dragging = false

		dragRect := h.main.Mouse.DragLeftRect()

		if dragRect.W < DragThreshold || dragRect.H < DragThreshold {
			break
		}

		dragRatio := dragRect.W / dragRect.H
		windowRatio := h.main.View.W / h.main.View.H

		var scale float64
		if dragRatio > windowRatio {
			scale = h.main.View.W / dragRect.W
		} else {
			scale = h.main.View.H / dragRect.H
		}

		// zoom in
		h.main.View.X += ((dragRect.X + dragRect.W/2) - h.main.View.X) * (1 - scale)
		h.main.View.Y += ((dragRect.Y + dragRect.H/2) - h.main.View.Y) * (1 - scale)

		// move to center
		h.main.View.X += h.main.View.W/2 - (dragRect.X + dragRect.W/2)
		h.main.View.Y += h.main.View.H/2 - (dragRect.Y + dragRect.H/2)

		h.main.View.Scale *= scale

	case StartDragRightCommand:
		h.main.Mouse.DragRight = MouseDrag{
			Dragging: true,
			X:        h.main.Mouse.X,
			Y:        h.main.Mouse.Y,
		}
		waitForCommand = true

	case StopDragRightCommand:
		h.main.Mouse.DragRight.Dragging = false
		waitForCommand = true

	case MoveViewCommand:
		c := command.(MoveViewCommand)
		h.main.View.X += c.X
		h.main.View.Y += c.Y

	default:
		log.Printf("unexpected command: %#v", command)
	}

	return
}

func (h *CommandHandler) HandleBlocking() {
	waitForCommand := true
	for waitForCommand {
		select {
		case command := <-h.commandChannel:
			waitForCommand = h.HandleCommand(command)
		}
	}
}

func (h *CommandHandler) HandleBlockingOrAtLeast(duration time.Duration) {
	timeout := time.After(duration)
	waitForCommand := true
	breakOutHit := false
	timeoutHit := false
	for !timeoutHit || !breakOutHit {
		select {
		case command := <-h.commandChannel:
			waitForCommand = h.HandleCommand(command)
			if !waitForCommand {
				breakOutHit = true
			}
		case <-timeout:
			timeoutHit = true
		}
	}
}

func (h *CommandHandler) HandleTimeout(timeout time.Duration) {
	waitForCommand := true
	for waitForCommand {
		select {
		case command := <-h.commandChannel:
			waitForCommand = h.HandleCommand(command)
		case <-time.After(timeout):
			waitForCommand = false
		}
	}
}
