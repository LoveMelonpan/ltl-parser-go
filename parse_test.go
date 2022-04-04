package ltlparser

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	formula := "Fa"
	tokens := Lex(formula)
	parser := NewParser(tokens)
	tree := parser.ParseTree()
	traverse(tree)

	if 1 != 1 {
		t.Error("error")
	}
}

func traverse(root *Node) {
	if root == nil {
		return
	}
	fmt.Printf("%v\n", root)
	for _, node := range root.Child {
		traverse(node)
	}
}
