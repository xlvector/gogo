package gogo

var pat3x3 = []string{
	"XOX...???",
	"X.X......",
	"XO....?.?",
	"XO?X..?.?",
	"XO?O.o?o?",
	"XO?O.X???",
	"?X?O.Oooo",
	"?X?o.O???",
	"OX?X.O   ",
	"?XOx.x   ",
	"?OXX.O   ",
}

var pat3x3Dict = PatExpand()

var rotateMatrix = [][]int{
	[]int{3, 2, 1, 6, 5, 4, 9, 8, 7},
	[]int{7, 8, 9, 4, 5, 6, 1, 2, 3},
	[]int{7, 4, 1, 8, 5, 2, 9, 6, 3},
	[]int{9, 8, 7, 6, 5, 4, 3, 2, 1},
	[]int{3, 6, 9, 2, 5, 8, 1, 4, 7},
	[]int{9, 6, 3, 8, 5, 2, 7, 4, 1},
	[]int{1, 4, 7, 2, 5, 8, 3, 6, 9},
}

func ColorRotatePat(pat string) string {
	ret := ""
	for _, c := range pat {
		if c == 'X' {
			ret += "O"
		} else if c == 'O' {
			ret += "X"
		} else if c == 'x' {
			ret += "o"
		} else if c == 'o' {
			ret += "x"
		} else {
			ret += string(c)
		}
	}
	return ret
}

func RotatePats(pat string) []string {
	ret := []string{}
	for _, r := range rotateMatrix {
		tmp := ""
		for _, k := range r {
			tmp += string(pat[k-1])
		}
		ret = append(ret, tmp)
	}
	return ret
}

func SinglePatExpand(pat string) []string {
	ret := []string{""}
	for _, c := range pat {
		ret2 := []string{}
		if c == '?' {
			for _, buf := range ret {
				ret2 = append(ret2, buf+"X")
				ret2 = append(ret2, buf+"O")
				ret2 = append(ret2, buf+".")
				ret2 = append(ret2, buf+" ")
			}
		} else if c == 'X' || c == 'O' || c == '.' || c == ' ' {
			for _, buf := range ret {
				ret2 = append(ret2, buf+string(c))
			}
		} else if c == 'x' {
			for _, buf := range ret {
				ret2 = append(ret2, buf+"O")
				ret2 = append(ret2, buf+" ")
				ret2 = append(ret2, buf+".")
			}
		} else if c == 'o' {
			for _, buf := range ret {
				ret2 = append(ret2, buf+"X")
				ret2 = append(ret2, buf+" ")
				ret2 = append(ret2, buf+".")
			}
		}
		ret = ret2
	}
	return ret
}

func PatExpand() map[string]byte {
	ret := make(map[string]byte)
	for _, pat := range pat3x3 {
		rpats := RotatePats(pat)
		for _, rpat := range rpats {
			epats := SinglePatExpand(rpat)
			for _, e := range epats {
				ret[e] = 1
			}

			epats2 := SinglePatExpand(ColorRotatePat(rpat))
			for _, e := range epats2 {
				ret[e] = 1
			}
		}
	}
	return ret
}

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
	for n4 > 0 {
		nk := int(n4 & 0x1ff)
		n4 = (n4 >> 9)
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

func FeatureEncode(key string, v ...int) int64 {
	h := int64(0)
	for _, c := range key {
		h *= 16777619
		h ^= int64(c)
	}

	for _, c := range v {
		h *= 16777619
		h ^= int64(c)
	}
	if h < 0 {
		h *= -1
	}
	return h
}

func (b *Board) LocalFeature(k int, c Color) []int64 {
	myNWorms := b.NeighWorms(k, c, c, 2)
	opNWorms := b.NeighWorms(k, c, OpColor(c), 2)
	worm := b.WormFromPoint(k, c, 2)
	ret := make([]int64, 0, 5)

	cm1 := 0
	cm2 := 0
	cm3 := 0
	for _, w := range myNWorms {
		if w.Liberty == 1 {
			cm1 += 1
		} else if w.Liberty == 2 {
			cm2 += 1
		} else if w.Liberty == 3 {
			cm3 += 1
		}
	}

	co1 := 0
	co2 := 0
	co3 := 0
	com1 := 0
	for _, w := range opNWorms {
		if w.Liberty == 1 {
			co1 += 1
			opNeighs := b.WormNeighWorms(w, c, 2)
			for _, w2 := range opNeighs {
				if w2.Liberty == 1 {
					com1 += 1
				}
			}
		} else if w.Liberty == 2 {
			co2 += 1
		} else if w.Liberty == 3 {
			co3 += 1
		}
	}
	ret = append(ret, FeatureEncode("liberty0", co1, com1, worm.Liberty))
	ret = append(ret, FeatureEncode("liberty1", cm1, cm2, cm3, co1, co2, co3, com1, worm.Liberty))
	pat3 := b.Pattern3x3String(k, c)
	ret = append(ret, FeatureEncode("pat3x3"+pat3))
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

func ColorMark2(c1, c2 Color) string {
	if c2 == INVALID_COLOR {
		return ColorMark(c1)
	} else {
		if c1 == c2 {
			return "X"
		} else if c1 == OpColor(c2) {
			return "O"
		} else {
			return ColorMark(c1)
		}
	}
}

func (b *Board) Pattern3x3String(p int, stone Color) string {
	x, y := IndexPos(p)
	ret := ""
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			x1, y1 := x+dx, y+dy
			if !PosOutBoard(x1, y1) {
				c := b.Points[PosIndex(x1, y1)]
				ret += ColorMark2(c, stone)
			} else {
				ret += " "
			}
		}
	}
	return ret
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

func (b *Board) NeighWorms(k int, c, wc Color, stopLiberty int) []Worm {
	n4 := Neigh4(k)
	ret := []Worm{}
	for n4 > 0 {
		nk := int(n4 & 0x1ff)
		n4 = (n4 >> 9)
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
