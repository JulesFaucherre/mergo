package gitlab

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"gitlab.com/jfaucherre/mergo/git"
	"gitlab.com/jfaucherre/mergo/models"
	"gitlab.com/jfaucherre/mergo/tools"
)

const (
	glPersonalAccessTokenURL = "https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html#creating-a-personal-access-token"
)

var (
	baseContent = []byte(`# Enter the content of your pull request
# Every line starting with a '#' will be considered as a comment and not treated
# Please note that new lines in gitlab's descriptions are <br />

`)
)

type gitlab struct {
	token string
}

func NewGitlab() (models.Host, error) {
	creds, err := tools.GetHostConfig("gitlab")
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	if len(creds) == 0 {
		return askForGlCredentials()
	}
	return gitlab{
		string(creds),
	}, nil
}

func (me gitlab) GetOwnerAndRepo(repository string) (string, string) {
	return tools.DefaultOwnerAndRepo(repository)
}

func (me gitlab) SubmitPr(opts *models.Opts) error {
	var err error
	body := struct {
		SourceBranch string `json:"source_branch"`
		TargetBranch string `json:"target_branch"`
		Title        string `json:"title"`
		Description  string `json:"description"`
	}{
		SourceBranch: opts.Head,
		TargetBranch: opts.Base,
	}

	stdin := bufio.NewReader(os.Stdin)
	fmt.Println("Enter the pull request's title:")
	if body.Title, err = stdin.ReadString('\n'); err != nil {
		return err
	}
	body.Title = strings.Trim(body.Title, "\n")

	if body.Description, err = git.EditText(baseContent); err != nil {
		return err
	}

	url := fmt.Sprintf(
		"https://gitlab.com/api/v4/projects/%s%%2f%s/merge_requests",
		opts.Owner,
		opts.Repo,
	)

	marshaled, err := json.Marshal(body)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(marshaled)

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, reader)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Private-Token", me.token)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("Request failed with status %s", resp.Status)
	}

	res := struct {
		WebURL string `json:"web_url"`
	}{}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(b, &res); err != nil {
		return err
	}

	fmt.Printf("Your pull request is available at the following URL:\n%s", res.WebURL)

	return nil
}

func askForGlCredentials() (models.Host, error) {
	g := &gitlab{}
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("We need your credentials for gitlab")

	for {
		fmt.Printf("Please enter your gitlab private token\nif you have none, you can get one here: %s\n", glPersonalAccessTokenURL)
		v, err := reader.ReadString('\n')
		g.token = strings.Trim(v, "\n")
		if err == nil && len(g.token) > 0 {
			break
		}
	}

	keep := ""
	for {
		fmt.Println("Do you want these credentials to be kept for next times ([y]/n)?")
		keep, _ = reader.ReadString('\n')
		keep = strings.Trim(keep, "\n")
		if keep == "" || keep == "y" || keep == "n" {
			break
		} else {
			fmt.Printf("Invalid input: %s\n", keep)
		}
	}

	if keep != "n" {
		content := []byte(g.token)
		if err := tools.WriteHostConfig("gitlab", content); err != nil {
			return nil, err
		}
	}

	return g, nil
}
