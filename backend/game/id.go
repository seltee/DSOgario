package game

var NextID uint32 = 1

func GetNextID() uint32 {
	out := NextID
	NextID++
	return out
}
