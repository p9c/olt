package olt

import (
	"errors"
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

// Ctx is a wrapper around layout.Context and app.Window and embeds an error so its methods can be chained
type Ctx struct {
	*layout.Context
	W    app.Window
	err  error
	Log  func(err string)
	LogC func(err error)
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

// Err returns the underlying error variable
func (c *Ctx) Err() error {
	return c.err
}

// Error implements the Error interface and returns a string
func (c *Ctx) Error() string {
	return c.err.Error()
}

// GetHFlexed returns a layout.FlexChild in horizontal orientation embedded in a layout.FlexChild
func (c *Ctx) GetHFlexed(weight float32, children ...layout.FlexChild) layout.FlexChild {
	return layout.Flexed(weight, func() { c.HorizontalFlexBox().Layout(c.Context, children...) })
}

// GetVFlexed returns a layout.FlexChild in vertical orientation embedded in a layout.FlexChild
func (c *Ctx) GetVFlexed(weight float32, children ...layout.FlexChild) layout.FlexChild {
	return layout.Flexed(weight, func() { c.VerticalFlexBox().Layout(c.Context, children...) })
}

// HorizontalFlexBox returns an empty layout.Flex set to horizontal
func (c *Ctx) HorizontalFlexBox() layout.Flex {
	return layout.Flex{Axis: layout.Horizontal}
}

// New returns a new context. This is an initializer, invoke thus:
//
// 		ctx := olt.Ctx{}.New()
func (c Ctx) New() *Ctx {
	cp := &Ctx{layout.NewContext(c.W.Queue()), c.W, nil, func(string) {}, func(error) {}}
	return cp
}

// NewFlexChildren creates a new FlexChildren and binds itself to it
func (c *Ctx) NewFlexChildren() FlexChildren {
	return FlexChildren{Ctx: c}
}

// Seterror sets the underlying error value directly and logs it if the closure is loaded
func (c *Ctx) Seterror(err error) *Ctx {
	c.err = err
	c.LogC(err)
	return c
}

// SetError sets the error to a new string and logs it
func (c *Ctx) SetError(err string) *Ctx {
	c.err = errors.New(err)
	c.LogC(c.err)
	return c
}

// SetErrorLogger loads a function that is used to print errors when they are set
func (c *Ctx) SetErrorLogger(logger func(err string)) *Ctx {
	c.Log = logger
	c.LogC = func(err error) {
		c.Log(err.Error())
	}
	return c
}

// VerticalFlexBox returns an empty layout.Flex set to vertical
func (c *Ctx) VerticalFlexBox() layout.Flex {
	return layout.Flex{Axis: layout.Vertical}
}

// FlexChildren is a struct to manage a list of layout.FlexChild(s) and provides a collection of editing functions
type FlexChildren struct {
	C      []layout.FlexChild
	Weight float32
	*Ctx
}

// AddHFlex adds a FlexChildren in horizontal orientation
func (f *FlexChildren) AddHFlex(weight float32, children FlexChildren) {
	f.C = append(f.C, f.GetHFlexed(weight, children.C...))
}

// AddVFlex adds a FlexChildren in vertical orientation
func (f *FlexChildren) AddVFlex(weight float32, children FlexChildren) {
	f.C = append(f.C, f.GetVFlexed(weight, children.C...))
}

// AddWidgets allows you to add widgets directly to a FlexChildren
func (f *FlexChildren) AddWidgets(weight float32, w ...layout.Widget) {
	for i := range w {
		f.C = append(f.C, layout.Flexed(weight, w[i]))
	}
}

// Append adds more FlexChildren to the end of a FlexChildren and returns it
func (f *FlexChildren) Append(a *FlexChildren) *FlexChildren {
	f.C = append(f.C, a.C...)
	return f
}

// Delete removes a specified set of elements in a FlexChildren
func (f *FlexChildren) Delete(start, end int) *FlexChildren {
	switch {
	case start < 0:
		f.err = errors.New("negative start")
		fallthrough
	case end < 0:
		f.err = errors.New("negative end")
		fallthrough
	case start > end:
		f.err = errors.New("region ends before it starts")
		fallthrough
	case end < len(f.C):
		f.err = errors.New("cannot delete outside of slice")
		fallthrough
	case start == end:
		f.err = errors.New("no elements will be deleted")
		break
	default:
		f.C = append(f.C[:start], f.C[end:]...)
	}
	return f
}

// FlexChildSlice returns the underlying []layout.FlexChild
func (f *FlexChildren) FlexChildSlice() []layout.FlexChild {
	return f.C
}

// GetHFlex returns a horizontal Layout.Flex with its contents inside
func (f *FlexChildren) GetHFlex() layout.Flex {
	out := f.HorizontalFlexBox()
	out.Layout(f.Context, f.C...)
	return out

}

// GetVFlex returns a vertical layout.Flex with its contents inside
func (f *FlexChildren) GetVFlex() layout.Flex {
	out := f.VerticalFlexBox()
	out.Layout(f.Context, f.C...)
	return out

}

// Insert inserts a given FlexChildren inside another FlexChildren and returns it
func (f *FlexChildren) Insert(index int, a *FlexChildren) *FlexChildren {
	switch {
	case index < 0:
		f.err = errors.New("negative index")
	case len(f.C) < index:
		f.err = errors.New("cannot insert beyond end of slice")
		break
	default:
		f.C = append(append(f.C[:index], a.C...), f.C[index:]...)
	}
	return f
}

// Prepend inserts a given FlexChildren before the existing contents and returns it
func (f *FlexChildren) Prepend(a *FlexChildren) *FlexChildren {
	f.C = append(a.C, f.C...)
	return f
}

// LayoutFunc is the function signature for a function that creates a new layout.Flex. This is used to create the root
// widget of a window where all the rest of the render tree is placed
type LayoutFunc func(c *Ctx) layout.Flex

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

// NewWindow encapsulates a window and binds it to a LayoutFunc which specifies how the interface works
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
