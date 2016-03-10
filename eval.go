package gogo

import "fmt"

type SGDNode struct {
	pos, dis int
}

func (b *Board) StoneGrayDistance(color Color) []int {
	q := make([]SGDNode, 0, len(b.w)/2)
	ret := make([]int, len(b.w), len(b.w))
	for i, p := range b.w {
		if p.color == color {
			q = append(q, SGDNode{i, 0})
			ret[i] = 0
		} else {
			ret[i] = -1
		}
	}
	qstart := 0
	visited := NewBoardBitmap()
	for len(q) > qstart {
		nv := q[qstart]
		qstart += 1
		if visited.IsSet(nv.pos) {
			continue
		}
		visited.Set(nv.pos)
		ret[nv.pos] = nv.dis
		pf := b.w[nv.pos]
		n4 := b.Neighbor4Color(pf.x, pf.y, GRAY)
		for _, p := range n4 {
			pi := b.index(p.x, p.y)
			if visited.IsSet(pi) {
				continue
			}
			q = append(q, SGDNode{pi, nv.dis + 1})
		}
	}
	return ret
}

func (b *Board) Dilation(seeds map[int]int) map[int]int {
	expand := make(map[int]int)
	for i, _ := range seeds {
		p := b.w[i]
		n4 := b.Neighbor4Color(p.x, p.y, GRAY)
		for _, pn := range n4 {
			expand[b.index(pn.x, pn.y)] = 0
		}
	}
	for i, v := range seeds {
		expand[i] = v
	}

	for i, v := range expand {
		if v == 64 || v == -64 {
			continue
		}
		p := b.w[i]
		n4 := b.Neighbor4(p.x, p.y)
		add := 0
		minus := 0
		for _, pn := range n4 {
			pni := b.index(pn.x, pn.y)
			vn, _ := expand[pni]
			if vn > 0 {
				add += 1
			} else if vn < 0 {
				minus += 1
			}
		}
		if v >= 0 && add > 0 && minus == 0 {
			seeds[i] = v + add
		}
		if v <= 0 && minus > 0 && add == 0 {
			seeds[i] = v - minus
		}
	}
	return seeds
}

func (b *Board) Erase(seeds map[int]int) map[int]int {
	ret := make(map[int]int)
	for i, v := range seeds {
		p := b.w[i]
		n4 := b.Neighbor4(p.x, p.y)
		add := 0
		for _, pn := range n4 {
			pni := b.index(pn.x, pn.y)
			vn, _ := seeds[pni]
			if vn*v <= 0 {
				add += 1
			}
		}
		if v > 0 {
			v = v - add
			if v < 0 {
				v = 0
			}
		} else if v < 0 {
			v = v + add
			if v > 0 {
				v = 0
			}
		}
		if v != 0 {
			ret[i] = v
		}
	}
	return ret
}

func (b *Board) PrintWithValue(vals map[int]int, ch bool) {
	for y := b.size - 1; y >= 0; y-- {
		for x := 0; x < b.size; x++ {
			p := b.index(x, y)
			v, _ := vals[p]
			if !ch {
				fmt.Printf("%2.0f ", float64(v))
			} else {
				if v > 0 && v < 64 {
					fmt.Print("x ")
				} else if v >= 64 {
					fmt.Print("X ")
				} else if v == 0 {
					fmt.Print(". ")
				} else if v < 0 && v > -64 {
					fmt.Print("o ")
				} else {
					fmt.Print("O ")
				}
			}
		}
		fmt.Println()
	}
}

func (b *Board) Influence() map[int]int {
	seeds := make(map[int]int)
	for i, p := range b.w {
		if p.color == BLACK {
			seeds[i] = 64
		} else if p.color == WHITE {
			seeds[i] = -64
		}
	}

	for i := 0; i < 5; i++ {
		seeds = b.Dilation(seeds)
	}

	return seeds
}

func (b *Board) Territory() map[int]int {
	seeds := make(map[int]int)
	for i, p := range b.w {
		if p.color == BLACK {
			seeds[i] = 64
		} else if p.color == WHITE {
			seeds[i] = -64
		}
	}

	for i := 0; i < 5; i++ {
		seeds = b.Dilation(seeds)
	}

	for i := 0; i < 21; i++ {
		seeds = b.Erase(seeds)
	}
	return seeds
}
