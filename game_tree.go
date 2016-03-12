package gogo

import (
	"container/list"
	"fmt"
	"strconv"
	"strings"
)

const (
	SGF_LABELS = "abcdefghijklmnopqrstuvwxyz"
)

type GameTreeNode struct {
	Father     *GameTreeNode
	Children   []*GameTreeNode
	x, y       int
	stone      Color
	attr       map[string]string
	visit, win int
	hit        int
}

func (p *GameTreeNode) GetAttr(key string) (string, bool) {
	ret, ok := p.attr[key]
	return ret, ok
}

func (p *GameTreeNode) HasAttr(key string) bool {
	_, ok := p.attr[key]
	return ok
}

func (p *GameTreeNode) Point() Point {
	return Point{p.x, p.y, p.stone}
}

func NewGameTreeNode(stone Color, x, y int) *GameTreeNode {
	return &GameTreeNode{
		Father:   nil,
		Children: []*GameTreeNode{},
		x:        x,
		y:        y,
		stone:    stone,
		attr:     make(map[string]string),
		visit:    0,
		win:      0,
		hit:      0,
	}
}

func NewSGFGameTreeNode(buf string) *GameTreeNode {
	ret := &GameTreeNode{
		Father:   nil,
		Children: []*GameTreeNode{},
		attr:     make(map[string]string),
	}
	buf = strings.TrimLeft(buf, "; \n\t\r")
	left := 0
	tk := ""
	for i := 0; i < len(buf); i++ {
		tk += buf[i : i+1]
		if buf[i] == '[' {
			left += 1
		} else if buf[i] == ']' {
			left -= 1
			if left == 0 {
				key, val := SGFKVToken(tk)
				ret.AddAttr(key, val)
				tk = ""
			}
		}
	}
	tk = strings.TrimSpace(tk)
	if len(tk) > 0 {
		key, val := SGFKVToken(tk)
		ret.AddAttr(key, val)
	}

	if white, ok := ret.attr["W"]; ok {
		ret.x, ret.y = SGFPosition(white)
		ret.stone = WHITE
		delete(ret.attr, "W")
	}

	if black, ok := ret.attr["B"]; ok {
		ret.x, ret.y = SGFPosition(black)
		ret.stone = BLACK
		delete(ret.attr, "B")
	}
	return ret
}

func SGFPosition(buf string) (int, int) {
	if len(buf) != 2 {
		fmt.Println("invalid pos: " + buf)
		return -1, -1
	}
	tx, ty := buf[0:1], buf[1:2]
	return strings.Index(SGF_LABELS, tx), strings.Index(SGF_LABELS, ty)
}

func SGFKVToken(buf string) (string, string) {
	buf = strings.TrimSpace(buf)
	p := strings.Index(buf, "[")
	if p < 0 || p == len(buf)-1 {
		return "", ""
	}
	return strings.ToUpper(buf[0:p]), buf[p+1 : len(buf)-1]
}

func (p *GameTreeNode) AddAttr(k, v string) *GameTreeNode {
	p.attr[k] = v
	if k == "HIT" {
		p.hit, _ = strconv.Atoi(v)
	}
	return p
}

func (p *GameTreeNode) AddChild(v *GameTreeNode) *GameTreeNode {
	if p.Children == nil {
		p.Children = []*GameTreeNode{}
	}
	v.Father = p
	p.Children = append(p.Children, v)
	return p
}

func (p *GameTreeNode) Path2Root() []*GameTreeNode {
	path := []*GameTreeNode{}
	v := p
	for {
		if v == nil {
			break
		}
		path = append(path, v)
		v = v.Father
	}
	return path
}

func (p *GameTreeNode) SGFAttr() string {
	ret := ";"
	if p.stone != GRAY {
		if p.x < 0 || p.x >= len(SGF_LABELS) || p.y < 0 || p.y >= len(SGF_LABELS) {
			return ""
		}
		ret += fmt.Sprintf("%s[%s%s]", SGFColor(p.stone), SGF_LABELS[p.x:p.x+1], SGF_LABELS[p.y:p.y+1])
		if p.hit > 0 {
			ret += fmt.Sprintf("HIT[%d]", p.hit)
		}
	}
	for k, v := range p.attr {
		ret += fmt.Sprintf("%s[%s]", k, v)
	}
	return ret
}

func (p *GameTreeNode) SGF() string {
	ret := p.SGFAttr()
	if ret == "" || p.Children == nil || len(p.Children) == 0 {
		return ret
	} else if len(p.Children) == 1 {
		ret += p.Children[0].SGF()
	} else {
		for i := 0; i < len(p.Children); i++ {
			ret += "(" + p.Children[i].SGF() + ")"
		}
	}
	return ret
}

func (p *GameTreeNode) HasChild(v *GameTreeNode) *GameTreeNode {
	if p.Children == nil {
		return nil
	}
	for _, child := range p.Children {
		fmt.Println(child.Point().String(), child.hit)
		if child.x == v.x && child.y == v.y && child.stone == v.stone {
			return child
		}
	}
	return nil
}

type GameTree struct {
	Root    *GameTreeNode
	Current *GameTreeNode
}

func NewGameTree(size int) *GameTree {
	ret := &GameTree{}
	ret.Root = NewGameTreeNode(GRAY, -1, -1).AddAttr("FF", "4").AddAttr("GM", "1").AddAttr("SZ", strconv.Itoa(size))
	ret.Current = ret.Root
	return ret
}

func (t *GameTree) HasHandicap() bool {
	if t.Root.HasAttr("AB") || t.Root.HasAttr("AW") {
		return true
	}
	return false
}

func (t *GameTree) SGFSize() int {
	if t.Root == nil {
		return 0
	}
	sz, ok := t.Root.GetAttr("SZ")
	if !ok {
		return 19
	}
	ret, _ := strconv.Atoi(sz)
	return ret
}

func (t *GameTree) Add(v *GameTreeNode) {
	if t.Root == nil {
		t.Root = v
		t.Current = t.Root
		return
	}
	t.Current.AddChild(v)
	t.Current = v
}

func (t *GameTree) Combine(ta *GameTree, depth int) {
	pa := ta.Path2Root()
	v := t.Root
	for i := len(pa) - 2; i >= 0 && i >= len(pa)-depth; i-- {
		v.hit += 1
		va := pa[i]
		if i == len(pa)-2 {
			if va.stone != BLACK {
				break
			}
		} else {
			if va.stone != OppColor(v.stone) {
				break
			}
		}
		if child := v.HasChild(va); child != nil {
			v = child
		} else {
			child := NewGameTreeNode(va.stone, va.x, va.y)
			v.AddChild(child)
			v = child
		}
	}
}

func (t *GameTree) NextMoveByQipu(tq *GameTree) []*GameTreeNode {
	path := t.Path2Root()
	vq := tq.Root
	for i := len(path) - 2; i >= 0; i-- {
		v := path[i]
		fmt.Println("==>", v.Point(), len(vq.Children))
		cq := vq.HasChild(v)
		if cq == nil {
			return nil
		}
		vq = cq
	}
	return vq.Children
}

func (t *GameTree) Path2Root() []*GameTreeNode {
	return t.Current.Path2Root()
}

func (t *GameTree) Back() {
	if t.Current.Father != nil {
		t.Current = t.Current.Father
	}
}

func (t *GameTree) WriteSGF() string {
	if t.Root == nil {
		return ""
	}
	return "(" + t.Root.SGF() + ")"
}

func (t *GameTree) ParseSGF(buf string) {
	t.Root = nil
	t.Current = nil
	stack := list.New()
	token := ""
	var node *GameTreeNode
	for {
		token, buf = SGFNextToken(buf)
		if len(token) == 0 {
			break
		}
		if token == ")" {
			for {
				if stack.Len() == 0 {
					break
				}
				v := stack.Back()
				stack.Remove(v)
				str, _ := v.Value.(string)
				if str == "(" {
					break
				}
				t.Back()
			}
		} else if token == "(" {
			stack.PushBack(token)
		} else {
			stack.PushBack(token)
			node = NewSGFGameTreeNode(token)
			t.Add(node)
		}
	}
	t.Current = node
}

func SGFNextToken(buf string) (string, string) {
	//buf = strings.TrimLeft(buf, " \n\t\r")
	buf = strings.TrimSpace(buf)
	if len(buf) == 0 {
		return "", buf
	}
	if buf[0] == '(' || buf[0] == ')' {
		return buf[0:1], buf[1:]
	} else if buf[0] == ';' {
		left := 0
		for i := 1; i < len(buf); i++ {
			if buf[i] == '[' {
				left += 1
			} else if buf[i] == ']' {
				left -= 1
			} else if buf[i] == ';' || buf[i] == ')' || buf[i] == '(' {
				if left == 0 {
					return buf[0:i], buf[i:]
				}
			}
		}
		return "", buf
	} else {
		left := 0
		for i := 0; i < len(buf); i++ {
			if buf[i] == '[' {
				left += 1
			} else if buf[i] == ']' {
				left -= 1
			} else if buf[i] == ';' || buf[i] == ')' || buf[i] == '(' {
				if left == 0 {
					return buf[0:i], buf[i:]
				}
			}
		}
		return "", buf
	}
}

func SGFColor(color Color) string {
	if color == WHITE {
		return "W"
	} else if color == BLACK {
		return "B"
	} else {
		return ""
	}
}
