package gogo

import (
	"io/ioutil"
	"log"
	"strconv"
	"sync"
)

func GenDLDataset(sgfFile string, out chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	buf, _ := ioutil.ReadFile(sgfFile)
	log.Println(string(buf))
	gt := NewGameTree(SIZE)
	gt.ParseSGF(string(buf))

	path := gt.Path2Root()
	board := NewBoard()

	n := 0
	for i := len(path) - 2; i >= 0; i-- {
		cur := path[i]
		sample := board.DLFeature(cur.stone)
		line := strconv.Itoa(PosIndex(cur.x, cur.y))
		line += "\t"
		for _, v := range sample {
			line += strconv.Itoa(int(v))
		}
		out <- line
		n += 1
		if ok := board.Put(PosIndex(cur.x, cur.y), cur.stone); !ok {
			log.Println(cur.x, cur.y, cur.stone)
		}
	}
	log.Println(sgfFile, len(out), n, len(path))
}

func (b *Board) DLFeature(stone Color) []byte {
	ret := make([]byte, NPOINT*2)
	for k, c := range b.Points {
		if c == stone {
			ret[k] = 1
		} else if c == OpColor(stone) {
			ret[k+NPOINT] = 1
		}
	}
	return ret
}
