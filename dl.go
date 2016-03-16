package gogo

import (
	"io/ioutil"
	"strconv"
)

func GenDLDataset(sgfFile string) []string {
	buf, _ := ioutil.ReadFile(sgfFile)
	gt := NewGameTree(SIZE)
	gt.ParseSGF(string(buf))

	path := gt.Path2Root()
	board := NewBoard()

	ret := make([]string, 0, 300)
	for i := len(path) - 2; i >= 0; i-- {
		cur := path[i]
		sample := board.DLFeature(cur.stone)
		line := strconv.Itoa(PosIndex(cur.x, cur.y))
		for _, v := range sample {
			line += "\t"
			line += strconv.Itoa(int(v))
		}
		ret = append(ret, line)
	}
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
