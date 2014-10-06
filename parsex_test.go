package gisp

import (
	"fmt"
	px "github.com/Dwarfartisan/goparsec/parsex"
	"testing"
	"reflect"
)

func TestParsexBasic(t *testing.T) {
	g := NewGispWith(
		map[string]Toolbox{
			"axiom": Axiom, "props": Propositions, "time": Time},
		map[string]Toolbox{"time": Time, "px": Parsex})

	digit := px.Bind(px.Many1(px.Digit), px.ReturnString)
	data := "344932454094325"
	state := NewStringState(data)
	pxre, err := digit(state)
	if err != nil {
		t.Fatalf("except \"%v\" pass test many1 digit but error:%v", data, err)
	}
	code := fmt.Sprintf("(var st (px.state \"%v\"))", data)
	src := code + `
	(var data ((px.many1 px.digit) st))
	(px.s2str data)`
	gre, err := g.Parse(src)
	if err != nil {
		t.Fatalf("except \"%v\" pass gisp many1 digit but error:%v", src, err)
	}
	t.Logf("from gisp: %v", gre)
	t.Logf("from parsex: %v", pxre)
	if !reflect.DeepEqual(pxre, gre){
		t.Fatalf("except got \"%v\" from gisp equal \"%v\" from parsex", gre, pxre)
	}
}
