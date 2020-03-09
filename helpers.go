package olt

import (
	"image"
	"image/color"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

// ARGB returns a color.ARGB from a uint32, use like this: ARGB(0xRRGGBBAA)
func ARGB(rgba uint32) (c color.RGBA) {
	c = color.RGBA{
		A: byte(rgba >> 24),
		R: byte(rgba >> 16),
		G: byte(rgba >> 8),
		B: byte(rgba),
	}
	return
}

func DP(i int) unit.Value {
	return unit.Dp(float32(i))
}

type Ctx struct {
	*layout.Context
	W app.Window
}

func (c Ctx) New() *Ctx {
	return &Ctx{layout.NewContext(c.W.Queue()), c.W}
}

func (c Ctx) Ctx() *layout.Context {
	return c.Context
}

type Widget struct {
	layout.Widget
}

type Coord struct {
	image.Point
}

func (c Coord) New(x, y int) Coord {
	c.X, c.Y = x, y
	return c
}

type Radius struct {
	SE, SW, NE, NW float32
}

type Box struct {
	W, H int
	Radius
}

func (b Box) New(w, h int, r Radius) Box {
	return Box{w, h, r}
}

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

type LayoutFunc func(c *Ctx) layout.Flex

type Widgeter interface {
	Draw()
	Dimensions() image.Point
}

func NewWindow(title string, W, H int, lf LayoutFunc) {
	go func() {
		c := Ctx{
			W: *app.NewWindow(
				app.Title(title),
				app.Size(DP(W), DP(H)),
			)}.New()
		for e := range c.W.Events() {
			if e, ok := e.(system.FrameEvent); ok {
				c.Reset(e.Config, e.Size)
				lf(c)
				e.Frame(c.Ops)
			}
		}
	}()
	app.Main()
}

// EmptyFlexBox is just a box with a given colour
func (c *Ctx) EmptyFlexBox(color color.RGBA) layout.Widget {
	return func() {
		// layout.Flexed(1,
		cs := c.Constraints
		c.DrawRectangle(color, Box{}.New(cs.Width.Max, cs.Height.Max, Radius{}), unit.Dp(0))
	}
}

func (c *Ctx) DrawRectangle(color color.RGBA, b Box, inset unit.Value) {
	in := layout.UniformInset(inset)
	in.Layout(c.Context, func() {
		square := f32.Rectangle{
			Max: f32.Point{
				X: float32(b.W),
				Y: float32(b.H),
			},
		}
		paint.ColorOp{Color: color}.Add(c.Ops)
		clip.Rect{Rect: square,
			NE: b.NE, NW: b.NW, SE: b.SE, SW: b.SW}.Op(c.Ops).Add(c.Ops) // HLdraw
		paint.PaintOp{Rect: square}.Add(c.Ops)
		c.Dimensions = layout.Dimensions{Size: image.Point{X: b.W, Y: b.H}}
	})
}

func (c *Ctx) HorizontalFlexBox() layout.Flex {
	return layout.Flex{Axis: layout.Horizontal}
}

func (c *Ctx) VerticalFlexBox() layout.Flex {
	return layout.Flex{Axis: layout.Vertical}
}
