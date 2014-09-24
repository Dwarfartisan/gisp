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
	Content: map[string]Functor{
		"chan": SimpleBox{
			ParsexSignChecker(px.Binds_(
				TypeAs(reflect.TypeOf((*reflect.Type)(nil)).Elem()),
				TypeAs(INT),
				px.Eof)),
			func(args ...interface{}) Tasker {
				return func(env Env) (interface{}, error) {
					return MakeBothChan(args[0].(reflect.Type), args[1].(Int)), nil
				}
			}},
		"chan->": SimpleBox{
			ParsexSignChecker(px.Binds_(
				TypeAs(reflect.TypeOf((*reflect.Type)(nil)).Elem()),
				TypeAs(INT),
				px.Eof)),
			func(args ...interface{}) Tasker {
				return func(env Env) (interface{}, error) {
					return MakeRecvChan(args[0].(reflect.Type), args[1].(Int)), nil
				}
			}},
		"chan<-": SimpleBox{
			ParsexSignChecker(px.Binds_(
				TypeAs(reflect.TypeOf((*reflect.Type)(nil)).Elem()),
				TypeAs(INT),
				px.Eof)),
			func(args ...interface{}) Tasker {
				return func(env Env) (interface{}, error) {
					return MakeSendChan(args[0].(reflect.Type), args[1].(Int)), nil
				}
			}},
		"send": SimpleBox{
			ParsexSignChecker(px.Binds_(
				TypeAs(reflect.TypeOf((*Chan)(nil))),
				px.Either(px.Try(TypeAs(ANYOPTION)), TypeAs(ANYMUST)),
				px.Eof)),
			func(args ...interface{}) Tasker {
				return func(env Env) (interface{}, error) {
					ch := args[0].(*Chan)
					ch.Send(args[1])
					return nil, nil
				}
			}},
		"send?": SimpleBox{
			ParsexSignChecker(px.Binds_(
				TypeAs(reflect.TypeOf((*Chan)(nil))),
				px.Either(px.Try(TypeAs(ANYOPTION)), TypeAs(ANYMUST)),
				px.Eof)),
			func(args ...interface{}) Tasker {
				return func(env Env) (interface{}, error) {
					ch := args[0].(*Chan)
					ch.TrySend(args[1])
					return nil, nil
				}
			}},
		"recv": SimpleBox{
			ParsexSignChecker(px.Binds_(
				TypeAs(reflect.TypeOf((*Chan)(nil))),
				px.Either(px.Try(TypeAs(ANYOPTION)), TypeAs(ANYMUST)),
				px.Eof)),
			func(args ...interface{}) Tasker {
				return func(env Env) (interface{}, error) {
					ch := args[0].(*Chan)
					data, ok := ch.Recv()
					return List{data, ok}, nil
				}
			}},
		"recv?": SimpleBox{
			ParsexSignChecker(px.Binds_(
				TypeAs(reflect.TypeOf((*Chan)(nil))),
				px.Either(px.Try(TypeAs(ANYOPTION)), TypeAs(ANYMUST)),
				px.Eof)),
			func(args ...interface{}) Tasker {
				return func(env Env) (interface{}, error) {
					ch := args[0].(*Chan)
					data, ok := ch.TryRecv()
					return List{data, ok}, nil
				}
			}},
	},
}
