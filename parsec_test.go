package gisp

import (
    //"fmt"
    p "github.com/Dwarfartisan/goparsec"
    "testing"
    "reflect"
)

func TestParsecBasic(t *testing.T) {
    g := NewGispWith(
        map[string]Toolbox{
            "axiom": Axiom, "props": Propositions, "time": Time},
        map[string]Toolbox{"time": Time, "p": Parsec})

    digit := p.Bind(p.Many1(p.Digit), p.ReturnString)
    data := "344932454094325"
    state := p.MemoryParseState(data)
    pre, err := digit(state)
    if err != nil {
        t.Fatalf("except \"%v\" pass test many1 digit but error:%v", data, err)
    }

    src := "(let ((st (p.state \"" + data + `")))
    (var data ((p.many1 p.digit) st))
    (p.s2str data))
    `
    gre, err := g.Parse(src)
    if err != nil {
        t.Fatalf("except \"%v\" pass gisp many1 digit but error:%v", src, err)
    }
    t.Logf("from gisp: %v", gre)
    t.Logf("from parsec: %v", pre)
    if !reflect.DeepEqual(pre, gre){
        t.Fatalf("except got \"%v\" from gisp equal \"%v\" from parsec", gre, pre)
    }
}
