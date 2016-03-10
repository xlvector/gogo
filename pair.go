package gogo

type IntFloatPair struct {
	First  int
	Second float64
}

type IntFloatPairList []IntFloatPair

func (p IntFloatPairList) Len() int {
	return len(p)
}

func (p IntFloatPairList) Less(i, j int) bool {
	return p[i].Second < p[j].Second
}

func (p IntFloatPairList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

type StringIntPair struct {
	First  string
	Second int
}

type StringIntPairList []StringIntPair

func (p StringIntPairList) Len() int {
	return len(p)
}

func (p StringIntPairList) Less(i, j int) bool {
	return p[i].Second < p[j].Second
}

func (p StringIntPairList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
