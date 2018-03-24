package flaskgo

import (
	"fmt"
	"net/http"
	"reflect"
	"regexp"
)

func handleDynamicRule(rule *regexp.Regexp, path string, route Router, w http.ResponseWriter, r *http.Request) bool {
	matched := rule.FindStringSubmatch(path)
	if len(matched) > 0 {
		parameters, values := route.Parameters(matched, rule.SubexpNames())
		handler := route.handler
		funcType := reflect.TypeOf(handler)
		switch funcType {
		case dynamicStringFunc:
			v, ok := handler.(func(Parameters) string)
			if !ok {
				panic("The handler function is not valid")
			}
			runInContext(w, r, func() {
				s := v(parameters)
				w.Write([]byte(s))
			})
		case dynamicStringFuncWithoutReturn:
			fun, ok := handler.(func(Parameters))
			if !ok {
				panic("The handler function is not valid")
			}
			runInContext(w, r, func() {
				fun(parameters)
			})
		default:
			numIn := funcType.NumIn()
			var in = make([]reflect.Value, numIn)

			for i := 0; i < numIn; i++ {
				in[i] = reflect.ValueOf(values[i])
			}
			runInContext(w, r, func() {
				outArray := reflect.ValueOf(handler).Call(in)
				if len(outArray) == 1 {
					out := outArray[0].Interface()
					outType := reflect.TypeOf(out)
					if outType.Kind() == reflect.String {
						var s = out.(string)
						w.Write([]byte(s))
					}
				}
			})
		}
		return true
	} else {
		return false
	}
}

func handlerStaticRule(handler StaticRouteFunc, w http.ResponseWriter, r *http.Request) {
	funcType := reflect.TypeOf(handler)
	if funcType.Kind() == reflect.Struct {
		value := handler.(reflect.Value)
		handler = value.Interface()
		funcType = reflect.TypeOf(handler)
	}
	switch funcType {
	case staticStringFunc:
		v, ok := handler.(func() string)
		if !ok {
			panic("The handler function is not valid")
		}
		runInContext(w, r, func() {
			s := v()
			w.Write([]byte(s))
		})
	case staticStringFuncWithoutReturn:
		v, ok := handler.(func())
		if !ok {
			panic("The handler function is not valid")
		}
		runInContext(w, r, func() {
			v()
		})
	default:
		fmt.Println(funcType)
		panic("Not supported handle function type")
	}
}
