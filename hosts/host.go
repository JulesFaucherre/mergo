package hosts

import (
	"fmt"
	"strings"

	"gitlab.com/jfaucherre/mergo/models"
)

var hosts = map[string]func() Host{
	"github": newGithub,
}

type Host interface {
	SubmitPr(models.Opts) error
	GetOwnerAndRepo(string) (string, string)
}

func getHostFromHTTP(remote string) string {
	splitted := strings.Split(remote, "/")
	return splitted[2]
}

func getHostFromSSH(remote string) string {
	splitted := strings.Split(remote, ":")
	r := strings.TrimRight(splitted[1], "/")
	return r
}

func GetHostNameFromRemoteString(remote string) string {
	if strings.HasPrefix(remote, "http") {
		return getHostFromHTTP(remote)
	}
	return getHostFromSSH(remote)
}

func GetHost(host string) (Host, error) {
	for h, builder := range hosts {
		if h == host {
			return builder(), nil
		}
	}
	return nil, fmt.Errorf("Not host for %s", host)
}
