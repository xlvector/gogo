package gogo

import (
	"io/ioutil"
	"log"
	"math/rand"
	"sort"
	"strconv"

	"github.com/xlvector/hector/core"
)

type PatternSample struct {
	Pattern []int64
	Label   int
}

type PairwiseSample struct {
	Features map[int64]int
	Label    int
}

func NewPairwiseSample(a, b *PatternSample) *PairwiseSample {
	ret := &PairwiseSample{
		Features: make(map[int64]int),
		Label:    a.Label - b.Label,
	}
	for _, p := range a.Pattern {
		ret.Features[p] = 1
	}
	for _, p := range b.Pattern {
		if _, ok := ret.Features[p]; ok {
			delete(ret.Features, p)
		} else {
			ret.Features[p] = -1
		}
	}
	return ret
}

func (p *PairwiseSample) String() string {
	label := 1
	if p.Label < 0 {
		label = 0
	}
	ret := strconv.Itoa(label)
	for k, v := range p.Features {
		ret += "\t"
		ret += strconv.FormatInt(k, 10)
		ret += ":"
		ret += strconv.Itoa(v)
	}
	return ret
}

func (p *PatternSample) String() string {
	ret := strconv.Itoa(p.Label)
	for _, k := range p.Pattern {
		ret += "\t"
		ret += strconv.FormatInt(k, 10)
		ret += ":1"
	}
	return ret
}

func (b *Board) PatternFeature(k int, c Color, last, cur []int64) []int64 {
	ret := make([]int64, 0, len(last)+len(cur))
	for _, v := range last {
		ret = append(ret, v*1000+int64(k))
	}
	for _, v := range cur {
		ret = append(ret, v*1000+int64(k))
	}
	ret = append(ret, b.LocalFeature(k, c)...)
	return ret
}

func (b *Board) RandomSelectValidPoint(n int, c Color) map[int]byte {
	ret := make(map[int]byte)

	for i := 0; i < NPOINT*2 && len(ret) < n; i++ {
		k := rand.Intn(NPOINT)
		if ok, _ := b.CanPut(k, c); ok {
			ret[k] = 1
		}
	}
	return ret
}

func (b *Board) Rotate(x, y, r int) (int, int) {
	switch r {
	case 0:
		return x, y
	case 1:
		return SIZE - x - 1, y
	case 2:
		return x, SIZE - y - 1
	case 3:
		return SIZE - x - 1, SIZE - y - 1
	case 4:
		return y, x
	case 5:
		return SIZE - y - 1, x
	case 6:
		return y, SIZE - x - 1
	case 7:
		return SIZE - y - 1, SIZE - x - 1
	default:
		return x, y
	}
}

func (b *Board) GenPattern(sgf string, rotate int) []*PairwiseSample {
	buf, _ := ioutil.ReadFile(sgf)
	gt := NewGameTree(SIZE)
	gt.ParseSGF(string(buf))
	if gt.HasHandicap() {
		return []*PairwiseSample{}
	}
	path := gt.Path2Root()
	ret := make([]*PairwiseSample, 12)
	lastPat := []int64{}
	for i := len(path) - 2; i >= 0; i-- {
		cur := path[i]
		if PosOutBoard(cur.x, cur.y) {
			continue
		}
		cur.x, cur.y = b.Rotate(cur.x, cur.y, rotate)
		curK := PosIndex(cur.x, cur.y)
		curPat := b.FinalPatternHash(curK, cur.stone)
		curPf := &PatternSample{b.PatternFeature(curK, cur.stone, lastPat, curPat), 1}

		vps := b.RandomSelectValidPoint(5, cur.stone)
		for p, _ := range vps {
			if p == curK {
				continue
			}
			pat := b.FinalPatternHash(p, cur.stone)
			spat := &PatternSample{b.PatternFeature(p, cur.stone, lastPat, pat), 0}
			ret = append(ret, NewPairwiseSample(curPf, spat))
			ret = append(ret, NewPairwiseSample(spat, curPf))
		}
		lastPat = curPat
		ok := b.Put(PosIndex(cur.x, cur.y), cur.stone)
		if !ok {
			break
		}
		b.RefreshInfluenceAndTerritory()
	}
	return ret
}

func (b *Board) EvaluateModel(sgf string, withLog bool) (int, int) {
	buf, _ := ioutil.ReadFile(sgf)
	gt := NewGameTree(SIZE)
	gt.ParseSGF(string(buf))
	if gt.HasHandicap() {
		return 0, 0
	}
	path := gt.Path2Root()
	lastPat := []int64{}
	hit := 0
	total := 0
	for i := len(path) - 2; i >= 0; i-- {
		cur := path[i]
		if PosOutBoard(cur.x, cur.y) {
			continue
		}

		rank := make(IntFloatPairList, 0, 100)
		for p, _ := range b.Points {
			if ok, _ := b.CanPut(p, cur.stone); !ok {
				continue
			}
			pat := b.FinalPatternHash(p, cur.stone)
			spat := b.PatternFeature(p, cur.stone, lastPat, pat)
			sample := core.NewSample()
			for _, f := range spat {
				sample.AddFeature(core.Feature{f, 1.0})
			}
			pr := b.Model.Predict(sample)
			rank = append(rank, IntFloatPair{p, pr})
		}
		sort.Sort(sort.Reverse(rank))
		x1, y1 := IndexPos(rank[0].First)
		if withLog {
			log.Println(PointString(cur.x, cur.y, cur.stone), PointString(x1, y1, cur.stone), rank[0].Second)
			if i%10 == 0 {
				mark := make(map[int]string)
				for k := 0; k < 26 && k < len(rank); k++ {
					mark[rank[k].First] = string(97 + k)
				}
				log.Println(b.String(mark))
			}
		}
		if rank[0].First == PosIndex(cur.x, cur.y) {
			hit += 1
		}
		total += 1
		lastPat = b.FinalPatternHash(PosIndex(cur.x, cur.y), cur.stone)
		ok := b.Put(PosIndex(cur.x, cur.y), cur.stone)
		if !ok {
			break
		}
		b.RefreshInfluenceAndTerritory()
	}
	return hit, total
}
