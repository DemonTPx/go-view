package view

import (
	"log"
	"time"
)

type Command struct{}

type ZoomInCommand Command
type ZoomOutCommand Command
type ZoomOriginalSizeCommand Command
type ZoomFitToWindowCommand Command
type QuitCommand Command

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
