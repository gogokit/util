package util

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gogokit/tostr"

	"google.golang.org/protobuf/runtime/protoiface"
)

type stringer struct {
	obj interface{}
}

// 示例用法: logs.CtxDebug(ctx, "%v", xf.Stringer(obj))或logs.CtxDebug(ctx, "%s", xf.Stringer(obj))
func Stringer(obj interface{}) fmt.Stringer {
	return &stringer{obj: obj}
}

func (s *stringer) String() string {
	return toString(s.obj)
}

func toString(obj interface{}) string {
	str := tostr.StringByConf(obj, tostr.Config{
		FilterStructField: []func(reflect.Value, int) bool{
			func(obj reflect.Value, fieldIdx int) bool {
				field := obj.Type().Field(fieldIdx)
				if field.Type.PkgPath() != "" && toStringMap[field.Type] == nil {
					return strings.HasPrefix(field.Type.PkgPath(), "google.golang.org/protobuf")
				}

				if reflect.PtrTo(obj.Type()).Implements(reflect.TypeOf((*protoiface.MessageV1)(nil)).Elem()) {
					if _, inMap := map[string]struct{}{
						"state":                {},
						"sizeCache":            {},
						"unknownFields":        {},
						"XXX_NoUnkeyedLiteral": {},
						"XXX_unrecognized":     {},
						"XXX_sizecache":        {},
					}[field.Name]; inMap {
						return true
					}
				}
				return false
			},
		},
		ToString: func(o reflect.Value) (objStr string) {
			if f, inMap := toStringMap[o.Type()]; inMap {
				return f(o)
			}

			return "<严重错误：指定对象的格式化函数未注册>"
		},
		FastSpecifyToStringProbe: func(o reflect.Value) (hasSpecifyToString bool) {
			_, inMap := toStringMap[o.Type()]
			return inMap
		},
		WarnSize: func() *int {
			ret := 100000
			return &ret
		}(),
		ResultSizeWarnCallback: func(str string) (shouldContinue bool) {
			return false
		},
		DisableMapKeySort: true,
	})
	return str
}

var toStringMap = map[reflect.Type]func(obj reflect.Value) string{
	reflect.TypeOf(time.Time{}): func(obj reflect.Value) string {
		return "{time.Time:\"" + obj.Interface().(time.Time).Format("2006-01-02 15:04:05.000") + "\"}"
	},
	reflect.TypeOf([]byte{}): func(obj reflect.Value) string {
		return strconv.Quote(string(obj.Interface().([]byte)))
	},
	reflect.TypeOf(json.RawMessage{}): func(obj reflect.Value) string {
		var s []byte = obj.Interface().(json.RawMessage)
		return strconv.Quote(string(s))
	},
}
