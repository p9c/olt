package olt

import (
	"errors"

	"gioui.org/layout"
)

// FlexChildren is a struct to manage a list of layout.FlexChild(s) and provides a collection of editing functions
type FlexChildren struct {
	C      []layout.FlexChild
	Weight float32
	*Ctx
}

// AddHFlex adds a FlexChildren in horizontal orientation
func (f *FlexChildren) AddHFlex(weight float32, children FlexChildren) {
	f.C = append(f.C, f.GetHFlexed(weight, children.C...))
}

// AddVFlex adds a FlexChildren in vertical orientation
func (f *FlexChildren) AddVFlex(weight float32, children FlexChildren) {
	f.C = append(f.C, f.GetVFlexed(weight, children.C...))
}

// AddWidgets allows you to add widgets directly to a FlexChildren
func (f *FlexChildren) AddWidgets(weight float32, w ...layout.Widget) {
	for i := range w {
		f.C = append(f.C, layout.Flexed(weight, w[i]))
	}
}

// Append adds more FlexChildren to the end of a FlexChildren and returns it
func (f *FlexChildren) Append(a *FlexChildren) *FlexChildren {
	f.C = append(f.C, a.C...)
	return f
}

// Delete removes a specified set of elements in a FlexChildren
func (f *FlexChildren) Delete(start, end int) *FlexChildren {
	switch {
	case start < 0:
		f.err = errors.New("negative start")
		fallthrough
	case end < 0:
		f.err = errors.New("negative end")
		fallthrough
	case start > end:
		f.err = errors.New("region ends before it starts")
		fallthrough
	case end < len(f.C):
		f.err = errors.New("cannot delete outside of slice")
		fallthrough
	case start == end:
		f.err = errors.New("no elements will be deleted")
		break
	default:
		f.C = append(f.C[:start], f.C[end:]...)
	}
	return f
}

// FlexChildSlice returns the underlying []layout.FlexChild
func (f *FlexChildren) FlexChildSlice() []layout.FlexChild {
	return f.C
}

// GetHFlex returns a horizontal Layout.Flex with its contents inside
func (f *FlexChildren) GetHFlex() *layout.Flex {
	out := HorizontalFlexBox()
	out.Layout(f.Context, f.C...)
	return out

}

// GetVFlex returns a vertical layout.Flex with its contents inside
func (f *FlexChildren) GetVFlex() *layout.Flex {
	out := VerticalFlexBox()
	out.Layout(f.Context, f.C...)
	return out

}

// Insert inserts a given FlexChildren inside another FlexChildren and returns it
func (f *FlexChildren) Insert(index int, a *FlexChildren) *FlexChildren {
	switch {
	case index < 0:
		f.err = errors.New("negative index")
	case len(f.C) < index:
		f.err = errors.New("cannot insert beyond end of slice")
		break
	default:
		f.C = append(append(f.C[:index], a.C...), f.C[index:]...)
	}
	return f
}

// Prepend inserts a given FlexChildren before the existing contents and returns it
func (f *FlexChildren) Prepend(a *FlexChildren) *FlexChildren {
	f.C = append(a.C, f.C...)
	return f
}
