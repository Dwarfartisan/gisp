package gisp

import (
	"reflect"
	t "time"
)

var (
	// BOOL 类型
	BOOL = reflect.TypeOf((*bool)(nil)).Elem()
	// STRING 是字符串类型
	STRING = reflect.TypeOf((*string)(nil)).Elem()
	// INT 整型
	INT = reflect.TypeOf((*Int)(nil)).Elem()
	// FLOAT 浮点型
	FLOAT = reflect.TypeOf((*Float)(nil)).Elem()
	// TIME 时间类型
	TIME = reflect.TypeOf((*t.Time)(nil)).Elem()
	// DURATION 时段类型
	DURATION = reflect.TypeOf((*t.Duration)(nil)).Elem()
	// ANY 是 interface{} 的封装
	ANY = reflect.TypeOf((*interface{})(nil)).Elem()
	// ATOM 原子类型
	ATOM = reflect.TypeOf((*Atom)(nil)).Elem()
	// LIST 序列类型
	LIST = reflect.TypeOf((*List)(nil)).Elem()
	// QUOTE 是 lisp quote 类型
	QUOTE = reflect.TypeOf((*Quote)(nil)).Elem()
	// DICT 是 map[string]interface{} 的封装
	DICT = reflect.TypeOf((*map[string]interface{})(nil)).Elem()

	// BOOLOPTION 是可空的 BOOL
	BOOLOPTION = Type{BOOL, true}
	// INTOPTION 是可空的 INT
	INTOPTION = Type{INT, true}
	// FLOATOPTION 是可空的 FLOAT
	FLOATOPTION = Type{FLOAT, true}
	// STRINGOPTION 是可空的 STRING
	STRINGOPTION = Type{STRING, true}
	// TIMEOPTION 是可空的 TIME
	TIMEOPTION = Type{TIME, true}
	// DURATIONOPTION 是可空的 DURATIONOPTION
	DURATIONOPTION = Type{DURATION, true}
	// ANYOPTION 是可空的 ANY
	ANYOPTION = Type{ANY, true}
	// ATOMOPTION 是可空的 ATOM
	ATOMOPTION = Type{ATOM, true}
	// LISTOPTION 是可空的 LIST
	LISTOPTION = Type{LIST, true}
	// QUOTEOPTION 是可空的 QUOTE
	QUOTEOPTION = Type{QUOTE, true}
	// DICTOPTION 是可空的 DICT
	DICTOPTION = Type{DICT, true}

	// BOOLMUST 是不可空的 BOOL
	BOOLMUST = Type{BOOL, false}
	// INTMUST 是不可空的 INT
	INTMUST = Type{INT, false}
	// FLOATMUST 是不可空的 FLOAT
	FLOATMUST = Type{FLOAT, false}
	// STRINGMUST 是不可空的 STRING
	STRINGMUST = Type{STRING, false}
	// TIMEMUST 是不可空的 TIME
	TIMEMUST = Type{TIME, false}
	// DURATIONMUST 是不可空的 DURATION
	DURATIONMUST = Type{DURATION, false}
	// ANYMUST 是不可空的 ANY
	ANYMUST = Type{ANY, false}
	// ATOMMUST 是不可空的 ATOM
	ATOMMUST = Type{ATOM, false}
	// LISTMUST 是不可空的 LISTMUST
	LISTMUST = Type{LIST, false}
	// QUOTEMUST 是不可空的 QUOTE
	QUOTEMUST = Type{QUOTE, false}
	// DICTMUST 是不可空的 DICT
	DICTMUST = Type{DICT, false}
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
