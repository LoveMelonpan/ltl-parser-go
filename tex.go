package ltlparser

import (
	"fmt"
	"strings"
	"unicode"
)

const (
	PROP = 257 + iota
	ALWAYS
	FUTURE
	UNTIL
	RELEASE
	AND
	OR
	IMPILCIT
	EQ
	NEXT
	NEG
	LEFT
	RIGHT
)

type Token struct {
	Value    string
	T_type   int
	position int
}

func Tex(formula string) []Token {
	tokens := make([]Token, 0)
	formula = strings.TrimSpace(formula)
	rformula := []rune(formula)
	i := 0
	for {
		if i >= len(rformula) {
			return tokens
		}

		s := rformula[i]

		if unicode.IsLower(s) {
			start := i
			var text string = string(s)
			i = i + 1
			var s rune
			for i < len(rformula) && (unicode.IsLower(s) || unicode.IsDigit(s) || s == '_') {
				s = rformula[i]
				text += string(s)
				i++
			}
			token := Token{Value: text, T_type: PROP, position: start}
			tokens = append(tokens, token)
			continue
		} else if s == '|' {
			if i >= len(rformula)-1 {
				err := fmt.Errorf("except | in forluma %s's end", formula)
				panic(err)
			}

			i++
			if rformula[i] == '|' {
				token := Token{Value: "||", T_type: OR, position: i - 1}
				tokens = append(tokens, token)
				i++
				continue
			} else {
				err := fmt.Errorf("except | but got %s in formula %s", string(rformula[i]), formula)
				panic(err)

			}
		} else if s == '&' {
			if i >= len(rformula)-1 {
				err := fmt.Errorf("except & in forluma %s's end", formula)
				panic(err)
			}
			i++

			if rformula[i] == '&' {
				token := Token{Value: "&&", T_type: AND, position: i - 1}
				tokens = append(tokens, token)
				i++
				continue
			} else {
				err := fmt.Errorf("except & but got %s in formula %s", string(rformula[i]), formula)
				panic(err)

			}
		} else if s == ' ' || s == '\r' || s == '\t' {
			i++
			continue
		} else if s == '(' {
			token := Token{Value: string(s), T_type: LEFT, position: i}
			tokens = append(tokens, token)
			i++
			continue
		} else if s == ')' {
			token := Token{Value: string(s), T_type: RIGHT, position: i}
			tokens = append(tokens, token)
			i++
			continue
		} else if s == 'G' {
			token := Token{Value: string(s), T_type: ALWAYS, position: i}
			tokens = append(tokens, token)
			i++
			continue
		} else if s == 'F' {
			token := Token{Value: string(s), T_type: FUTURE, position: i}
			tokens = append(tokens, token)
			i++
			continue
		} else if s == 'U' {
			token := Token{Value: string(s), T_type: UNTIL, position: i}
			tokens = append(tokens, token)
			i++
			continue
		} else if s == 'R' {
			token := Token{Value: string(s), T_type: RELEASE, position: i}
			tokens = append(tokens, token)
			i++
			continue
		} else if s == 'X' {
			token := Token{Value: string(s), T_type: NEXT, position: i}
			tokens = append(tokens, token)
			i++
			continue
		} else if i == '(' || i == ')' {
			i++
			continue
		} else if s == '!' {
			token := Token{Value: string(s), T_type: NEG, position: i}
			tokens = append(tokens, token)
			i++
			continue
		} else if s == '[' {
			if i >= len(rformula)-1 {
				err := fmt.Errorf("except ] in forluma %s's end", formula)
				panic(err)
			}

			if rformula[i+1] == ']' {
				token := Token{Value: "[]", T_type: ALWAYS, position: i}
				tokens = append(tokens, token)
				i += 2
				continue
			} else {
				err := fmt.Errorf("except ] but got %s in formula %s", string(rformula[i+1]), formula)
				panic(err)
			}
		} else if s == '-' {
			if i >= len(rformula)-1 {
				err := fmt.Errorf("except > in forluma %s's end", formula)
				panic(err)
			}

			if rformula[i+1] == '>' {
				token := Token{Value: "->", T_type: IMPILCIT, position: i}
				tokens = append(tokens, token)
				i += 2
				continue
			} else {
				err := fmt.Errorf("except > but got %s in formula %s", string(rformula[i+1]), formula)
				panic(err)
			}
		} else if s == '<' {
			if i >= len(rformula)-1 {
				err := fmt.Errorf("except > in forluma %s's end", formula)
				panic(err)
			}

			if rformula[i+1] == '>' {
				token := Token{Value: "->", T_type: IMPILCIT, position: i}
				tokens = append(tokens, token)
				i += 2
				continue
			} else if i >= len(rformula)-2 {
				err := fmt.Errorf("except -> in forluma %s's end", formula)
				panic(err)
			} else if rformula[i+1] == '-' && rformula[i+2] == '>' {
				token := Token{Value: "<->", T_type: EQ, position: i}
				tokens = append(tokens, token)
				i += 3
				continue
			} else {
				err := fmt.Errorf("except -> but got %s in formula %s", string(rformula[i+1]), formula)
				panic(err)
			}
		} else {
			err := fmt.Errorf("unexcepted token %v", rformula[i])
			panic(err)
		}

	}
}
