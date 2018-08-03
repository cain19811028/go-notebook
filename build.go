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
	"time"
)

var md = "framework.md"
var api = "https://api.github.com/repos/"
var layout = "1900-01-01 00:00:00"
var head = "# Framework\n"
var tail = "\n*Update Date: {}*"
var table = `
## {} Framework

| Project Name | Stars | Forks | Last Commit |
| ------------ | ----- | ----- | ----------- |
`

type Repo struct {
	Name string `json:"name"`
	HtmlUrl string `json:"html_url"`
	StargazersCount int `json:"stargazers_count"`
	ForksCount int `json:"forks_count"`
    DefaultBranch string `json:"default_branch"`
    LastCommitDate string
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

        	path = api + url[19:] + "/commits/" + repo.DefaultBranch + "?access_token=" + token
        	response, _ = http.Get(path)
        	defer response.Body.Close()
        	content, _ = ioutil.ReadAll(response.Body)
        	var commits Commits
        	json.Unmarshal([]byte(content), &commits)

        	t, _ := time.Parse("2006-01-02T15:04:05Z", commits.Commit.Committer.LastCommitDate)
        	repo.LastCommitDate = t.Format("2006-01-02 15:04:05")
        	fmt.Println(repo)
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