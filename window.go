package olt

import (
	"gioui.org/app"
	"gioui.org/io/system"
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
// 			Log:     c.Log,
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
	w.Println("new window")
	return w
}

func (w *Window) Size(H, W int) *Window {
	w.Println("h", H, "w", W)
	w.h, w.w = H, W
	return w
}

func (w *Window) Title(title string) *Window {
	w.Println("title", title)
	w.title = title
	return w
}

func (w *Window) Open() *Window {
	w.Println(w)
	w.Window = app.NewWindow(
		app.Title(w.title),
		app.Size(DP(w.w), DP(w.h)),
	)
	w.Println(w.Window)
	w.W = w.Window
	for e := range w.Window.Events() {
		if e, ok := e.(system.FrameEvent); ok {
			w.Println(w, e, e.Config, e.Size)
			w.Reset(e.Config, e.Size)
			w.lf(w.Ctx)
			e.Frame(w.Ops)
		}
	}
	app.Main()
	return w
}
