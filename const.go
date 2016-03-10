package gogo

type Color byte

const (
	GRAY    Color = 0
	BLACK   Color = 1
	WHITE   Color = 2
	UNKNOWN Color = 100
	LX            = "ABCDEFGHJKLMNOPQRSTUVWXYZ"
)

func OppColor(c Color) Color {
	if c == BLACK {
		return WHITE
	} else if c == WHITE {
		return BLACK
	}
	return c
}
