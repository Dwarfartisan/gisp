package gisp

import (
	"fmt"
	"reflect"
)

func adds(args ...interface{}) (interface{}, error) {
	switch len(args) {
	case 0:
		return 0, nil
	default:
		var fixed interface{} = 0
		var add = func(x interface{}) (interface{}, error) {
			switch val := x.(type) {
			case int:
				return fixed.(int) + val, nil
			case int8, int16, int32, int64:
				return fixed.(int) + int(reflect.ValueOf(val).Int()), nil
			case float32, float64:
				return float64(fixed.(int)) + reflect.ValueOf(val).Float(),
					fmt.Errorf("%v:%t is a float", x, x)
			default:
				return nil, fmt.Errorf("%v:%t is not a number", x, x)
			}
		}
		for _, arg := range args {
			var err error
			fixed, err = add(arg)
			if err != nil {
				if fixed != nil {
					add = func(x interface{}) (interface{}, error) {
						switch val := x.(type) {
						case int, int8, int16, int32, int64:
							return fixed.(float64) + float64(reflect.ValueOf(val).Int()), nil
						case float32:
							return fixed.(float64) + float64(val), nil
						case float64:
							return fixed.(float64) + val, nil
						default:
							return nil, fmt.Errorf("%v:%t is not a number", x, x)
						}
					}
				} else {
					return nil, err
				}
			}
		}
		return fixed, nil
	}
}
