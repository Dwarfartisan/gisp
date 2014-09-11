package gisp

import (
	t "time"
)

var time = Toolkit{
	Meta: map[string]interface{}{
		"category": "toolkit",
		"name":     "time",
	},
	Content: map[string]Expr{
		"now": func(env Env) Element {
			return func(args ...interface{}) (interface{}, error) {
				return t.Now(), nil
			}
		},
		// "Date": func(env Env) Element {
		// 	return func(args ...interface{}) (interface{}, error) {
		//
		// 	}
		// },
	},
}
