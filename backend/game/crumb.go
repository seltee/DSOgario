package game

const crumbsPerSecond = 3

const TypeCrumb = 2

type Crumb struct {
	Position Position
	ID       uint32
	Size     uint16
	Radius   float64
}
