package gisp

import (
	"reflect"
	t "time"
)

var (
	BOOL     = reflect.TypeOf((*bool)(nil)).Elem()
	STRING   = reflect.TypeOf((*string)(nil)).Elem()
	INT      = reflect.TypeOf((*Int)(nil)).Elem()
	FLOAT    = reflect.TypeOf((*Float)(nil)).Elem()
	TIME     = reflect.TypeOf((*t.Time)(nil)).Elem()
	DURATION = reflect.TypeOf((*t.Duration)(nil)).Elem()
	ANY      = reflect.TypeOf((*interface{})(nil)).Elem()
	ATOM     = reflect.TypeOf((*Atom)(nil)).Elem()
	QUOTE    = reflect.TypeOf((*Quote)(nil)).Elem()

	BOOLOPTION     = Type{BOOL, true}
	INTOPTION      = Type{INT, true}
	FLOATOPTION    = Type{FLOAT, true}
	STRINGOPTION   = Type{STRING, true}
	TIMEOPTION     = Type{TIME, true}
	DURATIONOPTION = Type{DURATION, true}
	ANYOPTION      = Type{ANY, true}
	ATOMOPTION     = Type{ATOM, true}
	QUOTEOPTION    = Type{QUOTE, true}

	BOOLMUST     = Type{BOOL, false}
	INTMUST      = Type{INT, false}
	FLOATMUST    = Type{FLOAT, false}
	STRINGMUST   = Type{STRING, false}
	TIMEMUST     = Type{TIME, false}
	DURATIONMUST = Type{DURATION, false}
	ANYMUST      = Type{ANY, false}
	ATOMMUST     = Type{ATOM, false}
	QUOTEMUST    = Type{QUOTE, false}
)
