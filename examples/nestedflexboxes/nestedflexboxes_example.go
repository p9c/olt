package main

import (
	"errors"
	"fmt"

	"github.com/p9c/olt"
)

func main() {
	c := olt.New()
	c.L.SetLevel("trace", true, "olt")
	c.L.Trace("testing")
	c.L.Debug("testing")
	fmt.Println("'", c.L.Check(errors.New("this is a test")), "'")
	c.L.Check(nil)
	c.L.Info("testing")
	c.L.Warn("testing")
	c.L.Error("testing")
	c.L.Fatal("testing")
	// c.Window().Title("nested flexboxes example").Dimensions(640, 480).Window()
	// c.Window("test", 640, 480, func(c *olt.Ctx) {
	// 	n := uint32(10)
	// 	nn := float32(n)
	// 	h := c.NewFlexChildren()
	// 	for i := n - n; i < n; i++ {
	// 		v := c.NewFlexChildren()
	// 		for j := n - n; j < n; j++ {
	// 			v.AddWidgets(1/nn, c.EmptyFlexBox(color.RGBA{
	// 				R: 0,
	// 				G: byte(255 - 256/n*(j+1)),
	// 				B: byte(255 - 256/n*(i+1)),
	// 				A: 255,
	// 			}))
	// 		}
	// 		h.AddVFlex(1/nn, v)
	// 	}
	// 	h.GetHFlex()
	// })
}
