package gisp

import (
	tm "time"

	px "github.com/Dwarfartisan/goparsec/parsex"
)

// Time 包引入了go的time包功能
var Time = Toolkit{
	Meta: map[string]interface{}{
		"category": "toolkit",
		"name":     "time",
	},
	Content: map[string]interface{}{
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
