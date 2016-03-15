package gogo

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/xlvector/hector/lr"
)

type Color byte

const (
	SIZE                = 19
	NPOINT              = 361
	GRAY          Color = 0
	BLACK         Color = 1
	WHITE         Color = 2
	INVALID_COLOR Color = 3
	PATTERN_SIZE        = 12
	LX                  = "ABCDEFGHJKLMNOPQRSTUVWXYZ"
)

//point, distance, other points
var PointDisMap [][][]int

func init() {
	PointDisMap = make([][][]int, NPOINT)
	for i := 0; i < NPOINT; i++ {
		PointDisMap[i] = make([][]int, PATTERN_SIZE)
		for j := 0; j < PATTERN_SIZE; j++ {
			PointDisMap[i][j] = make([]int, 0, j*4+5)
		}
	}
	for i := 0; i < NPOINT; i++ {
		for j := 0; j < NPOINT; j++ {
			d := Distance(i, j)
			if d < PATTERN_SIZE {
				PointDisMap[i][d] = append(PointDisMap[i][d], j)
			}
		}
	}
}

func PointString(x, y int, c Color) string {
	ret := ColorMark(c)
	ret += "["
	ret += string(LX[x])
	ret += strconv.Itoa(y + 1)
	ret += "]"
	return ret
}

func ParseAction(k int) (int, int, Color) {
	x, y := IndexPos(k / 10)
	c := Color(k % 10)
	return x, y, c
}

func ParseIndexAction(k int) (int, Color) {
	return k / 10, Color(k % 10)
}

func FormatAction(x, y int, c Color) int {
	return PosIndex(x, y)*10 + int(c)
}

func FormatIndexAction(k int, c Color) int {
	return k*10 + int(c)
}

func OpColor(c Color) Color {
	if c == BLACK {
		return WHITE
	} else if c == WHITE {
		return BLACK
	}
	return INVALID_COLOR
}

func PosIndex(x, y int) int {
	return y*SIZE + x
}

func IndexPos(k int) (int, int) {
	return k % SIZE, k / SIZE
}

func Neigh4(k int) []int {
	x, y := IndexPos(k)
	ret := make([]int, 0, 4)
	if !PosOutBoard(x-1, y) {
		ret = append(ret, PosIndex(x-1, y))
	}
	if !PosOutBoard(x+1, y) {
		ret = append(ret, PosIndex(x+1, y))
	}
	if !PosOutBoard(x, y-1) {
		ret = append(ret, PosIndex(x, y-1))
	}
	if !PosOutBoard(x, y+1) {
		ret = append(ret, PosIndex(x, y+1))
	}
	return ret
}

func NeighD(k int) []int {
	x, y := IndexPos(k)
	ret := make([]int, 0, 4)
	if !PosOutBoard(x-1, y-1) {
		ret = append(ret, PosIndex(x-1, y-1))
	}
	if !PosOutBoard(x+1, y-1) {
		ret = append(ret, PosIndex(x+1, y-1))
	}
	if !PosOutBoard(x-1, y+1) {
		ret = append(ret, PosIndex(x-1, y+1))
	}
	if !PosOutBoard(x+1, y+1) {
		ret = append(ret, PosIndex(x+1, y+1))
	}
	return ret
}

func PosOutBoard(x, y int) bool {
	return x < 0 || y < 0 || x >= SIZE || y >= SIZE
}

func EdgeDis(x, y int) int {
	ret := x
	if ret > SIZE-x-1 {
		ret = SIZE - x - 1
	}
	if ret > y {
		ret = y
	}
	if ret > SIZE-y-1 {
		ret = SIZE - y - 1
	}
	return ret
}

func Abs(n int) int {
	if n >= 0 {
		return n
	}
	return -1 * n
}

func Distance(a, b int) int {
	xa, ya := IndexPos(a)
	xb, yb := IndexPos(b)
	return Abs(xa-xb) + Abs(ya-yb)
}

func IndexEdgeDis(k int) int {
	x, y := IndexPos(k)
	return EdgeDis(x, y)
}

func IndexOutBoard(k int) bool {
	x, y := IndexPos(k)
	return PosOutBoard(x, y)
}

type Board struct {
	Points      []Color
	KoIndex     int
	Model       *lr.LogisticRegression
	PointHash   []int64
	PatternHash [][]int64
	Actions     []int
	LastPattern []int64
}

func NewBoard() *Board {
	ret := &Board{
		Points:  make([]Color, NPOINT),
		KoIndex: -1,
		Actions: make([]int, 0, 10),
	}
	ret.InitHash()
	return ret
}

func (b *Board) Clear() {
	for i := 0; i < NPOINT; i++ {
		b.Points[i] = GRAY
	}
	b.KoIndex = -1
	b.Model = nil
	b.InitHash()
}

func (b *Board) Copy() *Board {
	ret := &Board{
		Points:      make([]Color, NPOINT),
		KoIndex:     b.KoIndex,
		PointHash:   make([]int64, len(b.PointHash)),
		PatternHash: make([][]int64, len(b.PatternHash)),
		Actions:     make([]int, len(b.Actions)),
		LastPattern: make([]int64, len(b.LastPattern)),
		Model:       b.Model,
	}
	for i, v := range b.Points {
		ret.Points[i] = v
	}
	for i, v := range b.PointHash {
		ret.PointHash[i] = v
	}
	for i, v := range b.Actions {
		ret.Actions[i] = v
	}
	for i, v := range b.LastPattern {
		ret.LastPattern[i] = v
	}
	for i, a := range b.PatternHash {
		tmp := make([]int64, len(a))
		for j, v := range a {
			tmp[j] = v
		}
		ret.PatternHash[i] = tmp
	}
	return ret
}

func (b *Board) StableEye(k int, c Color) bool {
	n4 := Neigh4(k)
	for _, p := range n4 {
		if b.Points[p] != c {
			return false
		}
	}
	nd := NeighD(k)
	oc := OpColor(c)
	n := 0
	for _, p := range nd {
		if b.Points[p] == oc {
			n += 1
		}
	}
	if len(n4) == 4 && n < 2 {
		return true
	}
	if len(n4) < 4 && n < 1 {
		return true
	}
	return false
}

func (b *Board) CanPut(k int, c Color) (bool, map[int]Color) {
	if k < 0 || k >= NPOINT {
		return false, nil
	}
	if b.Points[k] != GRAY {
		return false, nil
	}

	if k == b.KoIndex {
		return false, nil
	}

	if b.StableEye(k, c) {
		return false, nil
	}

	take := make(map[int]Color)
	oc := OpColor(c)
	nworms := b.NeighWorms(k, c, oc, 2)
	for _, nw := range nworms {
		if nw.Liberty == 1 {
			for p, c1 := range nw.Points {
				take[p] = c1
			}
		}
	}

	if len(take) > 0 {
		return true, take
	}

	worm := b.WormFromPoint(k, c, 1)

	if worm.BorderColor == oc {
		return false, nil
	}

	return true, nil
}

func (b *Board) PutLabel(buf string) bool {
	c := ParseColor(buf[0:1])
	x := strings.Index(LX, buf[1:2])
	y, _ := strconv.Atoi(buf[2:])
	y -= 1
	return b.Put(PosIndex(x, y), c)
}

func (b *Board) Put(k int, c Color) bool {
	ok, take := b.CanPut(k, c)
	if !ok {
		return false
	}
	b.LastPattern = b.FinalPatternHash(k, c)
	b.KoIndex = -1
	if len(take) > 0 {
		for p, _ := range take {
			b.UpdateHash(p, b.Points[p], GRAY)
			b.Points[p] = GRAY
		}
		if len(take) == 1 {
			for p, _ := range take {
				b.KoIndex = p
			}
		}
	}
	b.Points[k] = c
	b.Actions = append(b.Actions, FormatIndexAction(k, c))
	b.UpdateHash(k, GRAY, c)
	return true
}

func (b *Board) LastMove() (int, Color) {
	if len(b.Actions) == 0 {
		return -1, INVALID_COLOR
	}
	a := b.Actions[len(b.Actions)-1]
	return ParseIndexAction(a)
}

type Worm struct {
	Points      map[int]Color
	Liberty     int
	Color       Color
	BorderColor Color
}

func NewWorm() *Worm {
	return &Worm{
		Points:      make(map[int]Color),
		Liberty:     0,
		Color:       INVALID_COLOR,
		BorderColor: INVALID_COLOR,
	}
}

func (w *Worm) AddPoint(p int, c Color) {
	w.Points[p] = c
}

func (w *Worm) IncludePoint(p int) bool {
	_, ok := w.Points[p]
	return ok
}

func (b *Board) WormFromPoint(k int, c Color, stopLiberty int) *Worm {
	// if pass invalid color, means use color in point k of board, otherwise, use specified color c
	if c == INVALID_COLOR {
		c = b.Points[k]
	}
	worm := NewWorm()
	worm.Color = c
	queue := make([]int, 0, 10)
	start := 0
	queue = append(queue, k)
	lb := make(map[int]byte)
	for {
		if start >= len(queue) {
			break
		}
		if stopLiberty > 0 && len(lb) > stopLiberty {
			break
		}
		v := queue[start]
		start += 1
		if worm.IncludePoint(v) {
			continue
		}
		worm.AddPoint(v, c)
		n4 := Neigh4(v)
		for _, nv := range n4 {
			if worm.IncludePoint(nv) {
				continue
			}
			if b.Points[nv] == c {
				queue = append(queue, nv)
			} else {
				if b.Points[nv] == GRAY {
					lb[nv] = 1
				}
				if worm.BorderColor == INVALID_COLOR {
					worm.BorderColor = b.Points[nv]
				} else if worm.BorderColor == GRAY {
					worm.BorderColor = GRAY
				} else if worm.BorderColor != b.Points[nv] {
					worm.BorderColor = GRAY
				}
			}
		}
	}
	worm.Liberty = len(lb)
	return worm
}

func ColorMark(c Color) string {
	if c == BLACK {
		return "X"
	} else if c == WHITE {
		return "O"
	}
	return "."
}

func ParseColor(c string) Color {
	if c == "B" {
		return BLACK
	} else if c == "W" {
		return WHITE
	}
	return GRAY
}

func FormatColor(c Color) string {
	if c == BLACK {
		return "B"
	} else if c == WHITE {
		return "W"
	}
	return ""
}

func (b *Board) String(mark map[int]string) string {
	if mark == nil {
		mark = make(map[int]string)
	}
	ret := "\n"
	ret += "   "
	for _, ch := range LX[0:SIZE] {
		ret += " "
		ret += string(ch)
	}
	ret += "\n"
	for y := 0; y < SIZE; y++ {
		for x := 0; x < SIZE; x++ {
			if x == 0 {
				ret += fmt.Sprintf("%2.f ", float64(SIZE-y))
			}
			ret += " "
			k := PosIndex(x, SIZE-y-1)
			if mk, ok := mark[k]; ok {
				ret += mk
			} else {
				ret += ColorMark(b.Points[PosIndex(x, SIZE-y-1)])
			}
		}
		ret += "\n"
	}
	ret += "   "
	for _, ch := range LX[0:SIZE] {
		ret += " "
		ret += string(ch)
	}
	ret += "\n"
	return ret
}
