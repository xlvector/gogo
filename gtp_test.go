package gogo

import "testing"

func TestGnuGo(t *testing.T) {
	prog := NewGTPProgram("gnugo", "--mode", "gtp", "--level", "10")
	prog.Put(3, 3, BLACK)
	prog.GenMove(WHITE)
}
