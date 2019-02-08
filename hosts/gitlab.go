package hosts

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"gitlab.com/jfaucherre/mergo/models"
	"gitlab.com/jfaucherre/mergo/tools"
)

const (
	glPersonalAccessTokenURL = "https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html#creating-a-personal-access-token"
)

type gitlab struct {
	token string
}

func newGitlab() (Host, error) {
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
	return ownerAndRepo(repository)
}

func (me gitlab) SubmitPr(opts *models.Opts) error {
	return nil
}

func askForGlCredentials() (Host, error) {
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
