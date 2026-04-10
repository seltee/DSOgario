package game

import "math/rand/v2"

const Precision = 100.0

type Position struct {
	X float64
	Y float64
}

func randRange(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func sizeToRadius(size uint16) float64 {
	return float64(size+10) * 0.02
}

func sizeToSpeed(size int) float64 {
	return 100.0 / float64(size)
}

func genFieldPosition(size float64) Position {
	return Position{
		X: randRange(-size, size),
		Y: randRange(-size, size),
	}
}
