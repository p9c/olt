package olt

import (
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/widget/material"
)

type List struct {
	*Ctx
	Theme *material.Theme
	*layout.List
}

func init() {
	gofont.Register()
}
