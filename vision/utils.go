package vision

import (
	"fmt"
	"image"
)

func IsFrameStatic(first, second image.Image, threshold float64) (bool, error) {
	firstObj, err := newObject(first)
	if err != nil {
		return false, fmt.Errorf("create first vision object: %w", err)
	}

	secondObj, err := newObject(second)
	if err != nil {
		return false, fmt.Errorf("create second vision object: %w", err)
	}

	change := firstObj.frameChangeRatio(secondObj)

	return change < threshold, nil
}

func FindPoint(source, template image.Image) (image.Point, bool, error) {
	sourceObj, err := newObject(source)
	if err != nil {
		return image.Pt(0, 0), false, fmt.Errorf("create source vision object: %w", err)
	}

	defer sourceObj.close()

	templateObj, err := newObject(template)
	if err != nil {
		return image.Pt(0, 0), false, fmt.Errorf("create template vision object: %w", err)
	}

	defer templateObj.close()

	points := sourceObj.findTopKMatchesInImage(templateObj)

	if len(points) == 0 {
		return image.Pt(0, 0), false, nil
	}

	clusters, noise := dbscan(points, 128, 3)

	if len(clusters) == 0 {
		return image.Pt(0, 0), false, nil
	}

	abs := image.Pt(0, 0)

	for _, point := range clusters[0] {
		abs.X += point.X
		abs.Y += point.Y
	}

	abs.X /= len(clusters[0])
	abs.Y /= len(clusters[0])

	if len(noise)*2 > len(clusters[0]) {
		return image.Pt(0, 0), false, nil
	}

	return abs, true, nil
}

// func saveTempResult(sourceObj *object, good, bad []image.Point, avg image.Point) error {
// 	mat := sourceObj.getMat()
// 	name := "_temp.jpeg"

// 	for _, point := range good {
// 		drowPoint(&mat, point, color.RGBA{255, 0, 255, 32})
// 	}

// 	for _, point := range bad {
// 		drowPoint(&mat, point, color.RGBA{0, 0, 255, 32})
// 	}

// 	drowPoint(&mat, avg, color.RGBA{0, 255, 0, 255})

// 	return saveMat(&mat, name)
// }

// func drowPoint(mat *gocv.Mat, point image.Point, cl color.RGBA) {
// 	if point.X > 0 && point.Y > 0 {
// 		startH := image.Pt(point.X-20, point.Y)
// 		endH := image.Pt(point.X+20, point.Y)

// 		_ = gocv.Line(mat, startH, endH, cl, 4)

// 		startV := image.Pt(point.X, point.Y-20)
// 		endV := image.Pt(point.X, point.Y+20)

// 		_ = gocv.Line(mat, startV, endV, cl, 4)
// 	}
// }

// func saveMat(mat *gocv.Mat, path string) error {
// 	img, err := mat.ToImage()
// 	if err != nil {
// 		return fmt.Errorf("convert mat: %w", err)
// 	}

// 	out, err := os.Create(path)
// 	if err != nil {
// 		return fmt.Errorf("create file: %w", err)
// 	}

// 	err = jpeg.Encode(out, img, &jpeg.Options{Quality: 95})
// 	if err != nil {
// 		return fmt.Errorf("encode image: %w", err)
// 	}

// 	if err := out.Close(); err != nil {
// 		return fmt.Errorf("close file: %w", err)
// 	}

// 	return nil
// }
