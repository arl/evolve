package main

import (
	"fmt"
	"image"
	"image/color"

	"github.com/arl/evolve/pkg/tsp"
	"golang.org/x/exp/constraints"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

var (
	pathColor = color.NRGBA{A: 255}
	dotColor  = color.NRGBA{R: 200, A: 255}
)

type pathWidget[T constraints.Integer] struct {
	citymax f32.Point
	cities  []tsp.Point2D

	zoomable Zoomable
}

func newPathWidget[T constraints.Integer](cities []tsp.Point2D) *pathWidget[T] {
	max := func(a, b float32) float32 {
		if a > b {
			return a
		}
		return b
	}

	// Compute world bounds
	var citymax f32.Point
	for _, c := range cities {
		citymax.X = max(citymax.X, float32(c.X))
		citymax.Y = max(citymax.Y, float32(c.Y))
	}
	fmt.Println("world bounds", citymax)

	return &pathWidget[T]{cities: cities, citymax: citymax}
}

func (pw *pathWidget[T]) Layout(onlyCities bool, sol []T, gtx C) D {
	return pw.zoomable.Layout(gtx, func(gtx C) D {
		// Draw cities as red dots
		const cityRadius = 5
		for i := range pw.cities {
			city := pw.cities[i]
			circle := clip.Ellipse{
				Min: image.Pt(int(city.X-cityRadius), int(city.Y-cityRadius)),
				Max: image.Pt(int(city.X+cityRadius), int(city.Y+cityRadius)),
			}.Op(gtx.Ops)
			paint.FillShape(gtx.Ops, dotColor, circle)
		}

		if !onlyCities && len(sol) != 0 {
			// At start we may not have received the first solution yet.
			p := clip.Path{}
			p.Begin(gtx.Ops)
			pt := f32.Pt(float32(pw.cities[sol[0]].X), float32(pw.cities[sol[0]].Y))
			p.MoveTo(pt)
			for i := 1; i < len(sol); i++ {
				pt := f32.Pt(float32(pw.cities[sol[i]].X), float32(pw.cities[sol[i]].Y))
				p.LineTo(pt)
			}
			p.LineTo(pt)
			paint.FillShape(gtx.Ops, pathColor, clip.Stroke{Path: p.End(), Width: 1}.Op())
		}

		op.InvalidateOp{}.Add(gtx.Ops)
		return layout.Dimensions{Size: gtx.Constraints.Max}
	})
}
