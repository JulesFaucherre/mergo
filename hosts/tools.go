package hosts

import "strings"

func ownerAndRepoSSH(remote string) (string, string) {
	splitted := strings.Split(remote, "/")

	return strings.TrimLeft(splitted[0], ":"), strings.TrimRight(splitted[1], ".")
}

func ownerAndRepoHTTP(remote string) (string, string) {
	splitted := strings.Split(remote, "/")

	return splitted[4], strings.TrimRight(splitted[5], ".")
}

func ownerAndRepo(remote string) (string, string) {
	if strings.HasPrefix("http") {
		return ownerAndRepoHTTP(remote)
	}
	return ownerAndRepoSSH(remote)
}
