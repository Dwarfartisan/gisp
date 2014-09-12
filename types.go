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
	LIST     = reflect.TypeOf((*List)(nil)).Elem()
	QUOTE    = reflect.TypeOf((*Quote)(nil)).Elem()

	BOOLOPTION     = Type{BOOL, true}
	INTOPTION      = Type{INT, true}
	FLOATOPTION    = Type{FLOAT, true}
	STRINGOPTION   = Type{STRING, true}
	TIMEOPTION     = Type{TIME, true}
	DURATIONOPTION = Type{DURATION, true}
	ANYOPTION      = Type{ANY, true}
	ATOMOPTION     = Type{ATOM, true}
	LISTOPTION     = Type{LIST, true}
	QUOTEOPTION    = Type{QUOTE, true}

	BOOLMUST     = Type{BOOL, false}
	INTMUST      = Type{INT, false}
	FLOATMUST    = Type{FLOAT, false}
	STRINGMUST   = Type{STRING, false}
	TIMEMUST     = Type{TIME, false}
	DURATIONMUST = Type{DURATION, false}
	ANYMUST      = Type{ANY, false}
	ATOMMUST     = Type{ATOM, false}
	LISTMUST     = Type{LIST, false}
	QUOTEMUST    = Type{QUOTE, false}
)

var TypeBox = Gearbox{
	Meta: map[string]interface{}{
		"category": "package",
		"name":     "types",
	},
	Content: map[string]interface{}{
		"int":       INT,
		"int?":      INTOPTION,
		"int!":      INTMUST,
		"float":     FLOAT,
		"float?":    FLOATOPTION,
		"float!":    FLOATMUST,
		"string":    STRING,
		"string?":   STRINGOPTION,
		"string!":   STRINGMUST,
		"time":      TIME,
		"time?":     TIMEOPTION,
		"time!":     TIMEMUST,
		"duration":  DURATION,
		"duration?": DURATIONOPTION,
		"duration!": DURATIONMUST,
		"list":      LIST,
		"list?":     LISTOPTION,
		"list!":     LISTMUST,
		"atom":      ATOM,
		"atom?":     ATOMOPTION,
		"atom!":     ATOMMUST,
		"quote":     QUOTE,
		"quote!":    QUOTEOPTION,
		"quote?":    QUOTEMUST,
		"any":       ANY,
		"any?":      ANYOPTION,
		"any!":      ANYMUST,
	},
}
