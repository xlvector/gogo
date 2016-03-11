package dataset

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"

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
				if rand.Float64() > 0.03 {
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

func SimpleFeatureString(label, index int, f []int64) string {
	ret := fmt.Sprintf("%d", label)
	for _, v := range f {
		ret += fmt.Sprintf("\t%d:1", v*1000+int64(index))
	}
	return ret
}

func genSimpleSamples(gt *gogo.GameTree) []string {
	path := gt.Path2Root()
	ret := []string{}
	board := gogo.NewBoard(gt.SGFSize())
	pdm := gogo.NewPointDistanceMap(board, gogo.PATTERN_SIZE)
	board.SetPointDistanceMap(pdm)
	lastPattern := []int64{}

	for i := len(path) - 2; i >= 1; i-- {
		cur := path[i].Point()
		if !cur.Valid() {
			break
		}
		fs := board.GenSimpleFeatures(lastPattern, cur)
		for p, pat := range fs {
			if p == board.Index(cur) {
				ret = append(ret, SimpleFeatureString(1, p, pat))
			} else {
				ret = append(ret, SimpleFeatureString(0, p, pat))
			}
		}
		lastPattern = board.PointSimpleFeature(cur, cur.Color())
		board.Put(cur.X(), cur.Y(), cur.Color())
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

func loadPatternModel(pat string) []map[int64]float64 {
	f, _ := os.Open(pat)
	reader := bufio.NewReader(f)
	ret := make([]map[int64]float64, gogo.PATTERN_SIZE)
	for i := 0; i < gogo.PATTERN_SIZE; i++ {
		ret[i] = make(map[int64]float64)
	}
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		tks := strings.Split(line, "\t")
		k, _ := strconv.Atoi(tks[0])
		p, _ := strconv.ParseInt(tks[1], 10, 64)
		v, _ := strconv.ParseFloat(tks[2], 64)
		ret[k][p] = v
	}
	return ret
}

func EvaluatePattern(sgfPath, modelPath string) (int, int) {
	gt := &gogo.GameTree{}
	buf, _ := ioutil.ReadFile(sgfPath)
	gt.ParseSGF(string(buf))
	return evalPattern(gt, modelPath)
}

func evalPattern(gt *gogo.GameTree, pat string) (int, int) {
	path := gt.Path2Root()
	board := gogo.NewBoard(gt.SGFSize())
	pdm := gogo.NewPointDistanceMap(board, gogo.PATTERN_SIZE)
	board.SetPointDistanceMap(pdm)
	patModel := loadPatternModel(pat)
	hit := 0
	total := 0
	for i := len(path) - 1; i >= 1; i-- {
		cur := path[i].Point()
		if !cur.Valid() {
			break
		}
		st := make([]gogo.IntFloatPairList, gogo.PATTERN_SIZE)
		for k := 0; k < gogo.PATTERN_SIZE; k++ {
			st[k] = make(gogo.IntFloatPairList, 0, 20)
		}
		for k, p := range board.W() {
			if p.Color() != gogo.GRAY {
				continue
			}
			pat := board.GetPatternHash(k)
			h := board.FeatureHash(gogo.MakePoint(p.X(), p.Y(), cur.Color()))
			for j := gogo.PATTERN_SIZE - 1; j >= 0; j-- {
				hh := pat[j] ^ h
				if v, ok := patModel[j][hh]; ok {
					st[j] = append(st[j], gogo.IntFloatPair{k, v})
					break
				}
			}
		}
		for j := gogo.PATTERN_SIZE - 1; j >= 0; j-- {
			if len(st[j]) == 0 {
				continue
			}
			sort.Sort(sort.Reverse(st[j]))
			if st[j][0].First == board.Index(cur) {
				hit += 1
			}
			fmt.Println(j, cur.String(), board.W()[st[j][0].First].String(), st[j][0].Second, st[j][0:10])
			break
		}
		total += 1
		board.Put(cur.X(), cur.Y(), cur.Color())
	}
	fmt.Println(hit, total, float64(hit)/float64(total))
	return hit, total
}

func evaluateLRModel(gt *gogo.GameTree, model *lr.LogisticRegression) (int, int) {
	path := gt.Path2Root()
	board := gogo.NewBoard(gt.SGFSize())
	pdm := gogo.NewPointDistanceMap(board, gogo.PATTERN_SIZE)
	board.SetPointDistanceMap(pdm)
	lastPattern := []int64{}
	hit := 0
	total := 0
	for i := len(path) - 2; i >= 1; i-- {
		cur := path[i].Point()
		if !cur.Valid() {
			break
		}
		rank := make(gogo.IntFloatPairList, 0, 50)
		for j, p := range board.W() {
			if p.Color() != gogo.GRAY {
				continue
			}
			s := core.NewSample()
			pat := board.PointSimpleFeature(p, cur.Color())
			pat = append(pat, lastPattern...)
			for _, v := range pat {
				s.AddFeature(core.Feature{v*1000 + int64(j), 1.0})
			}
			prob := model.Predict(s)
			rank = append(rank, gogo.IntFloatPair{j, prob})
		}
		sort.Sort(sort.Reverse(rank))
		for k := 0; k < 1 && k < len(rank); k++ {
			if rank[k].First == board.Index(cur) {
				hit += 1
				break
			}
		}
		total += 1
		fmt.Println(cur.String(), board.W()[rank[0].First].String(), rank[0].Second)
		lastPattern = board.PointSimpleFeature(cur, cur.Color())
		board.Put(cur.X(), cur.Y(), cur.Color())
		if i%20 == 0 {
			fmt.Println(board.String(cur))
		}
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
