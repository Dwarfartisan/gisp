package gisp

import (
	px "github.com/Dwarfartisan/goparsec/parsex"
	tm "time"
)

var Time = Toolkit{
	Meta: map[string]interface{}{
		"category": "toolkit",
		"name":     "time",
	},
	Content: map[string]Expr{
		"now": func(env Env) Element {
			return func(args ...interface{}) (interface{}, error) {
				_, err := GetArgs(env, px.Eof, args)
				if err != nil {
					return nil, err
				}
				return tm.Now(), nil
			}
		},
		"parseDuration": func(env Env) Element {
			return func(args ...interface{}) (interface{}, error) {
				params, err := GetArgs(env, px.Binds_(px.StringVal, px.Eof), args)
				if err != nil {
					return nil, err
				}
				return tm.ParseDuration(params[0].(string))
			}
		},
		"parseTime": func(env Env) Element {
			return func(args ...interface{}) (interface{}, error) {
				params, err := GetArgs(env, px.Binds_(px.StringVal, px.StringVal,
					px.Eof),
					args)
				if err != nil {
					return nil, err
				}
				return tm.Parse(params[0].(string), params[1].(string))
			}
		},
	},
}
