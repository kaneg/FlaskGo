package flaskgo

import (
	"log"
	"html/template"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

const PS = `(?P<static>[^<]*)<(?:(?P<converter>[a-zA-Z_][a-zA-Z0-9_]*)(?:\((?P<args>.*?)\))?\:)?(?P<variable>[a-zA-Z_][a-zA-Z0-9_]*)>`

var ruleRegexp *regexp.Regexp

var fcp FileContentProvider = new(DefaultFileContentProvider)

type App struct {
	Prefix       string
	staticRules  map[string]StaticRouteFunc
	dynamicRules map[*regexp.Regexp]Router
	funcMap      template.FuncMap
}

func (app *App) initDefaultFuncMap() {
	if app.funcMap == nil {
		app.funcMap = template.FuncMap{}
	}
	app.funcMap["size"] = func(size int64) string {
		unit := "B"
		if size > 1024*1024*1024 {
			size /= 1024 * 1024 * 1024
			unit = "GB"

		} else if size > 1024*1024 {
			size /= 1024 * 1024
			unit = "MB"

		} else if size > 1024 {
			size /= 1024
			unit = "KB"
		}
		return strconv.FormatInt(size, 10) + " " + unit
	}
	app.funcMap["time_format"] = func(layout string, t time.Time) string {
		if layout == "" {
			layout = "2006-01-02 15:04:05"
		}
		return t.Format(layout)
	}
}

func (app *App) getFuncMap() template.FuncMap {
	app.initDefaultFuncMap()
	return app.funcMap
}

func (app *App) Handler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	log.Println("Request URL: ", r.URL)
	var handled = false
	if handler, ok := app.staticRules[path]; ok {
		handlerStaticRule(handler, w, r)
		handled = true
	} else {
		for reg, route := range app.dynamicRules {
			methods := route.methods
			if methods != nil {
				if !contains(methods, r.Method) {
					continue
				}
			}
			handled = handleDynamicRule(reg, path, route, w, r)
			if handled {
				break
			}
		}
	}
	if !handled {
		w.WriteHeader(404)
		w.Write([]byte("No handler found"))
	}
}

func (app *App) Run(address string) {
	http.HandleFunc("/", app.Handler)
	http.ListenAndServe(address, nil)
}

func (app *App) RenderTemplate(name string, context *Context) string {
	tpl := FileTemplate{TplName: name}
	return renderTpl(tpl.Content(), context, app.getFuncMap())
}

func init() {
	ruleRegexp = regexp.MustCompile(PS)
}

func initApp(app App) {
	app.addRoute("/", app.Redirect("/"))
	app.AddRoute("/static/<path:path>", staticResource)
}

func (router Router) Parameters(real []string, keywords []string) (Parameters, []Object) {
	ruleConverterMap := make(map[string]string)
	for _, r := range router.ruleParts {
		ruleConverterMap[r.variable] = r.converter
	}

	var result = make(Parameters)
	var values = make([]Object, 0, len(keywords))
	for i, name := range keywords {
		if name != "" {
			converterType := ruleConverterMap[name]
			value := GetConverter(converterType).Parse(real[i])
			result[name] = value
			values = append(values, value)
		}
	}

	return result, values
}

func CreateAppWithPrefix(prefix string) App {
	app := App{Prefix: prefix}
	app.staticRules = make(map[string]StaticRouteFunc)
	app.dynamicRules = make(map[*regexp.Regexp]Router)
	initApp(app)

	return app
}

func CreateApp() App {
	return CreateAppWithPrefix("")
}