package gogo

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type GTPCommand struct {
	id, name string
	args     []string
}

func NewGTPCommand(line string) *GTPCommand {
	tks := strings.Split(line, " ")
	if len(tks) == 0 {
		return nil
	}
	c := GTPCommand{}
	if _, err := strconv.Atoi(tks[0]); err == nil {
		c.id = tks[0]
		if len(tks) == 1 {
			return nil
		}
		c.name = tks[1]
		c.args = tks[2:]
	} else {
		c.name = tks[0]
		c.args = tks[1:]
	}
	return &c
}

func (p *GTPCommand) getInt(k int) int {
	if k >= len(p.args) {
		return 0
	}
	if n, err := strconv.Atoi(p.args[k]); err == nil {
		return n
	}
	return 0
}

func (p *GTPCommand) GetInt() int {
	return p.getInt(0)
}

func (p *GTPCommand) getFloat(k int) float64 {
	if k >= len(p.args) {
		return 0
	}
	if f, err := strconv.ParseFloat(p.args[k], 64); err == nil {
		return f
	}
	return 0
}

func (p *GTPCommand) GetFloat() float64 {
	return p.getFloat(0)
}

func (p *GTPCommand) Position() (Color, int, int) {
	if len(p.args) == 0 {
		return GRAY, -1, -1
	}
	s := p.args[1]
	if s == "pass" {
		return GRAY, -1, -1
	}
	x := strings.Index(LX, strings.ToUpper(s[0:1]))
	y, _ := strconv.Atoi(s[1:])
	y -= 1
	return p.Color(p.args[0]), x, y
	if p.args[0] == "black" {
		return BLACK, x, y
	} else if p.args[0] == "white" {
		return WHITE, x, y
	} else {
		return GRAY, x, y
	}
}

func (p *GTPCommand) Color(buf string) Color {
	if buf == "black" || buf == "b" {
		return BLACK
	} else if buf == "white" || buf == "w" {
		return WHITE
	} else {
		return GRAY
	}
}

func (p *GTPCommand) ParsePosition(x, y int) string {
	if x < 0 || y < 0 {
		return "pass"
	}
	return LX[x:x+1] + strconv.Itoa(y+1)
}

func (p *GTPCommand) Output(buf string) {
	if buf == "nil" {
		fmt.Printf("?%s ???\n\n", p.id)
	} else {
		fmt.Printf("=%s %s\n\n", p.id, buf)
	}
}

func (g *Game) GTP() {
	r := bufio.NewReader(os.Stdin)
	f, _ := os.Create("/Users/xiangliang/GoCode/src/github.com/xlvector/gogo/gogo.txt")
	defer f.Close()
	cmds := []string{"boardsize", "clear_board", "komi", "play", "genmove", "name", "version", "protocol_version", "quit"}
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.ToLower(strings.TrimRight(line, " \n"))
		if len(line) == 0 {
			continue
		}
		cmd := NewGTPCommand(line)
		if cmd == nil {
			continue
		}

		switch cmd.name {
		case "boardsize":
			g.Init(cmd.GetInt())
			cmd.Output("")
		case "clear_board":
			g.Clear()
			cmd.Output("")
		case "komi":
			g.SetKomi(cmd.GetFloat())
			cmd.Output("")
		case "play":
			stone, x, y := cmd.Position()
			g.Put(stone, x, y)
			cmd.Output("")
		case "genmove":
			x, y := g.GenMove(cmd.Color(cmd.args[0]))
			cmd.Output(cmd.ParsePosition(x, y))
		case "name":
			cmd.Output("gogo")
		case "version":
			cmd.Output("0.1")
		case "protocol_version":
			cmd.Output("2")
		case "list_commands":
			cmd.Output(strings.Join(cmds, "\n"))
		case "known_command":
			for _, c := range cmds {
				if c == cmd.args[1] {
					cmd.Output("true")
					break
				}
			}
			cmd.Output("false")
		case "quit":
			cmd.Output("")
			break
		default:
			fmt.Fprintln(os.Stderr, "unknown command: ", line)
			cmd.Output("")
		}
		f.WriteString(g.String())
	}
}
