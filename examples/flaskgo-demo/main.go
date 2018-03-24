package main

import (
    "fmt"
    "github.com/kaneg/flaskgo"
    "strconv"
)

func s1() {
    w := flaskgo.GetResponseWriter()
    w.Write([]byte("this is in s1\n"))
}
func index() string {
    w := flaskgo.GetResponseWriter()
    w.Write([]byte("from context\n"))
    return "this is in s1111111\n"
}
func indexTpl() string {
    c := make(flaskgo.Context)
    c["Name"] = "Tom"
    c["abc"] = "Jim"
    return app.RenderTemplate("index.html", &c)
}
func d2(parameters flaskgo.Parameters) {
    w := flaskgo.GetResponseWriter()
    w.Write([]byte("this is in d2\n"))

    w.Write([]byte("id=" + fmt.Sprint(parameters["id"]) + "\n"))
    w.Write([]byte("id=" + fmt.Sprint(parameters["path"]) + "\n"))
}
func addRoutes() {
    app.AddRoute("/", indexTpl)
    app.AddRoute("/indexStr", index)
    app.AddRoute("/indexTpl", indexTpl)
    app.AddRoute("/indexStatic", s1)
    app.AddRoute("/param", d2)
    app.AddRoute("/indexDyna1/<int:id>/<path:path>", func(flaskgo.Parameters) string { return "abc" })
    app.AddRoute("/indexDyna3/<id>/<int:x>/<path:path>", func(id string, x int, path string) string { return "id is " + id + strconv.Itoa(x) + "/" + path })
}

func init() {
    //flaskgo.UseMemoryFile(fileMap)
    addRoutes()
}

var app = flaskgo.CreateApp()

func main() {
    fmt.Println("http://localhost:8080")
    app.Run(":8080")
}
