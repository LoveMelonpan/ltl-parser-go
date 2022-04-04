package ltlparser

import (
	"fmt"
	"log"
)

type Node struct {
	Ntype    string //state
	Operator string // terminal sign
	Child    []*Node
	Val      string //prop's symbol
}

/*
	syntax
	F -> (B,F')
	F' -> (Until,B,F'), (Release,B,F'), nil
	B -> (U,B')
	B' -> (and,U,B'), (or,U,B'), (->,U,B'), (<->,U,B'), nil
	U -> (Future,U), (Global,U),(!,U),(Next,U),(T)
	T -> ((E)),p
*/

type State int

const (
	FState     = iota + 1024 //formula
	FTailState               //formula' tail
	BState                   //binary operator
	BTailState               //binary operator's tail
	UState                   //unary operator
	TState                   //term
	UnKnownState
)

var operate map[int]string = map[int]string{
	FUTURE:  "future",
	ALWAYS:  "global",
	UNTIL:   "until",
	RELEASE: "release",
	NEXT:    "next",
	NEG:     "not",
	AND:     "and",
	OR:      "or",
}

type Stack interface {
	Pop()
	Push()
	Peek()
	Len()
}

type stateStack struct {
	states       []State
	presentState int
}

func (s *stateStack) Push(state State) {
	s.states = append(s.states, state)
	s.presentState++
}

func (s *stateStack) Peek() (State, error) {
	if s.presentState > -1 {
		return s.states[s.presentState], nil
	}
	return UnKnownState, StackEmptyError
}

func (s *stateStack) Len() int {
	return len(s.states)
}
func (s *stateStack) Pop() (State, error) {
	if s.presentState > -1 {
		state := s.states[s.presentState]
		s.presentState--
		s.states = s.states[:len(s.states)-1]
		return state, nil
	}
	return UnKnownState, StackEmptyError
}

func NewStackState() *stateStack {
	return &stateStack{
		states:       []State{},
		presentState: -1,
	}
}

type Parser struct {
	tokens   []Token
	parsePos int //parse position
	states   *stateStack
}

func NewParser(tokens []Token) *Parser {
	states := NewStackState()
	states.Push(FState)
	return &Parser{
		tokens:   tokens,
		parsePos: 0,
		states:   states, //initial state:F
	}
}

/*get current parsing token*/
func (p *Parser) word() *Token {
	if p.parsePos >= len(p.tokens) {
		return nil
	}
	return &p.tokens[p.parsePos]
}

/*got current parsing token and parse_pos+1*/
func (p *Parser) nextWord() *Token {
	if p.parsePos > len(p.tokens)-1 {
		return nil
	}
	p.parsePos++
	return &p.tokens[p.parsePos-1]
}
func (p *Parser) ParseTree() *Node {

	tree := p.Parse()
	return tree
}

func (p *Parser) Parse() *Node {

	defer p.states.Pop()

	presentState, err := p.states.Peek()
	fmt.Printf("state %v\n", p.states)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	switch presentState {
	case FState:
		{

			if w := p.word(); w != nil {
				node := &Node{
					Ntype:    "F",
					Operator: "",
					Child:    nil,
					Val:      "",
				}
				log.Printf("F word %v \n", w.ttype)
				p.states.Push(BState)
				b_node := p.Parse()

				if b_node == nil {
					err := NewMissingExprError(w.Value, "F", p.parsePos)
					log.Fatal(err)
				}

				p.states.Push(FTailState)
				f_tail_node := p.Parse()

				node.Child = append(node.Child, b_node)
				if f_tail_node != nil {
					node.Child = append(node.Child, f_tail_node)
				}
				return node
			} else {
				err := NewMissingExprError("", "F", p.parsePos)
				log.Fatal(err)
			}
		}

	case FTailState:
		{

			if w := p.word(); w != nil {
				node := &Node{
					Ntype:    "FTail",
					Operator: "",
					Child:    nil,
					Val:      "",
				}
				log.Printf("F' word %v", w)
				if w.ttype == UNTIL || w.ttype == RELEASE {
					p.nextWord()
					p.states.Push(BState)
					b_node := p.Parse()

					p.states.Push(FTailState)
					f_tail_node := p.Parse()

					node.Operator = operate[w.ttype]
					node.Child = append(node.Child, b_node)
					if f_tail_node != nil {
						node.Child = append(node.Child, f_tail_node)
					}
					return node
				} else {
					return nil
				}

			} else {
				return nil
			}
		}

	case BState:
		{

			if w := p.word(); w != nil {
				node := &Node{
					Ntype:    "B",
					Operator: "",
					Child:    nil,
					Val:      "",
				}
				log.Printf("B word %v", w)
				p.states.Push(UState)
				u_node := p.Parse()

				if u_node == nil {
					err := NewMissingExprError(w.Value, "B", p.parsePos)
					log.Fatal(err)
				}

				log.Printf("B get node U %+v\n", u_node)

				p.states.Push(BTailState)
				b_tail_node := p.Parse()
				log.Printf("B get node B' %+v\n", b_tail_node)

				node.Child = append(node.Child, u_node)
				if b_tail_node != nil {
					node.Child = append(node.Child, b_tail_node)
				}
				log.Printf("B return node U %+v", node)
				return node
			} else {
				err := NewMissingExprError("", "B", p.parsePos)
				log.Fatal(err)
			}
		}

	case BTailState:
		{
			if w := p.word(); w != nil {
				node := &Node{
					Ntype:    "BTail",
					Operator: "",
					Child:    nil,
					Val:      "",
				}
				log.Printf("B_tail word %v\n", w)
				if w.ttype == AND || w.ttype == OR || w.ttype == IMPILCIT || w.ttype == EQ {
					p.nextWord()

					p.states.Push(UState)
					u_node := p.Parse()

					p.states.Push(BTailState)
					b_tail_node := p.Parse()

					node.Operator = operate[w.ttype]
					node.Child = append(node.Child, u_node)
					if b_tail_node != nil {
						node.Child = append(node.Child, b_tail_node)
					}
					return node

				} else {
					return nil
				}
			} else {
				return nil
			}
		}
	case UState:
		{

			if w := p.word(); w != nil {
				node := &Node{
					Ntype:    "U",
					Operator: "",
					Child:    nil,
					Val:      "",
				}
				log.Printf("U word %v\n", w)
				if w.ttype == FUTURE || w.ttype == ALWAYS || w.ttype == NEXT || w.ttype == NEG {
					p.nextWord()
					p.states.Push(UState)
					u_node := p.Parse()
					if u_node == nil {
						err := NewMissingExprError(w.Value, "U", p.parsePos)
						log.Fatal(err)
					}

					node.Operator = operate[w.ttype]
					node.Child = append(node.Child, u_node)
					log.Printf("U return node %+v\n", node)
					return node
				} else {
					p.states.Push(TState)
					node = p.Parse()
					log.Printf("U return node %+v\n", node)
					return node
				}
			} else {
				err := NewMissingExprError("", "U", p.parsePos)
				log.Fatal(err)
			}
		}
	case TState:
		{

			if w := p.word(); w != nil {
				node := &Node{
					Ntype:    "T",
					Operator: "",
					Child:    nil,
					Val:      "",
				}
				log.Printf("T word %+v", w)
				if w.ttype == LEFT {
					p.nextWord()
					/*parse E*/
					p.states.Push(FState)
					e_node := p.Parse()
					if e_node == nil {
						err := NewMissingExprError("(", "T", p.parsePos-1)
						log.Fatal(err)
					}

					/*check )*/
					v := p.nextWord()
					if v == nil {
						err := NewSyntaxError(")", v.Value)
						log.Fatal(err)
					}
					if v.ttype == RIGHT {
						node.Child = append(node.Child, e_node)
						node.Operator = "()"
						return node
					}
					err := NewSyntaxError(")", v.Value)
					log.Fatal(err)
				} else if w.ttype == PROP {
					v := p.nextWord()

					if v == nil {
						err := NewMissingExprError("", "T", p.parsePos)
						log.Fatal(err)
					}
					node.Val = v.Value
					fmt.Printf("T return node %+v\n", node)
					return node
				} else {
					err := NewSyntaxError("( or lower case", w.Value)
					log.Fatal(err)
				}
			} else {
				err := NewMissingExprError("", "T", p.parsePos)
				log.Fatal(err)
			}
		}
	}
	return nil
}
