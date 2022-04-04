package ltlparser

import (
	"errors"
	"fmt"
	"runtime"
	"strconv"
)

/*parse erros*/
//get unexcept syntax
type SyntaxError struct {
	exceptSyntax string
	gotSyntax    string
}

func (e *SyntaxError) Error() string {
	_, srcName, line, _ := runtime.Caller(1)
	stackMsg := "[" + srcName + ":" + strconv.Itoa(line) + "] "
	err := fmt.Sprintf("Syntax Error :except %s but got %s", e.exceptSyntax, e.gotSyntax)
	return stackMsg + err
}

func NewSyntaxError(exceptSyntax string, gotSyntax string) SyntaxError {
	return SyntaxError{
		exceptSyntax: exceptSyntax,
		gotSyntax:    gotSyntax,
	}
}

//missing except expression during parse
type MissingExprError struct {
	beforeSyntax string
	nodeType     string
	pos          int
}

func (e *MissingExprError) Error() string {
	_, srcName, line, _ := runtime.Caller(1)
	stackMsg := "[" + srcName + ":" + strconv.Itoa(line) + "] "
	err := fmt.Sprintf("Missing Expr Error: missing expr after %s at %d in %s", e.beforeSyntax, e.pos, e.nodeType)
	return stackMsg + err
}

func NewMissingExprError(beforeSyntax string, nodeType string, pos int) MissingExprError {
	return MissingExprError{
		beforeSyntax: beforeSyntax,
		nodeType:     nodeType,
		pos:          pos,
	}
}

/*stack errors*/
var (
	StackEmptyError = errors.New("Stack is empty")
)
