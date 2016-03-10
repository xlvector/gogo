package gogo

type BoardBitmap struct {
	n []uint64
}

func NewBoardBitmap() *BoardBitmap {
	return &BoardBitmap{
		n : []uint64{0,0,0,0,0,0},
	}
}

func (b *BoardBitmap) Set(k int) {
	a := uint64(k / 64)
	r := uint64(k % 64)
	b.n[a] = b.n[a] | (1 << r)
}

func (b *BoardBitmap) IsSet(k int) bool {
	a, r := uint64(k / 64), uint64(k % 64)
	return (b.n[a] & (1 << r)) > 0
}

type PointMap struct {
	bitmap *BoardBitmap
	Points []int
}

func NewPointMap(capacity int) *PointMap {
	return &PointMap{
		bitmap: NewBoardBitmap(),
		Points: make([]int, 0, capacity),
	}
}

func (p *PointMap) Add(k int) {
	if p.bitmap.IsSet(k) {
		return
	}
	p.Points = append(p.Points, k)
	p.bitmap.Set(k)
}

func (p *PointMap) Exist(k int) bool {
	return p.bitmap.IsSet(k)
}

func (p *PointMap) Size() int {
	return len(p.Points)
}

func (p *PointMap) First() int {
	return p.Points[0]
}