package gogo

import (
	"io/ioutil"
	"log"
	"strconv"
)

func GenDLDataset(sgfFile string) []string {
	buf, err := ioutil.ReadFile(sgfFile)
	if err != nil {
		log.Println(err)
	}
	gt := NewGameTree(SIZE)
	gt.ParseSGF(string(buf))

	path := gt.Path2Root()
	log.Println(len(path))
	board := NewBoard()

	ret := []string{}
	for i := len(path) - 2; i >= 0; i-- {
		cur := path[i]
		sample := board.DLFeature(cur.stone)
		line := strconv.Itoa(PosIndex(cur.x, cur.y))
		line += "\t"
		for _, v := range sample {
			line += strconv.Itoa(int(v))
		}
		ret = append(ret, line)
		if ok := board.Put(PosIndex(cur.x, cur.y), cur.stone); !ok {
			log.Println(cur.x, cur.y, cur.stone)
		}
	}
	log.Println(sgfFile, len(ret), len(path))
	return ret
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
