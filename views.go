package flaskgo

import (
	"path"
)

var statePath = "static"

func staticResource(parameters Parameters) {
	w := GetResponseWriter()
	filePath := parameters.String("path")
	fullFilePath := path.Join(statePath, filePath)
	content := fcp.GetFileContent(fullFilePath)
	if content != nil {
		w.Write(content)
	} else {
		w.Write([]byte("Not Found"))
		w.WriteHeader(404)
	}
}

func Redirect(location string) func() {
	return func() {
		w := GetResponseWriter()
		w.Header().Add("Location", location)
		w.WriteHeader(302)
	}
}
