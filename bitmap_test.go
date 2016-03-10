package gogo

import (
	"testing"
)

func TestBitmap(t *testing.T) {
	bm := NewBoardBitmap()
	bm.Set(17)
	if !bm.IsSet(17) {
		t.Error()
	}
}
