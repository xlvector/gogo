package gogo

import (
	"container/list"
	"os"
	"path/filepath"
	"strings"
)

func TreeDir(root, ext string) []string {
	q := list.New()
	q.PushBack(root)
	ret := []string{}
	for q.Len() > 0 {
		p := q.Front()
		q.Remove(p)
		v := p.Value.(string)
		filepath.Walk(v, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() && path != root {
				q.PushBack(path)
			} else {
				if strings.HasSuffix(path, ext) {
					ret = append(ret, path)
				}
			}
			return nil
		})
	}
	return ret
}
