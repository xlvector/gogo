package gogo

import (
	"log"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/xlvector/hector/core"
)

var MCTSLock = &sync.Mutex{}

func (b *Board) SelfBattle(c Color) int {
	rand.Seed(time.Now().UnixNano())
	rank := make(map[int]float64)
	p := -1
	n := 0
	for {
		pass := 0
		p, rank = b.GenMove(c, rank)
		if p < 0 {
			pass += 1
		}
		p, rank = b.GenMove(OpColor(c), rank)
		if p < 0 {
			pass += 1
		}
		if pass >= 2 {
			break
		}
		n += 1
	}
	return n
}

func (b *Board) MakeWorms() []*Worm {
	ret := []*Worm{}
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
			for p, _ := range worm.Points {
				delete(nodes, p)
			}
			break
		}
	}
	return ret
}

func (b *Board) Score() float64 {
	worms := b.MakeWorms()
	ret := 0.0
	for _, w := range worms {
		if w.Color == BLACK {
			ret += float64(len(w.Points))
		} else if w.Color == WHITE {
			ret -= float64(len(w.Points))
		} else {
			if w.BorderColor == BLACK {
				ret += float64(len(w.Points))
			} else if w.BorderColor == WHITE {
				ret -= float64(len(w.Points))
			}
		}
	}
	return ret
}

func (b *Board) CandidateMoves(c Color, rank map[int]float64) map[int]float64 {
	last, _ := b.LastMove()
	if rand.Float64() < 0.1 || rank == nil {
		rank = make(map[int]float64)
	}
	for k, _ := range b.Points {
		if len(rank) > 0 && last >= 0 && Distance(k, last) > 3 {
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

func (b *Board) GenMove(c Color, rank map[int]float64) (int, map[int]float64) {
	rank = b.CandidateMoves(c, rank)
	cands := TopN(rank, 20)
	psum := 0.0
	for _, v := range cands {
		psum += v.Second
	}

	pr := rand.Float64() * psum
	pf := -1
	for _, v := range cands {
		pr -= v.Second
		if pr <= 0.0 {
			pf = v.First
			break
		}
	}
	b.Put(pf, c)
	delete(rank, pf)
	return pf, rank
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
	ret += math.Sqrt(math.Log(np) / float64(p.visit))
	ret += 0.0001 * rand.Float64()
	return ret
}

func (b *Board) MCTSMove(c Color, gt *GameTree) bool {
	root := gt.Current
	for i := 0; i < 20; i++ {
		node := MCTSSelection(gt)
		MCTSExpand(node, c, b)
		log.Println(i, root.visit)
		time.Sleep(time.Second * 60)
	}
	var best *GameTreeNode
	robust := 0.0
	for _, child := range root.Children {
		winrate := float64(child.win) / float64(10+child.visit)
		log.Println(string(LX[child.x]), child.y+1, ColorMark(child.stone), winrate, child.win, child.visit)
		if robust < winrate {
			robust = winrate
			best = child
		}
	}
	gt.Current = best
	return b.Put(PosIndex(best.x, best.y), c)
}

func MCTSSelection(gt *GameTree) *GameTreeNode {
	root := gt.Current
	ret := root
	for {
		if ret.Children == nil || len(ret.Children) == 0 {
			return ret
		}
		maxUCT := 0.0
		var best *GameTreeNode
		for _, child := range ret.Children {
			uct := child.UCTValue()
			if maxUCT < uct {
				maxUCT = uct
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

func MCTSExpand(node *GameTreeNode, wc Color, oBoard *Board) {
	board := NewBoardFromPath(node.Path2Root())
	board.Model = oBoard.Model
	oc := OpColor(node.stone)
	rank := board.CandidateMoves(oc, nil)
	topn := TopN(rank, 30)
	n := 0
	sg := make(chan byte, 100)
	sum := 0.0
	for _, child := range topn {
		sum += child.Second
	}
	for _, child := range topn {
		x, y := IndexPos(child.First)
		log.Println("expand node: ", x, y)
		cnode := NewGameTreeNode(oc, x, y)
		node.AddChild(cnode)
		tt := int(100.0*(child.Second/sum) + 0.5)
		for s := 0; s < tt; s++ {
			go MCTSSimulation(board.Copy(), cnode, wc, sg)
			n += 1
		}
	}
	for _ = range sg {
		n -= 1
		if n == 0 {
			break
		}
	}
	close(sg)
}

func MCTSSimulation(b *Board, next *GameTreeNode, wc Color, sg chan byte) {
	b.Put(PosIndex(next.x, next.y), next.stone)
	b.SelfBattle(OpColor(next.stone))
	s := b.Score()
	if (s > 0 && wc == BLACK) || (s <= 0 && wc == WHITE) {
		MCTSBackProp(next, 1)
	} else {
		MCTSBackProp(next, 0)
	}
	sg <- 1
}

func MCTSBackProp(node *GameTreeNode, win int) {
	MCTSLock.Lock()
	defer MCTSLock.Unlock()
	v := node
	for {
		if v == nil {
			return
		}
		v.visit += 1
		v.win += win
		v = v.Father
	}
}
