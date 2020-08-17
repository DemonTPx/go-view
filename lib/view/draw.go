package view

import gl "github.com/chsc/gogl/gl21"

type Color struct {
	R, G, B, A float64
}

func NewColor(r, g, b, a float64) Color {
	return Color{R: r, G: g, B: b, A: a}
}

type Rect struct {
	X, Y, W, H float64
}

func NewRect(x, y, w, h float64) Rect {
	return Rect{X: x, Y: y, W: w, H: h}
}

func (r *Rect) X2() float64 {
	return r.X + r.W
}

func (r *Rect) Y2() float64 {
	return r.Y + r.H
}

func DrawQuad(rect Rect, color Color) {
	gl.Color4f(gl.Float(color.R), gl.Float(color.G), gl.Float(color.B), gl.Float(color.A))

	gl.Begin(gl.QUADS)
	gl.Vertex3f(gl.Float(rect.X), gl.Float(rect.Y), 0)
	gl.Vertex3f(gl.Float(rect.X2()), gl.Float(rect.Y), 0)
	gl.Vertex3f(gl.Float(rect.X2()), gl.Float(rect.Y2()), 0)
	gl.Vertex3f(gl.Float(rect.X), gl.Float(rect.Y2()), 0)
	gl.End()
}

func DrawQuadOutline(rect Rect, width float64, color Color) {
	DrawQuad(NewRect(rect.X, rect.Y, rect.W, width), color)
	DrawQuad(NewRect(rect.X2()-width, rect.Y, width, rect.H), color)
	DrawQuad(NewRect(rect.X, rect.Y2()-width, rect.W, width), color)
	DrawQuad(NewRect(rect.X, rect.Y, width, rect.H), color)
}

func DrawQuadBorder(rect Rect, color Color, borderWidth float64, borderColor Color) {
	DrawQuad(rect, color)
	DrawQuadOutline(rect, borderWidth, borderColor)
}
