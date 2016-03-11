package gogo

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

const (
	FS_ZERO = iota
	FS_PAT33
	FS_EDGE_DIS
	FS_EDGE_DIS_X
	FS_EDGE_DIS_Y
)

func (b *Board) PointSimpleFeature(p Point, stone Color) []int64 {
	fh := b.FeatureHash(MakePoint(p.x, p.y, stone))
	if fh < 0 {
		return nil
	}
	ph := b.GetPatternHash(b.index(p.x, p.y))
	for i := 0; i < len(ph); i++ {
		ph[i] ^= fh
	}
	return ph
}

func (b *Board) GenSimpleFeatures(lastPat []int64, cur Point) map[int][]int64 {
	ret := make(map[int][]int64)
	for i, p := range b.w {
		if p.color != GRAY {
			continue
		}
		if (p.x == cur.x && p.y == cur.y) || rand.Float64() < 0.03 {
			pat := b.PointSimpleFeature(p, cur.color)
			if pat == nil {
				continue
			}
			pat = append(pat, lastPat...)
			ret[i] = pat
		}
	}
	return ret
}

func NewBoardFromPath(size int, path []*GameTreeNode) *Board {
	ret := NewBoard(size)
	for i := len(path) - 1; i >= 0; i-- {
		v := path[i]
		ret.Put(v.x, v.y, v.stone)
	}
	return ret
}

func (b *Board) SelfBattle(lastMove Point, color Color) Color {
	rand.Seed(time.Now().UnixNano())

	for {
		pass := 0
		lastMove = b.GenQuickMove(lastMove, color)
		if !lastMove.Valid() {
			pass += 1
		}

		lastMove = b.GenQuickMove(lastMove, OppColor(color))
		if !lastMove.Valid() {
			pass += 1
		}
		if pass >= 2 {
			break
		}
	}
	score := b.Score()
	if score > 0 {
		return BLACK
	} else {
		return WHITE
	}
}

func (p *GameTreeNode) BackPropVisit(root *GameTreeNode) {
	b := p
	for {
		if b == nil {
			break
		}
		b.visit += 1
		if b == root {
			break
		}
		b = b.Father
	}
}

func (p *GameTreeNode) BackPropWin(root *GameTreeNode) {
	b := p
	for {
		if b == nil {
			break
		}
		b.win += 1
		if b == root {
			break
		}
		b = b.Father
	}
}

func (g *Game) MCTreePolicy() *GameTreeNode {
	root := g.GT.Current
	node := root
	for {
		if len(node.Children) > 0 {
			maxScore := 0.0
			for _, child := range node.Children {
				score := float64(child.win) / float64(child.visit)
				score += math.Sqrt(2.0 * math.Log(float64(node.visit)) / float64(child.visit))
				if maxScore < score {
					maxScore = score
					node = child
				}
			}
		} else {
			return node
		}
	}
	return nil
}

type SingleBattleResult struct {
	node     *GameTreeNode
	winColor Color
}

func (g *Game) singleSimulate(newBoard *Board, gn *GameTreeNode, pm Point, ch chan SingleBattleResult) {
	winColor := newBoard.Copy().SelfBattle(pm, OppColor(pm.color))
	ch <- SingleBattleResult{gn, winColor}
}

func (g *Game) MCTSMove(stone Color) {
	root := g.GT.Current
	for root.visit < 1000 {
		fmt.Println(root.visit)
		node := g.MCTreePolicy()
		board := NewBoardFromPath(g.B.size, node.Path2Root())
		//info := board.CollectBoardInfo(InvalidPoint())
		cand := board.QuickCandidateMoves(Point{node.x, node.y, node.stone}, OppColor(node.stone), 20)
		//cand := info.CandidateMoves(Point{node.x, node.y, node.stone}, OppColor(node.stone), g.B.Model(), 5)
		ch := make(chan SingleBattleResult, len(cand)+1)
		n := 0
		for m, v := range cand {
			pm := board.w[m]
			if node == root {
				fmt.Printf("%s%d[%f] ", string(LX[pm.x]), pm.y+1, v)
			}
			pm.color = OppColor(node.stone)
			newBoard := board.Copy()
			if err := newBoard.Put(pm.x, pm.y, pm.color); err != nil {
				continue
			}
			gn := NewGameTreeNode(pm.color, pm.x, pm.y)
			node.AddChild(gn)
			for i := 0; i < 5; i++ {
				n += 1
				go g.singleSimulate(newBoard, gn, pm, ch)
			}
		}
		for i := 0; i < n; i++ {
			sbr := <-ch
			if sbr.winColor == stone {
				sbr.node.BackPropWin(root)
			}
			sbr.node.BackPropVisit(root)
		}
		close(ch)
		if node == root {
			fmt.Println()
		}
		if len(node.Children) == 0 {
			break
		}
	}

	maxRate := 0.0
	bestMove := root.Children[0]

	for _, child := range root.Children {
		winrate := float64(child.win) / float64(child.visit)
		fmt.Println(string(LX[child.x]), child.y+1, child.win, child.visit, winrate)
		if winrate > maxRate {
			maxRate = winrate
			bestMove = child
		}
	}

	root.Children = []*GameTreeNode{}
	g.Put(bestMove.stone, bestMove.x, bestMove.y)
}
