package ataxi

import "fmt"

type Coord struct {
	X int32
	Y int32
}

func (coord *Coord) String() string {
	return fmt.Sprintf("(%d, %d)", coord.X, coord.Y)
}
