package gogo

type SparseMat struct {
	data map[int]map[int]byte
}

func NewSparseMat() *SparseMat {
	return &SparseMat{
		data: make(map[int]map[int]byte),
	}
}

func (p *SparseMat) Add(k1, k2 int, v byte) {
	row, ok := p.data[k1]
	if !ok {
		p.data[k1] = make(map[int]byte)
		row = p.data[k1]
	}
	cell, ok := row[k2]
	if !ok {
		row[k2] = v
	} else {
		row[k2] = v + cell
	}
}

func (p *SparseMat) Get(k1, k2 int) byte {
	if row, ok1 := p.data[k1]; ok1 {
		if v, ok2 := row[k2]; ok2 {
			return v
		}
	}
	return 0
}
