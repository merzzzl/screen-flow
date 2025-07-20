package vision

import (
	"image"
	"image/color"

	"gocv.io/x/gocv"
)

func (p *Pipe) showImage(obj gocv.Mat, good image.Point, all []image.Point) {
	for _, pt := range all {
		drowPoint(&obj, pt, color.RGBA{255, 0, 0, 63})
	}

	if good.X != 0 || good.Y != 0 {
		drowPoint(&obj, good, color.RGBA{0, 255, 0, 255})
	}

	p.window.Show(obj)
}

func drowPoint(mat *gocv.Mat, point image.Point, cl color.RGBA) {
	if point.X > 0 && point.Y > 0 {
		startH := image.Pt(point.X-40, point.Y)
		endH := image.Pt(point.X+40, point.Y)

		_ = gocv.Line(mat, startH, endH, cl, 8)

		startV := image.Pt(point.X, point.Y-40)
		endV := image.Pt(point.X, point.Y+40)

		_ = gocv.Line(mat, startV, endV, cl, 8)
	}
}
