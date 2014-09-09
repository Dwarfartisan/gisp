package gisp

import (
	"fmt"
	p "github.com/Dwarfartisan/goparsec"
)

type Bracket struct {
	Root interface{}
	Expr []interface{}
}

func BracketParser(st p.ParseState) (interface{}, error) {
	bracket := p.Between(p.Rune('['), p.Rune(']'),
		p.Choice(
			p.SepBy1(ValueParser, p.Rune(':')),
		))
	t, err := bracket(st)
	if err != nil {
		return nil, err
	}
	tokens := t.([]interface{})
	if len(tokens) > 3 {
		return nil, fmt.Errorf(
			"Bracket expr error:except [key] or [int] or [int:int] or [int:int:int], %v too long",
			tokens,
		)
	}
	return tokens, nil
}
