package olt

import (
	"errors"
	"image"
	"image/color"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"

	"github.com/p9c/logi"
)

// Ctx is a wrapper around layout.Context and app.Window and embeds an error so its methods can be chained
type Ctx struct {
	*layout.Context
	W   *app.Window
	err error
	L *logi.Logger
}

// ClearError nils the embedded error
func (c *Ctx) ClearError() *Ctx {
	c.err = nil
	return c
}

// Ctx returns the underlying layout.Context
func (c Ctx) Ctx() *layout.Context {
	return c.Context
}

// DrawRectangle draws a box with a given set of corner radii and a fill colour
func (c *Ctx) DrawRectangle(color color.RGBA, b Box, inset unit.Value) {
	if color.A == 0 {
		c.err = errors.New("not drawing an invisible rectangle")
	}
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

// EmptyFlexBox is just a box with a given colour
func (c *Ctx) EmptyRigid(box Box, col ...color.RGBA) layout.Widget {
	return func() {
		// layout.Flexed(1,
		cs := c.Constraints
		switch {
		case box.H == 0 || box.W == 0:
			c.L.Error("not drawing zero width box. There is no spoon")
		case box.H < 0:
			box.H = cs.Height.Max
		case box.W < 0:
			box.W = cs.Width.Max
		}
		cc := color.RGBA{}
		if len(col) == 1 {
			cc = col[0]
		}
		c.DrawRectangle(cc, box, unit.Dp(0))
	}
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
		c.DrawRectangle(cc, NewBox(cs.Width.Max, cs.Height.Max), unit.Dp(0))
	}
}

// Err returns the underlying error variable
func (c *Ctx) Err() error {
	return c.err
}

// Error implements the Error interface and returns a string
func (c *Ctx) Error() string {
	return c.err.Error()
}

// GetList returns a simple list composed of a vertical layout.Flex that can scroll up and down
func (c *Ctx) GetList(fn func(i int)) List {
	ll := &layout.List{}
	// ll.Layout(c.Context, len(items), fn)
	return List{c, material.NewTheme(), ll}
}

// GetHFlexed returns a layout.FlexChild in horizontal orientation embedded in a layout.FlexChild
func (c *Ctx) GetHFlexed(weight float32, children ...layout.FlexChild) layout.FlexChild {
	return layout.Flexed(weight, func() { HorizontalFlexBox().Layout(c.Context, children...) })
}

// GetVFlexed returns a layout.FlexChild in vertical orientation embedded in a layout.FlexChild
func (c *Ctx) GetVFlexed(weight float32, children ...layout.FlexChild) layout.FlexChild {
	return layout.Flexed(weight, func() { VerticalFlexBox().Layout(c.Context, children...) })
}

// Seterror sets the underlying error value directly and logs it if the closure is loaded
func (c *Ctx) Seterror(err error) *Ctx {
	c.err = err
	c.L.Error(err)
	return c
}

// NewFlexChildren creates a new FlexChildren and binds itself to it
func (c *Ctx) NewFlexChildren() *FlexChildren {
	return &FlexChildren{Ctx: c}
}

// SetError sets the error to a new string and logs it
func (c *Ctx) SetError(err string) *Ctx {
	c.err = errors.New(err)
	c.L.Error(c.err)
	return c
}

// SetErrorLogger loads a function that is used to print errors when they are set
func (c *Ctx) SetErrorLogger(logger func(err string)) *Ctx {
	c.L = logi.L
	return c
}

// HorizontalFlexBox returns an empty layout.Flex set to horizontal
func HorizontalFlexBox() *layout.Flex {
	return &layout.Flex{Axis: layout.Horizontal}
}

// New returns a new context. This is an initializer, invoke thus:
//
// 		ctx := olt.Ctx{}.New()
func New(w ...*app.Window) (c *Ctx) {
	c = &Ctx{
		L: logi.Empty().SetLevel("trace", true, "olt"),
	}
	if len(w) > 0 {
		c.W = w[0]
		c.Context = layout.NewContext(c.W.Queue())
	}
	return c
}

// NewList returns a new list layout box
func NewList(scrollToEnd bool) *layout.List {
	list := &layout.List{
		Axis:        layout.Vertical,
		ScrollToEnd: scrollToEnd,
	}
	return list
}

// Rigid returns a rigid layout box
func Rigid(w layout.Widget) (out layout.FlexChild) {
	out = layout.Rigid(w)
	return
}

// VerticalFlexBox returns an empty layout.Flex set to vertical
func VerticalFlexBox() *layout.Flex {
	return &layout.Flex{Axis: layout.Vertical}
}
