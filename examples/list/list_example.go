package main

import (
	"fmt"

	"gioui.org/font/gofont"
	"gioui.org/widget/material"

	"github.com/p9c/olt"
)

func main() {
	gofont.Register()
	th := material.NewTheme()
	const n = 1e6
	c := olt.New()
	list := olt.NewList(false)
	c.L.SetLevel("trace", true, "olt")
	_ = c.Window().Title("test").Size(640, 480).Load(func(c *olt.Ctx) {
		list.Layout(c.Context, n, func(i int) {
			th.H3(fmt.Sprintf("List element #%d", i)).Layout(c.Context)
		})
	}).Open()
}
