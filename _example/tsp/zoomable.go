package main

import (
	"image"

	"gioui.org/f32"
	"gioui.org/gesture"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
)

type Zoomable struct {
	scroll gesture.Scroll
	drag   gesture.Drag
	mouse  f32.Point

	dragOrg f32.Point
	dragCur f32.Point
	offset  f32.Point

	tr f32.Affine2D
}

func (z *Zoomable) Layout(gtx C, zoomed layout.Widget) D {
	r := clip.Rect{Max: gtx.Constraints.Max}
	r.Push(gtx.Ops)

	z.scroll.Add(gtx.Ops, image.Rect(0, -1, 0, 1))
	nscroll := z.scroll.Scroll(gtx.Metric, gtx, gtx.Now, gesture.Vertical)
	pointer.InputOp{Tag: z, Types: pointer.Move}.Add(gtx.Ops)

	for _, ev := range gtx.Events(z) {
		switch ev := ev.(type) {
		case pointer.Event:
			switch ev.Type {
			case pointer.Move:
				z.mouse = ev.Position
			}
		}
	}

	z.drag.Add(gtx.Ops)
	for _, ev := range z.drag.Events(gtx.Metric, gtx, gesture.Both) {
		switch ev.Type {
		case pointer.Press:
			z.dragOrg = ev.Position
		case pointer.Drag:
			z.dragCur = z.dragOrg.Sub(ev.Position)
		case pointer.Release:
			z.offset = z.offset.Sub(z.dragOrg).Add(ev.Position)
			z.dragCur = f32.Point{}
		}
	}

	if nscroll != 0 {
		var change float32
		if nscroll > 0 {
			change = 0.9
		} else {
			change = 1.1
		}
		z.tr = z.tr.Scale(z.mouse, f32.Pt(change, change))
	}

	op.Affine(z.tr).Add(gtx.Ops)

	// Divide offset by scaling factor (sx == sy)
	sx, _, _, _, _, _ := z.tr.Elems()
	off := z.offset.Sub(z.dragCur).Div(sx)
	op.Offset(off.Round()).Add(gtx.Ops)

	return zoomed(gtx)
}
