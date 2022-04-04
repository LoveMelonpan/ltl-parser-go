package ltlparser

import (
	"fmt"
	"log"
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
	ttype    int
	position int
}

func Lex(formula string) []Token {
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
			token := Token{Value: text, ttype: PROP, position: start}
			tokens = append(tokens, token)
			continue
		} else if s == '|' {
			if i >= len(rformula)-1 {
				err := fmt.Errorf("except | in forluma %s's end", formula)
				log.Fatal(err)
			}

			i++
			if rformula[i] == '|' {
				token := Token{Value: "||", ttype: OR, position: i - 1}
				tokens = append(tokens, token)
				i++
				continue
			} else {
				err := fmt.Errorf("except | but got %s in formula %s", string(rformula[i]), formula)
				log.Fatal(err)

			}
		} else if s == '&' {
			if i >= len(rformula)-1 {
				err := fmt.Errorf("except & in forluma %s's end", formula)
				log.Fatal(err)
			}
			i++

			if rformula[i] == '&' {
				token := Token{Value: "&&", ttype: AND, position: i - 1}
				tokens = append(tokens, token)
				i++
				continue
			} else {
				err := fmt.Errorf("except & but got %s in formula %s", string(rformula[i]), formula)
				log.Fatal(err)

			}
		} else if s == ' ' || s == '\r' || s == '\t' {
			i++
			continue
		} else if s == '(' {
			token := Token{Value: string(s), ttype: LEFT, position: i}
			tokens = append(tokens, token)
			i++
			continue
		} else if s == ')' {
			token := Token{Value: string(s), ttype: RIGHT, position: i}
			tokens = append(tokens, token)
			i++
			continue
		} else if s == 'G' {
			token := Token{Value: string(s), ttype: ALWAYS, position: i}
			tokens = append(tokens, token)
			i++
			continue
		} else if s == 'F' {
			token := Token{Value: string(s), ttype: FUTURE, position: i}
			tokens = append(tokens, token)
			i++
			continue
		} else if s == 'U' {
			token := Token{Value: string(s), ttype: UNTIL, position: i}
			tokens = append(tokens, token)
			i++
			continue
		} else if s == 'R' {
			token := Token{Value: string(s), ttype: RELEASE, position: i}
			tokens = append(tokens, token)
			i++
			continue
		} else if s == 'X' {
			token := Token{Value: string(s), ttype: NEXT, position: i}
			tokens = append(tokens, token)
			i++
			continue
		} else if i == '(' || i == ')' {
			i++
			continue
		} else if s == '!' {
			token := Token{Value: string(s), ttype: NEG, position: i}
			tokens = append(tokens, token)
			i++
			continue
		} else if s == '[' {
			if i >= len(rformula)-1 {
				err := fmt.Errorf("except ] in forluma %s's end", formula)
				log.Fatal(err)
			}

			if rformula[i+1] == ']' {
				token := Token{Value: "[]", ttype: ALWAYS, position: i}
				tokens = append(tokens, token)
				i += 2
				continue
			} else {
				err := fmt.Errorf("except ] but got %s in formula %s", string(rformula[i+1]), formula)
				log.Fatal(err)
			}
		} else if s == '-' {
			if i >= len(rformula)-1 {
				err := fmt.Errorf("except > in forluma %s's end", formula)
				log.Fatal(err)
			}

			if rformula[i+1] == '>' {
				token := Token{Value: "->", ttype: IMPILCIT, position: i}
				tokens = append(tokens, token)
				i += 2
				continue
			} else {
				err := fmt.Errorf("except > but got %s in formula %s", string(rformula[i+1]), formula)
				log.Fatal(err)
			}
		} else if s == '<' {
			if i >= len(rformula)-1 {
				err := fmt.Errorf("except > in forluma %s's end", formula)
				log.Fatal(err)
			}

			if rformula[i+1] == '>' {
				token := Token{Value: "->", ttype: IMPILCIT, position: i}
				tokens = append(tokens, token)
				i += 2
				continue
			} else if i >= len(rformula)-2 {
				err := fmt.Errorf("except -> in forluma %s's end", formula)
				log.Fatal(err)
			} else if rformula[i+1] == '-' && rformula[i+2] == '>' {
				token := Token{Value: "<->", ttype: EQ, position: i}
				tokens = append(tokens, token)
				i += 3
				continue
			} else {
				err := fmt.Errorf("except -> but got %s in formula %s", string(rformula[i+1]), formula)
				log.Fatal(err)
			}
		} else {
			err := fmt.Errorf("unexcepted token %v", rformula[i])
			log.Fatal(err)
		}

	}
}
