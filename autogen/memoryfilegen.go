package main

import (
    "encoding/base64"
    "fmt"
    "os"
    "path/filepath"
    "strings"

    "github.com/kaneg/flaskgo"
)

func main() {
    fmt.Println("Gen")
    fmt.Println(os.Args)
    outputFile := os.Args[1] + ".go"
    folders := os.Args[2:]
    fmt.Println("outputfile:" + outputFile)
    fmt.Println(folders)
    fileContentMap := make(map[string][]byte)
    for _, folder := range folders {
        fmt.Println("Process folder:", folder)
        var walkFun filepath.WalkFunc = func(path string, info os.FileInfo, err error) error {
            if !info.IsDir() {
                fmt.Println("\t", path)
                fileContentMap[path] = flaskgo.GetFileContent(path)
            }
            return nil
        }
        filepath.Walk(folder, walkFun)
    }

    fmt.Println("Total file:", len(fileContentMap))

    s := 
`//Auto generated file

package main

import "encoding/base64"
import "github.com/kaneg/flaskgo"

var fileMap = make(map[string][]byte)

func init() {
	flaskgo.UseMemoryFile(fileMap)

`
    if file, e := os.Create(outputFile); e == nil {
        file.WriteString(s)
        for fileName, fileContent := range fileContentMap {
            toStrings(fileName, fileContent, file)
        }
        file.WriteString("}\n")
    }
}
func toBytes(fileName string, fileContent []byte, file *os.File) {
    fileName = strings.Replace(fileName, "\\", "/", -1)
    c := fmt.Sprint(fileContent)
    split := strings.Split(c, " ")
    for i := 0; i < len(split); i++ {
        if i%20 == 0 {
            split[i] = "\n" + split[i]
        }
    }
    join := strings.Join(split, ",")
    s := join[2:len(join)-1]
    s0 := fmt.Sprintf("fileMap[\"%s\"] = []byte{%s}\n", fileName, s)
    file.WriteString(s0)
}
func toStrings(fileName string, fileContent []byte, file *os.File) {
    fileName = strings.Replace(fileName, "\\", "/", -1)
    s := base64.StdEncoding.EncodeToString([]byte(fileContent))
    s0 := fmt.Sprintf("    fileMap[\"%s\"], _ = base64.StdEncoding.DecodeString(`%s`)\n", fileName, s)
    file.WriteString(s0)
}
