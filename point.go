package gogo

import (
	"fmt"

	"github.com/xlvector/hector/core"
)

type Point struct {
	x, y  int
	color Color
}

func MakePoint(x, y int, c Color) Point {
	return Point{x, y, c}
}

func (p Point) X() int {
	return p.x
}

func (p Point) Y() int {
	return p.y
}

func (p Point) Color() Color {
	return p.color
}

func (p Point) String() string {
	return fmt.Sprintf("%s[%s, %d]", SGFColor(p.color), string(LX[p.x]), p.y+1)
}

func (p Point) Valid() bool {
	return p.x >= 0 && p.y >= 0
}

func (p Point) EdgeDis(size int) int {
	ret := p.x
	if ret > p.y {
		ret = p.y
	}
	if ret > size-p.x-1 {
		ret = size - p.x - 1
	}
	if ret > size-p.y-1 {
		ret = size - p.y - 1
	}
	return ret
}

func (p Point) EdgeDisX(size int) int {
	ret := p.x
	if ret > size-p.x-1 {
		ret = size - p.x - 1
	}
	return ret
}

func (p Point) EdgeDisY(size int) int {
	ret := p.y
	if ret > size-p.y-1 {
		ret = size - p.y - 1
	}
	return ret
}

func InvalidPoint() Point {
	return Point{-1, -1, UNKNOWN}
}

func (p Point) Distance(s Point) int {
	dx := s.x - p.x
	if dx < 0 {
		dx *= -1
	}
	dy := s.y - p.y
	if dy < 0 {
		dy *= -1
	}
	return dx + dy
}

const (
	F_CAPTURE_CAN_ESCAPE = iota
	F_CAPTURE_CANNOT_ESCAPE
	F_OP_CAPTURE
	F_OP_ATARI
	F_ATARI
	F_LIBERTY
	F_OP_LIBERTY
	F_DISTANCE_EDGE
	F_CUT_COUNT
	F_CUT_WORM_SIZE
	F_CONNECT_COUNT
	F_CONNECT_WORM_SIZE
	F_CONNECT_WORM_MAX_LIBERTY
	F_SURROUND_MY
	F_SURROUND_OP
	F_EDGE_DIS_X
	F_EDGE_DIS_Y
	F_ORIGIN_WORM_SIZE
	F_ORIGIN_WORM_LIBERTY
	F_PAT_3X3
	F_PAT_5X5
)

type PointFeature struct {
	Liberty      int
	CutPoint     int
	P            Point
	OriginWorm   *Worm
	OriginDragon *Dragon
	BoardWorms   map[int]*Worm
	Score        float64
	SelectProb   float64
	WhiteDis     int
	BlackDis     int
	Fc           *core.Sample
}

func NewPointFeature(point Point) *PointFeature {
	return &PointFeature{
		Liberty:    0,
		CutPoint:   0,
		P:          point,
		BoardWorms: make(map[int]*Worm),
		Score:      0.0,
		SelectProb: 0.0,
		Fc:         core.NewSample(),
		WhiteDis:   40,
		BlackDis:   40,
	}
}

func (p *PointFeature) SetSelectProb(v float64, force ...bool) {
	if (len(force) > 0 && force[0]) || v > p.SelectProb {
		p.SelectProb = v
	}
}

func (b *Board) PointLiberty(p Point) int {
	n4 := b.Neighbor4(p.x, p.y)
	ret := 0
	for _, np := range n4 {
		if np.color == GRAY {
			ret += 1
		}
	}
	return ret
}

func FeatureString(label int, f *core.Sample) string {
	f.Label = label
	return string(f.ToString(false))
}
