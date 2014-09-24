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
	Content: map[string]Functor{
		"now": SimpleBox{
			ParsexSignChecker(px.Eof),
			func(args ...interface{}) Tasker {
				return func(env Env) (interface{}, error) {
					return tm.Now(), nil
				}
			}},
		"parseDuration": SimpleBox{
			ParsexSignChecker(px.Bind_(StringValue, px.Eof)),
			func(args ...interface{}) Tasker {
				return func(env Env) (interface{}, error) {
					return tm.ParseDuration(args[0].(string))
				}
			}},
		"parseTime": SimpleBox{
			ParsexSignChecker(px.Binds_(StringValue, StringValue, px.Eof)),
			func(args ...interface{}) Tasker {
				return func(env Env) (interface{}, error) {
					return tm.Parse(args[0].(string), args[1].(string))
				}
			}},
	},
}
