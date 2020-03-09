package main

import (
	"image/color"
	"log"

	"gioui.org/layout"

	"github.com/p9c/olt"
)

func main() {
	olt.NewWindow("test", 640, 480, func(c *olt.Ctx) (lf layout.Flex) {
		n := uint32(10)
		nn := float32(n)
		var h []layout.FlexChild
		for i := n - n; i < n; i++ {
			var v []layout.FlexChild
			for j := n - n; j < n; j++ {
				log.Println(i, j)
				v = append(v, layout.Flexed(1/nn, c.EmptyFlexBox(color.RGBA{
					R: 0,
					G: byte(255 - 256/n*(j+1)),
					B: byte(255 - 256/n*(i+1)),
					A: 255,
				})))
			}
			h = append(h, layout.Flexed(1/nn, func() { c.VerticalFlexBox().Layout(c.Context, v...) }))
		}
		// spew.Dump(h)
		out := c.HorizontalFlexBox()
		out.Layout(c.Context, h...)
		return out
	})
}
