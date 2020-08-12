package main

import (
	"github.com/DemonTPx/go-view/lib/view"
	"os"
	"runtime"
)

func main() {
	runtime.LockOSThread()

	filename := ""
	if len(os.Args) >= 2 {
		filename = os.Args[1]
	}

	main := view.NewMain(filename)
	err := main.Run()
	if err != nil {
		panic(err)
	}
}
