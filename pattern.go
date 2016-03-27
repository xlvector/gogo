package gogo

import (
	"sort"
)

func ExtendPosIndex(x, y int) int {
	return (y+SIZE)*(SIZE*3) + x + SIZE
}

func ExtendIndex(k int) int {
	x, y := IndexPos(k)
	return ExtendPosIndex(x, y)
}

func (b *Board) ColorHash(c Color) int64 {
	if c == BLACK {
		return 27620385010153
	} else if c == WHITE {
		return 93457175051051
	} else if c == GRAY {
		return 14137501853429
	}
	return 78017175060911
}

func (b *Board) VertexHash(x, y int, c Color) int64 {
	return (int64(ExtendPosIndex(x, y)+17)*13481917391 + 18223) ^ b.ColorHash(c)
}

func (b *Board) EdgeDisHash(k int) int64 {
	d := IndexEdgeDis(k)
	if d == 0 {
		return 34325608508451
	} else if d == 1 {
		return 73492759157915
	} else if d == 2 {
		return 89832645519571
	} else if d == 3 {
		return 71481275917501
	} else if d == 4 {
		return 91347170185103
	} else if d == 5 {
		return 54183651475917
	} else {
		return 43051045804350
	}
}

func (b *Board) SelfWormHash(k int, c Color) int64 {
	worm := b.WormFromPoint(k, c, 2)
	if worm.Liberty == 0 {
		return 34025815894375
	} else if worm.Liberty == 1 {
		return 93705801385943
	} else if worm.Liberty == 2 {
		return 70984305897151
	} else {
		return 68169175401475
	}
}

func (b *Board) OpWormHash(k int, c Color) int64 {
	nworms := b.NeighWorms(k, c, OpColor(c), 2)
	if len(nworms) == 0 {
		return 0
	}
	minLiberty := 10000
	for _, w := range nworms {
		if minLiberty > w.Liberty {
			minLiberty = w.Liberty
		}
	}
	if minLiberty == 0 {
		return 78901815894375
	} else if minLiberty == 1 {
		return 83705875206943
	} else if minLiberty == 2 {
		return 50091585897151
	} else {
		return 17169175589245
	}
}

func (b *Board) PointLiberty(k int) int {
	n4 := Neigh4(k)
	ret := 0
	for _, nk := range n4 {
		if b.Points[nk] == GRAY {
			ret += 1
		}
	}
	return ret
}

func (b *Board) EscapeAtari(k int, c Color) bool {
	nworms := b.NeighWorms(k, c, c, 2)
	minLiberty := 10000
	for _, w := range nworms {
		if minLiberty > w.Liberty {
			minLiberty = w.Liberty
		}
	}
	if minLiberty > 2 {
		return false
	} else {
		worm := b.WormFromPoint(k, c, 2)
		if worm.Liberty > 2 {
			return true
		}
		return false
	}
}

func (b *Board) Urgency(k int, c Color) float64 {
	if b.PointLiberty(k) > 3 {
		return 1.0
	}
	myNWorms := b.NeighWorms(k, c, c, 3)
	opNWorms := b.NeighWorms(k, c, OpColor(c), 3)
	worm := b.WormFromPoint(k, c, 3)
	minLiberty := 10000
	for _, w := range myNWorms {
		if minLiberty > w.Liberty {
			minLiberty = w.Liberty
		}
	}

	if minLiberty == 1 {
		if worm.Liberty == 2 {
			return 10.0
		} else if worm.Liberty >= 3 {
			return 100.0
		}
	} else if minLiberty == 2 {
		if worm.Liberty >= 3 {
			return 20.0
		}
	}

	minLiberty = 10000
	for _, w := range opNWorms {
		if minLiberty > w.Liberty {
			minLiberty = w.Liberty
		}
	}

	if minLiberty == 1 {
		return 50.0
	} else if minLiberty == 2 {
		if worm.Liberty == 1 {
			return 0.1
		}
		return 20.0
	}

	if worm.Liberty == 1 {
		return 0.1
	}

	return 1.0
}

func (b *Board) LocalFeature(k int, c Color) []int64 {
	myNWorms := b.NeighWorms(k, c, c, 3)
	opNWorms := b.NeighWorms(k, c, OpColor(c), 3)
	worm := b.WormFromPoint(k, c, 3)
	ret := make([]int64, 0, 5)

	minLiberty := 10000
	for _, w := range myNWorms {
		if minLiberty > w.Liberty {
			minLiberty = w.Liberty
		}
	}

	if minLiberty == 1 {
		if worm.Liberty == 1 {
			ret = append(ret, 493570158105)
		} else if worm.Liberty == 2 {
			ret = append(ret, 159084081432)
		} else if worm.Liberty == 3 {
			ret = append(ret, 897325971018)
		} else if worm.Liberty > 3 {
			ret = append(ret, 291850148415)
		}
	} else if minLiberty == 2 {
		if worm.Liberty == 1 {
			ret = append(ret, 932759347016)
		} else if worm.Liberty == 2 {
			ret = append(ret, 758724359874)
		} else if worm.Liberty == 3 {
			ret = append(ret, 238146923179)
		} else if worm.Liberty > 3 {
			ret = append(ret, 945876927621)
		}
	}

	//op
	var minLibertyWorm *Worm
	liberties := []int{}
	minLibertySize := 0
	for _, w := range opNWorms {
		if minLibertyWorm == nil {
			minLibertyWorm = w
		} else {
			if minLibertyWorm.Liberty > w.Liberty {
				minLibertyWorm = w
			}
		}
		liberties = append(liberties, w.Liberty)
	}
	sort.Ints(liberties)
	fl := int64(0)
	for _, l := range liberties {
		fl *= 5
		fl += int64(l + 1)
	}
	fl += 809438508012
	ret = append(ret, fl)
	if minLibertyWorm != nil {
		if minLibertyWorm.Liberty == 1 {
			ret = append(ret, 787401927621+int64(minLibertySize))
		} else if minLibertyWorm.Liberty == 2 {
			if worm.Liberty == 1 {
				ret = append(ret, 304580158101)
			} else if worm.Liberty == 2 {
				ret = append(ret, 843759137519)
				ret = append(ret, 843759137519+int64(minLibertySize))
			} else if worm.Liberty >= 3 {
				ret = append(ret, 934571349579)
				ret = append(ret, 934571349579+int64(minLibertySize))
			}
		}

		if minLibertyWorm.Liberty == 2 && worm.Liberty > 1 {
			opMyWorms := b.WormNeighWorms(minLibertyWorm, c, 2)
			ml := 10000
			for _, w := range opMyWorms {
				if ml > w.Liberty {
					ml = w.Liberty
				}
			}
			if ml == 2 {
				ret = append(ret, 148759154791191)
			}
		}
	}

	ret = append(ret, b.EdgeDisHash(k))

	nMy := int64(0)
	nOp := int64(0)
	for _, p := range PointDisMap[k][1] {
		if b.Points[p] == c {
			nMy += 1
		} else if b.Points[p] == OpColor(c) {
			nOp += 1
		}
	}

	ret = append(ret, 257012801851+nMy*710517801)
	ret = append(ret, 314501851003+nOp*837104719)

	nMy = 0
	nOp = 0
	for _, p := range PointDisMap[k][2] {
		if b.Points[p] == c {
			nMy += 1
		} else if b.Points[p] == OpColor(c) {
			nOp += 1
		}
	}
	ret = append(ret, 839457015011+nMy*834561911)
	ret = append(ret, 954876687874+nOp*231971295)

	wormHash := b.EmptyWormFromPoint(k, 5)
	for _, h := range wormHash {
		ret = append(ret, h^78860975057501)
	}
	ret = append(ret, b.PatternDxd(k, c, 1))

	ret = append(ret, 970803460911+751*int64(len(b.Actions)/10))

	/*
		influ := b.InfluenceVal[k]
		if c == WHITE {
			influ *= -1
		}
		ret = append(ret, 491570147501+11*int64(influ))
	*/

	fret := make([]int64, 0, 3*len(ret))
	for _, f1 := range ret {
		fret = append(fret, f1)
		for _, f2 := range ret {
			if f1 < f2 {
				fret = append(fret, f1^f2)
				for _, f3 := range ret {
					if f2 < f3 {
						fret = append(fret, f1^f2^f3)
					}
				}
			}
		}
	}

	return ret
}

func (b *Board) RotateNeigh(x, y, dx, dy, r int) (int, int) {
	if r == 0 {
		return x + dx, y + dy
	} else if r == 1 {
		return x - dx, y + dy
	} else if r == 2 {
		return x + dx, y - dy
	} else if r == 3 {
		return x - dx, y - dy
	} else if r == 4 {
		return y + dy, x + dx
	} else if r == 5 {
		return y + dy, x - dx
	} else if r == 6 {
		return y - dy, x + dx
	} else {
		return y - dy, x - dx
	}
}

func (b *Board) PatternDxd(p int, c Color, d int) int64 {
	x, y := IndexPos(p)
	ret := int64(0)
	for r := 0; r < 8; r++ {
		f := int64(0)
		for dy := -1 * d; dy <= d; dy++ {
			for dx := -1 * d; dx <= d; dx++ {
				f *= 20
				x1, y1 := b.RotateNeigh(x, y, dx, dy, r)
				c1 := INVALID_COLOR
				if !PosOutBoard(x1, y1) {
					c1 = b.Points[PosIndex(x1, y1)]
				}
				pl := b.PointLiberty(PosIndex(x1, y1))
				if c1 == c {
					f += int64(pl * 4)
				} else if c1 == OpColor(c) {
					f += int64(pl*4 + 1)
				} else if c1 == GRAY {
					f += int64(pl*4 + 2)
				} else {
					f += int64(pl*4 + 3)
				}
			}
		}
		if ret < f {
			ret = f
		}
	}
	ret += 90000000000000
	return ret
}

func (b *Board) NeighWorms(k int, c, wc Color, stopLiberty int) []*Worm {
	n4 := Neigh4(k)
	ret := []*Worm{}
	for _, nk := range n4 {
		if b.Points[nk] != wc {
			continue
		}
		include := false
		for _, w := range ret {
			if w.IncludePoint(nk) {
				include = true
				break
			}
		}
		if include {
			continue
		}
		worm := b.WormFromPoint(nk, b.Points[nk], stopLiberty)
		ret = append(ret, worm)
	}
	return ret
}

func (b *Board) InitHash() {
	b.PointHash = make([]int64, NPOINT)
	b.PatternHash = make([][]int64, NPOINT)
	for i := 0; i < NPOINT; i++ {
		x, y := IndexPos(i)
		b.PointHash[i] = b.VertexHash(x, y, b.Points[i])
	}
	for i := 0; i < NPOINT; i++ {
		b.PatternHash[i] = make([]int64, PATTERN_SIZE)
		x, y := IndexPos(i)

		for dy := -1*PATTERN_SIZE + 1; dy < PATTERN_SIZE; dy++ {
			for dx := -1*PATTERN_SIZE + 1; dx < PATTERN_SIZE; dx++ {
				x1, y1 := x+dx, y+dy
				d := Abs(dx) + Abs(dy)
				if d >= PATTERN_SIZE {
					continue
				}
				c := INVALID_COLOR
				if !PosOutBoard(x1, y1) {
					c = b.Points[PosIndex(x1, y1)]
				}
				b.PatternHash[i][d] ^= b.VertexHash(x1, y1, c)
			}
		}
	}
}

func (b *Board) UpdateHash(k int, oc, nc Color) {
	x, y := IndexPos(k)
	b.PointHash[k] ^= b.VertexHash(x, y, oc)
	b.PointHash[k] ^= b.VertexHash(x, y, nc)

	for d := 0; d < PATTERN_SIZE; d++ {
		for _, j := range PointDisMap[k][d] {
			b.PatternHash[j][d] ^= b.VertexHash(x, y, oc)
			b.PatternHash[j][d] ^= b.VertexHash(x, y, nc)
		}
	}
}

func (b *Board) FeatureHash(k int, c Color) int64 {
	ret := int64(0)
	ret ^= b.EdgeDisHash(k)
	ret ^= b.SelfWormHash(k, c)
	ret ^= b.OpWormHash(k, c)
	return ret
}

//calc before put
func (b *Board) FinalPatternHash(k int, c Color) []int64 {
	ret := make([]int64, PATTERN_SIZE)
	h := int64(0)
	fh := b.FeatureHash(k, c)
	for d := 0; d < PATTERN_SIZE; d++ {
		h ^= b.PatternHash[k][d]
		ret[d] = h ^ fh
	}
	return ret
}
