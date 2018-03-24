package flaskgo

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path"
)

var templatesPath = "templates"

type FileTemplate struct {
	TplName string
}

func (t *FileTemplate) Name() string {
	return t.TplName
}

func (t *FileTemplate) Type() string {
	return "file"
}

func (t *FileTemplate) Content() string {
	fullFilePath := path.Join(templatesPath, t.TplName)
	fmt.Println(fullFilePath)
	content := fcp.GetFileContent(fullFilePath)
	if content != nil {
		return string(content)
	} else {
		panic("Can not load FileTemplate " + t.Name())
	}
}

//todo
func (t *FileTemplate) loadAll() {
	dir, err := os.Open(templatesPath)
	if err == nil {
		info, err := dir.Stat()
		if err == nil {
			if info.IsDir() {
				files, _ := dir.Readdir(0)
				tplFiles := make([]string, len(files))
				for i, file := range tplFiles {
					tplFiles[i] = path.Join(templatesPath, file)
				}
				t, _ := template.ParseFiles(tplFiles...)
				//t.Execute()
				fmt.Print(t)

			}
		}
	}
}

func renderTpl(tpl string, context *Context, funcMap template.FuncMap) string {
	t, e := template.New("tpl").Funcs(funcMap).Parse(tpl)
	var out bytes.Buffer
	if e == nil {
		t.Execute(&out, context)
	}
	return out.String()
}
