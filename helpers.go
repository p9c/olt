package olt

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/unit"
)

// Box is a rectangle that is that simplifies specifying drawing a rectangle
type Box struct {
	W, H int
	Radius
}

// ClipOp returns a clip.Rect based on a Box
func (b Box) ClipOp(c *Ctx) clip.Op {
	square := f32.Rectangle{
		Max: f32.Point{
			X: float32(b.W),
			Y: float32(b.H),
		},
	}
	return clip.Rect{
		Rect: square,
		SE:   b.SE,
		SW:   b.SW,
		NW:   b.NW,
		NE:   b.NE,
	}.Op(c.Ops)
}

// New creates a new Box
func (b Box) New(w, h int, r Radius) Box {
	return Box{w, h, r}
}

// Coord is a wrapper on image.Point so we can attach local methods to it
type Coord struct {
	image.Point
}

// New creates a new Coord
func (c Coord) New(x, y int) Coord {
	c.X, c.Y = x, y
	return c
}

// LayoutFunc is the function signature for a function that creates a new layout.Flex. This is used to create the root
// widget of a window where all the rest of the render tree is placed
type LayoutFunc func(c *Ctx)

// Radius is the corner radii for rounded rectangles. Zero means no rounding.
type Radius struct {
	SE, SW, NE, NW float32
}

// Widget is a wrapper on layout.Widget so we can create local methods on the same type
type Widget struct {
	layout.Widget
}

// Widgeter is an interface that provides automatic access to dimensions information for custom widgets
type Widgeter interface {
	Draw()
	Dimensions() image.Point
}

// ARGB returns a color.ARGB from a uint32, use like this: ARGB(0xAARRGGBB)
func ARGB(rgba uint32) (c color.RGBA) {
	c = color.RGBA{
		A: byte(rgba >> 24),
		R: byte(rgba >> 16),
		G: byte(rgba >> 8),
		B: byte(rgba),
	}
	return
}

// DP returns a unit.DP
func DP(i int) unit.Value {
	return unit.Dp(float32(i))
}
