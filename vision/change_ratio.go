package vision

import "gocv.io/x/gocv"

func calcChangeRatio(a, b gocv.Mat) float64 {
	diff := gocv.NewMat()
	defer diff.Close()

	gocv.AbsDiff(a, b, &diff)

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
