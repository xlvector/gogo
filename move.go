package gogo

import (
	"log"
	"math/rand"
	"sort"

	"github.com/xlvector/hector/core"
	"github.com/xlvector/hector/dt"
)

type BoardInfo struct {
	Board        *Board
	Worms        []*Worm
	Stones       [][]int
	PointFetures []*PointFeature
}

func NewBoardInfo(b *Board) *BoardInfo {
	ret := &BoardInfo{
		Board:        b,
		Worms:        []*Worm{},
		PointFetures: make([]*PointFeature, len(b.w), len(b.w)),
		Stones:       make([][]int, 3, 3),
	}
	ret.Stones[GRAY] = make([]int, 0, 361)
	ret.Stones[BLACK] = make([]int, 0, 361)
	ret.Stones[WHITE] = make([]int, 0, 361)
	return ret
}

func (b *BoardInfo) genMovesOfLastMoveNeigh(lastMove Point, stone Color) map[int]float64 {
	n4 := b.Board.Neighbor4(lastMove.x, lastMove.y)
	nd := b.Board.NeighDiamond(lastMove.x, lastMove.y)
	ret := make(map[int]float64)
	for _, p := range n4 {
		ret[b.Board.index(p.x, p.y)] = 1.0
	}
	for _, p := range nd {
		ret[b.Board.index(p.x, p.y)] = 1.5
	}
	return ret
}

func (b *BoardInfo) genMovesOfOpLiberty(stone Color) map[int]float64 {
	ret := make(map[int]float64)
	for _, w := range b.Worms {
		if w.Color != OppColor(stone) {
			continue
		}
		if w.Liberty > 3 {
			continue
		}
		for _, p := range w.LibertyPoints.Points {
			if w.Liberty == 1 {
				ret[p] = 10.0
			} else {
				ret[p] = 5.0
			}
		}
	}
	return ret
}

func (b *BoardInfo) genMovesOfMyLiberty(stone Color) map[int]float64 {
	ret := make(map[int]float64)
	for _, w := range b.Worms {
		if w.Color != stone {
			continue
		}
		if w.Liberty > 3 {
			continue
		}
		// 如果棋串的气只有1，那么要看看这个气点附近有没有别的气，或者连接着其他同色的棋串
		if w.Liberty == 1 {
			for _, p := range w.LibertyPoints.Points {
				pf := b.PointFetures[p]
				if pf.Liberty >= 2 {
					ret[p] = 20.0
				} else {
					for _, w2 := range pf.BoardWorms {
						if w2.Color == stone && w2.Liberty > 1 {
							ret[p] = 20.0
						}
					}
				}
			}
			break
		}
		for _, p := range w.LibertyPoints.Points {
			pf := b.PointFetures[p]
			if pf.Liberty < 2 {
				continue
			}
			ret[p] = 5.0
		}
	}
	return ret
}

func (b *BoardInfo) genFillEyeMoves(stone Color) map[int]float64 {
	ret := make(map[int]float64)
	for _, w := range b.Worms {
		if w.Color != GRAY || w.Points.Size() != 1 || w.BorderColor != stone {
			continue
		}
		p := w.Points.First()
		pf := b.PointFetures[p]
		maxLiberty := 0
		minLiberty := 1000
		for _, wh := range pf.BoardWorms {
			if maxLiberty < wh.Liberty {
				maxLiberty = wh.Liberty
			}
			if minLiberty > wh.Liberty {
				minLiberty = wh.Liberty
			}
		}

		if minLiberty == 1 && maxLiberty > 1 {
			ret[p] = 20.0
		}
	}
	return ret
}

func (b *BoardInfo) genMovesByGlobalInfo(stone Color) map[int]float64 {
	ret := []int{}
	for _, w := range b.Worms {
		if w.Color != GRAY {
			continue
		}
		if w.Points.Size() == 1 {
			p := b.Board.w[w.Points.First()]
			if b.Board.SingleEye(p.x, p.y, stone) {
				continue
			}
		}
		for _, p := range w.Points.Points {
			ret = append(ret, p)
		}
	}
	pts := make(map[int]float64)
	if len(ret) == 0 {
		return pts
	}
	for i := 0; i < 2; i++ {
		k := ret[rand.Intn(len(ret))]
		pts[k] = 0.1
	}
	return pts
}

func addMap(a, b map[int]float64) map[int]float64 {
	for kb, vb := range b {
		va, ok := a[kb]
		if ok {
			a[kb] = va + vb
		} else {
			a[kb] = vb
		}
	}
	return a
}

func (b *BoardInfo) CandidateMoves(lastMove Point, stone Color, model *dt.RandomForest, n int) map[int]float64 {
	cand := make(map[int]float64)
	b.GenComplexFeatures(stone)
	ifl := make(IntFloatPairList, 0, 100)
	for _, p := range b.Stones[GRAY] {
		pf := b.PointFetures[p]
		if pf.SelectProb <= 0.00001 {
			continue
		}
		prob := pf.SelectProb
		if model != nil {
			prob = model.Predict(pf.Fc)
		}
		ifl = append(ifl, IntFloatPair{p, prob})
	}
	sort.Sort(ifl)
	for i := len(ifl) - 1; i >= 0 && i >= len(ifl)-n; i-- {
		pr := ifl[i].Second
		cand[ifl[i].First] = pr
	}
	return cand
}

func (info *BoardInfo) WormBorderConnected(worm *Worm) bool {
	if worm.BorderColor == GRAY {
		return false
	}
	k := -1
	for _, p := range worm.BorderPoints.Points {
		pf := info.PointFetures[p]
		if k < 0 {
			k = pf.OriginWorm.OriginPoint
		} else {
			if k != pf.OriginWorm.OriginPoint {
				return false
			}
		}
	}
	return true
}

func (info *BoardInfo) GenComplexFeatures(stone Color) {
	for _, w := range info.Worms {
		//自己的棋串
		if w.Color == stone {
			// 如果棋串的气只有1，那么要看看这个气点附近有没有别的气，或者连接着其他同色的棋串
			if w.Liberty == 1 {
				for _, p := range w.LibertyPoints.Points {
					pf := info.PointFetures[p]
					if pf.Liberty >= 2 {
						pf.Fc.AddFeature(core.Feature{F_CAPTURE_CAN_ESCAPE, 1.0})
						pf.SetSelectProb(1.0)
					} else {
						escape := false
						for _, w2 := range pf.BoardWorms {
							if w2.Color == stone && w2.OriginPoint != w.OriginPoint && w2.Liberty > 2 {
								pf.Fc.AddFeature(core.Feature{F_CAPTURE_CAN_ESCAPE, 1.0})
								escape = true
								pf.SetSelectProb(1.0)
							}
						}
						if !escape {
							pf.Fc.AddFeature(core.Feature{F_CAPTURE_CANNOT_ESCAPE, 1.0})
							pf.SetSelectProb(1.0)
						}
					}
				}
			} else if w.Liberty == 2 {
				for _, p := range w.LibertyPoints.Points {
					pf := info.PointFetures[p]
					pf.Fc.AddFeature(core.Feature{F_ATARI, float64(pf.Liberty)})
					pf.SetSelectProb(1.0)
				}
			} else {
				for _, p := range w.LibertyPoints.Points {
					pf := info.PointFetures[p]
					pf.Fc.AddFeature(core.Feature{F_LIBERTY, float64(pf.Liberty)})
					pf.SetSelectProb(1.0)
				}
			}
		} else if w.Color == OppColor(stone) {
			if w.Liberty == 1 {
				for _, p := range w.LibertyPoints.Points {
					pf := info.PointFetures[p]
					pf.Fc.AddFeature(core.Feature{F_OP_CAPTURE, float64(pf.Liberty)})
					pf.SetSelectProb(1.0)
				}
			} else if w.Liberty == 2 {
				for _, p := range w.LibertyPoints.Points {
					pf := info.PointFetures[p]
					pf.Fc.AddFeature(core.Feature{F_OP_ATARI, float64(pf.Liberty)})
					pf.SetSelectProb(1.0)
				}
			} else {
				for _, p := range w.LibertyPoints.Points {
					pf := info.PointFetures[p]
					pf.Fc.AddFeature(core.Feature{F_OP_LIBERTY, float64(pf.Liberty)})
					pf.SetSelectProb(1.0)
				}
			}
		} else if w.Color == GRAY {
			if w.BorderColor == stone {
				for _, p := range w.Points.Points {
					pf := info.PointFetures[p]
					if pf.Liberty >= 3 {
						continue
					}
					pf.Fc.AddFeature(core.Feature{F_SURROUND_MY, float64(pf.Liberty)})
					pf.SetSelectProb(0.2)
				}
			} else if w.BorderColor == OppColor(stone) {
				if w.Points.Size() > 10 {
					continue
				}
				for _, p := range w.Points.Points {
					pf := info.PointFetures[p]
					pf.Fc.AddFeature(core.Feature{F_SURROUND_OP, float64(pf.Liberty)})
					pf.SetSelectProb(0.2)
				}
			}
		}
	}

	for _, p := range info.Stones[GRAY] {
		pf := info.PointFetures[p]
		pf.Fc.AddFeature(core.Feature{F_DISTANCE_EDGE, float64(pf.P.EdgeDis(info.Board.size))})
		pf.Fc.AddFeature(core.Feature{F_EDGE_DIS_X, float64(pf.P.EdgeDisX(info.Board.size))})
		pf.Fc.AddFeature(core.Feature{F_EDGE_DIS_Y, float64(pf.P.EdgeDisY(info.Board.size))})
		pf.SetSelectProb(0.01)

		sameWormMaxLiberty := 0
		sameWormMinLiberty := 1000
		diffWormMaxLiberty := 0
		diffWormMinLiberty := 1000
		connectCount := 0
		connectWormSize := 0
		cutCount := 0
		cutWormSize := 0
		for _, w := range pf.BoardWorms {
			if w.Color == stone {
				connectCount += 1
				connectWormSize += w.Points.Size()
				if sameWormMaxLiberty < w.Liberty {
					sameWormMaxLiberty = w.Liberty
				}
				if sameWormMinLiberty > w.Liberty {
					sameWormMinLiberty = w.Liberty
				}

				pf.SetSelectProb(1.0)
			} else if w.Color == OppColor(stone) {
				cutCount += 1
				cutWormSize += w.Points.Size()
				pf.SetSelectProb(1.0)
				if diffWormMaxLiberty < w.Liberty {
					diffWormMaxLiberty = w.Liberty
				}
				if diffWormMinLiberty > w.Liberty {
					diffWormMinLiberty = w.Liberty
				}
			}
		}
		if connectCount > 0 {
			pf.Fc.AddFeature(core.Feature{F_CONNECT_COUNT, float64(connectCount)})
			pf.Fc.AddFeature(core.Feature{F_CONNECT_WORM_SIZE, float64(connectWormSize)})
			pf.Fc.AddFeature(core.Feature{F_CONNECT_WORM_MAX_LIBERTY, float64(sameWormMaxLiberty)})
		}
		if cutCount > 0 {
			pf.Fc.AddFeature(core.Feature{F_CUT_COUNT, float64(cutCount)})
			pf.Fc.AddFeature(core.Feature{F_CUT_WORM_SIZE, float64(cutWormSize)})
		}

		pf.Fc.AddFeature(core.Feature{F_ORIGIN_WORM_SIZE, float64(pf.OriginWorm.Points.Size())})
		pf.Fc.AddFeature(core.Feature{F_ORIGIN_WORM_LIBERTY, float64(pf.OriginWorm.Liberty)})
		//如果这个worm只有一个点，说明是一个眼，同时他的边界是自己的颜色
		if pf.OriginWorm.Points.Size() == 1 && pf.OriginWorm.BorderColor == stone {
			//要看他的边界所属的worm，如果这些worm中最小的超过1口气，那这个眼就不要填
			if sameWormMinLiberty > 1 {
				pf.SelectProb = 0.0
			}
		}

		//如果这个worm只有一个点，是对方的一只眼，那除非放进去能提掉对方的棋，否则位置非法
		if pf.OriginWorm.Points.Size() == 1 && pf.OriginWorm.BorderColor == OppColor(stone) {
			if diffWormMinLiberty > 1 {
				pf.SelectProb = 0.0
			}
		}

		//如果这个worm只有一个点，同时边界中自己的棋串最少的只有一口气，非法
		if pf.OriginWorm.Points.Size() == 1 {
			if sameWormMaxLiberty == 1 {
				pf.SelectProb = 0.0
			}
		}

		pat3x3 := info.Board.Pattern(pf.P.x, pf.P.y, stone, 1)
		pf.Fc.AddFeature(core.Feature{F_PAT_3X3 + 100*pat3x3, 1.0})
		pat5x5 := info.Board.Pattern(pf.P.x, pf.P.y, stone, 2)
		pf.Fc.AddFeature(core.Feature{F_PAT_5X5 + 100*pat5x5, 1.0})
	}

	//对于所有的gray worm,如果超过5个点，同时边界属于同一个worm，就不要往里面填棋了。这样可以提高自我对局的速度
	for _, w := range info.Worms {
		if w.Color != GRAY {
			continue
		}
		if w.BorderColor == GRAY {
			continue
		}
		if w.BorderPoints.Size() < 6 {
			continue
		}
		if info.WormBorderConnected(w) {
			for _, p := range w.Points.Points {
				pf := info.PointFetures[p]
				pf.SelectProb = 0.0
			}
		}
	}

}

func (b *Board) SinglePoint(x, y int) bool {
	n4 := b.Neighbor4(x, y)
	for _, p := range n4 {
		if p.Color() != GRAY {
			return false
		}
	}
	nd := b.NeighDiamond(x, y)
	for _, p := range nd {
		if p.Color() != GRAY {
			return false
		}
	}
	return true
}

func (b *BoardInfo) removeWorm(worm *Worm) {
	p := -1
	for k, w := range b.Worms {
		if worm.OriginPoint == w.OriginPoint {
			p = k
			break
		}
	}
	if p < 0 {
		return
	}
	b.Worms = append(b.Worms[:p], b.Worms[p:]...)
}

func (b *Board) CollectBoardInfo(lastMove Point) *BoardInfo {
	info := NewBoardInfo(b)

	if b.info != nil && lastMove.Valid() && b.SinglePoint(lastMove.x, lastMove.y) {
		p := b.index(lastMove.x, lastMove.y)
		pworm := b.info.PointFetures[p].OriginWorm

		worm1 := b.WormContainsPoint(b.index(lastMove.x, lastMove.y))
		n4 := b.Neighbor4(lastMove.x, lastMove.y)
		worm2 := b.WormContainsPoint(b.index(n4[0].x, n4[0].y))

		b.info.removeWorm(pworm)
		b.info.Worms = append(b.info.Worms, worm1)
		b.info.Worms = append(b.info.Worms, worm2)
		info.Worms = b.info.Worms
	} else {
		info.Worms = b.MakeWorms()
	}

	for i, p := range b.w {
		pf := NewPointFeature(p)
		info.PointFetures[i] = pf
		info.Stones[p.color] = append(info.Stones[p.color], i)
		pf.Liberty = b.PointLiberty(pf.P)
	}

	for _, worm := range info.Worms {
		for _, i := range worm.Points.Points {
			info.PointFetures[i].OriginWorm = worm
		}

		for _, i := range worm.BorderPoints.Points {
			info.PointFetures[i].BoardWorms[worm.OriginPoint] = worm
		}
	}
	//info.StoneGrayDistance(BLACK)
	//info.StoneGrayDistance(WHITE)
	b.info = info
	return info
}

func (b *Board) GenMove(lastMove Point, stone Color) Point {
	maxPr := 0.0
	best := InvalidPoint()
	for i, p := range b.w {
		if p.color != GRAY {
			continue
		}
		fe := b.PointSimpleFeature(p, stone)
		if fe == nil || len(fe) == 0 {
			continue
		}
		sample := core.NewSample()
		for _, k := range fe {
			sample.AddFeature(core.Feature{k*1000 + int64(i), 1.0})
		}
		for _, k := range b.lastMoveHash {
			sample.AddFeature(core.Feature{k*1000 + int64(i), 1.0})
		}
		pr := 0.5
		if b.model != nil {
			pr = b.model.Predict(sample)
		}
		if maxPr < pr {
			maxPr = pr
			best = p
		}
	}
	log.Println(best.String(), maxPr)
	if err := b.Put(best.x, best.y, stone); err != nil {
		log.Println(err)
		return InvalidPoint()
	}
	return best
}

func (b *Board) StableEye(x, y int, stone Color) bool {
	n4 := b.Neighbor4(x, y)
	for _, p := range n4 {
		if p.Color() != stone {
			return false
		}
	}
	nd := b.NeighDiamond(x, y)
	n := 0
	for _, p := range nd {
		if p.Color() == OppColor(stone) {
			n += 1
		}
	}
	if len(nd) == 4 && n < 2 {
		return true
	}
	if len(nd) < 4 && n == 0 {
		return true
	}
	return false
}

func (b *Board) QuickCandidateMoves(lastMove Point, stone Color, n int) map[int]float64 {
	ret := make(map[int]float64)

	rank := make(IntFloatPairList, 0, 100)
	for i, p := range b.w {
		if p.color != GRAY {
			continue
		}
		fe := b.PointSimpleFeature(p, stone)
		if fe == nil || len(fe) == 0 {
			continue
		}
		sample := core.NewSample()
		for _, k := range fe {
			sample.AddFeature(core.Feature{k*1000 + int64(i), 1.0})
		}
		for _, k := range b.lastMoveHash {
			sample.AddFeature(core.Feature{k*1000 + int64(i), 1.0})
		}
		pr := 0.51
		if b.model != nil {
			pr = b.model.Predict(sample)
		}
		rank = append(rank, IntFloatPair{i, pr})
	}
	sort.Sort(sort.Reverse(rank))
	for i := 0; i < n && i < len(rank); i++ {
		ret[rank[i].First] = rank[i].Second
	}
	return ret
}

func (b *Board) GenQuickMove(lastMove Point, stone Color) Point {
	movs := b.QuickCandidateMoves(lastMove, stone, 10)
	psum := 0.0
	for _, v := range movs {
		psum += v
	}

	pr := rand.Float64() * psum
	for k, v := range movs {
		log.Println(k, v)
		pr -= v
		if psum <= 0.0 {
			p := b.w[k]
			if err := b.Put(p.x, p.y, stone); err == nil {
				return p
			}
		}
	}
	return InvalidPoint()
}
