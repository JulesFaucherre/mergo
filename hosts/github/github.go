package github

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	"gitlab.com/jfaucherre/mergo/git"
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

func (me github) GetOwnerAndRepo(remote string) (string, string) {
	return tools.DefaultOwnerAndRepo(remote)
}

func (me github) SubmitPr(opts *models.Opts) error {
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

	stdin := bufio.NewReader(os.Stdin)
	fmt.Println("Enter the pull request's title:")
	if body.Title, err = stdin.ReadString('\n'); err != nil {
		return err
	}
	body.Title = strings.Trim(body.Title, "\n")

	if body.Body, err = git.EditText(baseContent); err != nil {
		return err
	}

	url := fmt.Sprintf(
		"https://%s:%s@api.github.com/repos/%s/%s/pulls",
		me.user,
		me.token,
		opts.Owner,
		opts.Repo,
	)

	marshaled, err := json.Marshal(body)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(marshaled)
	resp, err := http.Post(url, "application/json", reader)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("Request failed with status %s", resp.Status)
	}
	return nil
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
