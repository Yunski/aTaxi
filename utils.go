package main

import (
	"math"

	geo "github.com/kellydunn/golang-geo"
)

func getDistConstraint(dist float64) float64 {
	if dist < 2.0 {
		return 1.0
	} else if dist < 10.0 {
		return 1.5
	} else if dist < 100.0 {
		return 2.5
	} else {
		return 5.0
	}
}

func getSuperPixel(coord Coord, category int) Coord {
	x := mapToSuperCoord(coord.X, category)
	y := mapToSuperCoord(coord.Y, category)
	return Coord{X: x, Y: y}
}

func mapToSuperCoord(x int, category int) int {
	return sign(x)*(2*category+1)*int(math.Floor(math.Abs(float64(x))/float64(2*category+1))) + 2
}

func sign(x int) int {
	if x < 0 {
		return -1
	} else if x > 0 {
		return 1
	}
	return 0
}

func getMaxWaitingTime(dist float64) int {
	if dist < 2 {
		return 300
	} else if dist < 10 {
		return 420
	} else if dist < 100 {
		return 600
	} else if dist < 400 {
		return 720
	} else {
		return 900
	}
}

func getTripDistance(latlon1 *geo.Point, latlon2 *geo.Point) float64 {
	return 1.2 * latlon1.GreatCircleDistance(latlon2) / 1.6
}

func getTripCategory(gcDist float64) int {
	if gcDist < 0.5 {
		return 0
	} else if gcDist < 10 {
		return 1
	} else if gcDist < 100 {
		return 2
	} else if gcDist < 400 {
		return 3
	}
	return 4
}

type Coord struct {
	X int
	Y int
}
