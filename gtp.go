package gogo

import (
	"bufio"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

type GTPProgram struct {
	cmd    *exec.Cmd
	reader *bufio.Reader
	writer *bufio.Writer
}

func NewGTPProgram(name string, args ...string) *GTPProgram {
	ret := &GTPProgram{
		cmd: exec.Command(name, args...),
	}

	out, err := ret.cmd.StdoutPipe()
	if err != nil {
		log.Fatalln(err)
	}
	ret.reader = bufio.NewReader(out)

	in, err := ret.cmd.StdinPipe()
	if err != nil {
		log.Fatalln(err)
	}
	ret.writer = bufio.NewWriter(in)
	if err := ret.cmd.Start(); err != nil {
		log.Fatalln(err)
	}
	return ret
}

func (p *GTPProgram) Put(x, y int, stone Color) {
	buf := "play"
	if stone == BLACK {
		buf += " black "
	} else if stone == WHITE {
		buf += " white "
	} else {
		return
	}
	buf += string(LX[x])
	buf += strconv.Itoa(y + 1)
	buf += "\n"
	log.Print(buf)
	p.writer.WriteString(buf)
	p.writer.Flush()
	line, _ := p.reader.ReadString('\n')
	p.reader.ReadString('\n')
	log.Println(line)
}

func (p *GTPProgram) GenMove(stone Color) string {
	buf := "genmove "
	if stone == BLACK {
		buf += "black"
	} else if stone == WHITE {
		buf += "white"
	} else {
		return ""
	}
	buf += "\n"
	log.Print(buf)
	p.writer.WriteString(buf)
	p.writer.Flush()
	line, _ := p.reader.ReadString('\n')
	p.reader.ReadString('\n')
	line = strings.TrimLeft(line, " = ")
	line = strings.TrimSpace(line)
	log.Println(line)
	return line
}
