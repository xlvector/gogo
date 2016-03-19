package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"
	"sync"
	"time"

	"github.com/xlvector/gogo"
	"github.com/xlvector/hector/lr"
)

func GenDLDataset(root, output string) {
	paths := gogo.TreeDir(root, "sgf")
	f, _ := os.Create(output)
	defer f.Close()
	writer := bufio.NewWriter(f)
	lock := &sync.Mutex{}
	var wg sync.WaitGroup
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func(k int, wg *sync.WaitGroup, paths []string) {
			for j, path := range paths {
				if j%32 == k {
					lines := gogo.GenDLDataset(path)
					lock.Lock()
					for _, line := range lines {
						writer.WriteString(line)
						writer.WriteString("\n")
					}
					lock.Unlock()
				}
			}
			wg.Done()
		}(i, &wg, paths)
	}
	wg.Wait()
}

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
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file")
	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

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
		hit, total := board.EvaluateModel(*input, true)
		log.Println(hit, total, float64(hit)/float64(total))
	} else if *mode == "eval-folder" {
		m := &lr.LogisticRegression{}
		m.LoadModel(*model)
		hit := 0
		total := 0
		paths := gogo.TreeDir(*input, "sgf")
		for _, path := range paths {
			log.Println(path)
			board := gogo.NewBoard()
			board.Model = m
			h, t := board.EvaluateModel(path, false)
			hit += h
			total += t
			log.Println(h, t, hit, total, float64(hit)/float64(total))
		}
	} else if *mode == "simple" {
		board := gogo.NewBoard()
		if len(*model) > 0 {
			board.Model = &lr.LogisticRegression{}
			board.Model.LoadModel(*model)
		}
		log.Println(board.String(nil))
		gt := gogo.NewGameTree(gogo.SIZE)
		reader := bufio.NewReader(os.Stdin)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				break
			}
			line = strings.TrimSpace(line)
			if line == "pass" {
				break
			}
			line = strings.ToUpper(line)
			if ok := board.PutLabel("B" + line); !ok {
				log.Println("invalid")
				continue
			}
			last, c := board.LastMove()
			x, y := gogo.IndexPos(last)
			gt.Add(gogo.NewGameTreeNode(c, x, y))
			log.Println(board.String(nil))

			if ok := board.MCTSMove(gogo.WHITE, gt, 16, 10000); !ok {
				break
			}
			log.Println(board.String(nil))
		}
	} else if *mode == "self" {
		board := gogo.NewBoard()
		if len(*model) > 0 {
			board.Model = &lr.LogisticRegression{}
			board.Model.LoadModel(*model)
		}
		log.Println(board.String(nil))
		gt1 := gogo.NewGameTree(gogo.SIZE)
		gt2 := gogo.NewGameTree(gogo.SIZE)
		for {
			gt1.CurrentChild()
			if ok := board.MCTSMove(gogo.BLACK, gt1, 10, 1000); !ok {
				break
			}
			log.Println(board.String(nil))
			{
				last, _ := board.LastMove()
				lastX, lastY := gogo.IndexPos(last)
				gt2.Add(gogo.NewGameTreeNode(gogo.BLACK, lastX, lastY))
			}

			gt2.CurrentChild()
			if ok := board.MCTSMove(gogo.WHITE, gt2, 20, 1000); !ok {
				break
			}
			log.Println(board.String(nil))
			{
				last, _ := board.LastMove()
				lastX, lastY := gogo.IndexPos(last)
				gt1.Add(gogo.NewGameTreeNode(gogo.WHITE, lastX, lastY))
			}
		}
		log.Println(board.Score())
	} else if *mode == "dl-data" {
		GenDLDataset(*input, *output)
	}
}
