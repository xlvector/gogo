package gogo

import "math/rand"

func (b *PointMap) Set(k int) {
	a := uint64(k / 64)
	r := uint64(k % 64)
	b.n[a] = b.n[a] | (1 << r)
}

func (b *PointMap) IsSet(k int) bool {
	a, r := uint64(k/64), uint64(k%64)
	return (b.n[a] & (1 << r)) > 0
}

type PointMap struct {
	n      []uint64
	Points []int
}

func ZeroPointMap() PointMap {
	return PointMap{
		n:      []uint64{0, 0, 0, 0, 0, 0},
		Points: make([]int, 0, 5),
	}
}

func (p *PointMap) Random() int {
	if len(p.Points) == 0 {
		return -1
	}
	return p.Points[rand.Intn(len(p.Points))]
}

func (p *PointMap) Add(k int) {
	if p.IsSet(k) {
		return
	}
	p.Points = append(p.Points, k)
	p.Set(k)
}

func (p *PointMap) Exist(k int) bool {
	return p.IsSet(k)
}

func (p *PointMap) Size() int {
	return len(p.Points)
}

func (p *PointMap) First() int {
	return p.Points[0]
}
