package vision

import (
	"fmt"
	"image"

	"gocv.io/x/gocv"
)

type object struct {
	raw  gocv.Mat
	sfit gocv.Mat
	kp   []gocv.KeyPoint
}

var sift = gocv.NewSIFT()
var matcher = gocv.NewBFMatcherWithParams(gocv.NormL2, false)

func newObject(img image.Image) (*object, error) {
	mat, err := toMat(img)
	if err != nil {
		return nil, fmt.Errorf("failed to convert: %w", err)
	}

	return &object{
		raw:  mat,
		sfit: gocv.NewMat(),
	}, nil
}

func toMat(img image.Image) (gocv.Mat, error) {
	bounds := img.Bounds()
	x := bounds.Dx()
	y := bounds.Dy()
	bytes := make([]byte, 0, x*y*3)

	for j := bounds.Min.Y; j < bounds.Max.Y; j++ {
		for i := bounds.Min.X; i < bounds.Max.X; i++ {
			r, g, b, _ := img.At(i, j).RGBA()
			bytes = append(bytes, byte(b>>8), byte(g>>8), byte(r>>8))
		}
	}

	rgb, err := gocv.NewMatFromBytes(y, x, gocv.MatTypeCV8UC3, bytes)
	if err != nil {
		return gocv.NewMat(), err
	}

	return rgb, nil
}

func (o *object) close() {
	_ = o.raw.Close()
	_ = o.sfit.Close()
	o.kp = nil
}

func (o *object) genSIFT() bool {
	if o.raw.Empty() {
		return false
	}

	if !o.sfit.Empty() {
		return true
	}

	kp, ds := sift.DetectAndCompute(o.raw, gocv.NewMat())
	o.sfit = ds.Clone()
	_ = ds.Close()
	o.kp = kp

	return true
}

func (o *object) findTopKMatchesInImage(template *object) []image.Point {
	if !o.genSIFT() || !template.genSIFT() {
		return []image.Point{}
	}

	matches := matcher.KnnMatch(template.sfit, o.sfit, 2)

	var good []gocv.DMatch

	for _, m := range matches {
		if len(m) == 2 && m[0].Distance < 0.75*m[1].Distance {
			good = append(good, m[0])
		}
	}

	if len(good) == 0 {
		return []image.Point{}
	}

	points := make([]image.Point, 0, len(good))

	for _, match := range good {
		pt := o.kp[match.TrainIdx]
		points = append(points, image.Pt(int(pt.X), int(pt.Y)))
	}

	return points
}

func (o *object) frameChangeRatio(next *object) float64 {
	if o.raw.Empty() || next.raw.Empty() {
		return 1
	}

	diff := gocv.NewMat()
	defer diff.Close()

	gocv.AbsDiff(o.raw, next.raw, &diff)

	gray := gocv.NewMat()
	defer gray.Close()

	gocv.CvtColor(diff, &gray, gocv.ColorBGRToGray)

	binary := gocv.NewMat()
	defer binary.Close()

	gocv.Threshold(gray, &binary, 25, 255, gocv.ThresholdBinary)

	changedPixels := gocv.CountNonZero(binary)
	totalPixels := binary.Rows() * binary.Cols()

	changeRatio := float64(changedPixels) / float64(totalPixels)

	return changeRatio
}
