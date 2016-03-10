package gogo

func SelectOne(d map[int]byte) int {
	for k, _ := range d {
		return k
	}
	return -1
}

type Worm struct {
	Points        *PointMap
	Color         Color
	BorderColor   Color
	BorderPoints  *PointMap
	LibertyPoints *PointMap
	OriginDragon  *Dragon
	OriginPoint   int
	Liberty       int
}

func NewWorm() *Worm {
	return &Worm{
		Points:        NewPointMap(16),
		Color:         UNKNOWN,
		OriginPoint:   10000,
		BorderPoints:  NewPointMap(16),
		LibertyPoints: NewPointMap(8),
	}
}

func (w *Worm) calcLiberty(b *Board) {
	w.Liberty = 0
	for _, p := range w.BorderPoints.Points {
		if b.w[p].color == GRAY {
			w.Liberty += 1
			w.LibertyPoints.Add(p)
		}
	}
}

func (w *Worm) Surround() bool {
	return (w.Color == WHITE && w.BorderColor == BLACK) || (w.Color == BLACK && w.BorderColor == WHITE)
}

func (w *Worm) Dead() bool {
	return w.Surround() && w.Liberty == 0
}

func (w *Worm) Exist(p int) bool {
	return w.Points.Exist(p)
}

func (w *Worm) Add(p int) {
	if w.OriginPoint > p {
		w.OriginPoint = p
	}
	w.Points.Add(p)
}

func (b *Board) WormContainsPointBeforePut(p int, stone Color) *Worm {
	oldColor := b.w[p].color
	b.w[p].color = stone
	ret := b.WormContainsPoint(p)
	b.w[p].color = oldColor
	return ret
}

func (b *Board) WormContainsPoint(p int) *Worm {
	worm := NewWorm()
	qe := make([]int, 0, 5)
	qe = append(qe, p)
	qstart := 0
	worm.Color = b.w[p].color
	borderColors := 0
	for {
		if len(qe) == qstart {
			break
		}
		v := qe[qstart]
		qstart += 1
		if worm.Exist(v) {
			continue
		}
		worm.Add(v)
		vx, vy := b.pos(v)
		neigh := b.Neighbor4(vx, vy)
		for _, u := range neigh {
			pu := b.index(u.x, u.y)
			if b.Get(vx, vy).color == u.color {
				if !worm.Exist(pu) {
					qe = append(qe, pu)
				}
			} else {
				borderColors = borderColors | (1 << u.color)
				worm.BorderPoints.Add(pu)
			}
		}
	}
	if borderColors == (1 << BLACK) {
		worm.BorderColor = BLACK
	} else if borderColors == (1 << WHITE) {
		worm.BorderColor = WHITE
	} else {
		worm.BorderColor = GRAY
	}
	worm.calcLiberty(b)
	return worm
}

func (b *Board) WormsFromPoints(total map[int]byte) []*Worm {
	ret := make([]*Worm, 0, 20)
	for {
		if len(total) == 0 {
			break
		}
		v := SelectOne(total)
		worm := b.WormContainsPoint(v)
		ret = append(ret, worm)
		for _, p := range worm.Points.Points {
			delete(total, p)
		}
	}
	return ret
}

func (b *Board) MakeWorms() []*Worm {
	total := make(map[int]byte)
	for p := 0; p < len(b.w); p++ {
		total[p] = 1
	}
	return b.WormsFromPoints(total)
}

func (b *Board) points2Map(ps []Point) map[int]byte {
	ret := make(map[int]byte)
	for _, p := range ps {
		ret[b.index(p.x, p.y)] = 1
	}
	return ret
}

//stone should already be putted in x, y
func (b *Board) GetTakeWorms(stone Color, x, y int) []*Worm {
	n4 := b.Neighbor4(x, y)
	mn4 := b.points2Map(n4)
	cws := b.WormsFromPoints(mn4)
	ret := []*Worm{}
	for _, worm := range cws {
		if worm.Dead() && worm.Color != stone {
			ret = append(ret, worm)
			if worm.Color == WHITE {
				b.takeWhite += 1
			} else if worm.Color == BLACK {
				b.takeBlack += 1
			}
		}
	}
	return ret
}

func (b *Board) koPositionOfDeadWorms(cur Point, worms []*Worm) Point {
	if b.PointLiberty(cur) > 0 {
		return InvalidPoint()
	}
	for _, worm := range worms {
		if worm.Points.Size() != 1 {
			continue
		}
		return b.w[worm.Points.First()]
	}
	return InvalidPoint()
}

func (b *Board) TakeWorm(worm *Worm) {
	for _, p := range worm.Points.Points {
		b.w[p].color = GRAY
	}
}

func (b *Board) TakeWorms(worms []*Worm) {
	for _, worm := range worms {
		b.TakeWorm(worm)
	}
}

func (b *Board) AdjustByWorms() {
	worms := b.MakeWorms()
	for _, worm := range worms {
		if worm.Dead() {
			for _, p := range worm.Points.Points {
				b.w[p].color = GRAY
			}
		}
	}
}
