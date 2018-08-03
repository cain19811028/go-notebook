package main

import (
	"bufio"
	// "container/list"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
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

type Repo struct {
    DefaultBranch string `json:"default_branch"`
}

type Committer struct {
	LastCommitDate string `json:"date"`
}

type Commit struct {
    Committer Committer
}

type Commits struct {
    Commit Commit
}

func load(token string, title string, file string) {
	// l := list.New()
	f, _ := os.Open("list/" + file)
	scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        url := scanner.Text()
        if strings.HasPrefix(url, "https://github.com/") {
        	path := api + url[19:] + "?access_token=" + token
        	response, _ := http.Get(path)
        	defer response.Body.Close()
        	content, _ := ioutil.ReadAll(response.Body)
        	var repo Repo
        	json.Unmarshal([]byte(content), &repo)
        	fmt.Println(repo.DefaultBranch)

        	path = api + url[19:] + "/commits/" + repo.DefaultBranch + "?access_token=" + token
        	response, _ = http.Get(path)
        	defer response.Body.Close()
        	content, _ = ioutil.ReadAll(response.Body)
        	var commits Commits
        	json.Unmarshal([]byte(content), &commits)
        	fmt.Println(commits.Commit.Committer.LastCommitDate)
        }
    }
}

func build_info() {

}

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
	token := get_token()
    build_head()
    load(token, "Web", "web_list.txt")
    build_tail()
}