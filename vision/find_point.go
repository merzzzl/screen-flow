package vision

import (
	"image"
	"sort"

	"gocv.io/x/gocv"
)

type Algorithm int

const (
	AlgorithmSIFT = iota
	AlgorithmTM
)

func findPointSIFT(src, tpl gocv.Mat) (image.Point, []image.Point, bool) {
	if src.Empty() || tpl.Empty() {
		return image.Point{}, nil, false
	}

	srcGray := gocv.NewMat()
	tplGray := gocv.NewMat()

	defer srcGray.Close()
	defer tplGray.Close()

	if src.Channels() == 3 {
		gocv.CvtColor(src, &srcGray, gocv.ColorBGRToGray)
	} else {
		srcGray = src.Clone()
	}

	if tpl.Channels() == 3 {
		gocv.CvtColor(tpl, &tplGray, gocv.ColorBGRToGray)
	} else {
		tplGray = tpl.Clone()
	}

	sift := gocv.NewSIFT()
	defer sift.Close()

	kpSrc, descSrc := sift.DetectAndCompute(srcGray, gocv.NewMat())
	_, descTpl := sift.DetectAndCompute(tplGray, gocv.NewMat())

	if descSrc.Empty() || descTpl.Empty() {
		return image.Point{}, nil, false
	}

	matcher := gocv.NewBFMatcherWithParams(gocv.NormL2, false)
	defer matcher.Close()

	knn := matcher.KnnMatch(descTpl, descSrc, 2)

	good := make([]gocv.DMatch, 0, len(knn))

	for _, m := range knn {
		if len(m) == 2 && m[0].Distance < 0.75*m[1].Distance {
			good = append(good, m[0])
		}
	}

	if len(good) == 0 {
		return image.Point{}, nil, false
	}

	sort.Slice(good, func(i, j int) bool { return good[i].Distance < good[j].Distance })

	if len(good) > 50 {
		good = good[:50]
	}

	pts := make([]image.Point, len(good))

	for i, g := range good {
		p := kpSrc[g.TrainIdx]
		pts[i] = image.Pt(int(float64(p.X)), int(float64(p.Y)+0.5))
	}

	return pts[0], pts, true
}

func findPointTM(src, tpl gocv.Mat) (image.Point, []image.Point, bool) {
	if src.Empty() || tpl.Empty() {
		return image.Point{}, nil, false
	}

	srcGray := gocv.NewMat()
	tplGray := gocv.NewMat()

	defer srcGray.Close()
	defer tplGray.Close()

	if src.Channels() == 3 {
		gocv.CvtColor(src, &srcGray, gocv.ColorBGRToGray)
	} else {
		srcGray = src.Clone()
	}

	if tpl.Channels() == 3 {
		gocv.CvtColor(tpl, &tplGray, gocv.ColorBGRToGray)
	} else {
		tplGray = tpl.Clone()
	}

	res := gocv.NewMatWithSize(srcGray.Rows()-tplGray.Rows()+1, srcGray.Cols()-tplGray.Cols()+1, gocv.MatTypeCV32F)
	defer res.Close()

	gocv.MatchTemplate(srcGray, tplGray, &res, gocv.TmCcoeffNormed, gocv.NewMat())

	_, maxVal, _, maxLoc := gocv.MinMaxLoc(res)
	center := image.Point{X: maxLoc.X + tplGray.Cols()/2, Y: maxLoc.Y + tplGray.Rows()/2}

	return center, []image.Point{center}, maxVal > 0.75
}

func findPoint(src, tpl gocv.Mat, algo Algorithm) (image.Point, []image.Point, bool) {
	if src.Empty() || tpl.Empty() {
		return image.Point{}, nil, false
	}

	var fn func(src gocv.Mat, tpl gocv.Mat) (image.Point, []image.Point, bool)

	switch algo {
	case AlgorithmSIFT:
		fn = findPointSIFT
	case AlgorithmTM:
		fn = findPointTM
	default:
		return image.Point{}, nil, false
	}

	maxSide := src.Cols()

	if src.Rows() > maxSide {
		maxSide = src.Rows()
	}

	if maxSide <= 640 {
		return fn(src, tpl)
	}

	scale := float64(640) / float64(maxSide)

	newW := int(float64(src.Cols()) * scale)
	newH := int(float64(src.Rows()) * scale)

	dstSrc := gocv.NewMat()
	defer dstSrc.Close()

	gocv.Resize(src, &dstSrc, image.Pt(newW, newH), 0, 0, gocv.InterpolationLinear)

	tW := int(float64(tpl.Cols()) * scale)
	tH := int(float64(tpl.Rows()) * scale)

	dstTpl := gocv.NewMat()
	defer dstTpl.Close()

	gocv.Resize(tpl, &dstTpl, image.Pt(tW, tH), 0, 0, gocv.InterpolationLinear)

	pt, pts, ok := fn(dstSrc, dstTpl)
	if !ok {
		return image.Point{}, nil, false
	}

	inv := 1 / scale

	pt.X = int(float64(pt.X) * inv)
	pt.Y = int(float64(pt.Y) * inv)

	for i := range pts {
		pts[i].X = int(float64(pts[i].X) * inv)
		pts[i].Y = int(float64(pts[i].Y) * inv)
	}

	return pt, pts, true
}
