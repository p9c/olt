package olt

import (
	"gioui.org/app"
	"gioui.org/io/system"
)

// NewWindow encapsulates a window and binds it to a LayoutFunc which specifies how the interface works
func NewWindow(title string, W, H int, lf LayoutFunc) {
	go func() {
		w := app.NewWindow(
			app.Title(title),
			app.Size(DP(W), DP(H)),
		)
		c := &Ctx{W: w}
		c.New()
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
