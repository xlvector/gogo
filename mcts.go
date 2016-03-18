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
	for n < 350 {
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
	ret := 0.0
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

func (b *Board) GenMove(c Color, rank map[int]float64) (int, map[int]float64) {
	rank = b.CandidateMoves(c, rank)
	cands := TopN(rank, 30)
	if len(cands) == 0 {
		return -1, rank
	}
	pf := cands[rand.Intn(len(cands))].First
	b.Put(pf, c)
	delete(rank, pf)
	return pf, rank
}

func (b *Board) GenBestMove(c Color, gt *GameTree) bool {
	rank := b.CandidateMoves(c, nil)
	cands := TopN(rank, 1)
	if len(cands) == 0 {
		return false
	}
	if ok := b.Put(cands[0].First, c); ok {
		x, y := IndexPos(cands[0].First)
		gt.Add(NewGameTreeNode(c, x, y))
		return true
	}
	return false
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
	ret += math.Sqrt(math.Log(np) / float64(1+p.visit))
	return ret
}

func (b *Board) MCTSMove(c Color, gt *GameTree, expand, n int) bool {
	wg := &sync.WaitGroup{}
	root := gt.Current
	log.Println(PointString(root.x, root.y, root.stone), root.win, root.visit, "next stone color: ", ColorMark(c))
	if len(root.Children) == 0 {
		for _, child := range root.Children {
			log.Println(PointString(child.x, child.y, child.stone), child.win, child.visit)
		}
	}
	for i := 0; i < n; i++ {
		node := MCTSSelection(gt)
		MCTSExpand(node, b, expand, wg)
	}
	wg.Wait()
	var best *GameTreeNode
	robust := 0
	for _, child := range root.Children {
		winrate := float64(child.win) / float64(child.visit)
		log.Println(PointString(child.x, child.y, child.stone), winrate, child.win, child.visit, child.prior)
		if robust < child.win {
			robust = child.win
			best = child
		}
	}
	gt.Current = best
	return b.Put(PosIndex(best.x, best.y), c)
}

func MCTSSelection(gt *GameTree) *GameTreeNode {
	root := gt.Current
	ret := root
	depth := 0
	for {
		if ret.Children == nil || len(ret.Children) == 0 {
			return ret
		}
		if len(ret.CandMoves) > 0 {
			return ret
		}
		depth += 1
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

func MCTSExpand(node *GameTreeNode, oBoard *Board, nLeaf int, wg *sync.WaitGroup) {
	board := NewBoardFromPath(node.Path2Root())
	board.Model = oBoard.Model
	oc := BLACK
	if node.stone == BLACK || node.stone == WHITE {
		oc = OpColor(node.stone)
	}

	if len(node.CandMoves) == 0 {
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
	wg.Add(1)
	go MCTSSimulation(board.Copy(), cnode, wg)
}

func MCTSSimulation(b *Board, next *GameTreeNode, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()
	b.Put(PosIndex(next.x, next.y), next.stone)
	b.SelfBattle(OpColor(next.stone))
	s := b.Score()

	if s > 0 {
		MCTSBackProp(next, BLACK)
	} else {
		MCTSBackProp(next, WHITE)
	}
}

func MCTSBackProp(node *GameTreeNode, wc Color) {
	MCTSLock.Lock()
	defer MCTSLock.Unlock()
	v := node
	for {
		if v == nil {
			return
		}
		v.visit += 1
		if v.stone == wc {
			v.win += 1
		}
		v = v.Father
	}
}
