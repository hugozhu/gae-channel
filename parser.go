package channel

import (
	"strconv"
	"strings"
)

type Message struct {
	Key   string
	Val   *[]Message
	Child *Message
}

func (e *Message) Id() int {
	id, _ := strconv.Atoi(e.Key)
	return id
}

func (e *Message) HasChild() bool {
	return e.Child != nil
}

func (e *Message) Level() int {
	level := 1
	a := e.Child
	for {
		if a == nil {
			break
		}
		level++
		a = a.Child
	}
	return level
}

func (e *Message) ToString() string {
	s := ""
	if e.Key != "" {
		if e.Child == nil && e.Val == nil {
			s = "[" + e.Key
		} else {
			s = "[\"" + e.Key + "\", "
		}

	}
	if e.Child != nil {
		s = s + e.Child.ToString()
	}
	if e.Val != nil {
		for _, ee := range *e.Val {
			s = s + ee.ToString()
		}
	}
	s = s + "]"
	return s
}

type Parser struct {
	parent *Message
	arr    *[]Message
}

func (p *Parser) Root() *Message {
	return p.parent
}

func (p *Parser) Parse(bytes []byte) *Message {
	p.parse(bytes, 0)
	return p.Root()
}

func (p *Parser) parse(bytes []byte, pos1 int) int {
	if p.arr == nil {
		p.arr = &[]Message{}
	}
	key := ""
	for i := pos1; i < len(bytes); i++ {
		b := bytes[i]
		switch b {
		case '[':
			if bytes[i+1] == ']' {
				i = i + 1
			} else {
				key = strings.Trim(string(bytes[pos1:i]), "\"',\n")
				pos1 = i
				i = p.parse(bytes, i+1)
			}
		case ']':
			e := Message{
				Key: string(bytes[pos1:i]),
			}

			if key != "" {
				if key[0:1] == "[" { //ignore
					return i
				}
				p.parent = &Message{
					Key:   key,
					Val:   p.arr,
					Child: p.parent,
				}
				p.arr = &[]Message{}
			} else {
				tmp := append(*p.arr, e)
				p.arr = &tmp
			}
			return i
		}
	}
	return len(bytes)
}
