package gitlab

import (
	"fmt"
	"regexp"
	"strings"

	"gitlab.com/jfaucherre/mergo/credentials"
	"gitlab.com/jfaucherre/mergo/logger"
	"gitlab.com/jfaucherre/mergo/models"
	"gitlab.com/jfaucherre/mergo/tools"
)

var (
	// IsHostRegexp is a regexp which matches URLs that are gitlab's repositories
	IsHostRegexp = regexp.MustCompile("(^(https?://)|(git@))gitlab.com[:/].*$")
)

type gitlab struct {
	Token []byte `parse-name:"private token" parse-info:"If you have none, you can get one here: https://gitlab.com/-/profile/personal_access_tokens. Note that you must give api rights" parse-reg:"^.+$"`
}

func (me gitlab) Marshal() ([]byte, error) {
	return me.Token, nil
}

func (me *gitlab) Unmarshal(content []byte) error {
	me.Token = content
	return nil
}

func (me gitlab) Name() string {
	return "gitlab"
}

func New() (models.Host, error) {
	gl := gitlab{}

	err := tools.LoadHostCredentials(&gl)
	if err == credentials.ErrNoHostConfig {
		err = tools.AskForHostCredentials(&gl)
	}
	if err != nil {
		return nil, err
	}

	return &gl, nil
}

func (me gitlab) SubmitPr(params *models.MRParams) (*models.MRInfo, error) {
	message := strings.Split(params.Message, "\n")
	title := message[0]
	description := strings.Join(message[1:], "<br />")
	body := struct {
		SourceBranch string `json:"source_branch"`
		TargetBranch string `json:"target_branch"`
		Title        string `json:"title"`
		Description  string `json:"description"`
	}{
		SourceBranch: params.Head,
		TargetBranch: params.Base,
		Title:        title,
		Description:  description,
	}

	infos, err := tools.RepoInfoFromURL(params.URL)
	logger.Info("repo infos %+v\n", infos)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf(
		"https://gitlab.com/api/v4/projects/%s%%2F%s/merge_requests",
		infos.Owner,
		infos.Repo,
	)

	res := new(struct {
		Error   string `json:"error"`
		Message string `json:"message"`
		WebURL  string `json:"web_url"`
	})

	status, err := tools.Request(&tools.RequestParams{
		URL:    url,
		Method: "POST",
		Body:   &body,
		Headers: map[string]string{
			"Private-Token": string(me.Token),
		},
		Result: res,
	})

	if err != nil {
		return nil, err
	}
	if status >= 400 {
		logger.Info("res.Error = %+v\n", res.Error)
		logger.Info("res.Message = %+v\n", res.Message)
		return nil, fmt.Errorf("Request failed with status %d", status)
	}

	return &models.MRInfo{
		URL: res.WebURL,
	}, nil
}
