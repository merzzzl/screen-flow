package vision

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"gocv.io/x/gocv"
)

func (p *Pipe) showImage(obj gocv.Mat, res *Result) {
	var out gocv.Mat

	if res != nil {
		out = visualize(res)
	} else {
		srcW, srcH := obj.Cols(), obj.Rows()
		out = gocv.NewMatWithSize(srcH, srcW*2, gocv.MatTypeCV8UC3)
		roi := out.Region(image.Rect(srcW, 0, srcW*2, srcH))
		obj.CopyTo(&roi)
		roi.Close()
	}

	p.window.Resize(out.Cols(), out.Rows())
	p.window.Show(out)
}

func visualize(res *Result) gocv.Mat {
	if res != nil && res.algo == AlgorithmTM {
		res.src = heatmap(res.src, res.tpl)
	}

	tplW, tplH := 0, 0

	if !res.tpl.Empty() {
		tplW, tplH = res.tpl.Cols(), res.tpl.Rows()
	}

	srcW, srcH := res.src.Cols(), res.src.Rows()
	outH := int(math.Max(float64(srcH), float64(tplH)))
	outW := srcW * 2
	out := gocv.NewMatWithSize(outH, outW, gocv.MatTypeCV8UC3)

	if tplW > 0 && tplH > 0 {
		tplX := (srcW - tplW) / 2
		tplY := (outH - tplH) / 2
		tplROI := out.Region(image.Rect(tplX, tplY, tplX+tplW, tplY+tplH))

		res.tpl.CopyTo(&tplROI)
		tplROI.Close()

		red := color.RGBA{R: 255, A: 255}

		gocv.Rectangle(&out, image.Rect(tplX, tplY, tplX+tplW, tplY+tplH), red, 2)

		txt := fmt.Sprintf("%dms", res.dur.Milliseconds())
		pos := image.Point{X: tplX, Y: tplY - 10}

		gocv.PutText(&out, txt, pos, gocv.FontHersheySimplex, 0.4, color.RGBA{R: 255, G: 255, A: 255}, 2)
	}

	srcROI := out.Region(image.Rect(srcW, 0, srcW*2, srcH))

	res.src.CopyTo(&srcROI)
	srcROI.Close()

	if len(res.matches) > 0 {
		green := color.RGBA{G: 255, A: 255}

		for _, m := range res.matches {
			ps := res.kpSrc[m.TrainIdx]
			xS := int(ps.X) + srcW
			yS := int(ps.Y)
			pt := res.kpTpl[m.QueryIdx]
			xT := int(pt.X) + (srcW-tplW)/2
			yT := int(pt.Y) + (outH-tplH)/2
			gocv.Line(&out, image.Pt(xT, yT), image.Pt(xS, yS), green, 1)
		}
	}

	if res.best.X != 0 || res.best.Y != 0 {
		w, h := tplW, tplH
		topLeft := image.Point{X: res.best.X - w/2 + srcW, Y: res.best.Y - h/2}
		bottomRight := image.Point{X: topLeft.X + w, Y: topLeft.Y + h}
		red := color.RGBA{R: 255, A: 255}

		gocv.Rectangle(&out, image.Rect(topLeft.X, topLeft.Y, bottomRight.X, bottomRight.Y), red, 2)
		gocv.Circle(&out, image.Point{X: res.best.X + srcW, Y: res.best.Y}, 4, color.RGBA{G: 255, A: 255}, -1)
	}

	return out
}

func heatmap(src, tpl gocv.Mat) gocv.Mat {
	if src.Empty() || tpl.Empty() {
		return gocv.NewMat()
	}

	srcRows, srcCols := src.Rows(), src.Cols()
	tplRows, tplCols := tpl.Rows(), tpl.Cols()

	if srcRows < tplRows || srcCols < tplCols {
		return gocv.NewMat()
	}

	corr := gocv.NewMatWithSize(srcRows-tplRows+1, srcCols-tplCols+1, gocv.MatTypeCV32F)
	srcGray := gocv.NewMat()
	tplGray := gocv.NewMat()

	gocv.CvtColor(src, &srcGray, gocv.ColorBGRToGray)
	gocv.CvtColor(tpl, &tplGray, gocv.ColorBGRToGray)
	gocv.MatchTemplate(srcGray, tplGray, &corr, gocv.TmCcoeffNormed, gocv.NewMat())
	gocv.Normalize(corr, &corr, 0, 255, gocv.NormMinMax)
	corr.ConvertTo(&corr, gocv.MatTypeCV8U)

	heatSmall := gocv.NewMat()
	gocv.ApplyColorMap(corr, &heatSmall, gocv.ColormapJet)

	heatFull := gocv.NewMatWithSize(srcRows, srcCols, heatSmall.Type())
	heatFull.SetTo(gocv.Scalar{Val1: 0, Val2: 0, Val3: 0, Val4: 0})

	roiRect := image.Rect(tplCols/2, tplRows/2, tplCols/2+heatSmall.Cols(), tplRows/2+heatSmall.Rows())
	roi := heatFull.Region(roiRect)
	heatSmall.CopyTo(&roi)
	roi.Close()

	out := gocv.NewMat()
	gocv.AddWeighted(src, 0.4, heatFull, 0.6, 0, &out)

	corr.Close()
	srcGray.Close()
	tplGray.Close()
	heatSmall.Close()
	heatFull.Close()

	return out
}
