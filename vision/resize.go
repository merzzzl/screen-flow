package vision

import (
	"image"

	"gocv.io/x/gocv"
)

func resizeSrc(src gocv.Mat) (gocv.Mat, float64) {
	if src.Empty() {
		return src, 1
	}

	maxSide := src.Cols()

	if src.Rows() > maxSide {
		maxSide = src.Rows()
	}

	if maxSide <= 640 {
		return src, 1
	}

	scale := float64(640) / float64(maxSide)

	newW := int(float64(src.Cols()) * scale)
	newH := int(float64(src.Rows()) * scale)
	dstSrc := gocv.NewMat()

	gocv.Resize(src, &dstSrc, image.Pt(newW, newH), 0, 0, gocv.InterpolationLinear)

	return dstSrc, scale
}

func resizeTpl(tpl gocv.Mat, scale float64) gocv.Mat {
	if tpl.Empty() || scale == 1 {
		return tpl
	}

	newW := int(float64(tpl.Cols()) * scale)
	newH := int(float64(tpl.Rows()) * scale)
	dstTpl := gocv.NewMat()

	gocv.Resize(tpl, &dstTpl, image.Pt(newW, newH), 0, 0, gocv.InterpolationLinear)

	return dstTpl
}

func restorePoint(pt image.Point, scale float64) image.Point {
	inv := 1 / scale

	pt.X = int(float64(pt.X) * inv)
	pt.Y = int(float64(pt.Y) * inv)

	return pt
}
