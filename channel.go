package gisp

import (
	px "github.com/Dwarfartisan/goparsec/parsex"
	"reflect"
)

type Chan struct {
	Type  reflect.Type
	Dir   reflect.ChanDir
	value reflect.Value
}

func MakeChan(typ reflect.Type, dir reflect.ChanDir, buf Int) *Chan {
	return &Chan{typ, dir, reflect.MakeChan(typ, int(buf))}
}

func MakeRecvChan(typ reflect.Type, buf Int) *Chan {
	return MakeChan(typ, reflect.RecvDir, buf)
}

func MakeSendChan(typ reflect.Type, buf Int) *Chan {
	return MakeChan(typ, reflect.SendDir, buf)
}

func MakeBothChan(typ reflect.Type, buf Int) *Chan {
	return MakeChan(typ, reflect.BothDir, buf)
}

func (ch *Chan) Send(x interface{}) {
	ch.value.Send(reflect.ValueOf(x))
}

func (ch *Chan) Recv() (x interface{}, ok bool) {
	val, ok := ch.value.Recv()
	if val.IsValid() {
		return val.Interface(), ok
	} else {
		return nil, ok
	}
}

func (ch *Chan) TrySend(x interface{}) {
	ch.value.TrySend(reflect.ValueOf(x))
}

func (ch *Chan) TryRecv() (x interface{}, ok bool) {
	val, ok := ch.value.TryRecv()
	if val.IsValid() {
		return val.Interface(), ok
	} else {
		return nil, ok
	}
}

var channel = Toolkit{
	Meta: map[string]interface{}{
		"category": "toolkit",
		"name":     "channel",
	},
	Content: map[string]Expr{
		"chan": func(env Env) Element {
			return func(args ...interface{}) (interface{}, error) {
				params, err := GetArgs(env, px.Binds_(
					TypeAs(reflect.TypeOf((*reflect.Type)(nil)).Elem()),
					TypeAs(INT),
					px.Eof), args)
				if err != nil {
					return nil, err
				}
				return MakeBothChan(params[0].(reflect.Type), params[1].(Int)), nil
			}
		},
		"chan->": func(env Env) Element {
			return func(args ...interface{}) (interface{}, error) {
				params, err := GetArgs(env, px.Binds_(
					TypeAs(reflect.TypeOf((*reflect.Type)(nil)).Elem()),
					TypeAs(INT),
					px.Eof), args)
				if err != nil {
					return nil, err
				}
				return MakeRecvChan(params[0].(reflect.Type), params[1].(Int)), nil
			}
		},
		"chan<-": func(env Env) Element {
			return func(args ...interface{}) (interface{}, error) {
				params, err := GetArgs(env, px.Binds_(
					TypeAs(reflect.TypeOf((*reflect.Type)(nil)).Elem()),
					TypeAs(INT),
					px.Eof), args)
				if err != nil {
					return nil, err
				}
				return MakeSendChan(params[0].(reflect.Type), params[1].(Int)), nil
			}
		},
		"send": func(env Env) Element {
			return func(args ...interface{}) (interface{}, error) {
				params, err := GetArgs(env, px.Binds_(
					TypeAs(reflect.TypeOf((*Chan)(nil))),
					px.Either(px.Try(TypeAs(ANYOPTION)), TypeAs(ANYMUST)),
					px.Eof), args)
				if err != nil {
					return nil, err
				}
				ch := params[0].(*Chan)
				ch.Send(params[1])
				return nil, nil
			}
		},
		"send?": func(env Env) Element {
			return func(args ...interface{}) (interface{}, error) {
				params, err := GetArgs(env, px.Binds_(
					TypeAs(reflect.TypeOf((*Chan)(nil))),
					px.Either(px.Try(TypeAs(ANYOPTION)), TypeAs(ANYMUST)),
					px.Eof), args)
				if err != nil {
					return nil, err
				}
				ch := params[0].(*Chan)
				ch.TrySend(params[1])
				return nil, nil
			}
		},
		"recv": func(env Env) Element {
			return func(args ...interface{}) (interface{}, error) {
				params, err := GetArgs(env, px.Binds_(
					TypeAs(reflect.TypeOf((*Chan)(nil))),
					px.Either(px.Try(TypeAs(ANYOPTION)), TypeAs(ANYMUST)),
					px.Eof), args)
				if err != nil {
					return nil, err
				}
				ch := params[0].(*Chan)
				data, ok := ch.Recv()
				return List{data, ok}, nil
			}
		},
		"recv?": func(env Env) Element {
			return func(args ...interface{}) (interface{}, error) {
				params, err := GetArgs(env, px.Binds_(
					TypeAs(reflect.TypeOf((*Chan)(nil))),
					px.Either(px.Try(TypeAs(ANYOPTION)), TypeAs(ANYMUST)),
					px.Eof), args)
				if err != nil {
					return nil, err
				}
				ch := params[0].(*Chan)
				data, ok := ch.TryRecv()
				return List{data, ok}, nil
			}
		},
	},
}
