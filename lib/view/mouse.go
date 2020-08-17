package view

import "math"

type Mouse struct {
	X, Y float64

	DragLeft  MouseDrag
	DragRight MouseDrag
}

type MouseDrag struct {
	Dragging bool
	X, Y     float64
}

func (m *Mouse) DragLeftRect() Rect {
	return m.dragRect(m.DragLeft)
}

func (m *Mouse) DragRightRect() Rect {
	return m.dragRect(m.DragRight)
}

func (m *Mouse) dragRect(drag MouseDrag) Rect {
	return NewRect(
		math.Min(m.X, drag.X),
		math.Min(m.Y, drag.Y),
		math.Abs(m.X-drag.X),
		math.Abs(m.Y-drag.Y),
	)
}
