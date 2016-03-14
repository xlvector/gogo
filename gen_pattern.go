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

func (p *PatternSample) String() string {
	ret := strconv.Itoa(p.Label)
	for _, k := range p.Pattern {
		ret += "\t"
		ret += strconv.FormatInt(k, 10)
		ret += ":1"
	}
	return ret
}

func (b *Board) PatternFeature(k int, last, cur []int64) []int64 {
	ret := make([]int64, 0, len(last)+len(cur))
	for _, v := range last {
		ret = append(ret, v*1000+int64(k))
	}
	for _, v := range cur {
		ret = append(ret, v*1000+int64(k))
	}
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

func (b *Board) GenPattern(sgf string) []PatternSample {
	buf, _ := ioutil.ReadFile(sgf)
	gt := NewGameTree(SIZE)
	gt.ParseSGF(string(buf))
	path := gt.Path2Root()
	lastPat := []int64{}
	ret := []PatternSample{}
	for i := len(path) - 2; i >= 0; i-- {
		cur := path[i]
		if PosOutBoard(cur.x, cur.y) {
			continue
		}
		curK := PosIndex(cur.x, cur.y)
		curPat := b.FinalPatternHash(curK, cur.stone)
		ret = append(ret, PatternSample{b.PatternFeature(curK, lastPat, curPat), 1})

		vps := b.RandomSelectValidPoint(5, cur.stone)
		for p, _ := range vps {
			if p == curK {
				continue
			}
			pat := b.FinalPatternHash(p, cur.stone)
			ret = append(ret, PatternSample{b.PatternFeature(p, lastPat, pat), 0})
		}
		lastPat = curPat
		ok := b.Put(PosIndex(cur.x, cur.y), cur.stone)
		if !ok {
			break
		}
	}
	return ret
}

func (b *Board) EvaluateModel(sgf string) (int, int) {
	buf, _ := ioutil.ReadFile(sgf)
	gt := NewGameTree(SIZE)
	gt.ParseSGF(string(buf))
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
			spat := b.PatternFeature(p, lastPat, pat)
			sample := core.NewSample()
			for _, f := range spat {
				sample.AddFeature(core.Feature{f, 1.0})
			}
			pr := b.Model.Predict(sample)
			rank = append(rank, IntFloatPair{p, pr})
		}
		sort.Sort(sort.Reverse(rank))
		x1, y1 := IndexPos(rank[0].First)
		log.Println(PointString(cur.x, cur.y, cur.stone), PointString(x1, y1, cur.stone), rank[0].Second)
		if rank[0].First == PosIndex(cur.x, cur.y) {
			hit += 1
		}
		total += 1
		lastPat = b.FinalPatternHash(PosIndex(cur.x, cur.y), cur.stone)
		ok := b.Put(PosIndex(cur.x, cur.y), cur.stone)
		if !ok {
			break
		}
	}
	return hit, total
}