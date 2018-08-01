package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
)

var md = "framework.md"
var api = "https://api.github.com/repos/"
var head = "# Framework\n"
var tail = "\n*Update Date: {}*"
var table = `
## {} Framework

| Project Name | Stars | Forks | Last Commit |
| ------------ | ----- | ----- | ----------- |
`

func build_head() {
	f, _ := os.Create("test.md")
	w := bufio.NewWriter(f)
	w.WriteString(head)
	w.Flush()
}

func build_tail() {
	f, _ := os.OpenFile("test.md", os.O_APPEND|os.O_WRONLY, 0600)
	f.WriteString(tail)
}

func get_token() string {
	token, _ := ioutil.ReadFile("github_token.txt")
	return string(token)
}

func main() {
	// token := get_token()
    build_head()
    build_tail()
    fmt.Println("finish !")
}