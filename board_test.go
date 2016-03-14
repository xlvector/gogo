package gogo

import (
	"testing"
)

func TestRules(t *testing.T) {
	b := NewBoard()
	t.Log(b.String())

	b.PutLabel("BC3")
	t.Log(b.String())

	if ok := b.PutLabel("WC3"); ok {
		t.Error()
	}

	b.PutLabel("BD2")
	b.PutLabel("BE3")
	b.PutLabel("BD4")
	t.Log(b.String())

	/*
		white D3 is invalid
				A B C D E F G H J K L M N O P Q R S T
			19  . . . . . . . . . . . . . . . . . . .
			18  . . . . . . . . . . . . . . . . . . .
			17  . . . . . . . . . . . . . . . . . . .
			16  . . . . . . . . . . . . . . . . . . .
			15  . . . . . . . . . . . . . . . . . . .
			14  . . . . . . . . . . . . . . . . . . .
			13  . . . . . . . . . . . . . . . . . . .
			12  . . . . . . . . . . . . . . . . . . .
			11  . . . . . . . . . . . . . . . . . . .
			10  . . . . . . . . . . . . . . . . . . .
			 9  . . . . . . . . . . . . . . . . . . .
			 8  . . . . . . . . . . . . . . . . . . .
			 7  . . . . . . . . . . . . . . . . . . .
			 6  . . . . . . . . . . . . . . . . . . .
			 5  . . . . . . . . . . . . . . . . . . .
			 4  . . . X . . . . . . . . . . . . . . .
			 3  . . X . X . . . . . . . . . . . . . .
			 2  . . . X . . . . . . . . . . . . . . .
			 1  . . . . . . . . . . . . . . . . . . .
			    A B C D E F G H J K L M N O P Q R S T
	*/

	if ok := b.PutLabel("WD3"); ok {
		t.Error()
	}

	/*
		white D3 is valid, but then BE3 is invalid because of ko
				A B C D E F G H J K L M N O P Q R S T
			19  . . . . . . . . . . . . . . . . . . .
			18  . . . . . . . . . . . . . . . . . . .
			17  . . . . . . . . . . . . . . . . . . .
			16  . . . . . . . . . . . . . . . . . . .
			15  . . . . . . . . . . . . . . . . . . .
			14  . . . . . . . . . . . . . . . . . . .
			13  . . . . . . . . . . . . . . . . . . .
			12  . . . . . . . . . . . . . . . . . . .
			11  . . . . . . . . . . . . . . . . . . .
			10  . . . . . . . . . . . . . . . . . . .
			 9  . . . . . . . . . . . . . . . . . . .
			 8  . . . . . . . . . . . . . . . . . . .
			 7  . . . . . . . . . . . . . . . . . . .
			 6  . . . . . . . . . . . . . . . . . . .
			 5  . . . . . . . . . . . . . . . . . . .
			 4  . . . X O . . . . . . . . . . . . . .
			 3  . . X . X O . . . . . . . . . . . . .
			 2  . . . X O . . . . . . . . . . . . . .
			 1  . . . . . . . . . . . . . . . . . . .
			    A B C D E F G H J K L M N O P Q R S T
	*/
	b.PutLabel("WE4")
	b.PutLabel("WF3")
	b.PutLabel("WE2")
	t.Log(b.String())
	if ok := b.PutLabel("WD3"); !ok {
		t.Error()
	}
	t.Log(b.String())

	if ok := b.PutLabel("BE3"); ok {
		t.Error()
	}
	t.Log(b.String())

	/*
		stable eye shoud not fill, D3 should not fill, but J3 can
				A B C D E F G H J K L M N O P Q R S T
			19  . . . . . . . . . . . . . . . . . . .
			18  . . . . . . . . . . . . . . . . . . .
			17  . . . . . . . . . . . . . . . . . . .
			16  . . . . . . . . . . . . . . . . . . .
			15  . . . . . . . . . . . . . . . . . . .
			14  . . . . . . . . . . . . . . . . . . .
			13  . . . . . . . . . . . . . . . . . . .
			12  . . . . . . . . . . . . . . . . . . .
			11  . . . . . . . . . . . . . . . . . . .
			10  . . . . . . . . . . . . . . . . . . .
			 9  . . . . . . . . . . . . . . . . . . .
			 8  . . . . . . . . . . . . . . . . . . .
			 7  . . . . . . . . . . . . . . . . . . .
			 6  . . . . . . . . . . . . . . . . . . .
			 5  . . . . . . . . . . . . . . . . . . .
			 4  . . . X . . . O X . . . . . . . . . .
			 3  . . X . X . . X . X . . . . . . . . .
			 2  . . . X . . . . X O . . . . . . . . .
			 1  . . . . . . . . . . . . . . . . . . .
			    A B C D E F G H J K L M N O P Q R S T
	*/
	b.Clear()
	b.PutLabel("BC3")
	b.PutLabel("BD2")
	b.PutLabel("BD4")
	b.PutLabel("BE3")
	b.PutLabel("BH3")
	b.PutLabel("BJ2")
	b.PutLabel("BJ4")
	b.PutLabel("BK3")
	b.PutLabel("WH4")
	b.PutLabel("WK2")
	t.Log(b.String())
	if ok := b.StableEye(PosIndex(3, 2), BLACK); !ok {
		t.Error()
	}
	if ok := b.PutLabel("BD3"); ok {
		t.Error()
	}

	if ok := b.StableEye(PosIndex(8, 2), BLACK); ok {
		t.Error()
	}
	if ok := b.PutLabel("BJ3"); !ok {
		t.Error()
	}
}

func TestAtari(t *testing.T) {
	/*
		W[E3] true, W[F3] false
				A B C D E F G H J K L M N O P Q R S T
			19  . . . . . . . . . . . . . . . . . . .
			18  . . . . . . . . . . . . . . . . . . .
			17  . . . . . . . . . . . . . . . . . . .
			16  . . . . . . . . . . . . . . . . . . .
			15  . . . . . . . . . . . . . . . . . . .
			14  . . . . . . . . . . . . . . . . . . .
			13  . . . . . . . . . . . . . . . . . . .
			12  . . . . . . . . . . . . . . . . . . .
			11  . . . . . . . . . . . . . . . . . . .
			10  . . . . . . . . . . . . . . . . . . .
			 9  . . . . . . . . . . . . . . . . . . .
			 8  . . . . . . . . . . . . . . . . . . .
			 7  . . . . . . . . . . . . . . . . . . .
			 6  . . . . . . . . . . . . . . . . . . .
			 5  . . . . . . . . . . . . . . . . . . .
			 4  . . . X . . . . . . . . . . . . . . .
			 3  . . X O O . . . . . . . . . . . . . .
			 2  . . . X . . . . . . . . . . . . . . .
			 1  . . . . . . . . . . . . . . . . . . .
			    A B C D E F G H J K L M N O P Q R S T
	*/
	b := NewBoard()
	b.PutLabel("BC3")
	b.PutLabel("BD2")
	b.PutLabel("BD4")
	b.PutLabel("WD3")

	t.Log(b.String())

	if !b.EscapeAtari(PosIndex(4, 2), WHITE) {
		t.Error()
	}

	if b.EscapeAtari(PosIndex(5, 2), WHITE) {
		t.Error()
	}

	/*
		W[E1] false
				A B C D E F G H J K L M N O P Q R S T
			19  . . . . . . . . . . . . . . . . . . .
			18  . . . . . . . . . . . . . . . . . . .
			17  . . . . . . . . . . . . . . . . . . .
			16  . . . . . . . . . . . . . . . . . . .
			15  . . . . . . . . . . . . . . . . . . .
			14  . . . . . . . . . . . . . . . . . . .
			13  . . . . . . . . . . . . . . . . . . .
			12  . . . . . . . . . . . . . . . . . . .
			11  . . . . . . . . . . . . . . . . . . .
			10  . . . . . . . . . . . . . . . . . . .
			 9  . . . . . . . . . . . . . . . . . . .
			 8  . . . . . . . . . . . . . . . . . . .
			 7  . . . . . . . . . . . . . . . . . . .
			 6  . . . . . . . . . . . . . . . . . . .
			 5  . . . . . . . . . . . . . . . . . . .
			 4  . . . X X . . . . . . . . . . . . . .
			 3  . . X O O X . . . . . . . . . . . . .
			 2  . . . X O X . . . . . . . . . . . . .
			 1  . . . . . . . . . . . . . . . . . . .
			    A B C D E F G H J K L M N O P Q R S T
	*/
	b.PutLabel("BE4")
	b.PutLabel("BF2")
	b.PutLabel("BF3")
	b.PutLabel("WE3")
	b.PutLabel("WE2")
	t.Log(b.String())
	if b.EscapeAtari(PosIndex(4, 0), WHITE) {
		t.Error()
	}

	nworms := b.NeighWorms(PosIndex(4, 0), WHITE, WHITE)
	if len(nworms) != 1 {
		t.Error()
	}
	if len(nworms[0].Points) != 3 {
		t.Error()
	}
	if nworms[0].Liberty != 1 {
		t.Error()
	}
}

func TestHash(t *testing.T) {
	b := NewBoard()
	m := make(map[int64]byte)
	for _, h := range b.PointHash {
		m[h] = 1
	}
	t.Log(len(m))

	hp := make(map[int64]map[int]byte)
	for d := 0; d < PATTERN_SIZE; d++ {
		for k := 0; k < NPOINT; k++ {
			h := b.PatternHash[k][d]
			if _, ok := hp[h]; !ok {
				hp[h] = make(map[int]byte)
			}
			hp[h][k*100+d] = 1
		}
	}

	for h, v := range hp {
		if len(v) > 1 {
			t.Log(h)
			for k, _ := range v {
				x, y := IndexPos(k / 100)
				d := k % 100
				t.Log(d, x, y)
			}
		}
	}

	/*
		W[E1] false
				A B C D E F G H J K L M N O P Q R S T
			19  . . . . . . . . . . . . . . . . . . .
			18  . . . . . . . . . . . . . . . . . . .
			17  . . . . . . . . . . . . . . . . . . .
			16  . . . . . . . . . . . . . . . . . . .
			15  . . . . . . . . . . . . . . . . . . .
			14  . . . . . . . . . . . . . . . . . . .
			13  . . . . . . . . . . . . . . . . . . .
			12  . . . . . . . . . . . . . . . . . . .
			11  . . . . . . . . . . . . . . . . . . .
			10  . . . . . . . . . . . . . . . . . . .
			 9  . . . . . . . . . . . . . . . . . . .
			 8  . . . . . . . . . . . . . . . . . . .
			 7  . . . . . . . . . . . . . . . . . . .
			 6  . . . . . . . . . . . . . . . . . . .
			 5  . . . . . . . . . . . . . . . . . . .
			 4  . . . X X . . . . . . . . . . . . . .
			 3  . . X O O X . . . . . . . . . . . . .
			 2  . . . X X X . . . . . . . . . . . . .
			 1  O . . . . . . . . . . . . . . . . . .
			    A B C D E F G H J K L M N O P Q R S T
	*/

	b.PutLabel("BC3")
	b.PutLabel("WA1")
	b.PutLabel("BD2")
	b.PutLabel("BE2")
	b.PutLabel("BF2")
	b.PutLabel("BF3")
	b.PutLabel("BD4")
	b.PutLabel("WD3")
	b.PutLabel("WE3")
	t.Log(b.String())
	b.PutLabel("BE4")
	t.Log(b.String())
	b.PutLabel("BA2")
	b.PutLabel("BB1")
	t.Log(b.String())
	tmpPatternHash := make([][]int64, NPOINT)

	b2 := b.Copy()

	for i := 0; i < NPOINT; i++ {
		tmpPatternHash[i] = make([]int64, PATTERN_SIZE)
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
				tmpPatternHash[i][d] ^= b.VertexHash(x1, y1, c)
			}
		}

		for d := 0; d < PATTERN_SIZE; d++ {
			if tmpPatternHash[i][d] != b.PatternHash[i][d] {
				t.Error(i, d, tmpPatternHash[i][d], b.PatternHash[i][d])
			}
			if tmpPatternHash[i][d] != b2.PatternHash[i][d] {
				t.Error(i, d, tmpPatternHash[i][d], b2.PatternHash[i][d])
			}
		}
	}

}