package gogo

import (
	"container/list"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func TreeDir(root, ext string) []string {
	q := list.New()
	q.PushBack(root)
	ret := []string{}
	visit := make(map[string]byte)
	for q.Len() > 0 {
		p := q.Front()
		q.Remove(p)
		v := p.Value.(string)
		if _, ok := visit[v]; ok {
			continue
		}
		visit[v] = 1
		filepath.Walk(v, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() && path != root {
				q.PushBack(path)
			} else {
				if strings.HasSuffix(path, ext) {
					log.Println(path)
					ret = append(ret, path)
				}
			}
			return nil
		})
	}
	return ret
}

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
	p[j], p[i] = p[i], p[j]
}
