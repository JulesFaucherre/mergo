package hosts

import (
	"fmt"
	"regexp"
	"strings"

	"gitlab.com/jfaucherre/mergo/hosts/github"
	"gitlab.com/jfaucherre/mergo/hosts/gitlab"
	"gitlab.com/jfaucherre/mergo/models"
)

var hosts = map[*regexp.Regexp]func() (models.Host, error){
	regexp.MustCompile("[http://|https://]?[www.]?github[.com]?"): github.NewGithub,
	regexp.MustCompile("[http://|https://]?[www.]?gitlab[.com]?"): gitlab.NewGitlab,
}

func GetHost(host string) (models.Host, error) {
	for h, builder := range hosts {
		if h.MatchString(host) {
			return builder()
		}
	}
	return nil, fmt.Errorf("Not host for %s", host)
}

func GetHostNameFromRemoteString(remote string) string {
	fmt.Println(remote)
	if strings.HasPrefix(remote, "http") {
		return getHostFromHTTP(remote)
	}
	return getHostFromSSH(remote)
}

func getHostFromHTTP(remote string) string {
	splitted := strings.Split(remote, "/")
	return splitted[2]
}

func getHostFromSSH(remote string) string {
	return strings.FieldsFunc(remote, func(r rune) bool {
		return r == ':' || r == '@'
	})[1]
}
