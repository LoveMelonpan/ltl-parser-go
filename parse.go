package ltlparser

import (
	"fmt"
)

type Node struct {
	N_type string
	Value  []Token
	Child  []*Node
}

/*
	syntax
	E -> (T,E')
	E' -> (U,T,E'), (R,T,E'), nil
	T -> (M,T')
	T' -> (and,M,T'), (or,M,T'), (->,M,T'), (<->,M,T'), nil
	M -> (F,M), (G,M),(!,M'),(X,M),(B)
	B -> ((E)),p
*/
func Parse(tokens []Token) *Node {
	i := 0
	var E func() *Node
	var E_tail func() *Node
	var T func() *Node
	var T_tail func() *Node
	var M func() *Node
	var B func() *Node

	//next(word)
	next_word := func() *Token {
		i = i + 1
		if i >= len(tokens) {
			return nil
		}
		return &(tokens[i])

	}
	//get current word
	get_word := func() *Token {
		return &tokens[i]
	}

	/*B -> ((E)),p*/
	B = func() *Node {
		v := next_word()

		node := &Node{
			N_type: "B",
			Value:  nil,
			Child:  nil,
		}

		if v.T_type == LEFT {
			node.Value = []Token{*v}

			v = next_word()
			left := &Node{N_type: "(", Value: node.Value}

			var expr *Node
			if expr = E(); expr != nil {
				vals := expr.Value
				node.Value = append(node.Value, vals...)
			} else {
				err := fmt.Errorf("syntax error %v", v.Value)
				panic(err)
			}

			v = next_word()
			if v.T_type == RIGHT {
				node.Value = append(node.Value, *v)
				right_val := []Token{*v}
				right := &Node{N_type: ")", Value: right_val}
				node.Child = []*Node{left, expr, right}
				return node
			} else {
				err := fmt.Errorf("except ] but got %v", v.Value)
				panic(err)
			}

		}

		if v.T_type == PROP {
			node.Value = []Token{*v}
			return node
		}

		return nil
	}
	/*M -> (F,M), (G,M),(!,M'),(X,M),(B)*/
	M = func() *Node {
		v := next_word()

		node := &Node{
			N_type: "M",
			Value:  nil,
			Child:  nil,
		}

		switch v.T_type {
		case ALWAYS, FUTURE, NEXT, NEG:
			node.Value = []Token{*v}
			v = next_word()
			var expr *Node
			if expr = M(); expr != nil {
				node.Child = []*Node{expr}
				return node
			}
			err := fmt.Errorf("missing expression after %v in M", v.Value)
			panic(err)

		default:
			return B()
		}
		return nil
	}

	/*T -> (M,T')*/
	T = func() *Node {
		v := get_word()
		node := &Node{
			N_type: "T",
			Value:  nil,
			Child:  nil,
		}
		var m, t_tail *Node
		if m = M(); m != nil {
			node.Value = append(node.Value, m.Value...)
			node.Child = append(node.Child, m)

			next_word()
			t_tail = T_tail()
			if t_tail != nil {
				node.Value = append(node.Value, t_tail.Value...)
				node.Child = append(node.Child, t_tail)
			}
		} else {
			err := fmt.Errorf("missing expression after %v in T", v.Value)
			panic(err)
		}
		return node
	}

	/*T' -> (and,M,T'), (or,M,T'), (->,M,T'), (<->,M,T'), nil*/
	T_tail = func() *Node {
		v := next_word()
		node := &Node{N_type: "T_tail", Value: nil, Child: nil}

		if v.T_type == AND || v.T_type == OR || v.T_type == EQ || v.T_type == IMPILCIT {
			var first_node, second_node *Node
			if first_node = M(); first_node != nil {
				node.Value = append(node.Value, first_node.Value...)
				node.Child = append(node.Child, first_node)

				next_word()
				second_node = T_tail()
				node.Value = append(node.Value, second_node.Value...)
				node.Child = append(node.Child, second_node)
				return node

			}
		}
		return nil
	}

	/*E -> (T,E')*/
	E = func() *Node {
		v := get_word()
		node := &Node{
			N_type: "E",
			Value:  nil,
			Child:  nil,
		}
		var t, e_tail *Node
		if t = T(); t != nil {
			node.Value = append(node.Value, t.Value...)
			node.Child = append(node.Child, t)
			if e_tail = E_tail(); e_tail != nil {
				node.Value = append(node.Value, e_tail.Value...)
				node.Child = append(node.Child, e_tail)
			}
		} else {
			err := fmt.Errorf("missing expression afrer %v in T_tail", v.Value)
			panic(err)
		}
		return node
	}

	/*	E' -> (U,T,E'), (R,T,E'), nil*/
	E_tail = func() *Node {
		v := next_word()
		node := &Node{
			N_type: "E_tail",
			Value:  nil,
			Child:  nil,
		}
		if v.T_type == UNTIL || v.T_type == RELEASE {
			var t, e_tail *Node
			if t = T(); t != nil {
				node.Value = append(node.Value, t.Value...)
				node.Child = append(node.Child, t)
				e_tail = E_tail()
				if e_tail != nil {
					node.Value = append(node.Value, e_tail.Value...)
					node.Child = append(node.Child, e_tail)
				}
			}
		}
		return nil
	}

	return E()
}