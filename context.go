package flaskgo

import (
	"net/http"

	"github.com/jtolds/gls"
)

var (
	contextMgr        = gls.NewContextManager()
	requestKey        = gls.GenSym()
	responseWriterKey = gls.GenSym()
)

func GetRequest() *http.Request {
	x, ok := contextMgr.GetValue(requestKey)
	if ok {
		return x.(*http.Request)
	} else {
		panic("Failed to get Request from context")
	}
}

func GetResponseWriter() http.ResponseWriter {
	x, ok := contextMgr.GetValue(responseWriterKey)
	if ok {
		return x.(http.ResponseWriter)
	} else {
		panic("Failed to get ResponseWriter from context")
	}
}

func runInContext(w http.ResponseWriter, r *http.Request, fun func()) {
	contextMgr.SetValues(gls.Values{requestKey: r, responseWriterKey: w}, fun)
}
