package hosts

import (
	"fmt"
	"regexp"

	"gitlab.com/jfaucherre/mergo/hosts/github"
	"gitlab.com/jfaucherre/mergo/hosts/gitlab"
	"gitlab.com/jfaucherre/mergo/models"
)

var hosts = map[*regexp.Regexp]func() (models.Host, error){
	regexp.MustCompile("[http://|https://]?[www.]?github[.com]?"): github.NewGithub,
	regexp.MustCompile("[http://|https://]?[www.]?gitlab[.com]?"): gitlab.NewGitlab,
}

// GetHost returns a models.Host corresponding to the host 'host'
func GetHost(host string) (models.Host, error) {
	for h, builder := range hosts {
		if h.MatchString(host) {
			return builder()
		}
	}
	return nil, fmt.Errorf("Not host for %s", host)
}
