package olt

import (
	"image/color"

	"gioui.org/layout"
)

// ExampleNewWindow demonstrates a simple empty flexbox function
func ExampleNewWindow() {
	NewWindow("test", 640, 480, func(c *Ctx) (lf layout.Flex) {
		n := uint32(10)
		nn := float32(n)
		h := c.NewFlexChildren()
		for i := n - n; i < n; i++ {
			v := c.NewFlexChildren()
			for j := n - n; j < n; j++ {
				v.AddWidgets(1/nn, c.EmptyFlexBox(color.RGBA{
					R: 0,
					G: byte(255 - 256/n*(j+1)),
					B: byte(255 - 256/n*(i+1)),
					A: 255,
				}))
			}
			h.AddVFlex(1/nn, v)
		}
		return h.GetHFlex()
	})
}
