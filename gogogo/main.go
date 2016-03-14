package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/xlvector/gogo"
	"github.com/xlvector/hector/lr"
)

func GenPatterns(path string, ch chan string) {
	log.Println(path)
	for r := 0; r < 8; r++ {
		board := gogo.NewBoard()
		pats := board.GenPattern(path, r)
		for _, pat := range pats {
			ch <- pat.String()
		}
	}
}

func GenPatternsThread(paths []string, total, split int, ch chan string) {
	for i, path := range paths {
		if i%total != split {
			continue
		}
		GenPatterns(path, ch)
	}
}

func EvalModel(sgf string, model *lr.LogisticRegression) {
	board := gogo.NewBoard()
	board.Model = model
}

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	mode := flag.String("mode", "", "mode")
	input := flag.String("input", "", "input")
	output := flag.String("output", "", "output")
	model := flag.String("model", "", "model path")
	flag.Parse()

	runtime.GOMAXPROCS(runtime.NumCPU())
	if *mode == "gen-pattern" {
		paths := gogo.TreeDir(*input, "sgf")
		patCh := make(chan string, 1000)
		go GenPatternsThread(paths, 16, 0, patCh)
		go GenPatternsThread(paths, 16, 1, patCh)
		go GenPatternsThread(paths, 16, 2, patCh)
		go GenPatternsThread(paths, 16, 3, patCh)
		go GenPatternsThread(paths, 16, 4, patCh)
		go GenPatternsThread(paths, 16, 5, patCh)
		go GenPatternsThread(paths, 16, 6, patCh)
		go GenPatternsThread(paths, 16, 7, patCh)
		go GenPatternsThread(paths, 16, 8, patCh)
		go GenPatternsThread(paths, 16, 9, patCh)
		go GenPatternsThread(paths, 16, 10, patCh)
		go GenPatternsThread(paths, 16, 11, patCh)
		go GenPatternsThread(paths, 16, 12, patCh)
		go GenPatternsThread(paths, 16, 13, patCh)
		go GenPatternsThread(paths, 16, 14, patCh)
		go GenPatternsThread(paths, 16, 15, patCh)

		go func() {
			tc := time.NewTicker(time.Second)
			n := 0
			for _ = range tc.C {
				if len(patCh) == 0 {
					n += 1
					log.Println("==>", n)
				} else {
					n = 0
				}
				if n > 10 {
					close(patCh)
				}
			}
		}()

		f, _ := os.Create(*output)
		defer f.Close()
		w := bufio.NewWriter(f)
		for pat := range patCh {
			w.WriteString(pat)
			w.WriteString("\n")
		}
	} else if *mode == "eval" {
		board := gogo.NewBoard()
		board.Model = &lr.LogisticRegression{}
		board.Model.LoadModel(*model)
		hit, total := board.EvaluateModel(*input)
		log.Println(hit, total, float64(hit)/float64(total))
	}
}
