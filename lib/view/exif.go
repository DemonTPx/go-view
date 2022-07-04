package view

import (
	"fmt"
	"os"

	"github.com/rwcarlsen/goexif/exif"
)

type Orientation struct {
	mirrored     bool
	numRotations int
}

var DefaultOrientation = Orientation{false, 0}

var orientationMap = map[int]Orientation{
	1: {false, 0},
	2: {true, 0},
	3: {true, 2},
	4: {false, 2},
	5: {true, 1},
	6: {false, 1},
	7: {true, 3},
	8: {false, 3},
}

func ReadExifOrientation(filename string) (Orientation, error) {
	f, err := os.Open(filename)
	if err != nil {
		return DefaultOrientation, fmt.Errorf("error while opening file: %s", err)
	}
	defer f.Close()

	x, err := exif.Decode(f)
	if err != nil {
		return DefaultOrientation, nil
	}

	orientation, err := x.Get(exif.Orientation)
	if err != nil {
		return DefaultOrientation, nil
	}

	i, _ := orientation.Int(0)

	o, ok := orientationMap[i]
	if !ok {
		return DefaultOrientation, nil
	}

	return o, nil
}
