package gogo

import (
	"fmt"
	"testing"
)

func TestBoard(t *testing.T) {
	g := &Game{}
	g.Init(19)
	g.Put(BLACK, 3, 3)
	g.Put(WHITE, 3, 4)
	g.Put(BLACK, 4, 3)
	g.Print()
	fmt.Println(g.GT.WriteSGF())

	g.Undo()
	g.Undo()
	g.Put(WHITE, 4, 4)
	g.Put(BLACK, 5, 5)
	g.Put(WHITE, 6, 6)
	g.Put(BLACK, 7, 7)
	g.Print()
	fmt.Println(g.GT.WriteSGF())

	g.Undo()
	g.Undo()
	g.Put(WHITE, 6, 7)
	g.Put(BLACK, 7, 8)
	g.Print()
	fmt.Println(g.GT.WriteSGF())

	fmt.Println("parse sgf")
	g1 := &Game{}
	g1.Init(19)
	g1.GT.ParseSGF(g.GT.WriteSGF())
	g1.ResetBoardFromGT()

	g1.Print()
}

func TestWorm(t *testing.T) {
	g := &Game{}
	g.Init(9)
	g.Put(BLACK, 2, 2)
	g.Put(BLACK, 3, 1)
	g.Put(BLACK, 4, 2)
	g.Put(WHITE, 3, 2)
	g.Print()

	g.Put(BLACK, 3, 3)
	g.Print()

	if g.B.Get(3, 2).color != GRAY {
		t.Error()
	}

	g.Undo()
	g.Print()
	if g.B.Get(3, 2).color != WHITE {
		t.Error()
	}
}

func TestValid(t *testing.T) {

	/*

	   	32 is not valid

	        0 1 2 3 4 5 6 7 8
	      0 . . . . . . . . .
	      1 . . . X . . . . .
	      2 . . X . X . . . .
	      3 . . . X . . . . .
	      4 . . . . . . . . .
	      5 . . . . . . . . .
	      6 . . . . . . . . .
	      7 . . . . . . . . .
	      8 . . . . . . . . .


	*/

	g := &Game{}
	g.Init(9)
	g.Put(BLACK, 2, 2)
	g.Put(BLACK, 3, 1)
	g.Put(BLACK, 4, 2)
	g.Put(BLACK, 3, 3)
	g.Print()

	if err := g.Put(WHITE, 3, 2); err == nil {
		t.Error()
	}
	g.Print()

	/*

	   	00 is valid

	        0 1 2 3 4 5 6 7 8
	      0 . X O . . . . . .
	      1 X O . . . . . . .
	      2 O . . . . . . . .
	      3 . . . . . . . . .
	      4 . . . . . . . . .
	      5 . . . . . . . . .
	      6 . . . . . . . . .
	      7 . . . . . . . . .
	      8 . . . . . . . . .
	*/
	g.Clear()
	g.Put(BLACK, 1, 0)
	g.Put(BLACK, 0, 1)
	g.Put(WHITE, 2, 0)
	g.Put(WHITE, 1, 1)
	g.Put(WHITE, 0, 2)
	g.Print()

	if err := g.Put(WHITE, 0, 0); err != nil {
		t.Error(err)
	}
	g.Print()
}

func TestSelfBattle(t *testing.T) {
	g := &Game{}
	g.Init(19)
	g.B.SelfBattle(InvalidPoint(), BLACK)
	g.Print()
}

func TestKoRule(t *testing.T) {
	g := &Game{}
	g.Init(9)

	/*

	   	00 is valid

	        0 1 2 3 4 5 6 7 8
	      0 . X . . . . . . .
	      1 X O X . . . . . .
	      2 O . O . . . . . .
	      3 . O . . . . . . .
	      4 . . . . . . . . .
	      5 . . . . . . . . .
	      6 . . . . . . . . .
	      7 . . . . . . . . .
	      8 . . . . . . . . .
	*/
	g.Put(WHITE, 2, 2)
	g.Put(BLACK, 1, 0)
	g.Put(WHITE, 1, 1)
	g.Put(BLACK, 0, 1)
	g.Put(WHITE, 0, 2)
	g.Put(BLACK, 2, 1)
	g.Put(WHITE, 1, 3)
	g.Print()

	if err := g.Put(BLACK, 1, 2); err != nil {
		t.Error(err)
	}
	g.Print()

	if err := g.Put(WHITE, 1, 1); err != nil {
		if err.Error() != "ko position" {
			t.Error(err)
		}
	}
	g.Print()
}

func TestGenMovesOfMyLiberty(t *testing.T) {
	g := &Game{}
	g.Init(9)

	/*

	   	13 is important

	        0 1 2 3 4 5 6 7 8
	      0 . X . . . . . . .
	      1 X O X . . . . . .
	      2 X O X . . . . . .
	      3 . . . . . . . . .
	      4 . . . . . . . . .
	      5 . . . . . . . . .
	      6 . . . . . . . . .
	      7 . . . . . . . . .
	      8 . . . . . . . . .
	*/

	g.Put(BLACK, 1, 0)
	g.Put(BLACK, 0, 1)
	g.Put(BLACK, 0, 2)
	g.Put(WHITE, 1, 1)
	g.Put(WHITE, 1, 2)
	g.Put(BLACK, 2, 1)
	g.Put(BLACK, 2, 2)

	g.Print()

	info := g.B.CollectBoardInfo(InvalidPoint())
	pos := info.genMovesOfMyLiberty(WHITE)
	if len(pos) <= 0 {
		t.Error()
	}
	for p, _ := range pos {
		c := g.B.w[p]
		if c.x != 1 || c.y != 3 {
			t.Error()
		}
		g.Put(WHITE, c.x, c.y)
		break
	}
	g.Print()
}

func TestGenMovesOfMyLiberty2(t *testing.T) {
	g := &Game{}
	g.Init(9)

	/*

	   	14 is important

	        0 1 2 3 4 5 6 7 8
	      0 . . . . . . . . .
	      1 . . . . . . . . .
	      2 . . X X X . . . .
	      3 . X O O X . . . .
	      4 . . O X . . . . .
	      5 . O X X . . . . .
	      6 . O O O . . . . .
	      7 . . . . . . . . .
	      8 . . . . . . . . .
	*/

	g.Put(BLACK, 1, 3)
	g.Put(BLACK, 2, 2)
	g.Put(BLACK, 2, 5)
	g.Put(BLACK, 3, 2)
	g.Put(BLACK, 3, 4)
	g.Put(BLACK, 3, 5)
	g.Put(BLACK, 4, 2)
	g.Put(BLACK, 4, 3)
	g.Put(WHITE, 2, 3)
	g.Put(WHITE, 2, 4)
	g.Put(WHITE, 1, 5)
	g.Put(WHITE, 1, 6)
	g.Put(WHITE, 2, 6)
	g.Put(WHITE, 3, 3)
	g.Put(WHITE, 3, 6)

	g.Print()

	info := g.B.CollectBoardInfo(InvalidPoint())
	pos := info.genMovesOfMyLiberty(WHITE)
	if len(pos) <= 0 {
		t.Error()
	}
	for p, _ := range pos {
		c := g.B.w[p]
		if c.x != 1 || c.y != 4 {
			t.Error()
		}
		g.Put(WHITE, c.x, c.y)
		break
	}
	g.Print()
}

func TestGenMovesOfMyLiberty3(t *testing.T) {
	g := &Game{}
	g.Init(9)

	/*

	   	20 is important

	        0 1 2 3 4 5 6 7 8
	      0 O O . O O . O . .
	      1 O O O X O O O . .
	      2 X X X X X X . . .
	      3 . . . . . . . . .
	      4 . . . . . . . . .
	      5 . . . . . . . . .
	      6 . . . . . . . . .
	      7 . . . . . . . . .
	      8 . . . . . . . . .
	*/

	g.Put(BLACK, 0, 2)
	g.Put(BLACK, 1, 2)
	g.Put(BLACK, 2, 2)
	g.Put(BLACK, 3, 1)
	g.Put(BLACK, 3, 2)
	g.Put(BLACK, 4, 2)
	g.Put(BLACK, 5, 2)
	g.Put(WHITE, 0, 0)
	g.Put(WHITE, 0, 1)
	g.Put(WHITE, 1, 0)
	g.Put(WHITE, 1, 1)
	g.Put(WHITE, 2, 1)
	g.Put(WHITE, 3, 0)
	g.Put(WHITE, 4, 0)
	g.Put(WHITE, 4, 1)
	g.Put(WHITE, 5, 1)
	g.Put(WHITE, 6, 1)
	g.Put(WHITE, 6, 0)
	g.Print()

	info := g.B.CollectBoardInfo(InvalidPoint())
	pos := info.genMovesOfMyLiberty(WHITE)
	if len(pos) <= 0 {
		t.Error()
	}
	for p, _ := range pos {
		c := g.B.w[p]
		if c.x != 2 || c.y != 0 {
			t.Error()
		}
		g.Put(WHITE, c.x, c.y)
		break
	}
	g.Print()
}

func TestGenFillEyeMoves(t *testing.T) {
	g := &Game{}
	g.Init(9)

	/*

	   	11 should fill

	        0 1 2 3 4 5 6 7 8
	      0 . X . . . . . . .
	      1 X . X . . . . . .
	      2 O X O . . . . . .
	      3 . O . . . . . . .
	      4 . . . . . . . . .
	      5 . . . . . . . . .
	      6 . . . . . . . . .
	      7 . . . . . . . . .
	      8 . . . . . . . . .
	*/

	g.Put(BLACK, 1, 0)
	g.Put(BLACK, 0, 1)
	g.Put(BLACK, 2, 1)
	g.Put(BLACK, 1, 2)
	g.Put(WHITE, 0, 2)
	g.Put(WHITE, 1, 3)
	g.Put(WHITE, 2, 2)

	g.Print()

	info := g.B.CollectBoardInfo(InvalidPoint())
	pos := info.genFillEyeMoves(BLACK)
	if len(pos) <= 0 {
		t.Error()
	}
	for p, _ := range pos {
		c := g.B.w[p]
		if c.x != 1 || c.y != 1 {
			t.Error(c.String())
		}
		g.Put(BLACK, c.x, c.y)
		break
	}
	g.Print()
}

func TestDragon(t *testing.T) {
	g := &Game{}
	g.Init(9)

	/*

	   	1 dragon

	        0 1 2 3 4 5 6 7 8
	      0 . X . . . . . . .
	      1 X . X . . . . . .
	      2 . . . . X . . . .
	      3 X X X X . X . . .
	      4 . . . . X . . . .
	      5 . . . . . . . . .
	      6 . . . . . . . . .
	      7 . . . . . . . . .
	      8 . . . . . . . . .
	*/

	g.Put(BLACK, 1, 0)
	g.Put(BLACK, 0, 1)
	g.Put(BLACK, 2, 1)
	g.Put(BLACK, 0, 3)
	g.Put(BLACK, 1, 3)
	g.Put(BLACK, 2, 3)
	g.Put(BLACK, 3, 3)
	g.Put(BLACK, 4, 2)
	g.Put(BLACK, 4, 4)
	g.Put(BLACK, 5, 3)

	g.Print()

	info := g.B.CollectBoardInfo(InvalidPoint())
	dragons := info.BuildDragon()

	if len(dragons) != 1 {
		t.Error(len(dragons))
	}

	if len(dragons[0].Worms) != 7 {
		t.Error(len(dragons[0].Worms))
	}
}

func TestPatternHash(t *testing.T) {
	g := &Game{}
	g.Init(19)
	pdm := NewPointDistanceMap(g.B, PATTERN_SIZE)
	g.B.SetPointDistanceMap(pdm)
	g.Put(BLACK, 2, 2)
	g.Put(WHITE, 4, 5)
	g.Put(BLACK, 9, 11)

	for p, dh := range g.B.patternHash {
		newDh := make([]int64, len(dh))
		for d := 0; d < len(dh); d++ {
			qd := pdm.PointDistance(p, d)
			for _, q := range qd {
				newDh[d] ^= g.B.PointHash(g.B.w[q])
			}
			if dh[d] != newDh[d] {
				t.Error(dh[d], newDh[d])
			}

			//fmt.Println(g.B.w[p].String(), d, dh[d])
		}
	}

}
