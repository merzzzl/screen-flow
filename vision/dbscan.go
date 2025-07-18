package vision

import (
	"image"
	"math"
)

func distance(a, b image.Point) float64 {
	dx := a.X - b.X
	dy := a.Y - b.Y

	return math.Sqrt(float64(dx*dx + dy*dy))
}

func regionQuery(points []image.Point, p image.Point, eps int) []int {
	var neighbors []int

	for i, other := range points {
		if distance(p, other) <= float64(eps) {
			neighbors = append(neighbors, i)
		}
	}

	return neighbors
}

func expandCluster(points []image.Point, labels []int, pointIdx, clusterID, eps, minPts int) bool {
	seeds := regionQuery(points, points[pointIdx], eps)

	if len(seeds) < minPts {
		labels[pointIdx] = -1

		return false
	}

	for _, idx := range seeds {
		labels[idx] = clusterID
	}

	seeds = remove(seeds, pointIdx)

	for i := 0; i < len(seeds); i++ {
		n := seeds[i]

		if labels[n] == -1 {
			labels[n] = clusterID
		}

		if labels[n] != 0 {
			continue
		}

		labels[n] = clusterID
		neighbors := regionQuery(points, points[n], eps)

		if len(neighbors) >= minPts {
			seeds = append(seeds, neighbors...)
		}
	}

	return true
}

func remove(slice []int, value int) []int {
	var out []int

	for _, v := range slice {
		if v != value {
			out = append(out, v)
		}
	}

	return out
}

func dbscan(points []image.Point, eps, minPts int) ([][]image.Point, []image.Point) {
	clusterID := 1
	labels := make([]int, len(points))

	for i := range points {
		if labels[i] != 0 {
			continue
		}

		if expandCluster(points, labels, i, clusterID, eps, minPts) {
			clusterID++
		}
	}

	clustersMap := make(map[int][]image.Point)

	var noise []image.Point

	for i, label := range labels {
		switch {
		case label == -1:
			noise = append(noise, points[i])
		case label > 0:
			clustersMap[label] = append(clustersMap[label], points[i])
		}
	}

	clusters := make([][]image.Point, 0, len(clustersMap))

	for _, cluster := range clustersMap {
		clusters = append(clusters, cluster)
	}

	return clusters, noise
}
