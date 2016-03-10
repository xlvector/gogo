package gogo

import (
	"fmt"
	"testing"
)

/*
   19 . . . . . . . . . . . . . . . . . . .
   18 . . . . . . . . . . . . . . . . . . .
   17 . . . . . . . . . . . . . . . . . . .
   16 . . . X . . . . . . . . . . . O . . .
   15 . . . . . . . . . . . . . . . . . . .
   14 . . O . . . . . . . . . . . . X . . .
   13 . . . . . . . . . . . . . . . . . . .
   12 . . . . . . . . . . . . . . . . . . .
   11 . . . X . . . . . . . . . . . . . . .
   10 . . . . . . . . . . . . . . . . . . .
    9 . . . . . . . . . . . . . . . . . . .
    8 . . . . . . . . . . . . . . . . . . .
    7 . . . . . . . . . . . . . . . . . . .
    6 . . X . . O . . . . . . . . . O . . .
    5 . . . . . . . . . . . . O . . . . . .
    4 . . . X . O X . . . . . . . O . . . .
    3 . . . . . O . . . . . . X . . O . . .
    2 . . . . . . . . . . . . . . X . . . .
    1 . . . . . . . . . . . . . . . . . . .
      A B C D E F G H J K L M N O P Q R S T

*/

func TestStoneGrayDis(t *testing.T) {
	g := &Game{}
	g.Init(19)
	g.B.PutLabel("BC6")
	g.B.PutLabel("BD4")
	g.B.PutLabel("WF4")
	g.B.PutLabel("WF3")
	g.B.PutLabel("BG4")
	g.B.PutLabel("BN3")
	g.B.PutLabel("BP2")
	g.B.PutLabel("WP4")
	g.B.PutLabel("WQ3")
	g.B.PutLabel("WQ6")
	g.B.PutLabel("WN5")
	g.B.PutLabel("BD16")
	g.B.PutLabel("WQ16")
	g.B.PutLabel("BQ14")
	g.B.PutLabel("WC14")
	g.B.PutLabel("BD11")
	g.B.PutLabel("WF6")
	g.Print()

	black := g.B.StoneGrayDistance(BLACK)
	white := g.B.StoneGrayDistance(WHITE)

	fmt.Println(len(black))
	lines := []string{}
	maxDis := 10
	line := ""
	for i := 0; i < len(black); i++ {
		db := black[i]
		dw := white[i]
		label := "."
		if db <= maxDis || dw <= maxDis {
			if db > dw {
				label = "o"
			} else if db < dw {
				label = "x"
			}
			if db == 0 {
				label = "X"
			}
			if dw == 0 {
				label = "O"
			}
		}
		line += label
		line += " "
		x, _ := g.B.pos(i)
		if x == g.B.size-1 {
			lines = append(lines, line)
			line = ""
		}
	}
	for i := len(lines) - 1; i >= 0; i-- {
		fmt.Println(lines[i])
	}
}

/*
   19 . . . . . . . . . . . . . . . . . . .
   18 . . . . . . . . . . . . . . . . . . .
   17 . . . . . . . . . . . . . . . . . . .
   16 . . . . . . . . . . . . . . . . . . .
   15 . . . . . . . . . . . . . . . . . . .
   14 . . . . . . . . . . . . . . . . . . .
   13 . . . . . . . . . . . . . . . . . . .
   12 . . . . . . . . . . . . . . . . . . .
   11 . . . . . . . . . . . . . . . . . . .
   10 . . . . . . . . X . X . . . . . . . .
    9 . . . . . . . . . . . . . . . . . . .
    8 . . . . . . . . . . . . . . . . . . .
    7 . . . . . . . . . . . . . . . . . . .
    6 . . . . . . . . . . . . . . . . . . .
    5 . . . . . . . . . . . . . . . . . . .
    4 . . . . . . . . . . . . . . . . . . .
    3 . . . . . . . . . . . . . . . . . . .
    2 . . . . . . . . . . . . . . . . . . .
    1 . . . . . . . . . . . . . . . . . . .
      A B C D E F G H J K L M N O P Q R S T

*/

func TestDilation(t *testing.T) {
	g := &Game{}
	g.Init(19)
	g.B.PutLabel("BJ10")
	g.B.PutLabel("BM10")
	g.Print()

	seeds := make(map[int]int)
	seeds[g.B.index(8, 9)] = 64
	seeds[g.B.index(11, 9)] = 64
	g.B.PrintWithValue(seeds, false)

	for i := 0; i < 3; i++ {
		seeds = g.B.Dilation(seeds)
		fmt.Println()
		g.B.PrintWithValue(seeds, false)
	}

	for i := 0; i < 7; i++ {
		seeds = g.B.Erase(seeds)
		fmt.Println()
		g.B.PrintWithValue(seeds, false)
	}
	g.B.PrintWithValue(seeds, true)
}

/*
   19 . . . . . . . . . . . . . . . . . . .
   18 . . O O X . . . . . . . . . . . . . .
   17 . . O X . X . . . X . O . O . . O . .
   16 . . O X . . . . . . . . . . . O . X .
   15 . O X . . . . . . X . . . . . . . . .
   14 . O X . . . . . . . . . . . . . X . .
   13 . . X . . . . . . . . . . . . . . . .
   12 . . . . . . . . . . . . . . . X . . .
   11 . . . X . . . . . . . . . . . . . . .
   10 . . . . . . . . . . . . . . . . X . .
    9 . . . . . . . . . . . . . . . . . . .
    8 . . . . . . . . . . . . . . . . . . .
    7 . . . . . . . . . . . . . . . . . . .
    6 . . . . . . . . . . . . . . . . . . .
    5 . . . . . . . . . . . . . . . . . . .
    4 . . . . . . . . . . . . . . . . . . .
    3 . . . . . . . . . . . . . . . . . . .
    2 . . . . . . . . . . . . . . . . . . .
    1 . . . . . . . . . . . . . . . . . . .
      A B C D E F G H J K L M N O P Q R S T

*/

func TestTerritory(t *testing.T) {
	g := &Game{}
	g.Init(19)
	g.B.PutLabel("WC18")
	g.B.PutLabel("WD18")
	g.B.PutLabel("BE18")
	g.B.PutLabel("WC17")
	g.B.PutLabel("BD17")
	g.B.PutLabel("BF17")
	g.B.PutLabel("BK17")
	g.B.PutLabel("WM17")
	g.B.PutLabel("WO17")
	g.B.PutLabel("WR17")
	g.B.PutLabel("WC16")
	g.B.PutLabel("BD16")
	g.B.PutLabel("WQ16")
	g.B.PutLabel("BS16")
	g.B.PutLabel("WB15")
	g.B.PutLabel("BC15")
	g.B.PutLabel("BK15")
	g.B.PutLabel("WB14")
	g.B.PutLabel("BC14")
	g.B.PutLabel("BR14")
	g.B.PutLabel("BC13")
	g.B.PutLabel("BQ12")
	g.B.PutLabel("BD11")
	g.B.PutLabel("BR10")
	g.Print()

	territory := g.B.Territory()
	g.B.PrintWithValue(territory, true)

	influence := g.B.Influence()
	g.B.PrintWithValue(influence, true)
}
