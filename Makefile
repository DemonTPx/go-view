VERSION=0.1

.PHONY: all clean

all: build/go-view

clean:
	rm -rf build/*

build/go-view:
	go build -o build/go-view .
