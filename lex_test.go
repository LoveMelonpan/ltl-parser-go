package ltlparser

import (
	"fmt"
	"testing"
)

func TestLex(t *testing.T) {
	formula := "Fa"
	res := Lex(formula)

	f1 := "F((p||q)U(q&&r)))"
	fmt.Println(Lex(f1))

	except := []Token{
		{
			Value:    "F",
			ttype:    FUTURE,
			position: 0,
		},
		{
			Value:    "a",
			ttype:    PROP,
			position: 1,
		},
	}
	for k, v := range res {
		if v != except[k] {
			t.Error("should got {},but got {}", except[k], v)
		}
	}
}
