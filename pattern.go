package gogo

import "fmt"

type PointDistanceMap struct {
	data [][][]int
}

func NewPointDistanceMap(b *Board, ds int) *PointDistanceMap {
	size := len(b.w)
	data := make([][][]int, size, size)
	for i := 0; i < size; i++ {
		data[i] = make([][]int, ds, ds)
		for d := 0; d < ds; d++ {
			data[i][d] = make([]int, 0, ds*4+4)
		}
		pi := b.w[i]
		for j := 0; j < size; j++ {
			pj := b.w[j]
			d := pi.Distance(pj)
			if d >= ds {
				continue
			}
			data[i][d] = append(data[i][d], j)
		}
	}
	return &PointDistanceMap{data}
}

func (p *PointDistanceMap) PointDistance(q, d int) []int {
	return p.data[q][d]
}

func (b *Board) ColorHash(p Point) int64 {
	if p.color == GRAY {
		return 4367440637090516
	} else if p.color == BLACK {
		return 7534508431030358
	} else if p.color == WHITE {
		return 3401485023080138
	} else {
		return 6231701002348503
	}
}

func (b *Board) PointLibertyHash(p Point) int64 {
	n := b.PointLiberty(p)
	if n == 0 {
		return 4206701385013581
	} else if n == 1 {
		return 7305135150381401
	} else if n == 2 {
		return 2763013808356791
	} else {
		return 4150754907001058
	}
}

func (b *Board) EdgeDisHash(p Point) int64 {
	n := p.EdgeDis(b.size)
	if n <= 3 {
		return 4395792752175943
	} else if n == 4 {
		return 5672153418394501
	} else if n == 5 {
		return 9798437591759191
	} else {
		return 5788787957979347
	}
}

func (b *Board) KoHash(p Point) int64 {
	if p.x == b.ko.x && p.y == b.ko.y {
		return 6739755730475491
	} else {
		return 9486794256091760
	}
}

func (b *Board) PointHash(p Point) int64 {
	ret := b.ColorHash(p)
	ret ^= b.PointLibertyHash(p)
	ret ^= b.KoHash(p)
	ret ^= b.EdgeDisHash(p)
	return ret
}

func (b *Board) InitPatternHash(pdm *PointDistanceMap) {
	for p, dh := range b.patternHash {
		for d := 0; d < len(dh); d++ {
			qd := pdm.PointDistance(p, d)
			for _, q := range qd {
				dh[d] ^= b.pointHash[q]
			}
		}
	}
}

func (b *Board) UpdateHash(x, y int, c Color, pdm *PointDistanceMap) {
	if pdm == nil {
		return
	}
	p := b.index(x, y)
	oh := b.pointHash[p]
	nh := b.PointHash(Point{x, y, c})
	for d := 0; d < PATTERN_SIZE; d++ {
		dm := pdm.PointDistance(p, d)
		for _, q := range dm {
			b.patternHash[q][d] ^= oh
			b.patternHash[q][d] ^= nh
		}
	}
	b.pointHash[p] = nh
}

func (b *Board) Pattern(x, y int, stone Color, ds int) int64 {
	r := []int64{0, 0, 0, 0, 0, 0, 0, 0}
	mut := [][]int{[]int{1, 1}, []int{-1, 1}, []int{1, -1}, []int{-1, -1}}
	for dy := -1 * ds; dy <= ds; dy++ {
		for dx := -1 * ds; dx <= ds; dx++ {
			for k := 0; k < len(mut); k++ {
				if dx == 0 && dy == 0 {
					continue
				}
				{
					r[k] *= 16
					x1 := x + dx*mut[k][0]
					y1 := y + dy*mut[k][1]
					p1 := b.Get(x1, y1)
					lb := b.PointLiberty(p1)
					if lb > 3 {
						lb = 3
					}
					if p1.Valid() {
						if p1.Color() == stone {
							r[k] += int64(1 + lb*4)
						} else if p1.Color() == OppColor(stone) {
							r[k] += int64(2 + lb*4)
						} else {
							r[k] += int64(3 + lb*4)
						}
					}
				}
				{
					r[k+4] *= 16
					x1 := x + dy*mut[k][0]
					y1 := y + dx*mut[k][1]
					p1 := b.Get(x1, y1)
					lb := b.PointLiberty(p1)
					if lb > 3 {
						lb = 3
					}
					if p1.Valid() {
						if p1.Color() == stone {
							r[k+4] += int64(1 + lb*4)
						} else if p1.Color() == OppColor(stone) {
							r[k+4] += int64(2 + lb*4)
						} else {
							r[k+4] += int64(3 + lb*4)
						}
					}
				}
			}
		}
	}
	ret := r[0]
	for i := 1; i < len(r); i++ {
		if ret > r[i] {
			ret = r[i]
		}
	}
	return ret
}

func (b *Board) PatternString(x, y int, stone Color, ds int) string {
	r := []string{"", "", "", "", "", "", "", ""}
	mut := [][]int{[]int{1, 1}, []int{-1, 1}, []int{1, -1}, []int{-1, -1}}
	for dy := -1 * ds; dy <= ds; dy++ {
		for dx := -1 * ds; dx <= ds; dx++ {
			for k := 0; k < len(mut); k++ {
				{
					x1 := x + dx*mut[k][0]
					y1 := y + dy*mut[k][1]
					p1 := b.Get(x1, y1)
					if p1.Valid() {
						if p1.Color() == stone {
							r[k] += "O"
						} else if p1.Color() == OppColor(stone) {
							r[k] += "X"
						} else {
							r[k] += "."
						}
					} else {
						r[k] += "|"
					}
				}
				{
					x1 := x + dy*mut[k][0]
					y1 := y + dx*mut[k][1]
					p1 := b.Get(x1, y1)
					if p1.Valid() {
						if p1.Color() == stone {
							r[k+4] += "O"
						} else if p1.Color() == OppColor(stone) {
							r[k+4] += "X"
						} else {
							r[k+4] += "."
						}
					} else {
						r[k+4] += "|"
					}
				}
			}
		}
	}
	ret := r[0]
	for i := 1; i < len(r); i++ {
		if ret > r[i] {
			ret = r[i]
		}
	}
	return ret
}

func PrintPattern(pat string) {
	n := 3
	if len(pat) == 25 {
		n = 5
	}
	for i, ch := range pat {
		if i%n == 0 {
			fmt.Println()
		}
		fmt.Print(string(ch))
		fmt.Print(" ")

	}
	fmt.Println()
}
