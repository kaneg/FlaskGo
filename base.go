package flaskgo

import (
	"fmt"
	"os"
	"reflect"
)

type Object interface {
}

type RulePart struct {
	static,
	converter,
	variable,
	args string
}

type RPS []RulePart

func (r *RulePart) String() string {
	return fmt.Sprintf("{static: [%s], converter: [%s], var: [%s], args: [%s]}", r.static, r.converter, r.variable, r.args)
}

func (r RPS) String() string {
	s := "<"
	for _, x := range r {
		s += x.String()
		s += ", "
	}
	return s + ">"
}

type Converter interface {
	Pattern() string
	Parse(input string) Object
}

type Router struct {
	ruleParts []RulePart
	handler   DynamicRouteFunc
	methods   []string
}

type VarPair struct {
	name  string
	value interface{}
}

type Parameters map[string]Object

func (p Parameters) String(key string) string {
	object := p[key]
	return object.(string)
}

func (p Parameters) Int(key string) int {
	object := p[key]
	return object.(int)
}

type StaticRouteFunc interface{}

type DynamicRouteFunc interface{}

var staticStringFunc = reflect.TypeOf(func() string { return "" })
var staticStringFuncWithoutReturn = reflect.TypeOf(func() {})

var dynamicStringFunc = reflect.TypeOf(func(Parameters) string { return "" })
var dynamicStringFuncWithoutReturn = reflect.TypeOf(func(Parameters) {})

type Template interface {
	Name() string
	Type() string
	Content() string
}

type Context map[string]Object

type FileContentProvider interface {
	GetFileContent(fullFilePath string) []byte
}

type DefaultFileContentProvider struct {
}

func (f *DefaultFileContentProvider) GetFileContent(fullFilePath string) []byte {
	return GetFileContent(fullFilePath)
}

type MemoryFileContentProvider struct {
	FileMap map[string][]byte
}

func (m *MemoryFileContentProvider) GetFileContent(fullFilePath string) []byte {
	return m.FileMap[fullFilePath]
}

func SetFileContentProvider(f FileContentProvider) {
	fcp = f
}

func UseMemoryFile(fileMap map[string][]byte) {
	SetFileContentProvider(&MemoryFileContentProvider{FileMap: fileMap})

}

func GetFileContent(fullFilePath string) []byte {
	if file, e := os.Open(fullFilePath); e == nil {
		defer file.Close()
		if info, e := file.Stat(); e == nil {
			buffer := make([]byte, info.Size())
			file.Read(buffer)
			return buffer
		}
	}
	return nil
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
