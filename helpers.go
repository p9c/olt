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
	error
}

func (c Ctx) New() *Ctx {
	return &Ctx{layout.NewContext(c.W.Queue()), c.W, nil}
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
func (c *Ctx) EmptyFlexBox(col ...color.RGBA) layout.Widget {
	return func() {
		// layout.Flexed(1,
		cs := c.Constraints
		cc := color.RGBA{}
		if len(col) == 1 {
			cc = col[0]
		}
		c.DrawRectangle(cc, Box{}.New(cs.Width.Max, cs.Height.Max, Radius{}), unit.Dp(0))
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

type FlexChildren struct {
	C      []layout.FlexChild
	Weight float32
	*Ctx
}

func (f *FlexChildren) Append(a *FlexChildren) *FlexChildren {
	f.C = append(f.C, a.C...)
	return f
}

func (f *FlexChildren) Prepend(a *FlexChildren) *FlexChildren {
	f.C = append(a.C, f.C...)
	return f
}

func (f *FlexChildren) Insert(index int, a *FlexChildren) *FlexChildren {
	if len(f.C) < index {
		f.C = append(append(f.C[:index], a.C...), f.C[index:]...)
	}
	return f
}

func (f *FlexChildren) Delete(start, end int) *FlexChildren {
	if start < end && end < len(f.C) {
		f.C = append(f.C[:start], f.C[end:]...)
	}
	return f
}

func (f *FlexChildren) AddWidgets(weight float32, w ...layout.Widget) {
	for i := range w {
		f.C = append(f.C, layout.Flexed(weight, w[i]))
	}
}

func (f *FlexChildren) AddV(weight float32, children FlexChildren) {
	// for i := range children.C {
	f.C = append(f.C, f.GetVFlexed(weight, children.C...))
	// }
}

func (f *FlexChildren) AddH(weight float32, children FlexChildren) {
	// for i := range children.C {
	f.C = append(f.C, f.GetHFlexed(weight, children.C...))
	// }
}

func (f *FlexChildren) FlexChildSlice() []layout.FlexChild {
	return f.C
}

func (f *FlexChildren) GetHFlex() layout.Flex {
	out := f.HorizontalFlexBox()
	out.Layout(f.Context, f.C...)
	return out

}

func (f *FlexChildren) GetVFlex() layout.Flex {
	out := f.VerticalFlexBox()
	out.Layout(f.Context, f.C...)
	return out

}

func (c *Ctx) NewFlexChildren() FlexChildren {
	return FlexChildren{Ctx: c}
}

func (c *Ctx) GetHFlexed(weight float32, children ...layout.FlexChild) layout.FlexChild {
	return layout.Flexed(weight, func() { c.HorizontalFlexBox().Layout(c.Context, children...) })
}

func (c *Ctx) GetVFlexed(weight float32, children ...layout.FlexChild) layout.FlexChild {
	return layout.Flexed(weight, func() { c.VerticalFlexBox().Layout(c.Context, children...) })
}

func (c *Ctx) HorizontalFlexBox() layout.Flex {
	return layout.Flex{Axis: layout.Horizontal}
}

func (c *Ctx) VerticalFlexBox() layout.Flex {
	return layout.Flex{Axis: layout.Vertical}
}
