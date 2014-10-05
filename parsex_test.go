package gisp

import (
	"fmt"
	px "github.com/Dwarfartisan/goparsec/parsex"
	"testing"
)

func TestParsexBasic(t *testing.T) {
	g := NewGispWith(
		map[string]Toolbox{
			"axiom": Axiom, "props": Propositions, "time": Time},
		map[string]Toolbox{"time": Time, "px": Parsex})

	digit := px.Many1(px.Digit)
	data := "344932454094325"
	state := NewStringState(data)
	pxre, err := digit(state)
	if err != nil {
		t.Fatalf("except \"%v\" pass test many1 digit but error:%v", data, err)
	}
	code := fmt.Sprintf("(var st (px.state \"%v\"))", data)
	_, err = g.Parse(code)
	if err != nil {
		t.Fatalf("except \"%v\" create a state but error:%v", code, err)
	}
	src := "((px.many1 px.digit) st)"
	gre, err := g.Parse(src)
	if err != nil {
		t.Fatalf("except \"%v\" pass gisp many1 digit but error:%v", src, err)
	}
	t.Log(pxre)
	t.Log(gre)
}
