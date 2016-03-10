package gogo

import (
	"errors"
	"fmt"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/xlvector/hector/dt"
)

const (
	PATTERN_SIZE = 14
)

type Board struct {
	size        int
	w           []Point
	komi        float64
	ko          Point
	takeBlack   int
	takeWhite   int
	models      []*dt.RandomForest
	step        int
	info        *BoardInfo
	pdm         *PointDistanceMap
	pointHash   []int64
	patternHash [][]int64
}

func NewBoard(size int) *Board {
	b := &Board{}
	b.Clear(size)
	return b
}

func (b *Board) GetPatternHash(p int) []int64 {
	if p < 0 || p >= len(b.patternHash) {
		fmt.Println("GetPatternHash", p, len(b.patternHash))
		return nil
	}
	return b.patternHash[p]
}

func (b *Board) SetPointDistanceMap(pdm *PointDistanceMap) {
	b.pdm = pdm
	b.InitPatternHash(pdm)
}

func (b *Board) Clear(size int) {
	b.size = size
	b.w = make([]Point, size*size)
	b.pointHash = make([]int64, size*size)
	b.patternHash = make([][]int64, size*size)
	b.ko = InvalidPoint()
	i := 0
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			b.w[i] = Point{x, y, GRAY}
			b.patternHash[i] = make([]int64, PATTERN_SIZE)
			i += 1
		}
	}

	for i, p := range b.w {
		b.pointHash[i] = b.PointHash(p)
	}
}

func (b *Board) W() []Point {
	return b.w
}

func (b *Board) Model() *dt.RandomForest {
	step := b.step
	if b.models == nil || len(b.models) == 0 {
		return nil
	}
	if len(b.models) == 1 {
		return b.models[0]
	} else if len(b.models) == 2 {
		if step < 50 {
			return b.models[0]
		} else {
			return b.models[1]
		}
	} else if len(b.models) == 3 {
		if step < 20 {
			return b.models[0]
		} else if step >= 20 && step < 100 {
			return b.models[1]
		} else {
			return b.models[2]
		}
	} else {
		if step < 20 {
			return b.models[0]
		} else if step >= 20 && step < 60 {
			return b.models[1]
		} else if step >= 60 && step < 140 {
			return b.models[2]
		} else {
			return b.models[3]
		}
	}

}

func (b *Board) Copy() *Board {
	ret := &Board{
		size:      b.size,
		w:         make([]Point, 0, 1+len(b.w)),
		komi:      b.komi,
		takeBlack: b.takeBlack,
		takeWhite: b.takeWhite,
		ko:        b.ko,
		models:    b.models,
		step:      b.step,
		pdm:       b.pdm,
	}
	for _, p := range b.w {
		ret.w = append(ret.w, p)
	}
	return ret
}

func (b *Board) XMirror(p Point) Point {
	return b.Get(b.size-p.x, p.y)
}

func (b *Board) YMirror(p Point) Point {
	return b.Get(p.x, b.size-p.y)
}

func (b *Board) DMirror(p Point) Point {
	return b.Get(b.size-p.x, b.size-p.y)
}

func (b *Board) index(x, y int) int {
	return y*b.size + x
}

func (b *Board) Index(p Point) int {
	return b.index(p.x, p.y)
}

func (b *Board) pos(p int) (int, int) {
	return p % b.size, p / b.size
}

func (b *Board) PutLabel(buf string) error {
	c := buf[0:1]
	x := strings.Index(LX, strings.ToUpper(buf[1:2]))
	y, _ := strconv.Atoi(buf[2:])
	y -= 1
	if c == "B" {
		return b.Put(x, y, BLACK)
	} else if c == "W" {
		return b.Put(x, y, WHITE)
	} else {
		return errors.New("invalid input: " + buf)
	}
}

func (b *Board) Put(x, y int, stone Color) error {
	if !b.valid(stone, x, y) {
		return errors.New("invalid position")
	}
	if b.ko.x == x && b.ko.y == y {
		return errors.New("ko position")
	}
	b.ko = InvalidPoint()
	i := b.index(x, y)
	prevWi := b.w[i]
	b.w[i] = Point{x, y, stone}
	b.UpdateHash(x, y, stone, b.pdm)
	tworms := b.GetTakeWorms(stone, x, y)
	if len(tworms) == 0 {
		pworm := b.WormContainsPoint(i)
		if pworm.Dead() {
			b.w[i] = prevWi
			return errors.New("sucide position")
		}
	} else {
		b.ko = b.koPositionOfDeadWorms(b.w[i], tworms)
		b.TakeWorms(tworms)
	}
	b.step += 1
	return nil
}

func (b *Board) Size() int {
	return b.size
}

func (b *Board) Get(x, y int) Point {
	if x < 0 || y < 0 || x >= b.size || y >= b.size {
		return InvalidPoint()
	}
	return b.w[b.index(x, y)]
}

func (b *Board) SingleEye(x, y int, stone Color) bool {
	if b.Get(x, y).color != GRAY {
		return false
	}
	n4 := b.Neighbor4(x, y)
	for _, p := range n4 {
		if p.color != stone {
			return false
		}
	}
	return true
}

func (b *Board) Neighbor4(x, y int) []Point {
	ret := make([]Point, 0, 4)
	if p := b.Get(x-1, y); p.Valid() {
		ret = append(ret, p)
	}
	if p := b.Get(x+1, y); p.Valid() {
		ret = append(ret, p)
	}
	if p := b.Get(x, y-1); p.Valid() {
		ret = append(ret, p)
	}
	if p := b.Get(x, y+1); p.Valid() {
		ret = append(ret, p)
	}
	return ret
}

func (b *Board) Neighbor4Color(x, y int, color Color) []Point {
	ret := make([]Point, 0, 4)
	if p := b.Get(x-1, y); p.Valid() && p.color == color {
		ret = append(ret, p)
	}
	if p := b.Get(x+1, y); p.Valid() && p.color == color {
		ret = append(ret, p)
	}
	if p := b.Get(x, y-1); p.Valid() && p.color == color {
		ret = append(ret, p)
	}
	if p := b.Get(x, y+1); p.Valid() && p.color == color {
		ret = append(ret, p)
	}
	return ret
}

func (b *Board) NeighDiamond(x, y int) []Point {
	ret := make([]Point, 0, 4)
	if p := b.Get(x-1, y-1); p.Valid() {
		ret = append(ret, p)
	}
	if p := b.Get(x+1, y+1); p.Valid() {
		ret = append(ret, p)
	}
	if p := b.Get(x+1, y-1); p.Valid() {
		ret = append(ret, p)
	}
	if p := b.Get(x-1, y+1); p.Valid() {
		ret = append(ret, p)
	}
	return ret
}

func (b *Board) Valid(p Point) bool {
	return b.valid(p.color, p.x, p.y)
}

func (b *Board) valid(stone Color, x, y int) bool {
	if x < 0 || y < 0 || x >= b.size || y >= b.size {
		return false
	}

	p := b.Get(x, y)
	if p.color != GRAY {
		return false
	}
	return true
}

func (b *Board) Stone(x, y int, ls Point) string {
	i := b.index(x, y)
	if b.w[i].color == GRAY {
		return "."
	} else if b.w[i].color == BLACK {
		if ls.x == x && ls.y == y {
			return "#"
		} else {
			return "X"
		}
	} else if b.w[i].color == WHITE {
		if ls.x == x && ls.y == y {
			return "@"
		} else {
			return "O"
		}
	} else {
		return "?"
	}
}

func (b *Board) String(lastStep Point) string {
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
		}
	}()
	ret := ""
	if b == nil || b.size == 0 {
		return ret
	}
	for y := b.size; y >= 0; y-- {
		ret += fmt.Sprintf("    ")
		for x := 0; x <= b.size; x++ {
			if y == 0 {
				if x == 0 {
					ret += fmt.Sprintf("   ")
				} else {
					ret += fmt.Sprintf(string(LX[x-1]))
					ret += fmt.Sprintf(" ")
				}
			} else {
				if x == 0 {
					if y < 10 {
						ret += fmt.Sprintf(" %d", y)
					} else {
						ret += fmt.Sprintf("%d", y)
					}
					ret += fmt.Sprintf(" ")
				} else {
					ret += fmt.Sprintf(b.Stone(x-1, y-1, lastStep))
					ret += fmt.Sprintf(" ")
				}
			}
		}
		ret += fmt.Sprintf("\n")
	}
	ret += fmt.Sprintf("\n")
	return ret
}
