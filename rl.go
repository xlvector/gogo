package gogo

import (
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/xlvector/hector/core"
)

func BatchRLBattle(b *Board) {
	lock := &sync.Mutex{}
	for i := 0; i < 50; i++ {
		go func() {
			rank := b.Copy().RLBattle(BLACK)
			lock.Lock()
			defer lock.Unlock()
			for k, v := range rank {
				v1, _ := b.Model.Model[k]
				v1 += 0.001 * float64(v)
				b.Model.Model[k] = v1
			}
		}()
	}
	b.Model.SaveModel("expert.model.rl")
}

func (b *Board) RLBattle(c Color) map[int64]int {
	rand.Seed(time.Now().UnixNano())
	var fs []int64
	p := -1
	n := 0
	colorFs := make(map[Color]map[int64]int)
	colorFs[BLACK] = make(map[int64]int)
	colorFs[WHITE] = make(map[int64]int)
	ret := make(map[int64]int)
	for n < 350 {
		pass := 0
		p, fs = b.GenRLBattleMove(c)
		if p < 0 {
			pass += 1
		} else {
			for _, v := range fs {
				if v1, ok := colorFs[c][v]; ok {
					colorFs[c][v] = v1 + 1
				} else {
					colorFs[c][v] = 1
				}
			}
		}

		oc := OpColor(c)
		p, fs = b.GenRLBattleMove(oc)
		if p < 0 {
			pass += 1
		} else {
			for _, v := range fs {
				if v1, ok := colorFs[oc][v]; ok {
					colorFs[oc][v] = v1 + 1
				} else {
					colorFs[oc][v] = 1
				}
			}
		}
		if pass >= 2 {
			break
		}

		n += 1
	}
	s := b.Score()
	if s > 0 {
		for k, v := range colorFs[BLACK] {
			v1, _ := ret[k]
			ret[k] = v1 + v
		}
		for k, v := range colorFs[WHITE] {
			v1, _ := ret[k]
			ret[k] = v1 - v
		}
	} else {
		for k, v := range colorFs[WHITE] {
			v1, _ := ret[k]
			ret[k] = v1 + v
		}
		for k, v := range colorFs[BLACK] {
			v1, _ := ret[k]
			ret[k] = v1 - v
		}
	}
	log.Println(s)
	return ret
}

func (b *Board) GenRLBattleMove(c Color) (int, []int64) {
	rank := make(map[int]float64)
	fs := make(map[int][]int64)
	for k, _ := range b.Points {
		if ok, _ := b.CanPut(k, c); ok {
			pr := rand.Float64() * 0.1
			if b.Model != nil {
				pat := b.FinalPatternHash(k, c)
				smp := b.PatternFeature(k, c, b.LastPattern, pat)
				fs[k] = smp
				sample := core.NewSample()
				for _, v := range smp {
					sample.AddFeature(core.Feature{v, 1.0})
				}
				pr = b.Model.Predict(sample)
			}
			rank[k] = pr
		}
	}
	topn := TopN(rank, 10)
	if len(topn) == 0 {
		return -1, nil
	}
	k := rand.Intn(len(topn))
	b.Put(topn[k].First, c)
	return topn[k].First, fs[topn[k].First]
}
