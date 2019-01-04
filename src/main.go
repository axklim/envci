package main

import (
	"encoding/json"
	"fmt"
	"github.com/jessevdk/go-flags"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Opts struct {
	Token     string `long:"token" short:"t" required:"true"`
	ApiUrl    string `long:"url" short:"u" required:"true"`
	ProjectId string `long:"project" short:"p" required:"true"`
	Verbose   []bool `long:"verbose" short:"v" description:"Show verbose debug information"`
	Clear     []bool `long:"clear" short:"c" description:"Clear evn variables"`
	List      []bool `long:"list" short:"l" description:"List all available projects"`
}

type App struct {
	Options *Opts
	Client  *http.Client
}

var app = App{}

type Variables []struct {
	Key   string
	Value string
}

func (vars *Variables) ToMap(m map[string]string) {
	for _, v := range *vars {
		m[v.Key] = v.Value
	}
}

type Projects []Project

type Project struct {
	Id                int
	Name              string
	PathWithNamespace string `json:"path_with_namespace"`
	Namespace         struct {
		Id int
	}
}

func main() {
	var opts Opts
	if _, err := flags.ParseArgs(&opts, os.Args); err != nil {
		os.Exit(1)
	}
	client := &http.Client{}
	app.Options = &opts
	app.Client = client

	// https://docs.gitlab.com/ee/api/README.html#namespaced-path-encoding
	app.Options.ProjectId = strings.Replace(app.Options.ProjectId, "/", "%2F", -1)

	if len(app.Options.List) > 0 {
		app.ListProjects()
		return
	}

	vars := make(map[string]string)
	app.FeedGroupVars(vars)
	app.FeedProjectVars(vars)

	for i, v := range vars {
		if len(app.Options.Clear) > 0 {
			fmt.Printf("%s='%s'\n", i, "")
		} else {
			fmt.Printf("%s='%s'\n", i, v)
		}
	}
}

func (app *App) FeedProjectVars(vars map[string]string) {
	app.GetVars(fmt.Sprintf("%s/projects/%s/variables", app.Options.ApiUrl, app.Options.ProjectId)).ToMap(vars)
}

func (app *App) FeedGroupVars(vars map[string]string) {
	app.GetVars(fmt.Sprintf("%s/groups/%s/variables", app.Options.ApiUrl, app.GetGroupId())).ToMap(vars)
}

func (app *App) GetVars(apiUrl string) *Variables {
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		log.Panic(err)
	}
	req.Header.Add("PRIVATE-TOKEN", app.Options.Token)

	resp, err := app.Client.Do(req)
	defer resp.Body.Close()

	if len(app.Options.Verbose) > 0 {
		fmt.Printf("GET %s %d\n", apiUrl, resp.StatusCode)
	}

	if resp.StatusCode != http.StatusOK {
		return &Variables{}
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panic(err)
	}

	variables := &Variables{}
	if err := json.Unmarshal(body, &variables); err != nil {
		log.Panic(err)
	}
	return variables
}

func (app *App) GetGroupId() string {
	apiUrl := fmt.Sprintf("%s/projects/%s", app.Options.ApiUrl, app.Options.ProjectId)
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		log.Panic(err)
	}
	req.Header.Add("PRIVATE-TOKEN", app.Options.Token)

	resp, err := app.Client.Do(req)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panic(err)
	}

	project := &Project{}
	if err := json.Unmarshal(body, &project); err != nil {
		log.Panic(err)
	}
	return strconv.Itoa(project.Namespace.Id)
}

func (app *App) ListProjects() {
	apiUrl := fmt.Sprintf("%s/projects", app.Options.ApiUrl)
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		log.Panic(err)
	}
	req.Header.Add("PRIVATE-TOKEN", app.Options.Token)

	resp, err := app.Client.Do(req)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panic(err)
	}

	projects := &Projects{}
	if err := json.Unmarshal(body, &projects); err != nil {
		log.Panic(err)
	}

	for _, p := range *projects {
		fmt.Printf("%d\t%s\n", p.Id, p.PathWithNamespace)
	}
}
