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
	board := gogo.NewBoard()
	pats := board.GenPattern(path, 0)
	for _, pat := range pats {
		ch <- pat.String()
	}
}

func GenPatternsThread(paths []string, total, split int, ch chan string) {
	log.Println("begin split:", split)
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
		for i := 0; i < 32; i++ {
			go GenPatternsThread(paths, 32, i, patCh)
		}

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
				if n > 60 {
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
