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
		return dotSuffix(Dot{x, d.(Atom)})(st)
	}
}

func dotSuffix(x interface{}) p.Parser {
	return func(st p.ParseState) (interface{}, error) {
		d, err := p.Try(DotParser)(st)
		if err != nil {
			return x, nil
		}
		return dotSuffix(Dot{x, d.(Atom)})(st)
	}
}

func BracketSuffix(x interface{}) p.Parser {
	return func(st p.ParseState) (interface{}, error) {
		b, err := p.Try(BracketParser)(st)
		if err != nil {
			return nil, err
		}
		return bracketSuffix(Bracket{x, b.([]interface{})})(st)
	}
}

func bracketSuffix(x interface{}) p.Parser {
	return func(st p.ParseState) (interface{}, error) {
		b, err := BracketParser(st)
		if err != nil {
			return x, nil
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
	suffix := p.Either(p.Try(DotSuffix(prefix)), BracketSuffix(prefix))
	return func(st p.ParseState) (interface{}, error) {
		s, err := suffix(st)
		if err != nil {
			return prefix, nil
		}
		return SuffixParser(s)(st)
	}
}
