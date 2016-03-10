package gogo

type Dragon struct {
	Worms       []*Worm
	OriginPoint int
	Color       Color
}

func NewDragon() *Dragon {
	return &Dragon{
		Worms:       make([]*Worm, 0, 2),
		OriginPoint: 1000,
		Color:       UNKNOWN,
	}
}

func (d *Dragon) Add(w *Worm) {
	d.Worms = append(d.Worms, w)
	if d.OriginPoint > w.OriginPoint {
		d.OriginPoint = w.OriginPoint
	}
}

func (d *Dragon) WormSize() int {
	return len(d.Worms)
}

func (d *Dragon) PointSize() int {
	ret := 0
	for _, w := range d.Worms {
		ret += w.Points.Size()
	}
	return ret
}

func (b *BoardInfo) buildOriginDragon(dragon *Dragon) {
	for _, w := range dragon.Worms {
		for _, p := range w.Points.Points {
			b.PointFetures[p].OriginDragon = dragon
		}
	}
}

func (b *BoardInfo) BuildDragon() []*Dragon {
	co := NewSparseMat()
	for _, pf := range b.PointFetures {
		if pf.P.color != GRAY {
			continue
		}
		for _, w1 := range pf.BoardWorms {
			for _, w2 := range pf.BoardWorms {
				if w1.OriginPoint == w2.OriginPoint {
					continue
				}
				if w1.Color != w2.Color {
					continue
				}
				co.Add(w1.OriginPoint, w2.OriginPoint, 1)
			}
		}
	}

	for _, w := range b.Worms {
		if w.Color != GRAY {
			continue
		}
		if w.BorderColor == GRAY {
			continue
		}
		cws := NewPointMap(5)
		for _, p := range w.BorderPoints.Points {
			cws.Add(b.PointFetures[p].OriginWorm.OriginPoint)
		}

		for _, wo1 := range cws.Points {
			w1 := b.PointFetures[wo1].OriginWorm
			for _, wo2 := range cws.Points {
				if wo1 == wo2 {
					continue
				}
				w2 := b.PointFetures[wo2].OriginWorm
				if w1.Color == w2.Color {
					co.Add(wo1, wo2, 2)
				}
			}
		}
	}

	used := NewBoardBitmap()
	ret := []*Dragon{}
	for w1, row := range co.data {
		if used.IsSet(w1) {
			continue
		}
		dragon := NewDragon()
		dragon.Add(b.PointFetures[w1].OriginWorm)
		b.PointFetures[w1].OriginWorm.OriginDragon = dragon
		used.Set(w1)
		for w2, v := range row {
			if v < 2 {
				continue
			}
			dragon.Add(b.PointFetures[w2].OriginWorm)
			b.PointFetures[w2].OriginWorm.OriginDragon = dragon
			used.Set(w2)
		}
		b.buildOriginDragon(dragon)
		ret = append(ret, dragon)
	}
	return ret
}
