package gogo

import (
	"io/ioutil"
	"math/rand"
	"strconv"
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
			break
		}
		curK := PosIndex(cur.x, cur.y)
		curPat := b.FinalPatternHash(curK, cur.stone)
		ret = append(ret, PatternSample{b.PatternFeature(curK, lastPat, curPat), 1})

		vps := b.RandomSelectValidPoint(2, cur.stone)
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
