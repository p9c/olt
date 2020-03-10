package main

import (
	"image/color"

	"github.com/p9c/olt"
)

func main() {
	c := olt.New()
	c.L.SetLevel("trace", true, "olt")
	_ = c.Window().Size(640, 480).Title("nested flexboxes example").Load(func(c *olt.Ctx) {
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
		h.RenderHFlex()
	}).Open()
	// c.Window().Title("nested flexboxes example").Dimensions(640, 480).Window()
	// c.Window("test", 640, 480, func(c *olt.Ctx) {
	// })
}
