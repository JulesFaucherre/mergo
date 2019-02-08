package hosts

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	"gitlab.com/jfaucherre/mergo/models"
	"gitlab.com/jfaucherre/mergo/tools"
)

const ghSeparator = ';'

var (
	// Need to be implemented
	usernameRegex = regexp.MustCompile(".*")
)

type github struct {
	user string
	pswd string
}

func newGithub() (Host, error) {
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

	return askForGhCredentials()
}

func (me github) GetOwnerAndRepo(remote string) (string, string) {
	return ownerAndRepo(remote)
}

func (me github) SubmitPr(opts *models.Opts) error {
	body := struct {
		Head  string `json:"head"`
		Base  string `json:"base"`
		Title string `json:"title"`
		Body  string `json:"body"`
	}{
		opts.Head,
		opts.Base,
		opts.Title,
		opts.Body,
	}
	url := fmt.Sprintf(
		"https://%s:%s@api.github.com/repos/%s/%s/pulls",
		me.user,
		me.pswd,
		opts.Owner,
		opts.Repo,
	)

	marshaled, err := json.Marshal(body)
	fmt.Println(string(marshaled))
	fmt.Println(url)
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

func askForGhCredentials() (Host, error) {
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

	fmt.Println("Please enter your password")
	v, _ := reader.ReadString('\n')
	g.pswd = strings.Trim(v, "\n")

	keep := ""
	for {
		fmt.Println("Do you want these credentials to be kept for next times ([y]/n)?")
		v, _ = reader.ReadString('\n')
		keep = strings.Trim(v, "\n")
		if keep == "" || keep == "y" || keep == "n" {
			break
		}
	}

	if keep != "n" {
		sep := string(ghSeparator)
		content := []byte(g.user + sep + g.pswd)
		if err := tools.WriteHostConfig("github", content); err != nil {
			return nil, err
		}
	}

	return g, nil
}
