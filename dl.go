package gogo

import (
	"io/ioutil"
	"log"
	"strconv"
	"sync"
)

func GenDLDataset(sgfFile string, out chan string, wg *sync.WaitGroup) {
	buf, _ := ioutil.ReadFile(sgfFile)
	gt := NewGameTree(SIZE)
	gt.ParseSGF(string(buf))

	path := gt.Path2Root()
	board := NewBoard()

	for i := len(path) - 2; i >= 0; i-- {
		cur := path[i]
		sample := board.DLFeature(cur.stone)
		line := strconv.Itoa(PosIndex(cur.x, cur.y))
		line += "\t"
		for _, v := range sample {
			line += strconv.Itoa(int(v))
		}
		out <- line
		if ok := board.Put(PosIndex(cur.x, cur.y), cur.stone); !ok {
			break
		}
	}
	wg.Done()
	log.Println(sgfFile)
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
