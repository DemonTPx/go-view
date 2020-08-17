package view

import (
	"log"
	"time"
)

type Command struct{}

type QuitCommand Command
type ZoomInCommand Command
type ZoomOutCommand Command
type ZoomOriginalSizeCommand Command
type ZoomFitToWindowCommand Command
type FirstFileCommand Command
type LastFileCommand Command
type NextFileCommand Command
type PreviousFileCommand Command

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
	case ZoomInCommand:
		h.main.View.Scale *= 1.25
	case ZoomOutCommand:
		h.main.View.Scale *= 0.8
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
	default:
		log.Printf("unexpected command: %#v", command)
	}
}

func (h *CommandHandler) HandleBlocking() {
	select {
	case command := <-h.commandChannel:
		log.Printf("received command: %#v", command)
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
