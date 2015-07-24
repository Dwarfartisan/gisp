package gisp

import (
	"strconv"

	p "github.com/Dwarfartisan/goparsec"
)

// Float 是 gisp 系统的浮点数实现
type Float float64

// FloatParser 解析浮点数
func FloatParser(st p.ParseState) (interface{}, error) {
	f, err := p.Try(p.Float)(st)
	if err == nil {
		val, err := strconv.ParseFloat(f.(string), 64)
		if err == nil {
			return Float(val), nil
		}
		return nil, err
	}
	return nil, err
}
