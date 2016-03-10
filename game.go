package gogo

import (
	"fmt"

	"github.com/xlvector/hector/dt"
)

type Game struct {
	B      *Board
	GT     *GameTree
	QipuGT *GameTree
	Komi   float64
}

func NewGame(size int) *Game {
	return &Game{
		B:  NewBoard(size),
		GT: NewGameTree(size),
	}
}

func (g *Game) SetModel(models []*dt.RandomForest) {
	g.B.models = models
}

func (g *Game) Init(size int) {
	g.B = NewBoard(size)
	g.GT = NewGameTree(size)
}

func (g *Game) SetKomi(komi float64) {
	g.Komi = komi
	g.B.komi = komi
}

func (g *Game) Clear() {
	g.Init(g.B.size)
}

func (g *Game) Put(stone Color, x, y int) error {
	if err := g.B.Put(x, y, stone); err != nil {
		return err
	}
	g.GT.Add(NewGameTreeNode(stone, x, y))
	return nil
}

func (g *Game) GenMove(stone Color) (int, int) {
	//from qipu
	if g.QipuGT != nil {
		nexts := g.GT.NextMoveByQipu(g.QipuGT)
		for _, next := range nexts {
			if next.hit < 5 {
				continue
			}
			err := g.Put(next.stone, next.x, next.y)
			if err != nil {
				continue
			}
			fmt.Println("use qipu")
			return next.x, next.y
		}
	}
	lastStep := Point{g.GT.Current.x, g.GT.Current.y, g.GT.Current.stone}
	next := g.B.GenMove(lastStep, stone)
	g.GT.Add(NewGameTreeNode(stone, next.x, next.y))
	return next.x, next.y
}

func (g *Game) LastStep() Point {
	v := g.GT.Current
	return Point{v.x, v.y, v.stone}
}

func (g *Game) Undo() {
	g.GT.Back()
	g.ResetBoardFromGT()
}

func (g *Game) ResetBoardFromGT() error {
	path := g.GT.Path2Root()
	g.B = NewBoard(g.B.size)
	for i := len(path) - 1; i >= 0; i-- {
		v := path[i]
		if v.x < 0 || v.y < 0 {
			continue
		}
		err := g.B.Put(v.x, v.y, v.stone)
		if err != nil {
			return err
		}
	}
	return nil
}

//Print Board

func (g *Game) String() string {
	if g.B == nil {
		return ""
	}
	return g.B.String(g.LastStep())
}

func (g *Game) Print() {
	fmt.Println(g.String())
}
