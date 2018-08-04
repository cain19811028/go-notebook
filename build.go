package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

var md = "framework.md"
var api = "https://api.github.com/repos/"
var layout = "1900-01-01 00:00:00"
var head = "# Framework\n"
var tail = "\n*Update Date: %s*"
var table = `
## %s Framework

| Project Name | Stars | Forks | Last Commit |
| ------------ | ----- | ----- | ----------- |
`
var column = "| [%s](%s) | %s | %s | %s |\n"

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
	var repos []Repo
	var repo Repo
	var commits Commits
	f, _ := os.Open("list/" + file)
	scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        url := scanner.Text()
        if strings.HasPrefix(url, "https://github.com/") {
        	path := api + url[19:] + "?access_token=" + token
        	response, _ := http.Get(path)
        	defer response.Body.Close()
        	content, _ := ioutil.ReadAll(response.Body)
        	
        	json.Unmarshal([]byte(content), &repo)

        	path = api + url[19:] + "/commits/" + repo.DefaultBranch + "?access_token=" + token
        	response, _ = http.Get(path)
        	defer response.Body.Close()
        	content, _ = ioutil.ReadAll(response.Body)
        	
        	json.Unmarshal([]byte(content), &commits)

        	t, _ := time.Parse("2006-01-02T15:04:05Z", commits.Commit.Committer.LastCommitDate)
        	repo.LastCommitDate = t.Format("2006-01-02 15:04:05")
        	repos = append(repos, repo)
        	fmt.Println(repo)
        }
    }

    sort.SliceStable(repos, func(i, j int) bool {
		return repos[i].StargazersCount < repos[j].StargazersCount
	})

	build_info(title, repos)
}

func build_info(title string, repos []Repo) {
	f, _ := os.OpenFile(md, os.O_APPEND|os.O_WRONLY, 0600)
	f.WriteString(fmt.Sprintf(table, title))
	var result = ""
	for _, repo := range repos {
		result = fmt.Sprintf(column, repo.Name, repo.HtmlUrl, strconv.Itoa(repo.StargazersCount), strconv.Itoa(repo.ForksCount), repo.LastCommitDate)
		f.WriteString(result)
	}
}

func build_head() {
	f, _ := os.Create(md)
	w := bufio.NewWriter(f)
	w.WriteString(head)
	w.Flush()
}

func build_tail() {
	f, _ := os.OpenFile(md, os.O_APPEND|os.O_WRONLY, 0600)
	f.WriteString(fmt.Sprintf(tail, time.Now().Format("2006-01-02 15:04:05")))
}

func get_token() string {
	token, _ := ioutil.ReadFile("github_token.txt")
	return string(token)
}

func main() {
	token := get_token()
    build_head()
    load(token, "Web", "web_list.txt")
    load(token, "Testing", "test_list.txt")
    load(token, "IoT", "iot_list.txt")
    build_tail()
}