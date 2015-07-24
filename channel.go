package gisp

import (
	"reflect"

	px "github.com/Dwarfartisan/goparsec/parsex"
)

// Chan 封装 golang 的 channel 功能
type Chan struct {
	Type  reflect.Type
	Dir   reflect.ChanDir
	value reflect.Value
}

// MakeChan 实现 chan 的构造
func MakeChan(typ reflect.Type, dir reflect.ChanDir, buf Int) *Chan {
	return &Chan{typ, dir, reflect.MakeChan(typ, int(buf))}
}

// MakeRecvChan 实现一个单向的 recv chan
func MakeRecvChan(typ reflect.Type, buf Int) *Chan {
	return MakeChan(typ, reflect.RecvDir, buf)
}

// MakeSendChan 实现一个单向的 Send chan
func MakeSendChan(typ reflect.Type, buf Int) *Chan {
	return MakeChan(typ, reflect.SendDir, buf)
}

// MakeBothChan 构造一个双向 chan
func MakeBothChan(typ reflect.Type, buf Int) *Chan {
	return MakeChan(typ, reflect.BothDir, buf)
}

// Send 方法实现 chan x <- v
func (ch *Chan) Send(x interface{}) {
	ch.value.Send(reflect.ValueOf(x))
}

// Recv 方法实现 v <- chan x
func (ch *Chan) Recv() (x interface{}, ok bool) {
	val, ok := ch.value.Recv()
	if val.IsValid() {
		return val.Interface(), ok
	}
	return nil, ok
}

// TrySend 实现试写入（带状态返回）
func (ch *Chan) TrySend(x interface{}) {
	ch.value.TrySend(reflect.ValueOf(x))
}

// TryRecv 实现试接收（带状态返回）
func (ch *Chan) TryRecv() (x interface{}, ok bool) {
	val, ok := ch.value.TryRecv()
	if val.IsValid() {
		return val.Interface(), ok
	}
	return nil, ok
}

var channel = Toolkit{
	Meta: map[string]interface{}{
		"category": "toolkit",
		"name":     "channel",
	},
	Content: map[string]interface{}{
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
