package ataxi

import (
	"math"

	geo "github.com/kellydunn/golang-geo"
)

func GetSuperPixel(x int32, y int32, category uint32) (int32, int32) {
	X := mapToSuperCoord(x, category)
	Y := mapToSuperCoord(y, category)
	return X, Y
}

func mapToSuperCoord(x int32, category uint32) int32 {
	return sign(x)*(2*int32(category)+1)*int32(math.Floor(math.Abs(float64(x))/float64(2*int32(category)+1))) + 2
}

func sign(x int32) int32 {
	if x < 0 {
		return -1
	} else if x > 0 {
		return 1
	}
	return 0
}

func GetMaxWaitingTime(dist float64) uint32 {
	if dist < 2 {
		return 300
	} else if dist < 10 {
		return 420
	} else if dist < 100 {
		return 600
	} else if dist < 400 {
		return 900
	} else {
		return 1800
	}
}

func GetTripDistance(latlon1 *geo.Point, latlon2 *geo.Point) float64 {
	return 1.2 * latlon1.GreatCircleDistance(latlon2) / 1.6
}

func GetTripCategory(gcDist float64) uint32 {
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
