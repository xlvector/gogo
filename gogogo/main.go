package main

import (
	"bufio"
	"container/list"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"strconv"
	"strings"

	"github.com/xlvector/gogo"
	"github.com/xlvector/gogo/dataset"
	"github.com/xlvector/hector/dt"
)

func combineSGFs(root, ext string) *gogo.GameTree {
	q := list.New()
	q.PushBack(root)
	gt := gogo.NewGameTree(19)
	for q.Len() > 0 {
		p := q.Front()
		q.Remove(p)
		v := p.Value.(string)
		n := 0
		filepath.Walk(v, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() && path != root {
				fmt.Println("==>", path)
				q.PushBack(path)
			} else {
				if strings.HasSuffix(path, ext) {

					buf, _ := ioutil.ReadFile(path)
					newGt := &gogo.GameTree{}
					newGt.ParseSGF(string(buf))
					if gt == nil {
						gt = newGt
					} else {
						gt.Combine(newGt, 10)
					}
					n += 1
					if n%100 == 0 {
						fmt.Println(n, path)
					}
				}
			}
			return nil
		})
	}
	return gt
}

func getSimpleFeatures(root, ext string, ch chan string) {
	q := list.New()
	q.PushBack(root)
	n := 0
	sgfCh := make(chan string, 1000)
	ncpu := runtime.NumCPU()
	for i := 0; i < ncpu; i++ {
		go func() {
			for buf := range sgfCh {
				fs := dataset.GenSimpleSamplesFromSGF(buf)
				for _, f := range fs {
					ch <- f
				}
			}
		}()
	}
	for q.Len() > 0 && n < 1000 {
		p := q.Front()
		q.Remove(p)
		v := p.Value.(string)
		filepath.Walk(v, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				q.PushBack(path)
			} else {
				if strings.HasSuffix(path, ext) {
					fmt.Println(path)
					n += 1
					buf, _ := ioutil.ReadFile(path)
					sgfCh <- string(buf)
				}
			}
			return nil
		})
	}
	close(sgfCh)
}

func getAllFiles(root, ext string, ch chan string) {
	q := list.New()
	q.PushBack(root)
	n := 0
	sgfCh := make(chan string, 1000)
	ncpu := runtime.NumCPU()
	for i := 0; i < ncpu; i++ {
		go func() {
			for buf := range sgfCh {
				samples := dataset.GenSamplesFromSGF(buf, gogo.WHITE)
				for _, sample := range samples {
					fs := sample.FeatureString()
					for _, f := range fs {
						ch <- strconv.Itoa(sample.K) + "\t" + f
					}
				}
			}
		}()
	}
	for q.Len() > 0 && n < 1000 {
		p := q.Front()
		q.Remove(p)
		v := p.Value.(string)
		filepath.Walk(v, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				q.PushBack(path)
			} else {
				if strings.HasSuffix(path, ext) {
					fmt.Println(path)
					n += 1
					buf, _ := ioutil.ReadFile(path)
					sgfCh <- string(buf)
				}
			}
			return nil
		})
	}
	close(sgfCh)
}

func patternAllFiles(root, ext string, ch chan string) {
	q := list.New()
	q.PushBack(root)
	n := 0
	sgfCh := make(chan string, 1000)
	ncpu := runtime.NumCPU()
	for i := 0; i < ncpu; i++ {
		go func() {
			for buf := range sgfCh {
				patterns := dataset.GenPatternFromSGF(buf)
				for _, pat := range patterns {
					ch <- pat
				}
			}
		}()
	}
	for q.Len() > 0 && n < 1000 {
		p := q.Front()
		q.Remove(p)
		v := p.Value.(string)
		filepath.Walk(v, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				q.PushBack(path)
			} else {
				if strings.HasSuffix(path, ext) {
					fmt.Println(path)
					n += 1
					buf, _ := ioutil.ReadFile(path)
					sgfCh <- string(buf)
				}
			}
			return nil
		})
	}
	close(sgfCh)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	mode := flag.String("mode", "simple", "mode")
	sgfFile := flag.String("sgf-file", "", "sgf file path")
	sgfFolder := flag.String("sgf-folder", "", "sgf folder path")
	output := flag.String("output", "", "output path")
	model := flag.String("model", "", "model path")
	qipu := flag.String("qipu", "", "qipu path")
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file")
	mcts := flag.String("mcts", "", "use mcts")
	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	game := &gogo.Game{}
	if *mode == "gtp" {
		game.GTP()
	} else if *mode == "simple" {
		game.Init(19)
		if len(*model) > 0 {
			models := []*dt.RandomForest{}
			tks := strings.Split(*model, ",")
			for _, tk := range tks {
				rf := &dt.RandomForest{}
				rf.LoadModel(tk)
				models = append(models, rf)
			}
			game.SetModel(models)
		}
		if len(*qipu) > 0 {
			game.QipuGT = &gogo.GameTree{}
			buf, _ := ioutil.ReadFile(*qipu)
			game.QipuGT.ParseSGF(string(buf))
		}
		reader := bufio.NewReader(os.Stdin)
		game.Print()
		for {
			line, _ := reader.ReadString('\n')
			line = strings.TrimSpace(line)
			if line == "pass" {
				break
			}
			if len(line) > 3 || len(line) < 2 {
				fmt.Println("invalid line: ", line)
				continue
			}
			x := strings.Index(gogo.LX, strings.ToUpper(line[0:1]))
			y, _ := strconv.Atoi(line[1:])
			y -= 1
			if x < 0 || y < 0 {
				continue
			}
			fmt.Println(x, y)
			game.Put(gogo.BLACK, x, y)
			if *mcts != "1" {
				game.GenMove(gogo.WHITE)
			} else {
				game.MCTSMove(gogo.WHITE)
			}
			game.Print()
		}
	} else if *mode == "sgf" {
		if len(*sgfFile) > 0 {
			buf, _ := ioutil.ReadFile(*sgfFile)
			game.Init(19)
			game.GT.ParseSGF(string(buf))
			/*
				game.ResetBoardFromGT()
				game.Print()
				path := game.GT.Path2Root()
				for i := len(path) - 1; i >= 0; i-- {
					fmt.Println(path[i].Point().String())
				}

					samples := dataset.GenSamplesFromSGF(string(buf), gogo.WHITE)
					for _, smp := range samples {
						fmt.Println(smp.Board.String(gogo.InvalidPoint()))
					}
			*/
			return
		}
		if len(*sgfFolder) > 0 {
			ch := make(chan string, 1000)
			go func() {
				fout, _ := os.Create(*output)
				defer fout.Close()
				writer := bufio.NewWriter(fout)
				for line := range ch {
					writer.WriteString(line)
					writer.WriteString("\n")
				}
			}()
			getAllFiles(*sgfFolder, ".sgf", ch)
		}
	} else if *mode == "combine" {
		fmt.Println(*sgfFolder)
		gt := combineSGFs(*sgfFolder, ".sgf")
		ioutil.WriteFile(*output, []byte(gt.WriteSGF()), 0655)
	} else if *mode == "pattern" {
		ch := make(chan string, 1000)
		go func() {
			fout, _ := os.Create(*output)
			defer fout.Close()
			writer := bufio.NewWriter(fout)
			for line := range ch {
				writer.WriteString(line)
				writer.WriteString("\n")
			}
		}()
		patternAllFiles(*sgfFolder, ".sgf", ch)

	} else if *mode == "simple-feature" {
		ch := make(chan string, 1000)
		go func() {
			fout, _ := os.Create(*output)
			defer fout.Close()
			writer := bufio.NewWriter(fout)
			for line := range ch {
				writer.WriteString(line)
				writer.WriteString("\n")
			}
		}()
		getSimpleFeatures(*sgfFolder, ".sgf", ch)
	} else if *mode == "eval-lr" {
		dataset.EvaluateLRModel(*sgfFile, *model)
	}
}
