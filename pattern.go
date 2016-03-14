package gogo

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
	if d <= 3 {
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
	worm := b.WormFromPoint(k, c)
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
	nworms := b.NeighWorms(k, c, OpColor(c))
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
	nworms := b.NeighWorms(k, c, c)
	minLiberty := 10000
	for _, w := range nworms {
		if minLiberty > w.Liberty {
			minLiberty = w.Liberty
		}
	}
	if minLiberty > 2 {
		return false
	} else {
		worm := b.WormFromPoint(k, c)
		if worm.Liberty > 2 {
			return true
		}
		return false
	}
}

func (b *Board) LocalFeature(k int, c Color) int64 {
	myNWorms := b.NeighWorms(k, c, c)
	opNWorms := b.NeighWorms(k, c, OpColor(c))
	worm := b.WormFromPoint(k, c)

	f := int64(0)
	minLiberty := 10000
	for _, w := range myNWorms {
		if minLiberty > w.Liberty {
			minLiberty = w.Liberty
		}
	}

	if minLiberty == 1 {
		if worm.Liberty == 1 {
			//escape capture
			f ^= 493570158105
		} else if worm.Liberty == 2 {
			f ^= 159084081432
		} else if worm.Liberty == 3 {
			f ^= 897325971018
		} else if worm.Liberty > 3 {
			f ^= 291850148415
		}
	} else if minLiberty == 2 {
		if worm.Liberty == 1 {
			//escape capture
			f ^= 932759347016
		} else if worm.Liberty == 2 {
			f ^= 758724359874
		} else if worm.Liberty == 3 {
			f ^= 238146923179
		} else if worm.Liberty > 3 {
			f ^= 945876927621
		}
	}

	//op
	minLiberty = 10000
	for _, w := range opNWorms {
		if minLiberty > w.Liberty {
			minLiberty = w.Liberty
		}
	}
	if minLiberty == 1 {
		f ^= 787401927621
	} else if minLiberty == 2 {
		f ^= 304580158101
	}

	f ^= b.EdgeDisHash(k)
	return f
}

func (b *Board) NeighWorms(k int, c, wc Color) []*Worm {
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
		worm := b.WormFromPoint(nk, b.Points[nk])
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
	for d := 0; d < PATTERN_SIZE; d++ {
		h ^= b.PatternHash[k][d]
		ret[d] = h ^ b.FeatureHash(k, c)
	}
	return ret
}
