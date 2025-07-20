package vision

import (
	"gocv.io/x/gocv"
	"gocv.io/x/gocv/contrib"
)

func dSIFT(src, tpl gocv.Mat) (kpSrc, kpTpl []gocv.KeyPoint, descSrc, descTpl gocv.Mat) {
	sift := gocv.NewSIFT()
	defer sift.Close()

	kpSrc, descSrc = sift.DetectAndCompute(src, gocv.NewMat())
	kpTpl, descTpl = sift.DetectAndCompute(tpl, gocv.NewMat())

	return kpSrc, kpTpl, descSrc, descTpl
}

func dORB(srcGray, tplGray gocv.Mat) (kpSrc, kpTpl []gocv.KeyPoint, descSrc, descTpl gocv.Mat) {
	orb := gocv.NewORBWithParams(5000, 1.2, 12, 31, 0, 2, gocv.ORBScoreTypeHarris, 31, 10)
	defer orb.Close()

	kpSrc, descSrc = orb.DetectAndCompute(srcGray, gocv.NewMat())
	kpTpl, descTpl = orb.DetectAndCompute(tplGray, gocv.NewMat())

	return kpSrc, kpTpl, descSrc, descTpl
}

func dAKAZE(srcGray, tplGray gocv.Mat) (kpSrc, kpTpl []gocv.KeyPoint, descSrc, descTpl gocv.Mat) {
	akaze := gocv.NewAKAZE()
	defer akaze.Close()

	kpSrc, descSrc = akaze.DetectAndCompute(srcGray, gocv.NewMat())
	kpTpl, descTpl = akaze.DetectAndCompute(tplGray, gocv.NewMat())

	return kpSrc, kpTpl, descSrc, descTpl
}

func dBRISK(srcGray, tplGray gocv.Mat) (kpSrc, kpTpl []gocv.KeyPoint, descSrc, descTpl gocv.Mat) {
	brisk := gocv.NewBRISK()
	defer brisk.Close()

	kpSrc, descSrc = brisk.DetectAndCompute(srcGray, gocv.NewMat())
	kpTpl, descTpl = brisk.DetectAndCompute(tplGray, gocv.NewMat())

	return kpSrc, kpTpl, descSrc, descTpl
}

func dFAST(srcGray, tplGray gocv.Mat) (kpSrc, kpTpl []gocv.KeyPoint, descSrc, descTpl gocv.Mat) {
	fast := gocv.NewFastFeatureDetector()
	defer fast.Close()

	kpSrc = fast.Detect(srcGray)
	kpTpl = fast.Detect(tplGray)

	orb := gocv.NewORB()
	defer orb.Close()

	_, descSrc = orb.Compute(srcGray, gocv.NewMat(), kpSrc)
	_, descTpl = orb.Compute(tplGray, gocv.NewMat(), kpTpl)

	return kpSrc, kpTpl, descSrc, descTpl
}

func dKAZE(srcGray, tplGray gocv.Mat) (kpSrc, kpTpl []gocv.KeyPoint, descSrc, descTpl gocv.Mat) {
	kaze := gocv.NewKAZE()
	defer kaze.Close()

	kpSrc, descSrc = kaze.DetectAndCompute(srcGray, gocv.NewMat())
	kpTpl, descTpl = kaze.DetectAndCompute(tplGray, gocv.NewMat())

	return kpSrc, kpTpl, descSrc, descTpl
}

func dSURF(srcGray, tplGray gocv.Mat) (kpSrc, kpTpl []gocv.KeyPoint, descSrc, descTpl gocv.Mat) {
	surf := contrib.NewSURF()
	defer surf.Close()

	kpSrc, descSrc = surf.DetectAndCompute(srcGray, gocv.NewMat())
	kpTpl, descTpl = surf.DetectAndCompute(tplGray, gocv.NewMat())

	return kpSrc, kpTpl, descSrc, descTpl
}

func dAGAST(srcGray, tplGray gocv.Mat) (kpSrc, kpTpl []gocv.KeyPoint, descSrc, descTpl gocv.Mat) {
	agast := gocv.NewAgastFeatureDetector()
	defer agast.Close()

	kpSrc = agast.Detect(srcGray)
	kpTpl = agast.Detect(tplGray)

	orbDesc := gocv.NewORB()
	defer orbDesc.Close()

	_, descSrc = orbDesc.Compute(srcGray, gocv.NewMat(), kpSrc)
	_, descTpl = orbDesc.Compute(tplGray, gocv.NewMat(), kpTpl)

	return kpSrc, kpTpl, descSrc, descTpl
}

func dGFTT(srcGray, tplGray gocv.Mat) (kpSrc, kpTpl []gocv.KeyPoint, descSrc, descTpl gocv.Mat) {
	gftt := gocv.NewGFTTDetector()
	defer gftt.Close()

	kpSrc = gftt.Detect(srcGray)
	kpTpl = gftt.Detect(tplGray)

	orbDesc := gocv.NewORB()
	defer orbDesc.Close()

	_, descSrc = orbDesc.Compute(srcGray, gocv.NewMat(), kpSrc)
	_, descTpl = orbDesc.Compute(tplGray, gocv.NewMat(), kpTpl)

	return kpSrc, kpTpl, descSrc, descTpl
}

func dBRIEF(srcGray, tplGray gocv.Mat) (kpSrc, kpTpl []gocv.KeyPoint, descSrc, descTpl gocv.Mat) {
	fastDet := gocv.NewFastFeatureDetector()
	defer fastDet.Close()

	kpSrc = fastDet.Detect(srcGray)
	kpTpl = fastDet.Detect(tplGray)

	briefDesc := contrib.NewBriefDescriptorExtractor()
	defer briefDesc.Close()

	descSrc = briefDesc.Compute(kpSrc, srcGray)
	descTpl = briefDesc.Compute(kpTpl, tplGray)

	return kpSrc, kpTpl, descSrc, descTpl
}
