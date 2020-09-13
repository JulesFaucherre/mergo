package github

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"

	hostTools "gitlab.com/jfaucherre/mergo/hosts/tools"
	"gitlab.com/jfaucherre/mergo/models"
	"gitlab.com/jfaucherre/mergo/tools"
)

const (
	ghSeparator    = ';'
	createTokenURL = "https://github.com/settings/tokens"
)

var (
	// Need to be implemented
	usernameRegex = regexp.MustCompile(".*")
	baseContent   = []byte(`# Enter the content of your pull request
# Every line starting with a '#' will be considered as a comment and not treated
`)
)

type github struct {
	user  string
	token string
}

func NewGithub() (models.Host, error) {
	creds, err := tools.GetHostConfig("github")
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	for i, c := range creds {
		if c == ghSeparator && i != len(creds)-1 {
			return github{
				string(creds[0:i]),
				string(creds[i+1:]),
			}, nil
		}
	}

	return askForCredentials()
}

func (me github) SubmitPr(opts *models.Opts) (*models.MRInfo, error) {
	var err error
	body := struct {
		Head  string `json:"head"`
		Base  string `json:"base"`
		Title string `json:"title"`
		Body  string `json:"body"`
	}{
		Head: opts.Head,
		Base: opts.Base,
	}

	userInfo, err := hostTools.DefaultGetUserInfo(opts)
	if err != nil {
		return nil, err
	}

	body.Title, body.Body = userInfo.Title, userInfo.Body

	url := fmt.Sprintf(
		"https://%s:%s@api.github.com/repos/%s/%s/pulls",
		me.user,
		me.token,
		opts.Owner,
		opts.Repository,
	)

	marshaled, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(marshaled)
	resp, err := http.Post(url, "application/json", reader)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("request failed with status %s", resp.Status)
		}
		return nil, fmt.Errorf("request failed with\n\tstatus: %s\n\tbody: %s", resp.Status, string(body))
	}

	res := struct {
		HTMLURL string `json:"html_url"`
	}{}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(b, &res); err != nil {
		return nil, err
	}

	return &models.MRInfo{
		URL: res.HTMLURL,
	}, nil
}

func askForCredentials() (models.Host, error) {
	g := &github{}
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("We need your credentials for github")

	for {
		fmt.Println("Please enter your username")
		v, _ := reader.ReadString('\n')
		g.user = strings.Trim(v, "\n")
		if !usernameRegex.MatchString(g.user) {
			fmt.Println("Invalid github username")
		} else {
			break
		}
	}

	fmt.Printf(`Please enter an github API token.
To create a token, please go to this URL: %s
Note that you must give at least the "repo" rights
`, createTokenURL)
	v, _ := reader.ReadString('\n')
	g.token = strings.Trim(v, "\n")

	keep := ""
	for {
		fmt.Println("Do you want these credentials to be kept for next times ([y]/n)?")
		v, _ = reader.ReadString('\n')
		keep = strings.Trim(v, "\n")
		if keep == "" || keep == "y" || keep == "n" {
			break
		} else {
			fmt.Printf("Invalid input: %s\n", keep)
		}
	}

	if keep != "n" {
		sep := string(ghSeparator)
		content := []byte(g.user + sep + g.token)
		if err := tools.WriteHostConfig("github", content); err != nil {
			return nil, err
		}
	}

	return g, nil
}
