package ltlparser

import (
	"fmt"
	"testing"
)

func TestTex(t *testing.T) {
	formula := "Fa"
	res := Tex(formula)

	f1 := "F((p||q)U(q&&r)))"
	fmt.Println(Tex(f1))

	except := []Token{
		{
			Value:    "F",
			T_type:   FUTURE,
			position: 0,
		},
		{
			Value:    "a",
			T_type:   PROP,
			position: 1,
		},
	}
	for k, v := range res {
		if v != except[k] {
			t.Error("should got {},but got {}", except[k], v)
		}
	}
}
