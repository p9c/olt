package olt

import (
	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
)

// // Window encapsulates a window and binds it to a LayoutFunc which specifies how the interface works
// func Window(title string, W, H int, lf LayoutFunc) {
// 	go func() {
// 		c := New(app.NewWindow(app.Title(title), app.Size(DP(W), DP(H))))
// 		for e := range c.W.Events() {
// 			if e, ok := e.(system.FrameEvent); ok {
// 				c.Reset(e.Config, e.Size)
// 				lf(c)
// 				e.Frame(c.Ops)
// 			}
// 		}
// 	}()
// 	app.Main()
// }
//
// // Window encapsulates a window and binds it to a LayoutFunc which specifies how the interface works
// func (c *Ctx) Window(title string, W, H int, lf LayoutFunc) *Window {
// 	go func() {
// 		w := app.NewWindow(
// 			app.Title(title),
// 			app.Size(DP(W), DP(H)),
// 		)
// 		*c = Ctx{
// 			Context: layout.NewContext(w.Queue()),
// 			W:       w,
// 			err:     c.err,
// 			L:     c.L,
// 			LogC:    c.LogC,
// 		}
// 		for e := range c.W.Events() {
// 			if e, ok := e.(system.FrameEvent); ok {
// 				c.Reset(e.Config, e.Size)
// 				lf(c)
// 				e.Frame(c.Ops)
// 			}
// 		}
// 	}()
// 	app.Main()
// }

type Window struct {
	*Ctx
	*app.Window
	h, w  int
	title string
	lf    LayoutFunc
}

func (c *Ctx) Window() *Window {
	w := &Window{Ctx: c}
	w.L.Debug("new window")
	return w
}

func (w *Window) Size(W, H int) *Window {
	w.L.Debug("h", H, "w", W)
	w.h, w.w = H, W
	return w
}

func (w *Window) Title(title string) *Window {
	w.L.Debug("title", title)
	w.title = title
	return w
}

func (w *Window) Load(lf LayoutFunc) *Window {
	w.lf = lf
	return w
}

func (w *Window) Open() *Window {
	ww := app.NewWindow(
		app.Title(w.title),
		app.Size(DP(w.w), DP(w.h)),
	)
	w.Ctx.Context = layout.NewContext(ww.Queue())
	// w.L.Traces(ww)
	w.W = ww
	for e := range w.W.Events() {
		if e, ok := e.(system.FrameEvent); ok {
			w.Reset(e.Config, e.Size)
			w.lf(w.Ctx)
			e.Frame(w.Ops)
		}
	}
	app.Main()
	return w
}
