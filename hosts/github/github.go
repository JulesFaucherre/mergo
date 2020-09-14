package github

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"gitlab.com/jfaucherre/mergo/credentials"
	"gitlab.com/jfaucherre/mergo/logger"
	"gitlab.com/jfaucherre/mergo/models"
	"gitlab.com/jfaucherre/mergo/tools"
)

const (
	separator = ';'
	// IsHostRegexp is a regexp which matches URLs that are github's repositories
)

var (
	IsHostRegexp = regexp.MustCompile("(^(https?://)|(git@))github.com[:/].*$")

	ErrMalformedCredentials = errors.New("malformed github credntials")
)

type github struct {
	User  []byte `parse-name:"github username" parse-reg:"^.+$"`
	Token []byte `parse-name:"private token" parse-info:"To generate new tokens go here: https://github.com/settings/tokens" parse-reg:"^.+$"`
}

func (me github) Marshal() ([]byte, error) {
	return append(
		append(me.User, separator),
		me.Token...,
	), nil
}

func (me *github) Unmarshal(content []byte) error {
	index := bytes.IndexByte(content, separator)
	if index == -1 {
		return ErrMalformedCredentials
	}

	me.User, me.Token = content[:index], content[index+1:]
	return nil
}

func (me github) Name() string {
	return "github"
}

func New() (models.Host, error) {
	gh := github{}

	err := tools.LoadHostCredentials(&gh)
	if err == credentials.ErrNoHostConfig {
		err = tools.AskForHostCredentials(&gh)
	}
	if err != nil {
		return nil, err
	}

	return &gh, nil
}

func (me github) SubmitPr(params *models.MRParams) (*models.MRInfo, error) {
	message := strings.Split(params.Message, "\n")
	title := message[0]
	description := strings.Join(message[1:], "\n")
	body := struct {
		Head  string `json:"head"`
		Base  string `json:"base"`
		Title string `json:"title"`
		Body  string `json:"body"`
	}{
		Head:  params.Head,
		Base:  params.Base,
		Title: title,
		Body:  description,
	}

	infos, err := tools.RepoInfoFromURL(params.URL, "github.com")
	logger.Info("repo infos : %+v\n", infos)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf(
		"https://%s:%s@api.github.com/repos/%s/%s/pulls",
		me.User,
		me.Token,
		infos.Owner,
		infos.Repo,
	)

	res := new(struct {
		HTMLURL string `json:"html_url"`
	})

	status, err := tools.Request(&tools.RequestParams{
		URL:    url,
		Method: "POST",
		Body:   &body,
		Result: res,
	})
	logger.Debug("res = %+v\n", res)

	if status >= 400 {
		return nil, fmt.Errorf("Request failed with status %d", status)
	}

	return &models.MRInfo{
		URL: res.HTMLURL,
	}, nil
}
