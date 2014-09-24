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
	DICT     = reflect.TypeOf((*map[string]interface{})(nil)).Elem()

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
	DICTOPTION     = Type{DICT, true}

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
	DICTMUST     = Type{DICT, false}
)

// var TypeBox = Gearbox{
// 	Meta: map[string]interface{}{
// 		"category": "package",
// 		"name":     "types",
// 	},
// 	Content: map[string]interface{}{
// 		"int":       INTMUST,
// 		"int?":      INTOPTION,
// 		"float":     FLOATMUST,
// 		"float?":    FLOATOPTION,
// 		"string":    STRINGMUST,
// 		"string?":   STRINGOPTION,
// 		"time":      TIMEMUST,
// 		"time?":     TIMEOPTION,
// 		"dict":      DICTMUST,
// 		"dict?":     DICTMUST,
// 		"duration":  DURATIONMUST,
// 		"duration?": DURATIONOPTION,
// 		"list":      LISTMUST,
// 		"list?":     LISTOPTION,
// 		"atom":      ATOM,
// 		"atom?":     ATOMOPTION,
// 		"quote":     QUOTEMUST,
// 		"quote!":    QUOTEOPTION,
// 		"any":       ANYMUST,
// 		"any?":      ANYOPTION,
// 	},
// }
