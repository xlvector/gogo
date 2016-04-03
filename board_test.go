package gogo

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"
)

func TestCaseStudy(t *testing.T) {
	b := NewBoard()
	buf := "BR4,WD3,BC16,WQ17,BE17,WP3,BD5,WE5,BC3,WC4,BD4,WB3,BE3,WC2,BD2,WQ5,BC3,WR2,BQ4,WP4,BP2,WR5,BO3,WO4,BN3,WM5,BD3,WM16,BM15,WL16,BN15,WR15,BR13,WC10,BC8,WQ9,BN4,WN5,BL4,WO17,BQ14,WC13,BD9,WK6,BL5,WL6,BJ7,WE6,BJ6,WD6,BC6,WB1,BC5,WA2,BD10,WF4,BC9,WC11,BF3,WF12,BE11,WK5,BJ5,WK4,BG4,WS4,BS5,WS6,BF11,WG17,BF16,WG12,BG16,WH17,BL3,WH3,BH4,WJ3,BQ2,WQ3,BH16,WH10,BK10,WG3,BE4,WF5,BE7,WG5,BJ9,WH12,BJ12,WJ11,BK11,WE12,BD14,WF14,BK16,WK18,BL15,WB4,BD13,WD12,BK17,WF18,BL18,WC14,BJ13,WR11,BP10,WN7,BP9,WQ12,BQ10,WS12,BR9,WS10,BR10,WS9,BR7,WS13,BL8,WD15,BE15,WE14,BD16,WE18,BB18,WG19,BJ18,WH18,BC15,WQ15,BO16,WP14,BO14,WN17,BP13,WP15,BQ13,WP12,BO12,WN13,BO13,WP16,BR12,WG9,BG8,WF8,BF7,WQ11,BP11,WS11,BM7,WM6,BO8,WN8,BN9,WB14,BB10,WE9,BE8,WE13,BB11,WB12,BD11,WF9,BC12,WC11,BM8,WG7,BH8,WD15,BB5,WK2,BL2,WJ4,BH5,WO7,BP6,WO6,BP5,WO5,BQ6,WR6,BS7,WT5,BR1,WS1,BQ1,WO2,BO1,WN2,BM2,WP1,BD14,WA13,BB15,WG11,BG14,WQ7,BQ8,WP7,BS14,WR14,BT14,WT7,BT8,WS8,BT6,WL13,BR8,WH14,BG13,WH13,BG15,WJ10,BF10,WG10,BJ14,WM18,BT7,WD18,BC10,WJ17,BJ16,WJ19,BC12,WL17,BC11,WO15,BG6,WN16,BT9,WG2,BF2,WC18,BD15,WD13,BA11,WA12,BT10,WS15,BT15,WT13,BT16,WS16,BT17,WS17,BT11,WT18,BF6,WB17,BC17,WA17,BA16,WA14,BA18,WB19"
	tks := strings.Split(buf, ",")
	for _, tk := range tks {
		b.PutLabel(tk)
	}
	t.Log(b.String(nil))
	t.Log(b.LocalFeature(PosIndex(1, 15), BLACK))
}

func TestSelfBattleCase(t *testing.T) {
	b := NewBoard()
	buf := "BR4,WD3,BC16,WQ17,BE17,WP3,BD5,WE5,BC3,WC4,BD4,WB3,BE3,WC2,BD2,WQ5,BC3,WR2,BQ4,WP4,BP2,WR5,BO3,WO4,BN3,WM5,BD3,WM16,BM15,WL16,BN15,WR15,BR13,WC10,BC8,WQ9,BN4,WN5,BL4,WO17,BQ14,WC13,BD9,WK6,BL5,WL6,BJ7,WE6,BJ6,WD6,BC6,WB1,BC5,WA2,BD10,WF4,BC9,WC11,BF3,WF12,BE11,WK5,BJ5,WK4,BG4,WS4,BS5,WS6,BF11,WG17,BF16,WG12,BG16,WH17,BL3,WH3,BH4,WJ3,BQ2,WQ3,BH16,WH10,BK10,WG3,BE4,WF5,BE7,WG5,BJ9,WH12,BJ12,WJ11,BK11,WE12,BD14,WF14,BK16,WK18,BL15,WB4,BD13,WD12,BK17,WF18,BL18,WC14,BJ13,WR11,BP10,WN7,BP9,WQ12,BQ10,WS12,BR9,WS10,BR10,WS9,BR7,WS13,BL8,WD15,BE15,WE14,BD16,WE18,BB18,WG19,BJ18,WH18,BC15,WQ15,BO16,WP14,BO14,WN17,BP13,WP15,BQ13,WP12,BO12,WN13,BO13,WP16,BR12,WG9,BG8,WF8,BF7,WQ11,BP11,WS11,BM7,WM6,BO8,WN8,BN9,WB14,BB10,WE9,BE8,WE13,BB11,WB12,BD11,WF9,BC12,WC11,BM8,WG7,BH8,WD15,BB5,WK2,BL2,WJ4,BH5,WO7,BP6,WO6,BP5,WO5,BQ6,WR6,BS7,WT5,BR1,WS1,BQ1,WO2,BO1,WN2,BM2,WP1,BD14,WA13,BB15,WG11,BG14,WQ7,BQ8,WP7,BS14,WR14,BT14,WT7,BT8,WS8,BT6,WL13,BR8,WH14,BG13,WH13,BG15,WJ10,BF10,WG10,BJ14,WM18,BT7,WD18,BC10,WJ17,BJ16,WJ19,BC12,WL17"
	tks := strings.Split(buf, ",")
	for _, tk := range tks {
		b.PutLabel(tk)
	}
	t.Log(b.String(nil))
	win := 0
	for i := 0; i < 100; i++ {
		b2 := b.Copy()
		b2.SelfBattle(BLACK, nil)
		score := b2.Score()
		if score > 0 {
			win += 1
		}
		if i == 0 {
			t.Log(b2.String(nil))
		}
	}
	t.Log(win)
}

func TestRules(t *testing.T) {
	b := NewBoard()
	t.Log(b.String(nil))

	b.PutLabel("BC3")
	t.Log(b.String(nil))

	if ok := b.PutLabel("WC3"); ok {
		t.Error()
	}

	b.PutLabel("BD2")
	b.PutLabel("BE3")
	b.PutLabel("BD4")
	t.Log(b.String(nil))

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
	t.Log(b.String(nil))
	if ok := b.PutLabel("WD3"); !ok {
		t.Error()
	}
	t.Log(b.String(nil))

	if ok := b.PutLabel("BE3"); ok {
		t.Error()
	}
	t.Log(b.String(nil))

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
	t.Log(b.String(nil))
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

	t.Log(b.String(nil))

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
	t.Log(b.String(nil))
	if b.EscapeAtari(PosIndex(4, 0), WHITE) {
		t.Error()
	}

	nworms := b.NeighWorms(PosIndex(4, 0), WHITE, WHITE, -1)
	if len(nworms) != 1 {
		t.Error()
	}
	if nworms[0].Size() != 3 {
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
	t.Log(b.String(nil))
	b.PutLabel("BE4")
	t.Log(b.String(nil))
	b.PutLabel("BA2")
	b.PutLabel("BB1")
	t.Log(b.String(nil))
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

func TestSingleSelfBattle(t *testing.T) {
	/*
			A B C D E F G H J K L M N O P Q R S T
		19  . O . . . . . . . . . . . . . . . . .
		18  X . O . . X . . . . . . . . . . . . .
		17  . . . O . X . . . . . . . . . . . . .
		16  X . X . X . . . . . . . . . . . . . .
		15  . X . . . . . . . . . . . . . . . . .
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
		 4  . . . . . . . . . . . . . . . . . . .
		 3  . . . . . . . . . . . . . . . . . . .
		 2  . . . . . . . . . . . . . . . . . . .
		 1  . . . . . . . . . . . . . . . . . . .
		    A B C D E F G H J K L M N O P Q R S T
	*/
	b := NewBoard()
	b.PutLabel("BA16")
	b.PutLabel("WB19")
	b.PutLabel("BC16")
	b.PutLabel("WC18")
	b.PutLabel("BE16")
	b.PutLabel("WD17")
	b.PutLabel("BF17")
	b.PutLabel("BA18")
	b.PutLabel("BF18")
	b.PutLabel("BB15")

	n := 0
	rand.Seed(time.Now().UnixNano())
	f := 10
	for n < 350 {
		pass := 0
		p := b.GenSelfBattleMove(WHITE, nil)
		if p < 0 {
			pass += 1
		}
		if n < f {
			t.Log(b.String(nil))
		}

		p = b.GenSelfBattleMove(BLACK, nil)
		if p < 0 {
			pass += 1
		}
		if n < f {
			t.Log(b.String(nil))
		}
		if pass >= 2 {
			break
		}
		n += 1
	}
	t.Log(b.String(nil))
	t.Log(b.Score())
}

/*
func BenchmarkSelfBattle(t *testing.B) {
	for i := 0; i < t.N; i++ {
		b := NewBoard()
		b.SelfBattle(BLACK)
		b.Score()
	}
}
*/

func TestSelfBattle(t *testing.T) {
	win := 0
	lgr := NewLastGoodReply()
	b := NewBoard()
	/*
			A B C D E F G H J K L M N O P Q R S T
		19  . O . . . . . . . . . . . . . . . . .
		18  X . O . . X . . . . . . . . . . . . .
		17  . . . O . X . . . . . . . . . . . . .
		16  X . X . X X . . . . . . . . . . . . .
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
		 4  . . . . . . . . . . . . . . . . . . .
		 3  . . . . . . . . . . . . . . . . . . .
		 2  . . . . . . . . . . . . . . . . . . .
		 1  . . . . . . . . . . . . . . . . . . .
		    A B C D E F G H J K L M N O P Q R S T
	*/
	b.PutLabel("BA16")
	b.PutLabel("WB19")
	b.PutLabel("BC16")
	b.PutLabel("WC18")
	b.PutLabel("BE16")
	b.PutLabel("WD17")
	b.PutLabel("BF17")
	b.PutLabel("BA18")
	b.PutLabel("BF18")
	b.PutLabel("BF16")
	for i := 0; i < 500; i++ {
		b2 := b.Copy()
		b2.SelfBattle(WHITE, nil)
		s := b2.Score()
		if s > 0 {
			win += 1
			for j := 0; j < len(b2.Actions)-1; j++ {
				k1, c1 := ParseIndexAction(b2.Actions[j])
				k2, c2 := ParseIndexAction(b2.Actions[j+1])
				if c1 == WHITE && c2 == BLACK {
					lgr.Set(BLACK, k1, k2)
				}
			}
		} else {
			for j := 0; j < len(b2.Actions)-1; j++ {
				k1, c1 := ParseIndexAction(b2.Actions[j])
				k2, c2 := ParseIndexAction(b2.Actions[j+1])
				if c1 == BLACK && c2 == WHITE {
					lgr.Set(WHITE, k1, k2)
				}
			}
		}
		if i == 0 {
			t.Log(s)
			t.Log(b2.String(nil))
		}
	}
	t.Log(win * 100 / 500)
}

func TestSelfBattle2(t *testing.T) {
	win := 0
	lgr := NewLastGoodReply()
	b := NewBoard()
	/*
			A B C D E F G H J K L M N O P Q R S T
		19  . . . . . . . . . . . . . . . . . . .
		18  . O O O O O O O O O O . . . . . . . .
		17  . X X X X X X X X X X . . . . . . . .
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
		 4  . . . . . . . . . . . . . . . . . . .
		 3  . . . . . . . . . . . . . . . . . . .
		 2  . . . . . . . . . . . . . . . . . . .
		 1  . . . . . . . . . . . . . . . . . . .
		    A B C D E F G H J K L M N O P Q R S T
	*/
	b.PutLabel("BB18")
	b.PutLabel("WB17")
	b.PutLabel("BC18")
	b.PutLabel("WC17")
	b.PutLabel("BD18")
	b.PutLabel("WD17")
	b.PutLabel("BE18")
	b.PutLabel("WE17")
	b.PutLabel("BF18")
	b.PutLabel("WF17")
	b.PutLabel("BG18")
	b.PutLabel("WG17")
	b.PutLabel("BH18")
	b.PutLabel("WH17")
	b.PutLabel("BJ18")
	b.PutLabel("WJ17")
	b.PutLabel("BK18")
	b.PutLabel("WK17")
	b.PutLabel("BL18")
	b.PutLabel("WL17")
	for i := 0; i < 500; i++ {
		b2 := b.Copy()
		b2.SelfBattle(BLACK, nil)
		s := b2.Score()
		if s > 0 {
			win += 1
			for j := 0; j < len(b2.Actions)-1; j++ {
				k1, c1 := ParseIndexAction(b2.Actions[j])
				k2, c2 := ParseIndexAction(b2.Actions[j+1])
				if c1 == WHITE && c2 == BLACK {
					lgr.Set(BLACK, k1, k2)
				}
			}
		} else {
			for j := 0; j < len(b2.Actions)-1; j++ {
				k1, c1 := ParseIndexAction(b2.Actions[j])
				k2, c2 := ParseIndexAction(b2.Actions[j+1])
				if c1 == BLACK && c2 == WHITE {
					lgr.Set(WHITE, k1, k2)
				}
			}
		}
		if i == 0 {
			t.Log(s)
			t.Log(b2.String(nil))
		}
	}
	t.Log(win * 100 / 500)
}

func TestInfluence(t *testing.T) {
	/*
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
		 7  . . . . . . . . . . . . . . . X . . .
		 6  . . . . . . . . . . . . . . . . . . .
		 5  . . . . . . . . . . . . . . . . . X .
		 4  . . . . . . . . . . . . . . O . X . .
		 3  . . . . . . . . . . . . . . . O O X .
		 2  . . . . . . . . . . . . . . . . . . .
		 1  . . . . . . . . . . . . . . . . . . .
		    A B C D E F G H J K L M N O P Q R S T
	*/
	b := NewBoard()
	b.PutLabel("BS3")
	b.PutLabel("BS5")
	b.PutLabel("BR4")
	b.PutLabel("BQ7")
	b.PutLabel("WR3")
	b.PutLabel("WQ3")
	b.PutLabel("WP4")
	t.Log(b.String(nil))

	influence := b.Influence()
	mark := make(map[int]string)
	for k, v := range influence {
		tmp := fmt.Sprintf("%2.f", float64(v))
		if len(tmp) < 3 {
			tmp = " " + tmp
		}
		mark[k] = tmp
	}
	for k, c := range b.Points {
		if _, ok := influence[k]; !ok {
			mark[k] = "  " + ColorMark(c)
		}
	}
	t.Log(b.String(mark))
}

func TestTerritory(t *testing.T) {
	b := NewBoard()
	b.PutLabel("BJ10")
	b.PutLabel("BL10")
	b.PutLabel("BK9")
	t.Log(b.String(nil))

	territory := b.Territory()
	mark := make(map[int]string)
	for k, v := range territory {
		mark[k] = fmt.Sprintf("%2.f", float64(v))
	}
	for k, c := range b.Points {
		if _, ok := territory[k]; !ok {
			mark[k] = " " + ColorMark(c)
		}
	}
	t.Log(b.String(mark))
}

func TestLocalFeature(t *testing.T) {
	b := NewBoard()
	/*
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
		 5  . . . X . . . . . . . . . . . . . . .
		 4  . . X O X . . . . . . . . . . . . . .
		 3  . . . . . . . . . . . . . . . . . . .
		 2  . . . . . . . . . . . . . . . . . . .
		 1  . . . . . . . . . . . . . . . . . . .
		    A B C D E F G H J K L M N O P Q R S T
	*/
	b.PutLabel("BC4")
	b.PutLabel("BE4")
	b.PutLabel("BD5")
	b.PutLabel("WD4")
	k := PosIndex(3, 2)
	f := b.LocalFeature(k, WHITE)
	t.Log(f)
}

func TestPattern3x3(t *testing.T) {
	b := NewBoard()
	/*
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
		 6  . . . X O . . . . . . . . . . . . . .
		 5  . . X O X O . . . . . . . . . . . . .
		 4  . X O . . . . . . . . . . . . . . . .
		 3  . X . . . . . . . . . . . . . . . . .
		 2  . . . . . . . . . . . . . . . . . . .
		 1  . . . . . . . . . . . . . . . . . . .
		    A B C D E F G H J K L M N O P Q R S T
	*/
	b.PutLabel("BB3")
	b.PutLabel("BB4")
	b.PutLabel("BC5")
	b.PutLabel("WC4")
	b.PutLabel("WD5")
	b.PutLabel("BD6")
	b.PutLabel("BE5")
	b.PutLabel("WE6")
	b.PutLabel("WF5")
	t.Log(b.String(nil))
	t.Log(b.PatternDxd(PosIndex(3, 3), BLACK, 1))
	t.Log(b.PatternDxd(PosIndex(4, 3), BLACK, 1))
}

func TestSaveAtari(t *testing.T) {
	b := NewBoard()
	/*
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
		 6  . . . X O . . . . . . . . . . . . . .
		 5  . . X O X O . . . . . . . . . . . . .
		 4  . X O . . . . . . . . . . . . . . . .
		 3  . X . . . . . . . . . . . . . . . . .
		 2  . . . . . . . . . . . . . . . . . . .
		 1  . . . . . . . . . . . . . . . . . . .
		    A B C D E F G H J K L M N O P Q R S T
	*/
	b.PutLabel("BB3")
	b.PutLabel("BB4")
	b.PutLabel("BC5")
	b.PutLabel("WC4")
	b.PutLabel("WD5")
	b.PutLabel("BD6")
	b.PutLabel("BE5")
	b.PutLabel("WE6")
	b.PutLabel("WF5")
	worm := b.WormFromPoint(PosIndex(4, 4), BLACK, 3)
	t.Log(worm.Liberty)
	ret := b.SaveAtari(worm)
	for _, k := range ret {
		if k != PosIndex(3, 3) && k != PosIndex(4, 3) {
			t.Error(k)
		}
	}
}

func TestEmptyWormFromPoint(t *testing.T) {
	b := NewBoard()
	/*
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
		 5  . . . O . . . . . . . . . . . . . . .
		 4  . . O X X . . . . . . . . . . . . . .
		 3  . . X . . X . . . . . . . . . . . . .
		 2  . . . . . X . . . . . . . . . . . . .
		 1  X . . . . . . . . . . . . . . . . . .
		    A B C D E F G H J K L M N O P Q R S T
	*/

	b.PutLabel("BA1")
	b.PutLabel("BC3")
	b.PutLabel("BD4")
	b.PutLabel("BE4")
	b.PutLabel("BF2")
	b.PutLabel("BF3")
	b.PutLabel("WC4")
	b.PutLabel("WD5")

	b.EmptyWormFromPoint(PosIndex(2, 0), 1)
}

func BenchmarkPointMap(t *testing.B) {
	a := 361
	for i := 0; i < t.N; i++ {
		_ = make([]int, a)
	}
}

func TestSinglePatExpand(t *testing.T) {
	b := NewBoard()
	/*
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
		 5  . . O X . . . O X . . . . . . . . . .
		 4  . . X . . . . X . X . . . . . . . . .
		 3  . . . . . . . . . . . . . . . . . . .
		 2  . . . . . . . . . . . . . . . . . . .
		 1  . . . . . . . . . . . . . . . . . . .
		    A B C D E F G H J K L M N O P Q R S T
	*/
	//t.Log(pat3x3Dict)
	b.PutLabel("BC4")
	b.PutLabel("WC5")
	b.PutLabel("BD5")
	b.PutLabel("BH4")
	b.PutLabel("WH5")
	b.PutLabel("BJ5")
	b.PutLabel("BK4")

	buf := b.Pattern3x3String(PosIndex(3, 3), INVALID_COLOR)
	t.Log(buf)

	if _, ok := pat3x3Dict[buf]; !ok {
		t.Error()
	}

	buf = b.Pattern3x3String(PosIndex(8, 3), INVALID_COLOR)
	t.Log(buf)
	if _, ok := pat3x3Dict[buf]; ok {
		t.Error()
	}
}
