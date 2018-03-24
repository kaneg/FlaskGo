# FlaskGo - Go web framework inspired by Python Flask

# Usage

```go
var app = flaskgo.CreateApp()
app.AddRoute("/index", func() string { return "This is Index page" })
app.AddRoute("/indexDyna1/<int:id>/<path:path>", func(id int, path string) string { return "abc" })
app.Run(":8080")
```