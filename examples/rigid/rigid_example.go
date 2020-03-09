package main

import (
	"fmt"

	log "github.com/p9c/logi"

	"github.com/p9c/olt"
)

func main() {
	fmt.Println("testing")
	log.L.SetLevel("trace", true)
	log.DEBUG("logging")
	ctx := olt.New()
	_ = ctx
	// w := ctx.Window().Title("rigid example").Size(640,480)
	// w.Open()
	// h := ctx.NewFlexChildren()
	// ctx.Window("test", 640, 480, func(c *olt.Ctx) {
	// 	h.AddWidgets(1, c.EmptyFlexBox(c.ARGB(0xffff8800)))
	// 	h.GetHFlex()
	// })
}
