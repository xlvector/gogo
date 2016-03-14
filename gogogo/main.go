package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"runtime"

	"github.com/xlvector/gogo"
)

func GenPatterns(path string, ch chan string) {
	log.Println(path)
	board := gogo.NewBoard()
	pats := board.GenPattern(path)
	for _, pat := range pats {
		ch <- pat.String()
	}
}

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	mode := flag.String("mode", "", "mode")
	input := flag.String("input", "", "input")
	output := flag.String("output", "", "output")
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())
	if *mode == "gen-pattern" {
		paths := gogo.TreeDir(*input, "sgf")
		patCh := make(chan string, 1000)
		go func() {
			for _, path := range paths {
				go GenPatterns(path, patCh)
			}
		}()
		f, _ := os.Create(*output)
		defer f.Close()
		w := bufio.NewWriter(f)
		for pat := range patCh {
			w.WriteString(pat)
			w.WriteString("\n")
		}
	}
}
