package flaskgo

import (
	"fmt"
	"strconv"
)

type DefaultConverter struct {
}

func (c *DefaultConverter) Pattern() string {
	return "[^/]+"
}
func (c *DefaultConverter) Parse(input string) Object {
	return input
}

type IntConverter struct {
}

func (c *IntConverter) Pattern() string {
	return `\d+`
}

func (c *IntConverter) Parse(input string) Object {
	if i, err := strconv.Atoi(input); err == nil {
		return i
	} else {
		fmt.Println(err)
		panic(err)
	}
}

type BooleanConverter struct {
}

func (c *BooleanConverter) Pattern() string {
	return `true|false|1|0`
}

func (c *BooleanConverter) Parse(input string) Object {
	return input == "true" || input == "1"
}

type PathConverter struct {
}

func (c *PathConverter) Pattern() string {
	return ".*"
}

func (c *PathConverter) Parse(input string) Object {
	return input
}

var allConverters map[string]Converter

func init() {
	allConverters = make(map[string]Converter)
	allConverters["path"] = new(PathConverter)
	allConverters["int"] = new(IntConverter)
	allConverters["boolean"] = new(BooleanConverter)
	allConverters["default"] = new(DefaultConverter)
}

func GetConverter(name string) Converter {
	if name == "" {
		name = "default"
	}
	return allConverters[name]
}
