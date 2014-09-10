package gisp

import (
	p "github.com/Dwarfartisan/goparsec"
)

func DotSuffix(x interface{}) p.Parser {
	return func(st p.ParseState) (interface{}, error) {
		d, err := DotParser(st)
		if err != nil {
			return nil, err
		}
		return Dot{x, d.([]Atom)}, nil
	}
}

func BracketSuffix(x interface{}) p.Parser {
	return func(st p.ParseState) (interface{}, error) {
		b, err := BracketParser(st)
		if err != nil {
			return nil, err
		}
		return Bracket{x, b.([]interface{})}, nil
	}
}

func DotSuffixParser(x interface{}) p.Parser {
	return p.Either(p.Try(DotSuffix(x)), p.Return(x))
}

func BracketSuffixParser(x interface{}) p.Parser {
	return p.Either(p.Try(BracketSuffix(x)), p.Return(x))
}

func SuffixParser(prefix interface{}) p.Parser {
	suffix := p.Either(DotSuffix(prefix), BracketSuffix(prefix))
	return func(st p.ParseState) (interface{}, error) {
		s, err := suffix(st)
		if err != nil {
			return prefix, nil
		}
		return SuffixParser(s)(st)
	}
}
