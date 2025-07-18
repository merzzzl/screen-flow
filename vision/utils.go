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

	if len(points) < 5 {
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
