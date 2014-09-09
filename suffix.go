package gisp

import (
	p "github.com/Dwarfartisan/goparsec"
)

func DotSuffixParser(x interface{}) p.Parser {
	return p.Either(p.Try(DotSuffix(x)), p.Return(x))
}

func DotSuffix(x interface{}) p.Parser {
	return func(st p.ParseState) (interface{}, error) {
		d, err := p.Try(DotParser)(st)
		if err != nil {
			return nil, err
		}
		return Dot{x, d.([]Atom)}, nil
	}
}

func BracketSurffix(x interface{}) p.Parser {
	return func(st p.ParseState) (interface{}, error) {
		b, err := p.Try(BracketParser)(st)
		if err != nil {
			return nil, err
		}
		return Bracket{x, b.([]interface{})}, nil
	}
}

func SuffixParser(prefix interface{}) p.Parser {
	surffix := p.Try(p.Either(
		DotSuffix(prefix),
		BracketSurffix(prefix),
	))
	return func(st p.ParseState) (interface{}, error) {
		s, err := surffix(st)
		if err != nil {
			return prefix, nil
		}
		return SuffixParser(s)(st)
	}
}
