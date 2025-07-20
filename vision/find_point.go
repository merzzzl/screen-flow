package vision

import (
	"image"
	"sort"
	"time"

	"gocv.io/x/gocv"
)

type Algorithm int

const (
	AlgorithmSIFT = iota
	AlgorithmTM
	AlgorithmORB
	AlgorithmAKAZE
	AlgorithmBRISK
	AlgorithmFAST
	AlgorithmKAZE
	AlgorithmSURF
	AlgorithmAGAST
	AlgorithmGFTT
	AlgorithmBRIEF
)

type Result struct {
	matches []gocv.DMatch
	dur     time.Duration
	algo    Algorithm
	src     gocv.Mat
	tpl     gocv.Mat
	kpSrc   []gocv.KeyPoint
	kpTpl   []gocv.KeyPoint
	descSrc gocv.Mat
	descTpl gocv.Mat
	best    image.Point
}

func findPoint(src, tpl gocv.Mat, algo Algorithm) (*Result, bool) {
	if src.Empty() || tpl.Empty() {
		return nil, false
	}

	startTime := time.Now()

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

	var (
		kpSrc, kpTpl     []gocv.KeyPoint
		descSrc, descTpl gocv.Mat
	)

	switch algo {
	case AlgorithmTM:
		pt, ok := findPointTM(srcGray, tplGray)

		return &Result{
			algo:    algo,
			best:    pt,
			src:     src.Clone(),
			tpl:     tpl.Clone(),
			kpSrc:   []gocv.KeyPoint{},
			kpTpl:   []gocv.KeyPoint{},
			descSrc: gocv.NewMat(),
			descTpl: gocv.NewMat(),
			dur:     time.Since(startTime),
			matches: []gocv.DMatch{},
		}, ok
	case AlgorithmSIFT:
		kpSrc, kpTpl, descSrc, descTpl = dSIFT(srcGray, tplGray)
	case AlgorithmORB:
		kpSrc, kpTpl, descSrc, descTpl = dORB(srcGray, tplGray)
	case AlgorithmAKAZE:
		kpSrc, kpTpl, descSrc, descTpl = dAKAZE(srcGray, tplGray)
	case AlgorithmBRISK:
		kpSrc, kpTpl, descSrc, descTpl = dBRISK(srcGray, tplGray)
	case AlgorithmFAST:
		kpSrc, kpTpl, descSrc, descTpl = dFAST(srcGray, tplGray)
	case AlgorithmKAZE:
		kpSrc, kpTpl, descSrc, descTpl = dKAZE(srcGray, tplGray)
	case AlgorithmSURF:
		kpSrc, kpTpl, descSrc, descTpl = dSURF(srcGray, tplGray)
	case AlgorithmAGAST:
		kpSrc, kpTpl, descSrc, descTpl = dAGAST(srcGray, tplGray)
	case AlgorithmGFTT:
		kpSrc, kpTpl, descSrc, descTpl = dGFTT(srcGray, tplGray)
	case AlgorithmBRIEF:
		kpSrc, kpTpl, descSrc, descTpl = dBRIEF(srcGray, tplGray)
	}

	defer descSrc.Close()
	defer descTpl.Close()

	if descSrc.Empty() || descTpl.Empty() {
		return nil, false
	}

	matcher := gocv.NewBFMatcherWithParams(gocv.NormL2, false)
	defer matcher.Close()

	knn := matcher.KnnMatch(descTpl, descSrc, 2)
	_ = kpTpl
	good := make([]gocv.DMatch, 0, len(knn))

	for _, m := range knn {
		if len(m) == 2 && m[0].Distance < 0.75*m[1].Distance {
			good = append(good, m[0])
		}
	}

	if len(good) == 0 {
		return nil, false
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

	return &Result{
		algo:    algo,
		src:     src.Clone(),
		tpl:     tpl.Clone(),
		kpSrc:   kpSrc,
		kpTpl:   kpTpl,
		descSrc: descSrc.Clone(),
		descTpl: descTpl.Clone(),
		dur:     time.Since(startTime),
		matches: good,
		best:    pts[0],
	}, true
}

func findPointTM(srcGray, tplGray gocv.Mat) (image.Point, bool) {
	res := gocv.NewMatWithSize(srcGray.Rows()-tplGray.Rows()+1, srcGray.Cols()-tplGray.Cols()+1, gocv.MatTypeCV32F)
	defer res.Close()

	gocv.MatchTemplate(srcGray, tplGray, &res, gocv.TmCcoeffNormed, gocv.NewMat())

	_, maxVal, _, maxLoc := gocv.MinMaxLoc(res)
	center := image.Point{X: maxLoc.X + tplGray.Cols()/2, Y: maxLoc.Y + tplGray.Rows()/2}

	return center, maxVal > 0.75
}

func (r *Result) Close() {
	_ = r.descSrc.Close()
	_ = r.descTpl.Close()
	_ = r.src.Close()
	_ = r.tpl.Close()
}
