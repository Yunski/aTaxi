package ataxi

import (
	geo "github.com/kellydunn/golang-geo"
)

func GetSuperPixel(x int32, y int32, category uint32) (int32, int32) {
	var n int32
	switch category {
	case 0:
		n = 2
	case 1:
		n = 2
	case 2:
		n = 3
	case 3:
		n = 5
	case 4:
		n = 10
	default:
		panic("Unexpected trip category")
	}
	X := mapToSuperCoord(x, n)
	Y := mapToSuperCoord(y, n)
	return X, Y
}

func mapToSuperCoord(x int32, n int32) int32 {
	if x < 0 {
		return -1 + ((x+1)/n)*n
	}
	return (x / n) * n
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

func HashCode(x int32, y int32) int32 {
	x = bijection(x)
	y = bijection(y)
	if x >= y {
		return x*x + x + y
	}
	return x + y*y
}

func bijection(x int32) int32 {
	if x < 0 {
		return -x*2 - 1
	}
	return x * 2
}
