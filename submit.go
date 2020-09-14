package main

import (
	"errors"
	"regexp"

	"gitlab.com/jfaucherre/mergo/hosts/github"
	"gitlab.com/jfaucherre/mergo/hosts/gitlab"
	"gitlab.com/jfaucherre/mergo/logger"
	"gitlab.com/jfaucherre/mergo/models"
)

var (
	ErrHostNotFound = errors.New("host not found")

	hosts = map[*regexp.Regexp]func() (models.Host, error){
		gitlab.IsHostRegexp: gitlab.New,
		github.IsHostRegexp: github.New,
	}
)

func GetHost(url string) (models.Host, error) {
	logger.Info("searching host for URL: %s\n", url)

	for reg, constructor := range hosts {
		if reg.MatchString(url) {
			return constructor()
		}
	}

	return nil, ErrHostNotFound
}

func SubmitOnRemote(params *models.MRParams) (*models.MRInfo, error) {
	logger.Debug("Submitting PR with parameters: %+v\n", params)
	host, err := GetHost(params.URL)
	if err != nil {
		return nil, err
	}

	infos, err := host.SubmitPr(params)
	if err != nil {
		return nil, err
	}

	return infos, nil
}

func SubmitPr(params *models.MRParams, rmts []string) []string {
	mrUrls := []string{}

	for _, rmt := range rmts {
		params.URL = rmt
		infos, err := SubmitOnRemote(params)
		if err != nil {
			logger.Error("%+v\n", err)
			continue
		}

		mrUrls = append(mrUrls, infos.URL)
	}

	return mrUrls
}
