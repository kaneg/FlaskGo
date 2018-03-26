package flaskgo

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

func getRuleParts(rule string) (RPS, bool) {
	allMatched := ruleRegexp.FindAllStringSubmatch(rule, -1)
	size := len(allMatched)
	if size == 0 {
		return nil, false
	}
	allRouteRules := make(RPS, size)
	for i := 0; i < size; i++ {
		matched := allMatched[i]
		routeRule := &allRouteRules[i]
		for i, name := range ruleRegexp.SubexpNames() {
			value := matched[i]
			switch name {
			case "static":
				routeRule.static = value
			case "converter":
				routeRule.converter = value
			case "variable":
				routeRule.variable = value
			case "args":
				routeRule.args = value
			}
		}
	}
	return allRouteRules, true
}

func (app *App) AddRoute(rule string, handler interface{}, methods ...string) {
	rule = app.Prefix + rule
	app.addRoute(rule, handler, methods...)
}

func (app *App) addRoute(rule string, handler interface{}, methods ...string) {
	funcType := reflect.TypeOf(handler)
	if funcType.Kind() != reflect.Func && funcType.Kind() != reflect.Struct {
		panic("Second argument must be map function.")
	}
	if strings.Contains(rule, "<") {
		app.addDynamicRoute(rule, handler, methods...)
	} else {
		app.addStaticRoute(rule, handler)
	}
}

func (app *App) AddRouteAll(prefix string, target interface{}) {
	funcType := reflect.TypeOf(target)
	methodsType := reflect.ValueOf(target)
	for i := 0; i < methodsType.NumMethod(); i++ {
		method := methodsType.Method(i)
		app.AddRoute(prefix+"/"+strings.ToLower(funcType.Method(i).Name), method)
	}
}

func (app *App) addDynamicRoute(rule string, handler DynamicRouteFunc, methods ...string) {
	if ruleParts, ok := getRuleParts(rule); ok {
		var finalRegex = ""
		for _, routeRule := range ruleParts {
			finalRegex += routeRule.static
			converter := GetConverter(routeRule.converter)
			if converter == nil {
				panic("Converter is nil")
			}
			pattern := converter.Pattern()
			x := fmt.Sprintf("(?P<%s>%s)", routeRule.variable, pattern)
			finalRegex += x
		}
		app.dynamicRules[regexp.MustCompile(finalRegex)] = Router{ruleParts, handler, methods}
	}
}

func (app *App) addStaticRoute(rule string, handler StaticRouteFunc) {
	app.staticRules[rule] = handler
}

func (app *App) Redirect(location string) func() {
	return Redirect(app.Prefix+location)
}
