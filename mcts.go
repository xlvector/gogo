package gogo

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/xlvector/hector/core"
)

var MCTSLock = &sync.Mutex{}

func (b *Board) SelfBattle(c Color, lgr *LastGoodReply) map[int]Color {
	rand.Seed(time.Now().UnixNano())
	amaf := make(map[int]Color)
	//rank := make(map[int]float64)
	p := -1
	n := 0
	for n < 350 {
		pass := 0
		//p, rank = b.GenMove(c, rank)
		p = b.GenSelfBattleMove(c, lgr)
		if p < 0 {
			pass += 1
		} else {
			if _, ok := amaf[p]; !ok {
				amaf[p] = c
			}
		}
		//p, rank = b.GenMove(OpColor(c), rank)
		p = b.GenSelfBattleMove(OpColor(c), lgr)
		if p < 0 {
			pass += 1
		} else {
			if _, ok := amaf[p]; !ok {
				amaf[p] = OpColor(c)
			}
		}
		if pass >= 2 {
			break
		}
		n += 1
	}
	return amaf
}

func (b *Board) MakeWorms() []Worm {
	ret := []Worm{}
	nodes := make(map[int]byte)
	for k, _ := range b.Points {
		nodes[k] = 1
	}

	for {
		if len(nodes) == 0 {
			break
		}
		for k, _ := range nodes {
			worm := b.WormFromPoint(k, b.Points[k], -1)
			ret = append(ret, worm)
			for _, p := range worm.Points.Points {
				delete(nodes, p)
			}
			break
		}
	}
	return ret
}

func (b *Board) Score() float64 {
	worms := b.MakeWorms()
	ret := -7.5
	for _, w := range worms {
		if w.Color == BLACK {
			ret += float64(w.Size())
		} else if w.Color == WHITE {
			ret -= float64(w.Size())
		} else {
			if w.BorderColor == BLACK {
				ret += float64(w.Size())
			} else if w.BorderColor == WHITE {
				ret -= float64(w.Size())
			}
		}
	}
	return ret
}

func (b *Board) CandidateMoves(c Color, rank map[int]float64) map[int]float64 {
	last, _ := b.LastMove()
	if rank == nil {
		rank = make(map[int]float64)
	}
	calcAll := false
	if len(rank) == 0 {
		calcAll = true
	}
	for k, _ := range b.Points {
		calc := true
		if !calcAll {
			calc = false
			if last >= 0 && Distance(k, last) < 3 {
				calc = true
			} else {
				if pr, ok := rank[k]; ok && pr > 0.8 {
					calc = true
				}
			}
		}
		if !calc {
			continue
		}
		if ok, _ := b.CanPut(k, c); ok {
			pr := rand.Float64() * 0.1
			if b.Model != nil {
				pat := b.FinalPatternHash(k, c)
				smp := b.PatternFeature(k, c, b.LastPattern, pat)
				sample := core.NewSample()
				for _, v := range smp {
					sample.AddFeature(core.Feature{v, 1.0})
				}
				pr = b.Model.Predict(sample)
			}
			rank[k] = pr
		}
	}
	return rank
}

func (b *Board) SaveAtari(w Worm) []int {
	ret := make([]int, 0, 2)
	if w.Liberty != 1 {
		return ret
	}
	for _, p := range w.LibertyPoints.Points {
		if b.PointLiberty(p) > 2 {
			ret = append(ret, p)
		} else {
			w2 := b.WormFromPoint(p, w.Color, 2)
			if w2.Liberty > 1 {
				ret = append(ret, p)
			}
		}
	}
	n4 := b.WormNeighWorms(w, OpColor(w.Color), 2)
	for _, w2 := range n4 {
		if w2.Liberty == 1 {
			ret = append(ret, w2.LibertyPoints.First())
		}
	}
	ret2 := make([]int, 0, len(ret))
	for _, k := range ret {
		if ok, _ := b.CanPut(k, w.Color); ok {
			ret2 = append(ret2, k)
		}
	}
	return ret2
}

func (b *Board) LibertyWithColor(c Color, k int) int {
	n4 := Neigh4(k)
	ret := 0
	for n4 > 0 {
		p := int(n4 & 0x1ff)
		n4 = (n4 >> 9)
		if IndexOutBoard(p) {
			continue
		}
		if b.Points[p] == GRAY {
			ret += 1
			continue
		}
		if b.Points[p] != c {
			return -1
		}
	}
	return ret
}

func (b *Board) GenSelfBattleMove(c Color, lgr *LastGoodReply) int {
	rank := make(map[int]float64)
	last, _ := b.LastMove()
	if last >= 0 {
		if lgr != nil {
			p := lgr.Move(c, last)
			if p > 0 {
				if ok, _ := b.CanPut(p, c); ok {
					b.Put(p, c)
					return p
				}
			}
		}

		//fix atari
		ret := make([]int, 0, 2)
		worms := b.NeighWorms(last, OpColor(c), c, 2)
		for _, w := range worms {
			if w.Liberty == 1 {
				ret = append(ret, b.SaveAtari(w)...)
			}
		}
		if len(ret) > 0 {
			p := ret[rand.Intn(len(ret))]
			b.Put(p, c)
			return p
		}

		//pat3
		ret = make([]int, 0, 8)
		n4 := Neigh4(last)
		for n4 > 0 {
			p := int(n4 & 0x1ff)
			n4 = (n4 >> 9)
			pat := b.Pattern3x3String(p, INVALID_COLOR)
			if _, ok := pat3x3Dict[pat]; ok {
				if ok2, _ := b.CanPut(p, c); ok2 {
					ret = append(ret, p)
				}
			}
		}
		nd := NeighD(last)
		for nd > 0 {
			p := int(nd & 0x1ff)
			nd = (nd >> 9)
			pat := b.Pattern3x3String(p, INVALID_COLOR)
			if _, ok := pat3x3Dict[pat]; ok {
				if ok2, _ := b.CanPut(p, c); ok2 {
					ret = append(ret, p)
				}
			}
		}
		if len(ret) > 0 {
			p := ret[rand.Intn(len(ret))]
			b.Put(p, c)
			return p
		}
	}

	visited := ZeroPointMap()

	for j := len(b.Actions) - 1; j >= 0; j-- {
		a := b.Actions[j]
		k, ac := ParseIndexAction(a)
		if visited.Exist(k) {
			continue
		}
		worm := b.WormFromPoint(k, b.Points[k], 3)
		for _, p := range worm.Points.Points {
			visited.Add(p)
		}
		if ac == OpColor(c) {
			if worm.Liberty == 1 {
				rank[worm.LibertyPoints.First()] += 10000000.0 + float64(1000*worm.Size())
			} else if worm.Liberty == 2 {
				for _, p := range worm.LibertyPoints.Points {
					if b.PointLiberty(p) == 3 {
						rank[p] += 10000.0 + float64(worm.Size()*10)
					} else if b.PointLiberty(p) == 2 {
						rank[p] += 500.0 + float64(worm.Size())
					}
				}
			} else if worm.Liberty == 3 {
				ext := b.ExtendLiberty(worm.LibertyPoints)
				if ext.Size() <= 3 {
					for _, p := range worm.LibertyPoints.Points {
						rank[p] += 100.0 + float64(worm.Size())
					}
				}
			}

		} else if ac == c {
			if worm.Liberty == 1 {
				ps := b.SaveAtari(worm)
				for _, p := range ps {
					rank[p] += 100000.0 + float64(worm.Size()*100)
				}
			} else if worm.Liberty == 2 {
				for _, p := range worm.LibertyPoints.Points {
					if b.PointLiberty(p) > 2 {
						px, py := IndexPos(p)
						rank[p] += 10.0 + float64(EdgeDis(px, py))*10.0
					}
				}
			} else if worm.Liberty == 3 {
				ext := b.ExtendLiberty(worm.LibertyPoints)
				if ext.Size() <= 3 {
					for _, p := range worm.LibertyPoints.Points {
						px, py := IndexPos(p)
						rank[p] += 10.0 + float64(EdgeDis(px, py))*10.0
					}
				}
			}
		}

		for _, p := range worm.LibertyPoints.Points {
			if b.PointLiberty(p) == 2 {
				rank[p] += 1.0
			}
		}

		if worm.Liberty*3 < worm.Size() {
			for _, p := range worm.LibertyPoints.Points {
				rank[p] += 0.1
			}
		}
	}

	psum := 0.0
	for _, v := range rank {
		psum += v
	}

	psum *= rand.Float64()
	for k, v := range rank {
		psum -= v
		if psum <= 0.0 {
			if ok, _ := b.CanPut(k, c); ok {
				b.Put(k, c)
				return k
			}
		}
	}

	for i := 0; i < NPOINT*2; i++ {
		k := rand.Intn(NPOINT)
		/*
			if i < 10 && b.Model != nil {
				pat := b.FinalPatternHash(k, c)
				smp := b.PatternFeature(k, c, b.LastPattern, pat)
				sample := core.NewSample()
				for _, v := range smp {
					sample.AddFeature(core.Feature{v, 1.0})
				}
				pr := b.Model.Predict(sample)
				if pr < 0.1 {
					continue
				}
			}
		*/
		if ok, _ := b.CanPut(k, c); ok {
			b.Put(k, c)
			return k
		}
	}
	return -1
}

func (b *Board) EmptyWithinDis(k, d int) bool {
	for i := 0; i <= d; i++ {
		for _, p := range PointDisMap[k][i] {
			if b.Points[p] != GRAY {
				return false
			}
		}
	}
	return true
}

func (b *Board) GenMove(c Color, rank map[int]float64) (int, map[int]float64) {
	rank = b.CandidateMoves(c, rank)
	cands := TopN(rank, 10)
	if len(cands) == 0 {
		return -1, rank
	}
	psum := 0.0
	for _, cand := range cands {
		psum += cand.Second
	}
	pr := rand.Float64() * psum
	for _, cand := range cands {
		pr -= cand.Second
		if pr <= 0.0 {
			pf := cand.First
			b.Put(pf, c)
			delete(rank, pf)
			return pf, rank
		}
	}
	return -1, rank
}

func (b *Board) GenBestMove(c Color, gt *GameTree) (bool, int) {
	rank := b.CandidateMoves(c, nil)
	cands := TopN(rank, 1)
	if len(cands) == 0 {
		return false, -1
	}
	if ok := b.Put(cands[0].First, c); ok {
		x, y := IndexPos(cands[0].First)
		gt.Add(NewGameTreeNode(c, x, y))
		return true, cands[0].First
	}
	return false, -1
}

func (p *GameTreeNode) RaveValue() float64 {
	if p.visit == 0 {
		return p.prior
	}
	e := p.prior*0.1 + 0.9*float64(p.win)/float64(p.visit)
	if p.aVisit == 0 {
		return e
	}
	re := float64(p.aWin) / float64(p.aVisit)
	beta := float64(p.aVisit) / (float64(p.aVisit+p.visit) + float64(p.aVisit*p.visit)/3000.0)
	return beta*re + (1.0-beta)*e
}

func (p *GameTreeNode) UCTValue() float64 {
	if p.visit == 0 {
		return rand.Float64()
	}
	ret := float64(p.win) / float64(p.visit)
	np := 1.0
	if p.Father != nil && p.Father.visit > 0 {
		np = float64(p.Father.visit)
	}
	ret += p.prior * math.Sqrt(np) / float64(1+p.visit)
	return ret
}

type LastGoodReply struct {
	black map[int]int
	white map[int]int
	lock  *sync.RWMutex
}

func (p *LastGoodReply) Move(c Color, last int) int {
	p.lock.RLock()
	defer p.lock.RUnlock()
	if c == BLACK {
		k, ok := p.black[last]
		if ok {
			return k
		}
	} else if c == WHITE {
		k, ok := p.white[last]
		if ok {
			return k
		}
	}
	return -1
}

func (p *LastGoodReply) Set(c Color, k1, k2 int) {
	p.lock.Lock()
	defer p.lock.Unlock()
	if c == BLACK {
		p.black[k1] = k2
	} else if c == WHITE {
		p.white[k1] = k2
	}
}

func NewLastGoodReply() *LastGoodReply {
	return &LastGoodReply{
		black: make(map[int]int),
		white: make(map[int]int),
		lock:  &sync.RWMutex{},
	}
}

func (b *Board) MCTSMove(c Color, gt *GameTree, expand, n int) (bool, int) {
	wg := &sync.WaitGroup{}
	root := gt.Current
	log.Println(PointString(root.x, root.y, root.stone), root.win, root.visit, "next stone color: ", ColorMark(c))
	if len(root.Children) > 0 {
		for _, child := range root.Children {
			log.Println(PointString(child.x, child.y, child.stone), child.win, child.visit)
		}
	}
	lgr := NewLastGoodReply()
	for i := 0; i < n; i++ {
		if i%1000 == 0 {
			fmt.Print(".")
		}
		node := MCTSSelection(gt)
		MCTSExpand(node, b, expand, c, lgr, wg)
	}
	fmt.Println()
	wg.Wait()
	var best *GameTreeNode
	robust := 0
	for _, child := range root.Children {
		winrate := float64(child.win) / float64(child.visit)
		log.Println(PointString(child.x, child.y, child.stone), winrate, child.win, child.visit, child.prior)
		if robust < child.visit {
			robust = child.visit
			best = child
		}
	}
	gt.Current = best
	return b.Put(PosIndex(best.x, best.y), c), PosIndex(best.x, best.y)
}

func MCTSSelection(gt *GameTree) *GameTreeNode {
	root := gt.Current
	ret := root
	depth := 0
	for {
		ret.visit += 3
		if ret.Children == nil || len(ret.Children) == 0 {
			return ret
		}
		if len(ret.CandMoves) > 0 {
			return ret
		}
		depth += 1
		maxVal := 0.0
		var best *GameTreeNode
		for _, child := range ret.Children {
			val := child.UCTValue()
			if maxVal < val {
				maxVal = val
				best = child
			}
		}
		ret = best
	}

	return ret
}

func NewBoardFromPath(path []*GameTreeNode) *Board {
	ret := NewBoard()
	for i := len(path) - 2; i >= 0; i-- {
		cur := path[i]
		if PosOutBoard(cur.x, cur.y) {
			break
		}
		ret.Put(PosIndex(cur.x, cur.y), cur.stone)
	}
	return ret
}

func MCTSExpand(node *GameTreeNode, oBoard *Board, nLeaf int, wc Color, lgr *LastGoodReply, wg *sync.WaitGroup) {
	board := NewBoardFromPath(node.Path2Root())
	board.Model = oBoard.Model
	oc := BLACK
	if node.stone == BLACK || node.stone == WHITE {
		oc = OpColor(node.stone)
	}

	if len(node.Children) == 0 {
		rank := board.CandidateMoves(oc, nil)
		topn := TopN(rank, nLeaf)
		//line := PointString(node.x, node.y, node.stone) + ":"
		for _, child := range topn {
			x, y := IndexPos(child.First)
			cnode := NewGameTreeNode(oc, x, y)
			//line += PointString(x, y, oc) + ","
			cnode.prior = child.Second
			node.CandMoves = append(node.CandMoves, cnode)
		}
		//log.Println(line)
	}
	cnode := node.CandMoves[0]
	node.CandMoves = node.CandMoves[1:]
	_, cnode = node.AddChild(cnode)
	cnode.visit += 3
	wg.Add(1)
	go MCTSSimulation(board.Copy(), cnode, lgr, wg)
}

func MCTSSimulation(b *Board, next *GameTreeNode, lgr *LastGoodReply, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()
	b.Put(PosIndex(next.x, next.y), next.stone)
	amaf := b.SelfBattle(OpColor(next.stone), lgr)
	s := b.Score()

	if s > 0 {
		MCTSBackProp(next, BLACK, amaf, lgr)
	} else {
		MCTSBackProp(next, WHITE, amaf, lgr)
	}
}

func MCTSBackProp(node *GameTreeNode, wc Color, amaf map[int]Color, lgr *LastGoodReply) {
	MCTSLock.Lock()
	defer MCTSLock.Unlock()
	v := node
	for {
		if v == nil {
			return
		}
		if v.stone == wc {
			v.win += 1
		}
		v.visit -= 2

		if len(v.Children) != 0 {
			for _, child := range v.Children {
				pc := PosIndex(child.x, child.y)
				if cp, ok := amaf[pc]; ok && cp == child.stone {
					if cp == wc {
						child.aWin += 1
					}
					child.aVisit += 1
				}
			}
		}

		if lgr != nil {
			if v.stone == wc && v.Father != nil {
				lgr.Set(wc, PosIndex(v.Father.x, v.Father.y), PosIndex(v.x, v.y))
			}
		}
		v = v.Father
	}
}
