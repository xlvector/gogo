package dataset

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"strconv"

	"github.com/xlvector/gogo"
	"github.com/xlvector/hector/core"
	"github.com/xlvector/hector/lr"
)

type Sample struct {
	Board    *gogo.Board
	Info     *gogo.BoardInfo
	K        int
	CurStep  gogo.Point
	NextStep gogo.Point
}

func patternString(label int, pat []int64, h int64) string {
	ret := strconv.Itoa(label)
	for k, v := range pat {
		ret += "\t" + strconv.Itoa(k) + ":" + strconv.FormatInt(v^h, 10)
	}
	return ret
}

func genPatterns(gt *gogo.GameTree) []string {
	path := gt.Path2Root()
	board := gogo.NewBoard(gt.SGFSize())
	pdm := gogo.NewPointDistanceMap(board, gogo.PATTERN_SIZE)
	board.SetPointDistanceMap(pdm)
	ret := []string{}
	for i := len(path) - 1; i >= 0; i-- {
		cur := path[i].Point()
		if !board.Valid(cur) {
			break
		}

		for k, p := range board.W() {
			if p.Color() != gogo.GRAY {
				continue
			}
			label := 0
			if p.X() == cur.X() && p.Y() == cur.Y() {
				label = 1
			} else {
				if rand.Float64() > 0.1 {
					continue
				}
			}
			pat := board.GetPatternHash(k)
			h := board.FeatureHash(gogo.MakePoint(p.X(), p.Y(), cur.Color()))
			ret = append(ret, patternString(label, pat, h))
		}
		board.Put(cur.X(), cur.Y(), cur.Color())
	}
	return ret
}

func GenPatternFromSGF(buf string) []string {
	gt := &gogo.GameTree{}
	gt.ParseSGF(buf)
	return genPatterns(gt)
}

func genSamples(gt *gogo.GameTree, stone gogo.Color) []*Sample {
	path := gt.Path2Root()
	ret := []*Sample{}
	board := gogo.NewBoard(gt.SGFSize())
	for i := len(path) - 1; i >= 1; i-- {
		prob := float64(len(path) - i)
		prob = 100.0/100.0 + prob
		if rand.Float64() > prob {
			continue
		}
		cur := path[i].Point()
		if !cur.Valid() {
			break
		}
		next := path[i-1].Point()
		if !next.Valid() {
			break
		}
		board.Put(cur.X(), cur.Y(), cur.Color())
		sample := &Sample{
			K:        len(path) - i - 1,
			Board:    board.Copy(),
			NextStep: next,
			CurStep:  cur,
		}
		sample.GenComplexFeatures(next.Color())
		ret = append(ret, sample)
	}
	return ret
}

func FeatureString(label int, f map[int64]byte) string {
	ret := fmt.Sprintf("%d", label)
	for k, v := range f {
		ret += fmt.Sprintf("\t%d:%d", k, int(v))
	}
	return ret
}

func genSimpleSamples(gt *gogo.GameTree) []string {
	path := gt.Path2Root()
	ret := []string{}
	board := gogo.NewBoard(gt.SGFSize())
	pdm := gogo.NewPointDistanceMap(board, gogo.PATTERN_SIZE)
	board.SetPointDistanceMap(pdm)
	last := gogo.InvalidPoint()
	for i := len(path) - 1; i >= 1; i-- {
		cur := path[i].Point()
		if !cur.Valid() {
			break
		}
		fs := board.GenSimpleFeatures(last, cur)
		ratio := 15.0 / float64(len(fs))
		for p, f := range fs {
			if p == board.Index(cur) {
				ret = append(ret, FeatureString(1, f))
			} else {
				if rand.Float64() < ratio {
					ret = append(ret, FeatureString(0, f))
				}
			}
		}
		board.Put(cur.X(), cur.Y(), cur.Color())
		last = cur
	}
	return ret
}

func EvaluateLRModel(sgfPath, modelPath string) (int, int) {
	model := &lr.LogisticRegression{
		Model: make(map[int64]float64),
	}
	model.LoadModel(modelPath)

	gt := &gogo.GameTree{}
	buf, _ := ioutil.ReadFile(sgfPath)
	gt.ParseSGF(string(buf))
	return evaluateLRModel(gt, model)
}

func evaluateLRModel(gt *gogo.GameTree, model *lr.LogisticRegression) (int, int) {
	path := gt.Path2Root()
	board := gogo.NewBoard(gt.SGFSize())
	last := gogo.InvalidPoint()
	hit := 0
	total := 0
	for i := len(path) - 1; i >= 1; i-- {
		cur := path[i].Point()
		if !cur.Valid() {
			break
		}
		fs := board.GenSimpleFeatures(last, cur)
		best := -1
		maxProb := 0.0
		for p, f := range fs {
			s := core.NewSample()
			for k, v := range f {
				s.AddFeature(core.Feature{k, float64(v)})
			}
			prob := model.Predict(s)
			if maxProb < prob {
				best = p
			}
		}
		if best == board.Index(cur) {
			hit += 1
		}
		total += 1
		fmt.Println(cur.String(), board.W()[best].String())
		board.Put(cur.X(), cur.Y(), cur.Color())
		last = cur
	}
	fmt.Println(hit, total, float64(hit)/float64(total))
	return hit, total
}

func GenSimpleSamplesFromSGF(buf string) []string {
	gt := &gogo.GameTree{}
	gt.ParseSGF(buf)
	return genSimpleSamples(gt)
}

func GenSamplesFromSGF(buf string, stone gogo.Color) []*Sample {
	gt := &gogo.GameTree{}
	gt.ParseSGF(buf)
	return genSamples(gt, stone)
}

func (s *Sample) feature(i, t int) int {
	return t*1000 + i
}

func (s *Sample) GenComplexFeatures(stone gogo.Color) *Sample {
	s.Info = s.Board.CollectBoardInfo(gogo.InvalidPoint())
	s.Info.GenComplexFeatures(stone)
	return s
}

func (s *Sample) FeatureString() []string {
	ret := []string{}
	for p, pf := range s.Info.PointFetures {
		if pf.P.Color() != gogo.GRAY {
			continue
		}
		if p == s.Board.Index(s.NextStep) {
			ret = append(ret, gogo.FeatureString(1, pf.Fc))
		} else {
			if rand.Float64() < 0.05 {
				ret = append(ret, gogo.FeatureString(0, pf.Fc))
			}
		}
	}
	return ret
}
