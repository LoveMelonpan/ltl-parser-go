package ltlparser

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	formula := "Fa"
	tokens := Tex(formula)
	ast := Parse(tokens)
	traverse(ast)
	if formula != "Fa" {
		t.Error("error")
	}
}

func traverse(root *Node) {
	if root == nil {
		return
	}
	fmt.Println(root)
	for _, v := range root.Child {
		traverse(v)
	}
}
