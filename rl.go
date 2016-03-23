package gogo

import (
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/xlvector/hector/core"
	"github.com/xlvector/hector/lr"
)

func BatchRLBattle(b *Board) {
	b.Model2 = &lr.LogisticRegression{}
	b.Model2.Model = make(map[int64]float64)
	for k, v := range b.Model.Model {
		b.Model2.Model[k] = v
	}
	for k := 0; k < 100; k++ {
		wg := &sync.WaitGroup{}
		win := 0
		ch := make(chan map[int64]int, 100)
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				s, rank := b.Copy().RLBattle(BLACK)
				if s > 0 {
					win += 1
					ch <- rank
				}
				wg.Done()
			}()
		}
		wg.Wait()
		close(ch)
		for rank := range ch {
			for k, v := range rank {
				v1, _ := b.Model.Model[k]
				v1 += 0.0001 * float64(v)
				b.Model.Model[k] = v1
			}
		}
		log.Println(win)
	}
}

func (b *Board) RLBattle(c Color) (float64, map[int64]int) {
	rand.Seed(time.Now().UnixNano())
	var fs []int64
	p := -1
	n := 0
	colorFs := make(map[Color]map[int64]int)
	colorFs[BLACK] = make(map[int64]int)
	colorFs[WHITE] = make(map[int64]int)
	//ret := make(map[int64]int)
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
	/*
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
	*/
	return s, colorFs[BLACK]
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
				if c == BLACK {
					pr = b.Model.Predict(sample)
				} else {
					pr = b.Model2.Predict(sample)
				}
			}
			rank[k] = pr
		}
	}
	topn := TopN(rank, 2)
	if len(topn) == 0 {
		return -1, nil
	}
	k := rand.Intn(len(topn))
	b.Put(topn[k].First, c)
	return topn[k].First, fs[topn[k].First]
}
