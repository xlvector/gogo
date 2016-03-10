package gogo

import (
	"fmt"
	"strings"

	"github.com/axgle/mahonia"
	"github.com/saintfish/chardet"
)

func ToUTF8(buf string) string {
	dt := chardet.NewTextDetector()
	ret, err := dt.DetectBest([]byte(buf))
	if err != nil {
		return buf
	}
	fmt.Println(ret.Charset)
	charset := strings.ToLower(ret.Charset)
	if strings.HasPrefix(charset, "gb") {
		enc := mahonia.NewEncoder("gbk")
		return enc.ConvertString(buf)
	}
	return buf
}
