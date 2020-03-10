package main

import (
	"github.com/p9c/olt"
)

func main() {
	c := olt.New()
	_ = c.Window().
		Title("rigid example").
		Size(640, 480).
		Load(func(c *olt.Ctx) {
			c.NewFlexChildren().AddWidgets(1,
				// fills all behind body
				c.EmptyFlexBox(c.ARGB(0xffcccccc)),
			).RenderHFlex()
			c.NewFlexChildren().AddWidgets(1,
				// sidebar
				c.EmptyRigid(olt.NewBox(64, -1), c.ARGB(0xff333333)),
			).RenderVFlex()
			c.NewFlexChildren().AddWidgets(1,
				// top bar
				c.EmptyRigid(olt.NewBox(-1, 64), c.ARGB(0xff000000)),
			).RenderHFlex()
		}).Open()
}
