package view

import "math"

type Mouse struct {
	X, Y uint32

	dragging     bool
	dragX, dragY uint32
}

func (m *Mouse) DragRect() Rect {
	return NewRect(
		math.Min(float64(m.X), float64(m.dragX)),
		math.Min(float64(m.Y), float64(m.dragY)),
		math.Abs(float64(m.X)-float64(m.dragX)),
		math.Abs(float64(m.Y)-float64(m.dragY)),
	)
}
